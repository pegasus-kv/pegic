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

	HashKeyEnc, SortKeyEnc, ValueEnc util.PegicBytesEncoding
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
func readPegasusArgs(ctx *Context, args []string) ([]*util.PegicBytes, error) {
	// the first argument must be hashkey
	hashkey, err := util.CreateBytesFromString(args[0], ctx.HashKeyEnc)
	if err != nil {
		return nil, err
	}
	if len(args) == 1 {
		return []*util.PegicBytes{hashkey}, nil
	}

	sortkey, err := util.CreateBytesFromString(args[1], ctx.SortKeyEnc)
	if err != nil {
		return nil, err
	}
	if len(args) == 2 {
		return []*util.PegicBytes{hashkey, sortkey}, nil
	}

	value, err := util.CreateBytesFromString(args[1], ctx.SortKeyEnc)
	if err != nil {
		return nil, err
	}
	if len(args) == 3 {
		return []*util.PegicBytes{hashkey, sortkey, value}, nil
	}

	panic("more than 3 arguments are given")
}
