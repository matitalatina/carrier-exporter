package metrics

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"mattianatali.it/carrier-exporter/internal/config"
	"mattianatali.it/carrier-exporter/internal/tim"
	"mattianatali.it/carrier-exporter/internal/wind"
)

var (
	available = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "carrier_data_available_bytes",
		Help: "The available data",
	}, []string{
		"carrier",
	})
)

func HandleMetrics(config config.Config) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		status := make(chan error)

		go func(c chan<- error) {
			status <- registerWind(config)
		}(status)
		go func(c chan<- error) {
			status <- registerTim(config)
		}(status)

		for i := 0; i < 2; i++ {
			if err := <-status; err != nil {
				fmt.Printf("Error encountered: %+v", err)
			}
		}
		promhttp.Handler().ServeHTTP(w, r)
	}
}

func registerTim(config config.Config) error {
	availableDataBytes, err := tim.GetAvailableDataBytes(tim.Credentials{
		Username: config.Secrets.Tim.Username,
		Password: config.Secrets.Tim.Password,
	})

	if err != nil {
		return err
	}

	available.WithLabelValues("tim").Set(float64(availableDataBytes))
	return nil
}

func registerWind(config config.Config) error {
	container := wind.Container{}
	windService := container.GetService()

	insight, err := windService.GetInsights(wind.Credentials{
		Username: config.Secrets.Wind.Username,
		Password: config.Secrets.Wind.Password,
	},
		config.Secrets.Wind.LineID,
		config.Secrets.Wind.ContractID,
	)

	if err != nil {
		return err
	}

	available.WithLabelValues("wind").Set(insight.National.Data.Available)

	return nil
}
