FROM library/golang:1.25rc1-alpine3.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go mod download


COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /webhook-server ./cmd/


FROM alpine:3.18
RUN apk --no-cache add ca-certificates
COPY --from=builder /webhook-server /webhook-server

RUN mkdir -p /etc/webhook/certs
ENTRYPOINT ["/webhook-server"]
