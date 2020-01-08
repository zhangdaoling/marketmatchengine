all: build

prepare:
	cp ./app/kafka_demo/mysql/order.sql ./order.sql

sync:
	go build -o build/sync app/main.go

vendor:
	go mod vendor

build: sync

clean:
	rm -rf build

cleanAll: clean
	rm -rf vendor

.PHONY: test build clean cleanAll vendor docker
