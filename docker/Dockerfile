FROM golang:1.13-buster as build

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
ENV GOOS linux
ENV GARCH amd64
ENV CGO_ENABLED 0

RUN go build

FROM gcr.io/distroless/static-debian10:nonroot

LABEL maintainer="Marco Santos <marco.santos@loq9.com>" \
    org.label-schema.name="Drone Nomad" \
    org.label-schema.vendor="LOQ9" \
    org.label-schema.schema-version="1.0"

COPY --from=build /app/drone-nomad /bin/drone-nomad

ENTRYPOINT ["/bin/drone-nomad"]