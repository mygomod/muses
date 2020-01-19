package prom

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/goecology/muses/pkg/app"
	"github.com/goecology/muses/pkg/common"
)

var (
	// HTTPServer for http server
	HTTPServerTimer = NewPromTimer("goecology_http_server", []string{"app", "env", "method"})

	HTTPServerCounter = NewPromCounter("goecology_http_server_code", []string{"app", "env", "method", "code"})

	AppBuildInfo = NewPromCounter("goecology_app_build_info", []string{"app", "env", "version"})
)

// toodo prometheus也要初始化
var defaultCaller = &callerStore{
	Name: common.ModPromName,
}

type callerStore struct {
	Name string
	cfg  Cfg
}

func Register() common.Caller {
	return defaultCaller
}

func (c *callerStore) InitCfg(cfg []byte) error {
	return nil
}

func (c *callerStore) InitCaller() error {
	// 信息初始化到prometheus
	AppBuildInfo.Set(app.Config().Muses.App.Version)
	return nil
}

// Prom struct info
type PromTimer struct {
	histogram *prometheus.HistogramVec
	summary   *prometheus.SummaryVec
}

type PromCounter struct {
	counter *prometheus.GaugeVec
}

// New creates a Prom instance.
func NewPromTimer(name string, labels []string) *PromTimer {
	obj := &PromTimer{
		histogram: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: name,
				Help: name,
			}, labels),
	}
	prometheus.MustRegister(obj.histogram)
	return obj
}

func NewPromCounter(name string, labels []string) *PromCounter {
	obj := &PromCounter{
		counter: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: name,
				Help: name,
			}, labels),
	}
	prometheus.MustRegister(obj.counter)
	return obj

}

// Timing log timing information (in milliseconds) without sampling
func (p *PromTimer) Timing(name string, time int64, extra ...string) {
	label := append([]string{app.Config().Muses.App.Name, app.Config().Muses.App.Env, name}, extra...)
	p.histogram.WithLabelValues(label...).Observe(float64(time))
}

// Incr increments one stat counter without sampling
func (p *PromCounter) Incr(name string, extra ...string) {
	label := append([]string{app.Config().Muses.App.Name, app.Config().Muses.App.Env, name}, extra...)
	p.counter.WithLabelValues(label...).Inc()
}

// Decr decrements one stat counter without sampling
func (p *PromCounter) Decr(name string, extra ...string) {
	label := append([]string{app.Config().Muses.App.Name, app.Config().Muses.App.Env, name}, extra...)
	p.counter.WithLabelValues(label...).Dec()
}

// Add add count    v must > 0
func (p *PromCounter) Add(name string, v int64, extra ...string) {
	label := append([]string{app.Config().Muses.App.Name, app.Config().Muses.App.Env, name}, extra...)
	p.counter.WithLabelValues(label...).Add(float64(v))
}

func (p *PromCounter) Set(name string, extra ...string) {
	label := append([]string{app.Config().Muses.App.Name, app.Config().Muses.App.Env, name}, extra...)
	p.counter.WithLabelValues(label...).Set(1)
}
