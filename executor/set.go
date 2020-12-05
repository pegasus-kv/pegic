package executor

import (
	"context"
	"pegic/executor/util"
	"time"

	"github.com/XiaoMi/pegasus-go-client/pegasus"
)

func Set(rootCtx *Context, tb pegasus.TableConnector, hashKey, sortkey, value *util.PegicBytes) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	return tb.Set(ctx, hashKey.Bytes(), sortkey.Bytes(), value.Bytes())
}
