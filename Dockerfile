FROM golang:1.23-bookworm as build
ARG  OTEL_VERSION=0.117.0
ARG TARGETPLATFORM
RUN case "$TARGETPLATFORM" in \
      "linux/amd64") export GOARCH=amd64 ;; \
      "linux/arm64") export GOARCH=arm64 ;; \
      "linux/arm/v7") export GOARCH=arm ;; \
      *) echo "Unsupported platform: $TARGETPLATFORM"; exit 1 ;; \
    esac && echo "Building for GOARCH=$GOARCH"
WORKDIR /app
RUN export CGO_ENABLED=0
RUN go install go.opentelemetry.io/collector/cmd/builder@v${OTEL_VERSION}
COPY . .
RUN builder --config=builder.yaml

FROM scratch
COPY --from=build app/cmd/watcher-collector /
EXPOSE 4317/tcp 4318/tcp 8891/tcp 8899/tcp 13133/tcp

CMD ["/mdai-watcher-collector", "--config=/etc/watcher-collector-config/collector.yaml"]