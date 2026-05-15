// Package patcher provides a lightweight patch engine for env maps.
//
// A patch is a slice of [Op] values that are applied in order against a source
// env map.  Three operation kinds are supported:
//
//   - [OpSet]    – create or overwrite a key with a new value.
//   - [OpDelete] – remove a key (returns an error if the key is absent).
//   - [OpRename] – rename a key while preserving its value.
//
// The source map is never mutated; [Apply] always returns a fresh copy.
// [FormatResult] and [Summary] can be used to render human-readable output.
package patcher
