package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"logarda/internal/config"
)

func EncryptString(word string) string {
    text := []byte(word) // convert the word into each ASCII bits

    key := []byte(config.ENCRYPTION_KEY) // same here with encryption key

    c, err := aes.NewCipher(key) // make the AES engine with encryption key
    if err != nil {
        fmt.Println(err)
    }

    gcm, err := cipher.NewGCM(c) // use GCM for additional security and wrap the AES engine
    if err != nil {
        fmt.Println(err)
    }


    nonce := make([]byte, gcm.NonceSize()) // 12-bytes nonce (Number used once) 

    if _, err = io.ReadFull(rand.Reader, nonce); err != nil { // fill the nonce slice with random numbers (to allow different encryptions for different runs)
        fmt.Println(err)
    }

    ciphertext := gcm.Seal(nonce, nonce, text, nil) 
	encryptedString := base64.StdEncoding.EncodeToString(ciphertext) // convert to string

	return encryptedString
}


func HashString(word string) string {
    hash := sha256.Sum256([]byte(word)) // converts into 32 bytes array

    return hex.EncodeToString(hash[:])
}

