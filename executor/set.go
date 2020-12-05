package executor

import (
	"context"
	"pegic/executor/util"
	"time"
)

func Set(rootCtx *Context, hashKey, sortkey, value *util.PegicBytes) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	return rootCtx.UseTable.Set(ctx, hashKey.Bytes(), sortkey.Bytes(), value.Bytes())
}
