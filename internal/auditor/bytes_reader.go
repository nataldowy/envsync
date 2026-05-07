package auditor

import "bytes"

// newBytesReader wraps a byte slice in a bytes.Reader, satisfying
// the io.Reader interface required by json.NewDecoder.
func newBytesReader(b []byte) *bytes.Reader {
	return bytes.NewReader(b)
}
