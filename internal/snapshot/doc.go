// Package snapshot captures and compares .env file states over time.
//
// A Snapshot records all key-value pairs from an .env file together with
// a timestamp. Snapshots can be persisted to disk as JSON and reloaded
// later to compute a Delta — showing which keys were added, removed, or
// changed between two points in time.
//
// Typical usage:
//
//	s := snapshot.Capture(".env", entries)
//	_ = snapshot.Save(s, ".env.snap")
//
//	old, _ := snapshot.Load(".env.snap")
//	delta := snapshot.Compare(old, s)
package snapshot
