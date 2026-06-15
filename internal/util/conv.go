package util

import (
	"encoding/base64"
	"encoding/binary"
	"errors"
	"strconv"
	"strings"

	"github.com/btcsuite/btcutil/base58"
)

func B58_encode(msg []byte) string {
	return base58.Encode(msg)
}

func B58_decode(msg string) []byte {
	return base58.Decode(msg)
}

func B64_encode(msg []byte) string {
	encoded := base64.StdEncoding.EncodeToString(msg)
	return encoded
}

func B64_decode(msg string) []byte {
	decoded, _ := base64.StdEncoding.DecodeString(msg)
	return decoded
}

func Sats_to_btc(sats string) string {
	digits := len(sats)
	var ans string
	if digits <= 8 {
		zero_num := 8 - digits
		zeros := strings.Repeat("0", zero_num)
		ans = "0." + zeros + sats
	} else {
		sats_byte := []byte(sats)
		big_byte := sats_byte[:len(sats_byte)-8]
		small_byte := sats_byte[len(sats_byte)-8:]
		big := string(big_byte)
		small := string(small_byte)
		ans = big + "." + small
	}
	ans_byte := []byte(ans)
	prev_length := len(ans_byte)
	for i := prev_length - 1; ans_byte[i] == byte('0'); i-- {
		ans_byte = ans_byte[:len(ans_byte)-1]
	}
	if ans_byte[len(ans_byte)-1] == byte('.') {
		ans_byte = append(ans_byte, byte('0'))
	}
	ans = string(ans_byte)
	return ans
}

func Btc_to_sats(btc string) int {
	sats := 0
	if strings.Contains(btc, ".") {
		numl := strings.Split(btc, ".")
		big, _ := strconv.Atoi(numl[0])
		small_str := numl[1] + strings.Repeat("0", 8-len(numl[1]))
		small, _ := strconv.Atoi(small_str)
		sats += big * 100000000
		sats += small
	} else {
		am, _ := strconv.Atoi(btc)
		sats += am * 100000000
	}
	return sats
}

func Int_to_bytes(value int, length int) ([]byte, error) {
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

func Bytes_to_int(value []byte, length int) (int, error) {
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

func Int_to_str(i int) string {
	return strconv.Itoa(i)
}

func Str_to_int(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func Int_to_compactsize(i int) ([]byte, error) {
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
