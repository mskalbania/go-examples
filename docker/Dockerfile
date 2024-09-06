#using builder to first compile to target platform and then copy the binary to the final image
FROM golang:alpine AS builder

#optimizing build time (when dependencies are not changed layers are taken from cache)
COPY go.mod go.sum /source/
WORKDIR /source
RUN go mod download

ADD . /source
RUN go build -ldflags "-X main.version=1.0.0"

FROM alpine:latest
COPY --from=builder /source/go-examples /app/go-examples
COPY rest/config-local-docker.yaml /app/rest/config-local-docker.yaml
WORKDIR /app
CMD ["./go-examples"]