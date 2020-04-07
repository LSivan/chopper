/**
 * DateTime : 18-12-22 下午2:59
 * Author : liangxingwei
 */

package xtime

import (
	"fmt"
	"testing"
)

func TestGetTimeStampFromStr(t *testing.T) {
	fmt.Println(GetTimeStampFromStr("2018-12-22 14:23:12"))
}

func TestGetTimeIntervalStr(t *testing.T) {
	type args struct {
		interval int64
	}
	tests := []struct {
		name string
		args args
	}{
		{"normal test", args{3700}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetTimeIntervalStr(tt.args.interval)
			t.Logf("got:%v\n", got)
		})
	}
}
