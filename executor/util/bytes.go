package util

import (
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

type PegicBytesEncoding int

const (
	// EncodingUTF8 is the default encoding of string.
	EncodingUTF8  PegicBytesEncoding = 0
	EncodingASCII                    = 1
	EncodingInt                      = 2
	EncodingBytes                    = 3
)

type PegicBytes struct {
	encoding PegicBytesEncoding

	value []byte
}

func CreateBytesFromString(s string, enc PegicBytesEncoding) (*PegicBytes, error) {
	pb := &PegicBytes{}
	switch enc {
	case EncodingUTF8:
		if !utf8.ValidString(s) {
			return nil, errors.New("invalid utf8 string")
		}
		pb.value = []byte(s) // go uses utf8 by default.

	case EncodingInt:
		i, err := strconv.ParseUint(s, 10, 64 /*bits*/)
		if err != nil {
			return nil, errors.New("invalid integer")
		}
		pb.value = make([]byte, 8 /*bytes*/)
		binary.BigEndian.PutUint64(pb.value, i)

	case EncodingBytes:
		bytesInStrList := strings.Split(s, " ")

		pb.value = make([]byte, len(bytesInStrList))
		for i, byteStr := range bytesInStrList {
			b, err := strconv.Atoi(byteStr)
			if err != nil || b >= 128 || b < -128 { // byte ranges from [-128, 127]
				return nil, fmt.Errorf("invalid byte \"%s\"", byteStr)
			}
			pb.value[i] = byte(b)
		}

	default:
		panic(fmt.Sprintf("unsupported encoding %d", enc))
	}
	return pb, nil
}

func (*PegicBytes) String() string {
	// TODO(wutao)
	return ""
}

func (b *PegicBytes) Bytes() []byte {
	return b.value
}
