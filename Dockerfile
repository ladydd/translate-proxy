FROM golang:1.22-alpine AS builder
WORKDIR /app
ENV GOPROXY=https://goproxy.cn,direct
COPY go.mod go.sum ./
RUN go mod download
COPY main.go .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o translate-proxy .

FROM alpine:3.19
COPY --from=builder /app/translate-proxy /translate-proxy
COPY config.json /config.json
EXPOSE 8787
CMD ["/translate-proxy"]
