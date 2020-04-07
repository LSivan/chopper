package workgroup

import (
	"fmt"
	xerror "github.com/LSivan/hatchet/x-error"
	xlog "github.com/LSivan/hatchet/x-log"
	"math/rand"
	"testing"
	"time"
)

func TestAsyncWorkGroup(t *testing.T) {
	wg := NewAsyncWorkGroup(20, 512)
	wg.SetIdxGen(func() int {
		return 0
	})
	wg.Start()



	time.Sleep(time.Second)
	for i := 0; i <= 1000; i++ {
		xerror.DoIfErrorNotNil(wg.AddWork(func(j int) func() error {
			return func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(5)+1))
				fmt.Println("this is", j)
				return nil
			}
		}(i)), func(err error) {
			xlog.Sugar.Named("add work error").Errorw("", "err", err.Error())
		})
	}

	time.Sleep(time.Second * 5)
}
