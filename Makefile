.PHONY: build build-arm

BUILD_CMD=go build -o dist/carrier-exporter cmd/carrier-exporter/main.go

build-arm:
	GOOS=linux GOARCH=arm GOARM=7 $(BUILD_CMD)

build:
	$(BUILD_CMD)
