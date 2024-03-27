# Build stage
FROM golang:1.22.1 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-s -w" -o prometheus-mqtt-sd ./cmd/prometheus-mqtt-sd/main.go

# Final stage
FROM scratch
COPY --from=builder /app/prometheus-mqtt-sd /app/
WORKDIR /app
CMD ["./prometheus-mqtt-sd"]