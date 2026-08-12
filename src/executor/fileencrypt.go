// encdec
// Écrit par J.F. Gratton <jean-francois@famillegratton.net>
// Orininal name: src/executor/fileencrypt.go
// Original time: 2023/07/31 07:58

package executor

import (
	"fmt"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	hf "github.com/jeanfrancoisgratton/helperFunctions/v5"
	hftx "github.com/jeanfrancoisgratton/helperFunctions/v5/terminalfx"
)

// EncodeFile encrypts sourcefile. With no destfile it works in place: the
// ciphertext goes to a scratch file that then replaces the source, unless -k
// asks for the source to be kept. With a destfile, the ciphertext goes there
// and the source is left alone -- -k has nothing to do in that mode.
func EncodeFile(sourcefile, destfile string) *ce.CustomError {
	inPlace := destfile == ""
	if inPlace {
		destfile = sourcefile + ".enc"
	}

	if cerr := assertWritable(destfile); cerr != nil {
		return cerr
	}

	if cerr := encode(sourcefile, destfile); cerr != nil {
		return cerr
	}

	if !inPlace || Keep {
		return nil
	}
	return replaceSource(sourcefile, destfile)
}

func encode(source, dest string) *ce.CustomError {
	if !Quiet {
		fmt.Println(hftx.InProgressSign(fmt.Sprintf("Enconding %s", source)))
	}

	if err := hf.EncodeFile(source, dest, Passphrase); err != nil {
		return &ce.CustomError{Title: "Error encoding file", Message: err.Error()}
	}

	if !Quiet {
		fmt.Println(hftx.EnabledSign(fmt.Sprintf("Succesfully encoded %s as %s\n",
			hftx.Green(source), hftx.Green(dest))))
	}
	return nil
}
