FROM golang:1.22.6-alpine AS build
ENV CGO_ENABLED=1

RUN apk add --no-cache \
    # Important: required for go-sqlite3
    gcc \
    # Required for Alpine
    musl-dev

WORKDIR /app

COPY . .

RUN go build -o forum ./cmd/api/

CMD ["./forum"]