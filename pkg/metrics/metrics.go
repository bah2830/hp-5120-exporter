package metrics

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bah2830/switch-exporter/pkg/config"
	"github.com/prometheus/client_golang/prometheus"
	prom "github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	namespace = "switch_exporter"

	Temp = prom.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "temp_celcius",
			Help:      "Chassis temperature in celcius",
		},
		[]string{"name", "warning", "alarm", "critical"},
	)
)

func Serve(conf *config.Config) {
	log.Printf("serving metrics on 0.0.0.0:%d%s", conf.Metrics.Port, conf.Metrics.Path)
	http.Handle(conf.Metrics.Path, promhttp.Handler())
	http.ListenAndServe(fmt.Sprintf(":%d", conf.Metrics.Port), nil)
}
