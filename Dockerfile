FROM golang:1.22-alpine AS builder

WORKDIR /src
COPY go.mod go.sum* ./
RUN go mod download || true

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/cert-sentinel ./main.go

FROM alpine:3.19
RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=builder /app/cert-sentinel /usr/local/bin/cert-sentinel

EXPOSE 9109
ENTRYPOINT ["cert-sentinel"]
CMD ["-config=/etc/cert-sentinel/config.yaml", "-listen=:9109"]
