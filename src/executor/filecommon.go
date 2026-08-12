// encdec
// Écrit par J.F. Gratton <jean-francois@famillegratton.net>
// Orininal name: src/executor/filecommon.go

package executor

import (
	"fmt"
	"os"

	ce "github.com/jeanfrancoisgratton/customError/v3"
)

// assertWritable refuses to clobber an existing destination unless -F was given.
// It guards the scratch file of an in-place run just as much as a destination
// named on the command line: a leftover foo.enc from an aborted run is still
// somebody's data.
func assertWritable(destfile string) *ce.CustomError {
	if Force {
		return nil
	}

	_, err := os.Stat(destfile)
	switch {
	case err == nil:
		return &ce.CustomError{Title: "Destination file already exists",
			Message: fmt.Sprintf("%s already exists; pass -F (--force) to overwrite it", destfile)}
	case !os.IsNotExist(err):
		return &ce.CustomError{Title: "Error checking the destination file", Message: err.Error()}
	}
	return nil
}

// replaceSource completes an in-place run: the scratch file we just wrote takes
// the place of the original. This is only ever called when the destination was
// derived rather than given by the user -- an explicit destination is where the
// result belongs, and moving it elsewhere would be a surprise.
func replaceSource(sourcefile, scratchfile string) *ce.CustomError {
	if err := os.Remove(sourcefile); err != nil {
		return &ce.CustomError{Title: "Error removing the source file", Message: err.Error()}
	}
	if err := os.Rename(scratchfile, sourcefile); err != nil {
		return &ce.CustomError{Title: "Error renaming the destination file", Message: err.Error()}
	}
	return nil
}
