// encdec
// Written by Jean-François Gratton <jean-francois.gratton@aylo.com>
// Original timestamp : 2026.09.03 14:54:19
// Original filename : src/executor/directorycrypt.go

package executor

import ce "github.com/jeanfrancoisgratton/customError/v3"

// DecodeDirectory recursively decodes every regular file below rootdir using
// the same in-place semantics and flags as DecodeFile. When rootdir is empty,
// the walk starts in the current working directory.
func DecodeDirectory(rootdir string) *ce.CustomError {
	return processDirectory(rootdir, "decoding", DecodeFile)
}
