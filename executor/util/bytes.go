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
	EncodingUTF8 PegicBytesEncoding = iota
	EncodingInt32
	EncodingInt64
	EncodingBytes
)

func (e PegicBytesEncoding) String() string {
	switch e {
	case EncodingUTF8:
		return "UTF-8"
	case EncodingInt32:
		return "INT32"
	case EncodingInt64:
		return "INT64"
	case EncodingBytes:
		return "BYTES"
	default:
		panic(nil)
	}
}

type PegicBytes struct {
	encoding PegicBytesEncoding

	value []byte
}

func NewBytesFromString(s string, enc PegicBytesEncoding) (*PegicBytes, error) {
	pb := &PegicBytes{}
	switch enc {
	case EncodingUTF8:
		if !utf8.ValidString(s) {
			return nil, errors.New("invalid utf8 string")
		}
		pb.value = []byte(s) // go uses utf8 by default.

	case EncodingInt32:
		i, err := strconv.ParseInt(s, 10, 32 /*bits*/)
		if err != nil {
			return nil, errors.New("invalid INT32")
		}
		if enc == EncodingInt32 {
			pb.value = make([]byte, 4 /*bytes*/)
			binary.BigEndian.PutUint32(pb.value, uint32(i))
		} else {
			pb.value = make([]byte, 8 /*bytes*/)
			binary.BigEndian.PutUint64(pb.value, uint64(i))
		}

	case EncodingInt64:
		i, err := strconv.ParseInt(s, 10, 64 /*bits*/)
		if err != nil {
			return nil, errors.New("invalid INT64")
		}
		if enc == EncodingInt32 {
			pb.value = make([]byte, 4 /*bytes*/)
			binary.BigEndian.PutUint32(pb.value, uint32(i))
		} else {
			pb.value = make([]byte, 8 /*bytes*/)
			binary.BigEndian.PutUint64(pb.value, uint64(i))
		}

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

func NewBytes(s []byte, enc PegicBytesEncoding) (*PegicBytes, error) {
	pb := &PegicBytes{
		value: s,
	}

	switch enc {
	case EncodingUTF8:
		if !utf8.Valid(s) {
			return pb, errors.New("invalid utf8 bytes")
		}

	case EncodingInt32:
		if len(s) != 4 {
			return pb, fmt.Errorf("bytes is not a valid INT32")
		}

	case EncodingInt64:
		if len(s) != 8 {
			return pb, fmt.Errorf("bytes is not a valid INT64")
		}

	case EncodingBytes:

	default:
		panic(fmt.Sprintf("unsupported encoding %d", enc))
	}
	return pb, nil
}

func (b *PegicBytes) String() string {
	switch b.encoding {
	case EncodingUTF8:
		return string(b.value)

	case EncodingInt32:
		i := binary.BigEndian.Uint32(b.value)
		return fmt.Sprint(int32(i))

	case EncodingInt64:
		i := binary.BigEndian.Uint64(b.value)
		return fmt.Sprint(int64(i))

	case EncodingBytes:
		s := ""
		for _, c := range b.value {
			s += fmt.Sprint(int(c)) + " "
		}
		return strings.TrimSpace(s)

	default:
		panic(nil)
	}
}

func (b *PegicBytes) Bytes() []byte {
	return b.value
}
