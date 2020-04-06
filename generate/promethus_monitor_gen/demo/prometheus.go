package main

import (
	"github.com/labstack/echo"
	"github.com/prometheus/client_golang/prometheus"
	"os"
	//"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

var (
	MetricCounter *prometheus.CounterVec
	MetricTimer   *prometheus.HistogramVec
)

func Init(serviceName string) {
	MetricConf = metricConfig{
		NameSpace:         "base_service",
		ServiceName:       serviceName,
		TraceMethodMysql:  "mysql",
		TraceMethodRedis:  "redis",
		TraceMethodRouter: "router",
		TraceMethodRpc:    "rpc",
	}
	// 注册metric
	MetricCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: MetricConf.NameSpace,
			Subsystem: MetricConf.ServiceName,
			Name:      "counter",
			Help:      "count request",
		},
		[]string{"instance", "method", "endpoint", "code"},
	)
	MetricTimer = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: MetricConf.NameSpace,
			Subsystem: MetricConf.ServiceName,
			Name:      "latency",
			Help:      "consume time",
			Buckets:   []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5},
		},
		[]string{"instance", "method", "endpoint"},
	)

	prometheus.MustRegister(MetricCounter)
	prometheus.MustRegister(MetricTimer)

	//go baseMetrics()
}

func GetMonitorHandle() http.Handler {
	return promhttp.Handler()
}

func HostName() string {
	name, _ := os.Hostname()
	return name
}

func Timer(method, endpoint string) *prometheus.Timer {
	t := prometheus.NewTimer(MetricTimer.WithLabelValues(HostName(), method, endpoint))
	return t
}

func TimerMiddleware(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		//延时
		t := RouterTimeTrace(c.Path())
		defer t.Trace()
		return handlerFunc(c)
	}
}

type tracer struct {
	t *prometheus.Timer
}

func (t *tracer) Trace() {
	t.t.ObserveDuration()
}
func TimeTrace(method, endpoint string) *tracer {
	return &tracer{t: Timer(method, endpoint)}
}

func MysqlTimeTrace(endpoint string) *tracer {
	return TimeTrace(MetricConf.TraceMethodMysql, endpoint)
}
func RedisTimeTrace(endpoint string) *tracer {
	return TimeTrace(MetricConf.TraceMethodRedis, endpoint)
}
func RpcTimeTrace(endpoint string) *tracer {
	return TimeTrace(MetricConf.TraceMethodRpc, endpoint)
}
func RouterTimeTrace(endpoint string) *tracer {
	return TimeTrace(MetricConf.TraceMethodRouter, endpoint)
}


var MetricConf metricConfig

//go:generate pd_gen -service_name=api -namespace=base_service -title=test
type metricConfig struct {
	NameSpace   string
	ServiceName string

	TraceMethodMysql  string `trace_val:"mysql"`
	TraceMethodRedis  string `trace_val:"redis"`
	TraceMethodRouter string `trace_val:"router"`
	TraceMethodRpc    string `trace_val:"rpc"`
}
