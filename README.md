# watcher-collector

A purpose-built OpenTelemetry collector for watching and routing telemetry streams

# Usage

## Build the collector binary

> Prerequisite: Install the collector builder by following instructions [here](https://opentelemetry.io/docs/collector/custom-collector/#step-1---install-the-builder)

```sh
ocb --config=builder.yaml 
```

## Run the collector

```sh
./cmd/watcher-collector/mdai-watcher-collector --config=collector.yaml  
```
