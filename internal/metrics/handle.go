package metrics

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"mattianatali.it/carrier-exporter/internal/config"
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
			log.Printf("Error getting wind insights: %+v", err)
		}

		available.WithLabelValues("wind").Set(insight.National.Data.Available)

		promhttp.Handler().ServeHTTP(w, r)
	}
}
