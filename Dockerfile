FROM golang:1.23.3-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest

# RUN apk --no-cahce add ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .

ENV DB_HOST=localhost
ENV DB_USER=postgres
ENV DB_PASSWORD=postgres
ENV DB_NAME=golang
ENV DB_PORT=5432
ENV DB_SSLMODE=disable
ENV PORT=8000
ENV NATS_URL=localhost
ENV NATS_CHANNEL=golang

EXPOSE 8080

CMD ["./main"]