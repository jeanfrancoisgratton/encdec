package executor

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"

	ce "github.com/jeanfrancoisgratton/customError/v3"
)

// ref: https://www.golinuxcloud.com/golang-encrypt-decrypt/#Encryption

func Encode(string2encrypt string) (string, *ce.CustomError) {
	if len(SecretKey) != 32 {
		return "", &ce.CustomError{Title: "Unable to encode", Message: fmt.Sprintf("Current key is only %v bytes long. It needs to be of exactly 32 bytes.",
			len(SecretKey))}
	}

	key := []byte(SecretKey)
	plaintext := []byte(string2encrypt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", &ce.CustomError{Title: "Unable to encode", Message: err.Error()}
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", &ce.CustomError{Title: "Unable to encode", Message: err.Error()}
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}
