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
