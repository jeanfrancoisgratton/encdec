// encdec
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// src/executor/types.go :
// Original file timestamp: 2024.08.11 07:42:43

package executor

/*	For now, we only declare the global variables we need in this tool
	Eventually, all functions will be wrapped around type structs and function receivers
*/

var Quiet = false

// Passphrase is not a key: helperFunctions derives the actual 32-byte AES-256
// key from its SHA256 sum. Any length is therefore valid, and the empty
// passphrase -- the default, when -s is not given -- is a legitimate one.
var Passphrase = ""

// Keep only applies to in-place runs (no destination file given): it leaves the
// source alone instead of replacing it with the result.
var Keep = false

// Force allows an existing destination file to be overwritten.
var Force = false

var FileOps = false
var DirectoryOps = false
var DEBUG = false
