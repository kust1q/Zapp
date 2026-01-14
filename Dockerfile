FROM golang:1.24.4 AS builder

RUN apt-get update && apt-get install -y --no-install-recommends \
    gcc \
    libc-dev \
    libpq-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ .

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o main-app ./cmd/app

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o main-search ./cmd/search

FROM alpine:3.21

RUN apk add --no-cache \
    ca-certificates \
    libpq \
    libc6-compat \
    && rm -rf /var/cache/apk/*

RUN mkdir -p /app/certs \
    && mkdir -p /app/configs

COPY --from=builder /app/main-app /app/main-app
COPY --from=builder /app/main-search /app/main-search

COPY --from=builder /app/certs/ /app/certs/
COPY --from=builder /app/configs/ /app/configs/

RUN chmod 755 /app/main-app /app/main-search \
    && chmod 644 /app/certs/*.pem \
    && chmod 644 /app/configs/*.yml

WORKDIR /app

EXPOSE 8080 8081

CMD ["/app/main-app"]