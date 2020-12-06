package executor

import (
	"context"
	"fmt"
	"time"
)

func Get(rootCtx *Context, hashKeyStr, sortkeyStr string) error {
	pegasusArgs, err := readPegasusArgs(rootCtx, []string{hashKeyStr, sortkeyStr})
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	rawValue, err := rootCtx.UseTable.Get(ctx, pegasusArgs[0], pegasusArgs[1])
	if err != nil {
		return err
	}

	value, err := rootCtx.ValueEnc.DecodeAll(rawValue)
	if err != nil {
		return err
	}
	fmt.Fprintln(rootCtx, value)
	return nil
}
