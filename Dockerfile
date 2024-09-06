FROM golang:1.23.1-alpine3.20 AS builder

ENV RS_CONFIG_PATH=./config/dev.yaml

WORKDIR /go/src/reports

COPY . .

RUN go mod tidy \
    go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ./cmd/main.go

FROM alpine:latest AS runner

RUN apk --no-cache add ca-certificates

WORKDIR /root

ENV CONFIG_PATH=/root/config/dev.yaml

RUN mkdir -p /root/config

COPY --from=builder /go/src/reports/config ./config

COPY --from=builder /go/src/reports/cmd/main .

EXPOSE 10501

ENTRYPOINT [ "main" ]