package workgroup

import (
	xerror "github.com/LSivan/hatchet/x-error"
	xlog "github.com/LSivan/hatchet/x-log"
	"testing"
	"time"
)

func TestAsyncWorkGroup(t *testing.T) {
	wg := NewAsyncWorkGroup(20, 512)
	wg.Start()

	time.Sleep(time.Second * 3)

	for i := 0; i <= 1000; i++ {
		xerror.DoIfErrorNotNil(wg.AddWork(func() error {
			time.Sleep(time.Millisecond * 10)
			xlog.Sugar.Named("test").Debug("this is", i)
			return nil
		}), func(err error) {
			xlog.Sugar.Named("add work error").Errorw("", "err", err.Error())
		})
	}
	time.Sleep(time.Second)
}
