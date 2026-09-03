// encdec
// Written by Jean-François Gratton <jean-francois.gratton@aylo.com>
// Original timestamp : 2026.09.03 14:53:45
// Original filename : src/executor/directoryencrypt.go

package executor

import ce "github.com/jeanfrancoisgratton/customError/v3"

// EncodeDirectory recursively encodes every regular file below rootdir using
// the same in-place semantics and flags as EncodeFile. When rootdir is empty,
// the walk starts in the current working directory.
func EncodeDirectory(rootdir string) *ce.CustomError {
	return processDirectory(rootdir, "encoding", EncodeFile)
}
