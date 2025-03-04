FROM --platform=$BUILDPLATFORM golang:1.23-bookworm AS builder
WORKDIR /app
COPY . .
ARG OTEL_VERSION=0.117.0
ARG BUILDOS
ARG BUILDARCH
ARG TARGETOS
ARG TARGETARCH
RUN curl --proto '=https' --tlsv1.2 -fL https://github.com/open-telemetry/opentelemetry-collector-releases/releases/download/cmd%2Fbuilder%2Fv${OTEL_VERSION}/ocb_${OTEL_VERSION}_${BUILDOS}_${BUILDARCH} -o /app/builder && chmod +x /app/builder
COPY . .
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} /app/builder --config=builder.yaml

FROM scratch
COPY --from=builder app/cmd/watcher-collector /
EXPOSE 4317/tcp 4318/tcp 8891/tcp 8899/tcp 13133/tcp
USER 65533:65533
ENTRYPOINT ["/mdai-watcher-collector"]
