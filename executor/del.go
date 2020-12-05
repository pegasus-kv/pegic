package executor

import (
	"context"
	"pegic/executor/util"
	"time"
)

func Del(rootCtx *Context, hashKey *util.PegicBytes, sortkey *util.PegicBytes) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	return rootCtx.UseTable.Del(ctx, hashKey.Bytes(), sortkey.Bytes())
}
