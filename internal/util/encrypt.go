package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"log"

	"golang.org/x/crypto/scrypt"
)

func scryptHashNew(pw []byte) ([]byte, []byte) {
	salt := GenKey(16) // Always use a unique salt
	N := 16384         // CPU/memory cost parameter
	r := 8             // Block size
	p := 1             // Parallelization factor
	keyLen := 32
	key, err := scrypt.Key(pw, salt, N, r, p, keyLen)
	if err != nil {
		log.Fatalf("Error generating key: %v", err)
	}
	return salt, key
}

func scryptHashAgain(salt []byte, pw []byte) []byte {
	N := 16384 // CPU/memory cost parameter
	r := 8     // Block size
	p := 1     // Parallelization factor
	keyLen := 32
	key, err := scrypt.Key(pw, salt, N, r, p, keyLen)
	if err != nil {
		log.Fatalf("Error generating key: %v", err)
	}
	return key
}

func AesEncrypt(priv []byte, pw []byte) []byte {
	salt, key := scryptHashNew(pw)
	block, _ := aes.NewCipher(key)
	aesgcm, _ := cipher.NewGCM(block)

	nonce := GenKey(12)
	rand.Read(nonce)
	ciphertext := aesgcm.Seal(nil, nonce, priv, nil)
	ans := append(append(salt, nonce...), ciphertext...)
	return ans
}

func AesDecrypt(priv []byte, pw []byte) []byte {
	salt := priv[:16]
	nonce := priv[16:28]
	key := scryptHashAgain(salt, pw)
	block, _ := aes.NewCipher(key)
	aesgcm, _ := cipher.NewGCM(block)

	decrypted, _ := aesgcm.Open(nil, nonce, priv[28:], nil)
	return decrypted
}
