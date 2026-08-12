// encdec
// Écrit par J.F. Gratton <jean-francois@famillegratton.net>
// Orininal name: src/executor/stringdecrypt.go

package executor

import (
	"fmt"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	hf "github.com/jeanfrancoisgratton/helperFunctions/v5"
)

// Decode deciphers a string produced by Encode(); helperFunctions panics on
// malformed input (a ciphertext shorter than one AES block, typically), so we
// catch that and return a CustomError instead.
func Decode(cryptedString string) (decoded string, cerr *ce.CustomError) {
	defer func() {
		if r := recover(); r != nil {
			decoded = ""
			cerr = &ce.CustomError{Title: "Unable to decode", Message: fmt.Sprintf("%v", r)}
		}
	}()

	return hf.DecodeString(cryptedString, Passphrase), nil
}
