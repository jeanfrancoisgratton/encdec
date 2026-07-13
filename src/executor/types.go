// encdec
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// src/executor/types.go :
// Original file timestamp: 2024.08.11 07:42:43

package executor

/*	For now, we only declare the global variables we need in this tool
	Eventually, all functions will be wrapped around type structs and function receivers
*/

var Quiet = false
var SecretKey = "secret key 2 encrypt and decrypt"
var Keep = false
var FileOps = false
var DEBUG = false

const chunkSize = 64 * 1024 // 64 KB chunk size
