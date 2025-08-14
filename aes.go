package main

import (
	"crypto/aes"
	"encoding/hex"
)

// Key for AES 256 encryption/decryption: fdsl;mewrjope456fds4fbvfnjwaugfo
const keyHex = "6664736c3b6d6577726a6f706534353666647334666276666e6a77617567666f"

func aesEcb256Encrypt(input string) string {
	bs := 16
	nopadplaintext := []byte(input)
	pad := bs - len(nopadplaintext)%bs

	var plaintext []byte
	if pad != 16 {
		paddedlength := len(nopadplaintext) + pad
		paddedplaintext := make([]byte, paddedlength)
		copy(paddedplaintext, nopadplaintext)
		plaintext = paddedplaintext
	} else {
		plaintext = []byte(input)
	}

	key, _ := hex.DecodeString(keyHex)
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	ciphertext := make([]byte, len(plaintext))
	be := 0
	for len(plaintext) > 0 {
		block.Encrypt(ciphertext[be:], plaintext)
		plaintext = plaintext[bs:]
		be += bs
	}

	return hex.EncodeToString(ciphertext)
}

func aesEcb256Decrypt(input string) string {
	bs := 16
	nopadciphertext, _ := hex.DecodeString(input)
	pad := bs - len(nopadciphertext)%bs

	var ciphertext []byte
	if pad != 16 {
		paddedlength := len(nopadciphertext) + pad
		paddedciphertext := make([]byte, paddedlength)
		copy(paddedciphertext, nopadciphertext)
		ciphertext = paddedciphertext
	} else {
		ciphertext, _ = hex.DecodeString(input)
	}

	key, _ := hex.DecodeString(keyHex)
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	plaintext := make([]byte, len(ciphertext))
	be := 0
	for len(ciphertext) > 0 {
		block.Decrypt(plaintext[be:], ciphertext)
		ciphertext = ciphertext[bs:]
		be += bs
	}

	return string(plaintext)
}
