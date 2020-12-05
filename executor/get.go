package executor

import (
	"context"
	"fmt"
	"pegic/executor/util"
	"time"

	"github.com/XiaoMi/pegasus-go-client/pegasus"
)

func Get(rootCtx *Context, tb pegasus.TableConnector, hashKey, sortkey *util.PegicBytes) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	value, err := tb.Get(ctx, hashKey.Bytes(), sortkey.Bytes())
	if err != nil {
		return err
	}

	fmt.Println(value)
	return nil
}
