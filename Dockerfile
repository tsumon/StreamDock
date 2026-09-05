# Build stage
FROM golang:1.26.6-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
ARG VERSION=v1.4.0
RUN CGO_ENABLED=0 go build -trimpath -buildvcs=false -ldflags="-s -w -X main.appVersion=${VERSION}" -o streamdock .

# Runtime stage
FROM alpine:3.21

RUN apk add --no-cache ca-certificates su-exec tzdata

WORKDIR /app
RUN addgroup -S streamdock && \
    adduser -S -D -H -u 10001 -G streamdock streamdock && \
    mkdir -p /app/data && \
    chown streamdock:streamdock /app/data && \
    chmod 0700 /app/data
COPY --from=builder --chown=root:root --chmod=0555 /app/streamdock /app/streamdock
COPY --chown=root:root --chmod=0555 docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh

EXPOSE 9090

ENV PORT=9090
ENV DB_PATH=/app/data/streamdock.db
ENV PANEL_BIND_ADDR=0.0.0.0
ENV ALLOW_INSECURE_HTTP=true

VOLUME ["/app/data"]

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
  CMD ["/app/streamdock", "--healthcheck"]

ENTRYPOINT ["/usr/local/bin/docker-entrypoint.sh"]
CMD ["./streamdock"]
