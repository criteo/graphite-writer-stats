FROM golang:1.12.9 as builder
WORKDIR /workspace
COPY go.sum go.mod ./
RUN go mod download
COPY . .
RUN go build ./cmd/graphite-writer-stats/graphite-writer-stats.go

FROM golang:1.12.9
WORKDIR /app
COPY --from=builder /workspace/graphite-writer-stats .
COPY --from=builder /workspace/configs/ configs/
ENTRYPOINT ["./graphite-writer-stats"]
CMD ["--brokers","kafka:9092","--topic","metrics","--group","graphite-writer"]
