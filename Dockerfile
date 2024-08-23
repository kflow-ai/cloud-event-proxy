FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -v -o /app/cloud-event-proxy github.com/kflow-ai/cloud-event-proxy/cmd/cloudeventproxy

FROM gcr.io/distroless/static:nonroot

COPY --from=builder /app/cloud-event-proxy /app/cloud-event-proxy

ENTRYPOINT ["/app/cloud-event-proxy"]
