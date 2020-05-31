package metrics

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"mattianatali.it/carrier-exporter/internal/config"
	"mattianatali.it/carrier-exporter/internal/tim"
	"mattianatali.it/carrier-exporter/internal/vodafone"
	"mattianatali.it/carrier-exporter/internal/wind"
)

var (
	available = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "carrier_data_available_bytes",
		Help: "The available data",
	}, []string{
		"carrier",
	})
	windContainer = wind.Container{}
	timContainer  = tim.Container{}
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
		go func(c chan<- error) {
			status <- registerVodafone(config)
		}(status)

		for i := 0; i < 3; i++ {
			if err := <-status; err != nil {
				fmt.Printf("Error encountered: %+v", err)
			}
		}
		promhttp.Handler().ServeHTTP(w, r)
	}
}

func registerTim(config config.Config) error {
	service := timContainer.GetService()
	availableDataBytes, err := service.GetAvailableDataBytes(tim.Credentials{
		Username: config.Carriers.Tim.Username,
		Password: config.Carriers.Tim.Password,
	})

	if err != nil {
		return err
	}

	available.WithLabelValues("tim").Set(availableDataBytes)
	return nil
}

func registerWind(config config.Config) error {
	windService := windContainer.GetService()

	insight, err := windService.GetInsights(wind.Credentials{
		Username: config.Carriers.Wind.Username,
		Password: config.Carriers.Wind.Password,
	},
		config.Carriers.Wind.LineID,
		config.Carriers.Wind.ContractID,
	)

	if err != nil {
		return err
	}

	available.WithLabelValues("wind").Set(insight.National.Data.Available)

	return nil
}

func registerVodafone(config config.Config) error {
	container := vodafone.Container{}
	service := container.GetService()
	vodafoneConf := config.Carriers.Vodafone
	sim := vodafoneConf.Sims[0]
	resp, err := service.GetCounters(vodafone.Credentials{
		Username: vodafoneConf.Username,
		Password: vodafoneConf.Password,
	}, sim.Phone)

	if err != nil {
		return err
	}

	availableGb := float64(0)

	for _, c := range resp.Result.Counters {
		if contains(sim.Counters, c.ID) {
			for _, v := range c.Threshold.Values {
				if v.Unit == "GB" {
					availableGb += v.ResidualValue
				}
			}
		}
	}

	available.WithLabelValues("vodafone").Set(availableGb * 1024 * 1024 * 1024)

	return nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
