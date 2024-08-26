FROM golang:1.22 as builder

ENV CGO_ENABLED=0

RUN apt-get update \
    && apt-get install -y git make tzdata ca-certificates --no-install-recommends \
    && update-ca-certificates

WORKDIR /usr/src/build

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app ./cmd


FROM debian:stable-slim

RUN apt-get update \
    && apt-get install -y make tzdata ca-certificates curl --no-install-recommends \
    && update-ca-certificates \
    && apt-get -y clean && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY schema /app/schema
COPY config /app/config

COPY --from=builder /usr/src/build/app /app/

CMD ["./app"]


