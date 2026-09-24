FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o fleet-server ./cmd/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o fleet-simulator ./simulator/simulator.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/fleet-server .
COPY --from=builder /app/fleet-simulator .

EXPOSE 8080
CMD ["./fleet-server"]