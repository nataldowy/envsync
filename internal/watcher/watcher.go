// Package watcher monitors .env files for changes and triggers a callback
// when modifications are detected. It is useful for long-running processes
// that need to react to environment configuration updates without restarting.
package watcher

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// Event describes a file-change notification produced by the Watcher.
type Event struct {
	// Path is the absolute or relative path of the file that changed.
	Path string
	// OldHash is the SHA-256 hex digest of the file before the change.
	OldHash string
	// NewHash is the SHA-256 hex digest of the file after the change.
	NewHash string
	// At is the time the change was detected.
	At time.Time
}

// Handler is a function invoked whenever a watched file changes.
type Handler func(Event)

// Watcher polls one or more files at a configurable interval and calls
// the registered Handler when a file's content hash changes.
type Watcher struct {
	mu       sync.Mutex
	files    map[string]string // path -> last known hash
	interval time.Duration
	handler  Handler
	stop     chan struct{}
	wg       sync.WaitGroup
}

// New creates a Watcher that checks files every interval duration.
// interval must be positive; if it is zero or negative New defaults to 2s.
func New(interval time.Duration, handler Handler) *Watcher {
	if interval <= 0 {
		interval = 2 * time.Second
	}
	return &Watcher{
		files:    make(map[string]string),
		interval: interval,
		handler:  handler,
		stop:     make(chan struct{}),
	}
}

// Add registers a file path to be watched. It is safe to call Add before or
// after Start. If the file does not exist yet, it will be tracked once it
// appears.
func (w *Watcher) Add(path string) error {
	hash, err := hashFile(path)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("watcher: initial hash of %q: %w", path, err)
	}
	w.mu.Lock()
	w.files[path] = hash
	w.mu.Unlock()
	return nil
}

// Start begins the polling loop in a background goroutine.
// Calling Start more than once without a preceding Stop is a no-op.
func (w *Watcher) Start() {
	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		ticker := time.NewTicker(w.interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				w.poll()
			case <-w.stop:
				return
			}
		}
	}()
}

// Stop shuts down the polling loop and waits for it to exit.
func (w *Watcher) Stop() {
	close(w.stop)
	w.wg.Wait()
}

// poll checks every registered file and fires the handler on any change.
func (w *Watcher) poll() {
	w.mu.Lock()
	paths := make([]string, 0, len(w.files))
	for p := range w.files {
		paths = append(paths, p)
	}
	w.mu.Unlock()

	for _, path := range paths {
		newHash, err := hashFile(path)
		if err != nil && !os.IsNotExist(err) {
			// Transient read error – skip this cycle.
			continue
		}

		w.mu.Lock()
		oldHash := w.files[path]
		if newHash != oldHash {
			w.files[path] = newHash
			w.mu.Unlock()
			if w.handler != nil {
				w.handler(Event{
					Path:    path,
					OldHash: oldHash,
					NewHash: newHash,
					At:      time.Now(),
				})
			}
		} else {
			w.mu.Unlock()
		}
	}
}

// hashFile returns a hex-encoded SHA-256 digest of the named file's contents.
// It returns an empty string when the file does not exist.
func hashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
