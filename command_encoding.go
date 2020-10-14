package pegic

import (
	"encoding/binary"
	"errors"
	"fmt"
	"pegic/ast"
	"strconv"
	"strings"
	"unicode/utf8"
)

type pegicBytesEncoding int

const (
	encodingUTF8  pegicBytesEncoding = 1
	encodingASCII                    = 2
	encodingInt                      = 3
	encodingBytes                    = 4
)

type pegicBytes struct {
	encoding pegicBytesEncoding

	value []byte
}

func createBytesFromString(s string, enc pegicBytesEncoding) (*pegicBytes, error) {
	pb := &pegicBytes{}
	switch enc {
	case encodingUTF8:
		if !utf8.ValidString(s) {
			return nil, errors.New("invalid utf8 string")
		}
		pb.value = []byte(s) // go uses utf8 by default.

	case encodingInt:
		i, err := strconv.ParseUint(s, 10, 64 /*bits*/)
		if err != nil {
			return nil, errors.New("invalid integer")
		}
		pb.value = make([]byte, 8 /*bytes*/)
		binary.BigEndian.PutUint64(pb.value, i)

	case encodingBytes:
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

func (*pegicBytes) String() string {
	// TODO(wutao)
	return ""
}

type encodingCommand struct {
	algorithm string
}

func (*encodingCommand) execute(parsedCmd *ast.ParsedCommand) error {
	// TODO(wutao)
	return nil
}

func (*encodingCommand) astNode() *ast.CommandASTNode {
	// Possible inputs:
	//   ENCODING HASHKEY <UTF-8|ASCII|INT|BYTES>
	//   ENCODING SORTKEY <UTF-8|ASCII|INT|BYTES>
	//   ENCODING VALUE <UTF-8|ASCII|INT|BYTES>
	//   ENCODING RESET

	encNode := &ast.CommandASTNode{
		Arguments: []*ast.CommandArgument{
			{
				Name: "encoding",
				Selections: &ast.Select{
					Items: []string{
						"UTF-8",
						"ASCII",
						"INT",
						"BYTES",
					},
				},
			},
		},
	}
	node := &ast.CommandASTNode{
		CustomDescription: "Configure the encoding of hashkey, sortkey and value",
		SubNodes:          make(map[string]*ast.CommandASTNode),
	}
	node.SubNodes["HASHKEY"] = encNode
	node.SubNodes["SORTKEY"] = encNode
	node.SubNodes["VALUE"] = encNode
	node.SubNodes["RESET"] = &ast.CommandASTNode{
		CustomDescription: "Reset the encoding settings to UTF-8",
	}
	return node
}
