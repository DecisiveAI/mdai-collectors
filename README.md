# watcher-collector

A purpose-built OpenTelemetry collector for watching and routing telemetry streams

# Usage

## Setup: Install the OpenTelemetry Collector Builder

```sh
CGO_ENABLED=0 go install go.opentelemetry.io/collector/cmd/builder@v0.117.0
```

## Build the collector binary

```sh
builder --config=builder.yaml 
```

## Run the collector

```sh
./cmd/watcher-collector/mdai-watcher-collector --config=collector.yaml  
```
