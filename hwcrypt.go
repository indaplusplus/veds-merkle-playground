package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"log"
	"fmt"
	"io"
)

//https://golang.org/src/crypto/cipher/example_test.go
func Encrypt(data []byte, key []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatal(err)
	}

	cipher_text := make([]byte, aes.BlockSize + len(data))
	iv := cipher_text[:aes.BlockSize]
	_, err = io.ReadFull(rand.Reader, iv);
	if err != nil {
		log.Panic(err)
	}

	cip := cipher.NewCFBEncrypter(block, iv)
	cip.XORKeyStream(cipher_text[aes.BlockSize:], data)

	return cipher_text
}

func Decrypt(data []byte, key []byte) []byte{
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatal(err)
	}
	iv := data[:aes.BlockSize]
	cipher_text := data[aes.BlockSize:]

	cip := cipher.NewCFBDecrypter(block, iv)
	cip.XORKeyStream(cipher_text, cipher_text)
	return cipher_text
}

func main() {
	data := []byte{42, 69}
	key := []byte{11,22,33,44,11,22,33,44,11,22,33,44,11,22,33,44}
	ec := Encrypt(data, key)
	dc := Decrypt(ec, key)
	fmt.Println(dc)
}
