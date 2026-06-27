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

func GenKey(len int) []byte {
	buf := make([]byte, len)
	_, err := rand.Read(buf)
	if err != nil {
		panic(err)
	}
	return buf
}

func BytesToWif(priv []byte) string {
	head := []byte{0x80}
	var tail byte = 0x01
	privEx := append(append(head, priv...), tail)
	checksum := Hash256(privEx)[0:4]
	privCheck := append(privEx, checksum...)
	wif := B58Encode(privCheck)
	return wif
}

func NewKeypair() (*btcec.PrivateKey, *btcec.PublicKey) {
	privKey, _ := btcec.NewPrivateKey()
	pubKey := privKey.PubKey()
	return privKey, pubKey
}

func BytesToKeypair(priv []byte) (*btcec.PrivateKey, *btcec.PublicKey) {
	privKey, pubKey := btcec.PrivKeyFromBytes(priv)
	return privKey, pubKey
}

func PubkeyToAddressB32(pub []byte) string {
	pubBytesHashed := Hash160(pub)
	head := []byte{0}
	data5, err := bech32.ConvertBits(pubBytesHashed, 8, 5, true)
	if err != nil {
		log.Fatal(err)
	}
	pubBytesToEncode := append(head, data5...)
	addressBech32, err := bech32.Encode("bc", pubBytesToEncode)
	return addressBech32
}

func NewPrivAddressB32() ([]byte, string) {
	privkey, pubkey := NewKeypair()
	priv := privkey.Serialize()
	pub := pubkey.SerializeCompressed()
	address := PubkeyToAddressB32(pub)
	return priv, address
}

func AddressType(address string) (string, error) {
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
	} else if address[0] == '3' {
		_, _, err := base58.CheckDecode(address)
		if err != nil {
			return "", err
		}
		return "p2sh", nil
	}
	return "", errors.New("このアドレス形式には対応していません。")
}

func IsValidAddress(addr string) bool {
	_, err := AddressType(addr)
	if err != nil {
		return false
	} else {
		return true
	}
}
