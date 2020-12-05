package executor

import (
	"io"

	"github.com/XiaoMi/pegasus-go-client/pegasus"
)

type Client struct {
	// Every command should use Client as the fmt.Fprint's writer.
	io.Writer

	pegasus.Client
}

func NewClient(writer io.Writer, metaAddrs []string) *Client {
	c := &Client{
		Writer: writer,
		Client: pegasus.NewClient(pegasus.Config{
			MetaServers: metaAddrs,
		}),
	}
	return c
}
