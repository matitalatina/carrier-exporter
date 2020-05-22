package main

import (
	"fmt"

	"mattianatali.it/carrier-exporter/internal/config"
	"mattianatali.it/carrier-exporter/internal/wind"
)

func main() {
	container := wind.Container{}
	windService := container.GetService()

	config, err := config.ParseFile("../../config.yml")

	if err != nil {
		fmt.Println(err)
	}

	insight, err := windService.GetInsights(wind.Credentials{
		Username: config.Secrets.Wind.Username,
		Password: config.Secrets.Wind.Password,
	},
		config.Secrets.Wind.LineID,
		config.Secrets.Wind.ContractID,
	)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(insight)
}
