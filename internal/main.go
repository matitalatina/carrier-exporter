package main

import (
	"fmt"

	"mattianatali.it/carrier-exporter/internal/tim"
)

func main() {
	container := tim.Container{}
	client := container.GetClient()
	resp, err := client.Authorize()
	fmt.Printf("%+v", err)
	fmt.Printf("%+v", resp)
}
