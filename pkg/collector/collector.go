package collector

import (
	"fmt"
	"net/http"

	"github.com/giantswarm/microerror"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Config struct {
	Port uint
}

type Collector struct {
	port uint
}

func New(config Config) (*Collector, error) {
	if config.Port == 0 {
		return nil, microerror.Maskf(invalidConfigError, "%T.Port can't be 0", config)
	}

	return &Collector{
		port: config.Port,
	}, nil
}

func (c *Collector) StartAsync() {
	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(fmt.Sprintf(":%d", c.port), nil) //nolint
}
