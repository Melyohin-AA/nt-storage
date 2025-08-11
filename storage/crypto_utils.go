package storage

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

func checkKeySize(key []byte) error {
	if len(key) != aes.BlockSize {
		return fmt.Errorf("expected key to be %d bytes long but was %d", aes.BlockSize, len(key))
	}
	return nil
}

func EncryptR(key []byte, source io.Reader) (io.Reader, error) {
	if err := checkKeySize(key); err != nil {
		return nil, err
	}
	// IV
	iv := make([]byte, aes.BlockSize)
	rand.Read(iv)
	// Payload
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}
	stream := cipher.NewCTR(block, iv)
	return io.MultiReader(bytes.NewReader(iv), cipher.StreamReader{S: stream, R: source}), nil
}

func Encrypt(key []byte, dest io.Writer, source io.Reader) error {
	if err := checkKeySize(key); err != nil {
		return err
	}
	// IV
	iv := make([]byte, aes.BlockSize)
	rand.Read(iv)
	dest.Write(iv)
	// Payload
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return err
	}
	stream := cipher.NewCTR(block, iv)
	_, err = io.Copy(cipher.StreamWriter{S: stream, W: dest}, source)
	return err
}

func Decrypt(key []byte, dest io.Writer, source io.Reader) error {
	if err := checkKeySize(key); err != nil {
		return err
	}
	// IV
	iv := make([]byte, aes.BlockSize)
	if err := ReadFixed(source, iv); err != nil {
		return err
	}
	// Payload
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return err
	}
	stream := cipher.NewCTR(block, iv)
	_, err = io.Copy(dest, cipher.StreamReader{S: stream, R: source})
	return err
}
