FROM golang:1.23-bookworm as build
ARG  OTEL_VERSION=0.117.0
WORKDIR /app
RUN CGO_ENABLED=0 go install go.opentelemetry.io/collector/cmd/builder@v${OTEL_VERSION}
COPY . .
RUN CGO_ENABLED=0 builder --config=builder.yaml

FROM scratch
COPY --from=build app/cmd/watcher-collector /
EXPOSE 4317/tcp 4318/tcp 8891/tcp 8899/tcp 13133/tcp

CMD ["/mdai-watcher-collector", "--config=/etc/watcher-collector-config/collector.yaml"]