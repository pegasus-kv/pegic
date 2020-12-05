package executor

import (
	"context"
	"fmt"
	"pegic/executor/util"
	"time"
)

func Get(rootCtx *Context, hashKeyStr, sortkeyStr string) error {
	pegasusArgs, err := readPegasusArgs(rootCtx, []string{hashKeyStr, sortkeyStr})
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	rawValue, err := rootCtx.UseTable.Get(ctx, pegasusArgs[0].Bytes(), pegasusArgs[1].Bytes())
	if err != nil {
		return err
	}

	value, err := util.NewBytes(rawValue, rootCtx.ValueEnc)
	if err != nil {
		return err
	}
	fmt.Fprintln(rootCtx, value.String())
	return nil
}
