package executor

import (
	"io"
	"pegic/executor/util"

	"github.com/XiaoMi/pegasus-go-client/pegasus"
)

type Context struct {
	// Every command should use Context as the fmt.Fprint's writer.
	io.Writer

	pegasus.Client

	// default to nil
	UseTable pegasus.TableConnector

	// default to nil
	Compressor util.BytesCompression

	HashKeyEnc, SortKeyEnc, ValueEnc util.BytesEncoder
}

func NewContext(writer io.Writer, metaAddrs []string) *Context {
	c := &Context{
		Writer: writer,
		Client: pegasus.NewClient(pegasus.Config{
			MetaServers: metaAddrs,
		}),
	}
	return c
}

// readPegasusArgs returns exactly the same number of arguments of input `args` if no failure.
// The order of arguments are also preserved.
func readPegasusArgs(ctx *Context, args []string) ([][]byte, error) {
	// the first argument must be hashkey
	hashkey, err := ctx.HashKeyEnc.EncodeAll(args[0])
	if err != nil {
		return nil, err
	}
	if len(args) == 1 {
		return [][]byte{hashkey}, nil
	}

	sortkey, err := ctx.SortKeyEnc.EncodeAll(args[1])
	if err != nil {
		return nil, err
	}
	if len(args) == 2 {
		return [][]byte{hashkey, sortkey}, nil
	}

	value, err := ctx.SortKeyEnc.EncodeAll(args[1])
	if err != nil {
		return nil, err
	}
	if len(args) == 3 {
		return [][]byte{hashkey, sortkey, value}, nil
	}

	panic("more than 3 arguments are given")
}
