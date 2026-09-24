.PHONY: run test simulate build docker-build

run:
	go run cmd/main.go

test:
	go test -v ./...

simulate:
	go run simulator/simulator.go

build:
	go build -o bin/fleet-server cmd/main.go
	go build -o bin/fleet-simulator simulator/simulator.go

docker-build:
	docker build -t mini-fleet-monitor .