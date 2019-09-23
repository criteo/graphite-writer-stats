export GO111MODULE := on
all: fmt lint vet build-dev test run 
build-dev:
	go build ./cmd/graphite-writer-stats/graphite-writer-stats.go
fmt:
	go fmt ./...
vet:
	go vet ./...
lint:
	golint ./...
test:
	go test -v ./...
run:
	./graphite-writer-stats --brokers localhost:9092 --topic metrics --group graphite-writer-stats
docker-build:
	docker build . -t graphite-writer-stats -f build/Dockerfile
docker-run:
	docker run -p 8080:8080 graphite-writer-stats --brokers localhost:9092 --topic metrics -group graphite-writer-stats
