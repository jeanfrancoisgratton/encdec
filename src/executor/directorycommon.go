// encdec
// Written by Jean-François Gratton <jean-francois.gratton@aylo.com>
// Original timestamp : 2026.09.03 14:53:03
// Original filename : src/executor/directorycommon.go

package executor

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	ce "github.com/jeanfrancoisgratton/customError/v3"
)

// processDirectory walks rootdir recursively and applies process to every
// regular file. An empty rootdir means the current working directory.
// Symbolic links and other special files are deliberately left untouched.
func processDirectory(rootdir, operation string, process func(string, string) *ce.CustomError) *ce.CustomError {
	if rootdir == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return &ce.CustomError{Title: "Error getting the current directory", Message: err.Error()}
		}
		rootdir = cwd
	}

	resolvedRoot, err := filepath.EvalSymlinks(rootdir)
	if err != nil {
		return &ce.CustomError{Title: "Error checking the directory", Message: err.Error()}
	}
	rootdir = resolvedRoot

	info, err := os.Stat(rootdir)
	if err != nil {
		return &ce.CustomError{Title: "Error checking the directory", Message: err.Error()}
	}
	if !info.IsDir() {
		return &ce.CustomError{Title: "Invalid directory", Message: fmt.Sprintf("%s is not a directory", rootdir)}
	}

	var operationError *ce.CustomError
	walkError := filepath.WalkDir(rootdir, func(path string, entry fs.DirEntry, walkError error) error {
		if walkError != nil {
			return walkError
		}
		if entry.IsDir() {
			return nil
		}

		entryInfo, err := entry.Info()
		if err != nil {
			return err
		}
		if !entryInfo.Mode().IsRegular() {
			return nil
		}

		operationError = process(path, "")
		if operationError != nil {
			return operationError
		}
		return nil
	})
	if operationError != nil {
		return operationError
	}
	if walkError != nil {
		return &ce.CustomError{
			Title:   fmt.Sprintf("Error %s directory", operation),
			Message: walkError.Error(),
		}
	}
	return nil
}
