# drone-nomad

Drone plugin for deployment with Nomad

## Build

Build the binary with the following command:

```console
export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=0
export GO111MODULE=on

go build -v -a -tags netgo -o build/linux/amd64/drone-nomad
```

## Docker

Build the Docker image with the following command:

```console
docker build \
  --label org.label-schema.build-date=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
  --label org.label-schema.vcs-ref=$(git rev-parse --short HEAD) \
  --file docker/Dockerfile.linux.amd64 --tag plugins/nomad .
```

## Usage

```console
docker run --rm \
  -v $(pwd):$(pwd) \
  -w $(pwd) \
  plugins/nomad
```