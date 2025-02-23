# Этап сборки
FROM golang:1.24.0 AS builder
WORKDIR /go/src/availabilityMonitorWeb
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/checkerWeb ./cmd/

# Этап выполнения
FROM debian:buster-slim
WORKDIR /app
COPY --from=builder /go/src/availabilityMonitorWeb/bin/checkerWeb /app/checkerWeb
RUN chmod +x /app/
EXPOSE 8080/tcp
CMD ["/app/checkerWeb"]
