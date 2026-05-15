// Package patcher applies a set of patch operations (set, delete, rename) to an
// env map, returning a new map and a structured result describing every change.
package patcher

import "fmt"

// OpKind is the kind of patch operation.
type OpKind string

const (
	OpSet    OpKind = "set"
	OpDelete OpKind = "delete"
	OpRename OpKind = "rename"
)

// Op describes a single patch operation.
type Op struct {
	Kind    OpKind
	Key     string // source key
	Value   string // used by OpSet
	NewKey  string // used by OpRename
}

// Change records what happened to a single key.
type Change struct {
	Op       OpKind
	Key      string
	OldValue string
	NewValue string
}

// Result holds the patched env map and the list of applied changes.
type Result struct {
	Env     map[string]string
	Changes []Change
}

// Apply executes ops against src (which is never mutated) and returns a Result.
// Operations are applied in order; referencing a non-existent key is an error
// unless the operation is OpSet.
func Apply(src map[string]string, ops []Op) (Result, error) {
	out := make(map[string]string, len(src))
	for k, v := range src {
		out[k] = v
	}

	var changes []Change

	for _, op := range ops {
		switch op.Kind {
		case OpSet:
			old := out[op.Key]
			out[op.Key] = op.Value
			changes = append(changes, Change{Op: OpSet, Key: op.Key, OldValue: old, NewValue: op.Value})

		case OpDelete:
			old, ok := out[op.Key]
			if !ok {
				return Result{}, fmt.Errorf("patcher: delete: key %q not found", op.Key)
			}
			delete(out, op.Key)
			changes = append(changes, Change{Op: OpDelete, Key: op.Key, OldValue: old})

		case OpRename:
			if op.NewKey == "" {
				return Result{}, fmt.Errorf("patcher: rename: NewKey must not be empty")
			}
			old, ok := out[op.Key]
			if !ok {
				return Result{}, fmt.Errorf("patcher: rename: key %q not found", op.Key)
			}
			delete(out, op.Key)
			out[op.NewKey] = old
			changes = append(changes, Change{Op: OpRename, Key: op.Key, NewValue: op.NewKey, OldValue: old})

		default:
			return Result{}, fmt.Errorf("patcher: unknown op kind %q", op.Kind)
		}
	}

	return Result{Env: out, Changes: changes}, nil
}
