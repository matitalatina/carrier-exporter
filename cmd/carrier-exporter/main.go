package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"mattianatali.it/carrier-exporter/internal/config"
	"mattianatali.it/carrier-exporter/internal/metrics"
)

func main() {
	configPath := flag.String("config", "config.yml", "Config location")
	flag.Parse()
	config, err := config.ParseFile(*configPath)

	if err != nil {
		log.Panic(err)
	}

	port := config.App.Port
	http.HandleFunc("/metrics", metrics.HandleMetrics(config))
	log.Printf("Listening on port %d", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
