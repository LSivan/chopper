/**
 * DateTime : 18-11-29 上午11:42
 * Author : liangxingwei
 */
package xtime

import (
	"time"
)

// 获取时间差的小时差
func GetTimeIntervalStr(interval int) (res [3]int) {
	d := interval
	h := d / 3600

	d = d % 3600
	m := d / 60

	s := d % 60
	res[0] = h
	res[1] = m
	res[2] = s
	return
}

// 从
func GetTimeStampFromStr(day string) int64 {
	//  直接用parse的话默认会用UTC，时间可能会有问题
	t, err := time.ParseInLocation("2006-01-02 15:04:05", day, time.Local)
	if err != nil {
		return 0
	}
	return t.Unix()
}
