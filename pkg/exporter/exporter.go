package exporter

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/bah2830/switch-exporter/pkg/metrics"

	"github.com/bah2830/switch-exporter/pkg/config"
	"github.com/bah2830/switch-exporter/pkg/hpswitch"
)

type Exporter struct {
	conf   *config.Config
	cancel context.CancelFunc
}

func New(c *config.Config) *Exporter {
	return &Exporter{
		conf: c,
	}
}

func (e *Exporter) Start() error {
	log.Printf("starting switch exporter with %s poll", e.conf.Interval.String())

	// Run the first poll to check for errors
	if err := e.pollSwitch(); err != nil {
		return err
	}

	// Setup context for cancel signal later
	ctx, cancel := context.WithCancel(context.Background())
	e.cancel = cancel

	tick := time.NewTicker(e.conf.Interval)
	go func() {
		for {
			select {
			case <-tick.C:
				if err := e.pollSwitch(); err != nil {
					log.Println(err)
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}

func (e *Exporter) pollSwitch() error {
	log.Println("getting switch environment details")

	hpSwitch, err := hpswitch.NewWithPassword(e.conf.SSH.Host, e.conf.SSH.Port, e.conf.SSH.Username, e.conf.SSH.Password)
	if err != nil {
		return err
	}
	defer hpSwitch.Close()

	details, err := hpSwitch.GetEnvironmentDetails()
	if err != nil {
		return err
	}

	for _, sensor := range details.Sensors {
		labels := prometheus.Labels{
			"name":     sensor.Name,
			"warning":  strconv.Itoa(sensor.Limits.Warning),
			"alarm":    strconv.Itoa(sensor.Limits.Warning),
			"critical": strconv.Itoa(sensor.Limits.Warning),
		}

		metrics.Temp.With(labels).Set(float64(sensor.TempCelsius))
	}

	return nil
}

func (e *Exporter) Stop() {
	e.cancel()
}
