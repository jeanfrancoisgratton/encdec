// encdec
// Écrit par J.F. Gratton <jean-francois@famillegratton.net>
// Orininal name: src/executor/filedecrypt.go
// Original time: 2023/07/31 07:58

package executor

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"io"
	"os"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	hftx "github.com/jeanfrancoisgratton/helperFunctions/v5/terminalfx"
)

func DecodeFile(sourcefile, destfile string) *ce.CustomError {
	//var err error = nil
	var cerr *ce.CustomError = nil

	if destfile == "" {
		destfile = sourcefile + ".dec"
	}
	//if DEBUG {
	//	fmt.Printf("[DecodeFile] source file %s\n", sourcefile)
	//	fmt.Printf("[DecodeFile] output file %s\n", destfile)
	//	fmt.Printf("[DecodeFile] keep file ? %\n", Keep)
	//	fmt.Printf("[DecodeFile] keep file ? % %\n", Quiet)
	//}
	cerr = decode(sourcefile, destfile)

	if !Keep && cerr == nil {
		if err := os.Remove(sourcefile); err != nil {
			cerr = &ce.CustomError{Title: "Error removing the source file", Message: err.Error()}
		}
		if err := os.Rename(destfile, sourcefile); err != nil {
			cerr = &ce.CustomError{Title: "Error renaming the destination file", Message: err.Error()}
		}
	}
	return cerr
}

func decode(source, dest string) *ce.CustomError {
	if PromptForKeys {
		SecretKey = getSecretKey("Please enter a 32 bytes (characters) key: ")
	}
	if len(SecretKey) != 32 {
		fmt.Printf("Current key is only %v bytes long. It needs to be of exactly 32 bytes. Aborting.\n", len(SecretKey))
		fmt.Println(hftx.InfoSign(fmt.Sprintf("%s %s", hftx.Red("ATTEMPTING TO"),
			hftx.Yellow("decode the file with the default, hardcoded key"))))
	}
	if !Quiet {
		fmt.Println(hftx.InProgressSign("Decoding " + source))
	}
	key := []byte(SecretKey)

	// Create a new AES cipher block based on the provided encryption key
	block, err := aes.NewCipher(key)
	if err != nil {
		return &ce.CustomError{Title: "Error creating AES cipher", Message: err.Error()}
	}

	// Open the input file for reading
	inFile, err := os.Open(source)
	if err != nil {
		return &ce.CustomError{Title: "Error opening file", Message: err.Error()}
	}
	defer inFile.Close()

	// Create the output file for writing the decrypted data
	outFile, err := os.Create(dest)
	if err != nil {
		return &ce.CustomError{Title: "Error creating file", Message: err.Error()}
	}
	defer outFile.Close()

	// Read the IV (Initialization Vector) from the beginning of the input file
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(inFile, iv); err != nil {
		return &ce.CustomError{Title: "Error reading file", Message: err.Error()}
	}

	// Create a new CFB (Cipher Feedback) decrypter using the block cipher and IV
	stream := cipher.NewCFBDecrypter(block, iv)

	// Create a buffer to hold a chunk of data to be processed
	buf := make([]byte, chunkSize)
	for {
		// Read a chunk of data from the input file
		n, err := inFile.Read(buf)
		if n > 0 {
			// Decrypt the chunk of data in-place using XORKeyStream
			stream.XORKeyStream(buf[:n], buf[:n])

			// Write the decrypted chunk to the output file
			if _, err := outFile.Write(buf[:n]); err != nil {
				return &ce.CustomError{Title: "Error writing file", Message: err.Error()}
			}
		}
		// Check for the end of file
		if err == io.EOF {
			break
		}
		// Handle other read errors
		if err != nil {
			return &ce.CustomError{Title: "Error reading file", Message: err.Error()}
		}
	}

	if !Quiet {
		fmt.Println(hftx.InfoSign(fmt.Sprintf("Successfully decoded %s to %s",
			hftx.Green(source), hftx.Green(dest))))
	}
	return nil
}
