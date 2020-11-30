package pegic

import (
	"errors"
	"fmt"
	"strings"

	"github.com/XiaoMi/pegasus-go-client/admin"
	"github.com/XiaoMi/pegasus-go-client/pegasus"
)

type Executor struct {
	ctx *ExecContext
}

type ExecContext struct {
	client      pegasus.Client
	adminClient admin.Client
	table       pegasus.TableConnector
}

func NewExecutor(metaServerList []string) *Executor {
	cfg := pegasus.Config{
		MetaServers: metaServerList,
	}
	adminCfg := admin.Config{
		MetaServers: metaServerList,
	}
	client := pegasus.NewClient(cfg)
	adminClient := admin.NewClient(adminCfg)
	return &Executor{
		ctx: &ExecContext{client: client, adminClient: adminClient},
	}
}

var noTableError = errors.New("no table selected, please run `USE <\"table_name\">`")

// Executor is the pegic command executor in interactive mode.
func (e *Executor) Execute(s string) {
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return
	}

	parts := strings.SplitN(s, " ", 2)
	cmdStr := strings.ToUpper(parts[0])
	cmd, found := commandsTable[cmdStr]
	if !found {
		fmt.Printf("ERROR: unsupported command: \"%s\"\n", cmdStr)
		return
	}
	var subcommand string = ""
	if len(parts) != 1 {
		subcommand = parts[1]
	}
	err := cmd.parse(strings.TrimSpace(subcommand))
	if err != nil {
		fmt.Printf("ERROR: unable to parse command: %s\n", err)
		return
	}
	if err := cmd.execute(e.ctx); err != nil {
		fmt.Printf("ERROR: execution failed: %s\n", err)
		return
	}
}
