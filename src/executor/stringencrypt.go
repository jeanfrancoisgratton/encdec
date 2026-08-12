// encdec
// Écrit par J.F. Gratton <jean-francois@famillegratton.net>
// Orininal name: src/executor/stringencrypt.go

package executor

import (
	"fmt"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	hf "github.com/jeanfrancoisgratton/helperFunctions/v5"
)

// Encode ciphers a string with Passphrase; the actual work is done by
// helperFunctions, we only translate its panics into a CustomError.
func Encode(string2encrypt string) (encoded string, cerr *ce.CustomError) {
	defer func() {
		if r := recover(); r != nil {
			encoded = ""
			cerr = &ce.CustomError{Title: "Unable to encode", Message: fmt.Sprintf("%v", r)}
		}
	}()

	return hf.EncodeString(string2encrypt, Passphrase), nil
}
