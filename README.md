# drone-nomad

Drone plugin for deployment with Nomad

![Docker Pulls](https://img.shields.io/docker/pulls/loq9/drone-nomad?label=drone-nomad)

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
  plugins/drone-nomad
```

## Template Variables

The following variables could be configured on a nomad template with the following syntax `${VAR_NAME}`.

| Environment | Argument | Description |
|---|---|---|
| PLUGIN_ADDR | | nomad addr |
| PLUGIN_CONSUL_TOKEN | | consul token |
| PLUGIN_VAULT_TOKEN | | vault token |
| PLUGIN_TOKEN | | nomad token |
| PLUGIN_REGION | | nomad region |
| PLUGIN_NAMESPACE | | nomad namespace |
| PLUGIN_TEMPLATE | | nomad template |
| PLUGIN_PRESERVE_COUNTS | | preserve task counts when deploying (bool) |
| PLUGIN_WATCH_DEPLOYMENT | | trigger a deploy and wait till the deployment is complete (bool) |
| PLUGIN_WATCH_DEPLOYMENT_TIMEOUT | | if watch deployment is enabled, wait up to this time duration for the deploy to finish. Errors on timeout. Default: "5m" (duration string) |
| PLUGIN_TLS_CA_CERT | tls_ca_cert | nomad tls ca certificate file |
| PLUGIN_TLS_CA_PATH | tls_ca_path | nomad tls ca certificate file path |
| PLUGIN_TLS_CA_CERT_PEM | tls_ca_cert_pem | nomad tls ca certificate pem |
| PLUGIN_TLS_CLIENT_CERT | tls_client_cert | nomad tls client certificate |
| PLUGIN_TLS_CLIENT_CERT_PEM | tls_client_cert_pem | nomad tls client certificate pem |
| PLUGIN_TLS_CLIENT_KEY | tls_client_key | nomad tls client private key |
| PLUGIN_TLS_CLIENT_KEY_PEM | tls_client_key_pem | nomad tls client private key pem |
| PLUGIN_TLS_SERVERNAME | tls_servername | nomad tls server name |
| PLUGIN_TLS_INSECURE | tls_insecure | nomad tls insecure |
| DRONE_REPO_OWNER | | repository owner |
| DRONE_REPO_NAME | | repository name |
| DRONE_COMMIT_SHA | | git commit sha |
| DRONE_COMMIT_REF | | git commit ref |
| DRONE_COMMIT_BRANCH | | git commit branch |
| DRONE_COMMIT_AUTHOR | | git author name |
| DRONE_COMMIT_MESSAGE | | commit message |
| DRONE_BUILD_EVENT | | build event |
| DRONE_BUILD_NUMBER | | build number |
| DRONE_BUILD_PARENT | | build parent |
| DRONE_BUILD_STATUS | | build status |
| DRONE_BUILD_LINK | | build link |
| DRONE_BUILD_STARTED | | build started |
| DRONE_BUILD_CREATED | | build created |
| DRONE_TAG | | build tag |
| DRONE_JOB_STARTED | | job started |
