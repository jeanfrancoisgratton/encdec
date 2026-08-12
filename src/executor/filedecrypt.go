// encdec
// Écrit par J.F. Gratton <jean-francois@famillegratton.net>
// Orininal name: src/executor/filedecrypt.go
// Original time: 2023/07/31 07:58

package executor

import (
	"fmt"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	hf "github.com/jeanfrancoisgratton/helperFunctions/v5"
	hftx "github.com/jeanfrancoisgratton/helperFunctions/v5/terminalfx"
)

// DecodeFile mirrors EncodeFile: in place when no destfile is given, to the
// named destination -- source untouched -- when there is one.
func DecodeFile(sourcefile, destfile string) *ce.CustomError {
	inPlace := destfile == ""
	if inPlace {
		destfile = sourcefile + ".dec"
	}

	if cerr := assertWritable(destfile); cerr != nil {
		return cerr
	}

	if cerr := decode(sourcefile, destfile); cerr != nil {
		return cerr
	}

	if !inPlace || Keep {
		return nil
	}
	return replaceSource(sourcefile, destfile)
}

func decode(source, dest string) *ce.CustomError {
	if !Quiet {
		fmt.Println(hftx.InProgressSign("Decoding " + source))
	}

	if err := hf.DecodeFile(source, dest, Passphrase); err != nil {
		return &ce.CustomError{Title: "Error decoding file", Message: err.Error()}
	}

	if !Quiet {
		fmt.Println(hftx.InfoSign(fmt.Sprintf("Successfully decoded %s to %s",
			hftx.Green(source), hftx.Green(dest))))
	}
	return nil
}
