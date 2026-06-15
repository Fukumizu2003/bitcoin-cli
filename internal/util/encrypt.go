package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"log"

	"golang.org/x/crypto/scrypt"
)

func scrypt_hash_new(pw []byte) ([]byte, []byte) {
	salt := Gen_key(16) // Always use a unique salt
	N := 16384          // CPU/memory cost parameter
	r := 8              // Block size
	p := 1              // Parallelization factor
	keyLen := 32
	key, err := scrypt.Key(pw, salt, N, r, p, keyLen)
	if err != nil {
		log.Fatalf("Error generating key: %v", err)
	}
	return salt, key
}

func scrypt_hash_again(salt []byte, pw []byte) []byte {
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

func Aes_encrypt(priv []byte, pw []byte) []byte {
	salt, key := scrypt_hash_new(pw)
	block, _ := aes.NewCipher(key)
	aesgcm, _ := cipher.NewGCM(block)

	nonce := Gen_key(12)
	rand.Read(nonce)
	ciphertext := aesgcm.Seal(nil, nonce, priv, nil)
	ans := append(append(salt, nonce...), ciphertext...)
	return ans
}

func Aes_decrypt(priv []byte, pw []byte) []byte {
	salt := priv[:16]
	nonce := priv[16:28]
	key := scrypt_hash_again(salt, pw)
	block, _ := aes.NewCipher(key)
	aesgcm, _ := cipher.NewGCM(block)

	decrypted, _ := aesgcm.Open(nil, nonce, priv[28:], nil)
	return decrypted
}
