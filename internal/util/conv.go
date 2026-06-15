package util

import (
	"encoding/base64"
	"encoding/binary"
	"errors"
	"strconv"
	"strings"

	"github.com/btcsuite/btcutil/base58"
)

func B58Encode(msg []byte) string {
	return base58.Encode(msg)
}

func B58Decode(msg string) []byte {
	return base58.Decode(msg)
}

func B64Encode(msg []byte) string {
	encoded := base64.StdEncoding.EncodeToString(msg)
	return encoded
}

func B64Decode(msg string) []byte {
	decoded, _ := base64.StdEncoding.DecodeString(msg)
	return decoded
}

func SatsToBtc(sats string) string {
	digits := len(sats)
	var ans string
	if digits <= 8 {
		zeroNum := 8 - digits
		zeros := strings.Repeat("0", zeroNum)
		ans = "0." + zeros + sats
	} else {
		satsByte := []byte(sats)
		bigByte := satsByte[:len(satsByte)-8]
		smallByte := satsByte[len(satsByte)-8:]
		big := string(bigByte)
		small := string(smallByte)
		ans = big + "." + small
	}
	ansByte := []byte(ans)
	prevLength := len(ansByte)
	for i := prevLength - 1; ansByte[i] == byte('0'); i-- {
		ansByte = ansByte[:len(ansByte)-1]
	}
	if ansByte[len(ansByte)-1] == byte('.') {
		ansByte = append(ansByte, byte('0'))
	}
	ans = string(ansByte)
	return ans
}

func BtcToSats(btc string) int {
	sats := 0
	if strings.Contains(btc, ".") {
		numl := strings.Split(btc, ".")
		big, _ := strconv.Atoi(numl[0])
		smallStr := numl[1] + strings.Repeat("0", 8-len(numl[1]))
		small, _ := strconv.Atoi(smallStr)
		sats += big * 100000000
		sats += small
	} else {
		am, _ := strconv.Atoi(btc)
		sats += am * 100000000
	}
	return sats
}

func IntToBytes(value int, length int) ([]byte, error) {
	if length == 4 {
		buf := make([]byte, 4)
		binary.LittleEndian.PutUint32(buf, uint32(value))
		return buf, nil
	} else if length == 8 {
		buf := make([]byte, 8)
		binary.LittleEndian.PutUint64(buf, uint64(value))
		return buf, nil
	} else {
		return nil, errors.New("Length must be 4 or 8.")
	}
}

func BytesToInt(value []byte, length int) (int, error) {
	if length == 4 {
		ans := binary.LittleEndian.Uint32(value)
		return int(ans), nil
	} else if length == 8 {
		ans := binary.LittleEndian.Uint64(value)
		return int(ans), nil
	} else {
		return -1, errors.New("Length must be 4 or 8.")
	}
}

func IntToStr(i int) string {
	return strconv.Itoa(i)
}

func StrToInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func IntToCompactsize(i int) ([]byte, error) {
	if i < 253 {
		return []byte{byte(i)}, nil
	} else if i < 65536 {
		head := []byte{0xfd}
		buf := make([]byte, 2)
		binary.LittleEndian.PutUint16(buf, uint16(i))
		return append(head, buf...), nil
	} else if i < 4294967296 {
		head := []byte{0xfe}
		buf := make([]byte, 4)
		binary.LittleEndian.PutUint32(buf, uint32(i))
		return append(head, buf...), nil
	} else {
		return nil, errors.New("Too large to convert to conpact size")
	}
}
