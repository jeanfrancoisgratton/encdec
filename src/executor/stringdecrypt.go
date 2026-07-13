package executor

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"

	ce "github.com/jeanfrancoisgratton/customError/v3"
)

// ref: https://www.golinuxcloud.com/golang-encrypt-decrypt/#Encryption
func Decode(cryptedString string) (string, *ce.CustomError) {
	if len(SecretKey) != 32 {
		return "", &ce.CustomError{Title: "Unable to decode", Message: fmt.Sprintf("Current key is only %v bytes long. It needs to be of exactly 32 bytes.",
			len(SecretKey))}
	}

	key := []byte(SecretKey)
	ciphertext, _ := base64.URLEncoding.DecodeString(cryptedString)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", &ce.CustomError{Title: "Unable to decode", Message: err.Error()}
	}

	if len(ciphertext) < aes.BlockSize {
		return "", &ce.CustomError{Title: "Unable to decode", Message: "ciphertext too short"}
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return fmt.Sprintf("%s", ciphertext), nil
}
