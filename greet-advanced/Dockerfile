# docker build --progress=plain -t greetapi -f Dockerfile .
FROM golang:1.23 AS builder

WORKDIR /app

COPY . ./

RUN go mod download

# Set environment variables for Go build
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
# Build the Go application
RUN go build -v -o main ./cmd/main

FROM alpine:latest

WORKDIR /app

RUN set -eux; \
    \
    apk update; \
    apk upgrade; \
    apk add --no-cache --virtual .build-deps \
    zip \
    unzip \
    ; \
    apk add --no-cache \
    curl \
    ca-certificates \
    bash \
    ; \
    # curl -k http://certificates../file.zip > /tmp/file.zip;unzip /tmp/file.zip -d /usr/local/share/ca-certificates && \
    # update-ca-certificates && \
    apk del --no-network .build-deps;

COPY --from=builder /app/main .

RUN addgroup --system --gid 1001 appgroup && \
    adduser --system --uid 1001 appuser --ingroup appgroup && \
    chown -R appuser:appgroup /app && mkdir -p /home/appuser && chown -R appuser:appgroup /home/appuser && \
    chmod +x /app/main

# USER appuser
# Make port 50051 available to the world outside this container
EXPOSE 50051

# Run the command to start
CMD ["./main"]
