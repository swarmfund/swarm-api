all: build run

build:
	go build -i gitlab.com/swarmfund/api/cmd/api

run:
	./api run --config=local-config.yaml
