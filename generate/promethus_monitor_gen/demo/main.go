package main

import (
	"math/rand"
	"time"
)

func main() {
	Init("api")
	for range time.Tick(time.Second) {
		foo()
	}
}

func foo() {
	defer MysqlTimeTrace("foo").Trace()
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(250)))
	return
}
