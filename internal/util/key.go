package util

import (
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"log"

	"golang.org/x/crypto/ripemd160"

	"github.com/btcsuite/btcutil/bech32"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcutil/base58"
)

func Hash256(msg []byte) []byte {
	hasher1 := sha256.New()
	hasher2 := sha256.New()
	hasher1.Write(msg)
	hasher2.Write(hasher1.Sum(nil))
	return hasher2.Sum(nil)
}

func Hash160(msg []byte) []byte {
	hasher1 := sha256.New()
	hasher2 := ripemd160.New()
	hasher1.Write(msg)
	hasher2.Write(hasher1.Sum(nil))
	return hasher2.Sum(nil)
}

func Gen_key(len int) []byte {
	buf := make([]byte, len)
	_, err := rand.Read(buf)
	if err != nil {
		panic(err)
	}
	return buf
}

func Bytes_to_wif(priv []byte) string {
	head := []byte{0x80}
	var tail byte = 0x01
	priv_ex := append(append(head, priv...), tail)
	checksum := Hash256(priv_ex)[0:4]
	priv_check := append(priv_ex, checksum...)
	wif := B58_encode(priv_check)
	return wif
}

func New_keypair() (*btcec.PrivateKey, *btcec.PublicKey) {
	privKey, _ := btcec.NewPrivateKey()
	pubKey := privKey.PubKey()
	return privKey, pubKey
}

func Bytes_to_keypair(priv []byte) (*btcec.PrivateKey, *btcec.PublicKey) {
	privKey, pubKey := btcec.PrivKeyFromBytes(priv)
	return privKey, pubKey
}

func Pubkey_to_address_b32(pub []byte) string {
	pub_bytes_hashed := Hash160(pub)
	head := []byte{0}
	data5, err := bech32.ConvertBits(pub_bytes_hashed, 8, 5, true)
	if err != nil {
		log.Fatal(err)
	}
	pub_bytes_to_encode := append(head, data5...)
	address_bech32, err := bech32.Encode("bc", pub_bytes_to_encode)
	return address_bech32
}

func New_priv_address_b32() ([]byte, string) {
	privkey, pubkey := New_keypair()
	priv := privkey.Serialize()
	pub := pubkey.SerializeCompressed()
	address := Pubkey_to_address_b32(pub)
	return priv, address
}

func Address_type(address string) (string, error) {
	if address[:4] == "bc1q" {
		_, _, err := bech32.Decode(address)
		if err != nil {
			return "", err
		}
		return "p2wpkh", nil
	} else if address[0] == '1' {
		_, _, err := base58.CheckDecode(address)
		if err != nil {
			return "", err
		}
		return "p2pkh", nil
	}
	return "", errors.New("このアドレス形式には対応していません。")
}
