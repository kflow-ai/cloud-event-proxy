FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -v -o /app/cloud-event-proxy github.com/kflow-ai/cloud-event-proxy/cmd/cloudeventproxy

FROM alpine:3.18

COPY --from=builder /app/cloud-event-proxy /app/cloud-event-proxy

ENTRYPOINT ["/app/cloud-event-proxy"]
