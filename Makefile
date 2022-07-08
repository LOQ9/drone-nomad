.PHONY: clean build build-mac build-linux docker

PACKAGE = drone-nomad

export GO111MODULE=on
export CGO_ENABLED=0

all: build

build:
	go build -a -tags netgo -o ./dist/$(PACKAGE)

dep-install:
	go mod download

build-mac:
	mkdir -p build/mac
	env GOOS=darwin GOARCH=amd64 go build -a -tags netgo -o ./dist/$(PACKAGE)_darwin_amd64

build-linux:
	mkdir -p build/linux
	env GOOS=linux GOARCH=amd64 go build -a -tags netgo -o ./dist/$(PACKAGE)_linux_amd64

lint: check-lint
	golangci-lint run ./...

test:
	go test -v ./...

check-lint:
	@if ! [ -x "$$(command -v golangci-lint)" ]; then \
		echo "Downloading golangci-lint..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	fi;

docker:
	docker build \
	--label org.label-schema.build-date=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
	--label org.label-schema.vcs-ref=$(git rev-parse --short HEAD) \
	--file docker/Dockerfile --tag plugins/nomad .

clean:
	rm -rf dist