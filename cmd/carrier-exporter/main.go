package main

import (
	"fmt"
	"log"
	"net/http"

	"mattianatali.it/carrier-exporter/internal/config"
	"mattianatali.it/carrier-exporter/internal/metrics"
)

func main() {
	config, err := config.ParseFile("../../config.yml")

	if err != nil {
		log.Panic(err)
	}

	port := config.App.Port
	http.HandleFunc("/metrics", metrics.HandleMetrics(config))
	log.Printf("Listening on port %d", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
