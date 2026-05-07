// Package auditor provides a lightweight append-only audit log for envsync
// operations. Each sync or diff action can be recorded as a JSON-lines entry
// containing a timestamp, event type, source/target paths, and a change count.
//
// Usage:
//
//	a := auditor.New(".envsync-audit.log")
//	err := a.Record(auditor.Entry{
//		Event:   auditor.EventSync,
//		Source:  ".env.production",
//		Target:  ".env.local",
//		Changes: 3,
//	})
//
//	entries, err := a.ReadAll()
package auditor
