package utils

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "fmt"
    "io"
	"encoding/base64"
)

func EncryptString(word string) {
    fmt.Println("Encryption Program v0.01")

    text := []byte(word)
    key := []byte("passphrasewhichneedstobe32bytes!")

    c, err := aes.NewCipher(key) // make the AES engine
    if err != nil {
        fmt.Println(err)
    }

    gcm, err := cipher.NewGCM(c) // use GCM for additional security
    if err != nil {
        fmt.Println(err)
    }


    nonce := make([]byte, gcm.NonceSize())

    if _, err = io.ReadFull(rand.Reader, nonce); err != nil { // fill the nonce slice with random numbers
        fmt.Println(err)
    }

    ciphertext := gcm.Seal(nonce, nonce, text, nil) 
	encryptedString := base64.StdEncoding.EncodeToString(ciphertext) // convert to string

	fmt.Println(encryptedString)
}