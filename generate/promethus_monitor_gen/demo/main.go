package main

import (
	xprometheus "github.com/LSivan/hatchet/x-prometheus"
	"math/rand"
	"time"
)

func main() {
	xprometheus.InitDefault("api")
	for range time.Tick(time.Second) {
		foo()
	}
}

func foo() {
	defer xprometheus.MysqlTimeTrace("foo").Trace()
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(250)))
	return
}
