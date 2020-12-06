package executor

import (
	"context"
	"time"
)

func Set(rootCtx *Context, hashKeyStr, sortkeyStr, valueStr string) error {
	pegasusArgs, err := readPegasusArgs(rootCtx, []string{hashKeyStr, sortkeyStr, valueStr})
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	return rootCtx.UseTable.Set(ctx, pegasusArgs[0], pegasusArgs[1], pegasusArgs[2])
}
