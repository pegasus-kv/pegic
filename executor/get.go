package executor

import (
	"context"
	"fmt"
	"pegic/executor/util"
	"time"
)

func Get(rootCtx *Context, hashKey, sortkey *util.PegicBytes) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	value, err := rootCtx.UseTable.Get(ctx, hashKey.Bytes(), sortkey.Bytes())
	if err != nil {
		return err
	}

	fmt.Println(value)
	return nil
}
