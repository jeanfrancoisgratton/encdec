// encdec
// Écrit par J.F. Gratton <jean-francois@famillegratton.net>
// Orininal name: src/executor/fileencrypt.go
// Original time: 2023/07/31 07:58

package executor

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"os"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	hftx "github.com/jeanfrancoisgratton/helperFunctions/v4/terminalfx"
)

func EncodeFile(sourcefile, destfile string) error {
	var cerr *ce.CustomError = nil

	if destfile == "" {
		destfile = sourcefile + ".enc"
	}
	//if DEBUG {
	//	fmt.Println("[EncodeFile] source file:", sourcefile)
	//	fmt.Println("[EncodeFile] output file:", destfile)
	//	fmt.Println("[EncodeFile] keep file?", Keep)
	//	fmt.Println("[EncodeFile] quiesce output?", Quiet)
	//}

	cerr = encode(sourcefile, destfile)

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

func encode(source, dest string) *ce.CustomError {
	if PromptForKeys {
		SecretKey = getSecretKey("Please enter a 32 bytes (characters) key: ")
	}
	if len(SecretKey) != 32 {
		fmt.Printf("Current key is only %v bytes long. It needs to be of exactly 32 bytes. Aborting.\n", len(SecretKey))
		os.Exit(1)
	}

	if !Quiet {
		fmt.Println(hftx.InProgressSign(fmt.Sprintf("Enconding %s", source)))
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

	// Create the output file for writing the encrypted data
	outFile, err := os.Create(dest)
	if err != nil {
		return &ce.CustomError{Title: "Error creating the destination file", Message: err.Error()}
	}
	defer outFile.Close()

	// Generate a random IV (Initialization Vector) to use with CFB mode
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return &ce.CustomError{Title: "Error creating a random IV", Message: err.Error()}
	}

	// Write the IV to the beginning of the output file
	outFile.Write(iv)

	// Create a new CFB (Cipher Feedback) encrypter using the block cipher and IV
	stream := cipher.NewCFBEncrypter(block, iv)

	// Create a buffer to hold a chunk of data to be processed
	buf := make([]byte, chunkSize)
	for {
		// Read a chunk of data from the input file
		n, err := inFile.Read(buf)
		if n > 0 {
			// Encrypt the chunk of data in-place using XORKeyStream
			stream.XORKeyStream(buf[:n], buf[:n])

			// Write the encrypted chunk to the output file
			if _, err := outFile.Write(buf[:n]); err != nil {
				return &ce.CustomError{Title: "Error encoding file", Message: err.Error()}
			}
		}
		// Check for the end of file
		if err == io.EOF {
			break
		}
		// Handle other read errors
		if err != nil {
			return &ce.CustomError{Title: "Error encoding file", Message: err.Error()}
		}
	}

	if !Quiet {
		fmt.Println(hftx.EnabledSign(fmt.Sprintf("Succesfully encoded %s as %s\n",
			hftx.Green(source), hftx.Green(dest))))
	}
	return nil
}
