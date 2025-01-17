FROM golang:1.23-bookworm AS build
ARG OTEL_VERSION=0.117.0
ARG TARGETOS
ARG TARGETARCH
ENV CGO_ENABLED=0
WORKDIR /app
RUN curl https://github.com/open-telemetry/opentelemetry-collector-releases/releases/download/cmd%2Fbuilder%2Fv${OTEL_VERSION}/ocb_${OTEL_VERSION}_${TARGETOS}_${TARGETARCH} -o /app/builder && chmod +x /app/builder
COPY . .
RUN /app/builder --config=builder.yaml

FROM scratch
COPY --from=build app/cmd/watcher-collector /
EXPOSE 4317/tcp 4318/tcp 8891/tcp 8899/tcp 13133/tcp

CMD ["/mdai-watcher-collector", "--config=/etc/watcher-collector-config/collector.yaml"]
