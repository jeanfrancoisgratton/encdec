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
)

func DecodeFile(sourcefile, destfile string) error {
	var err error = nil

	if destfile == "" {
		destfile = sourcefile + ".dec"
	}
	//if DEBUG {
	//	fmt.Printf("[DecodeFile] source file %s\n", sourcefile)
	//	fmt.Printf("[DecodeFile] output file %s\n", destfile)
	//	fmt.Printf("[DecodeFile] keep file ? %\n", Keep)
	//	fmt.Printf("[DecodeFile] keep file ? % %\n", Quiet)
	//}
	if err = decode(sourcefile, destfile); err != nil {
		return err
	}

	if !Keep {
		if err = os.Remove(sourcefile); err != nil {
			return err
		}
		err = os.Rename(destfile, sourcefile)
	}

	return err
}

func decode(source, dest string) error {
	if PromptForKeys {
		SecretKey = getSecretKey("Please enter a 32 bytes (characters) key: ")
	}
	if len(SecretKey) != 32 {
		fmt.Printf("Current key is only %v bytes long. It needs to be of exactly 32 bytes. Aborting.\n", len(SecretKey))
		os.Exit(1)
	}
	if !Quiet {
		fmt.Println("Decoding ", source)
	}
	key := []byte(SecretKey)

	// Create a new AES cipher block based on the provided encryption key
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	// Open the input file for reading
	inFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer inFile.Close()

	// Create the output file for writing the decrypted data
	outFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// Read the IV (Initialization Vector) from the beginning of the input file
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(inFile, iv); err != nil {
		return err
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
				return err
			}
		}
		// Check for the end of file
		if err == io.EOF {
			break
		}
		// Handle other read errors
		if err != nil {
			return err
		}
	}

	if !Quiet {
		fmt.Printf("Succesfully decoded %s as %s\n", source, dest)
	}
	return nil
}
