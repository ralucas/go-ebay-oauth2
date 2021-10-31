.PHONY: build
build:
	go build -v ./...

.PHONY: test
test:
	go test -v --tags=unit ./...

.PHONY: test-integration
test-integration:
	go test -v --tags=integration ./...