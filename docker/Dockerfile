#using builder to first compile to target platform and then copy the binary to the final image
FROM golang:alpine as builder
ADD . /source
WORKDIR /source
RUN go build -ldflags "-X main.version=1.0.0"

FROM alpine:latest
COPY --from=builder /source/go-examples /app/go-examples
WORKDIR /app
CMD ["./go-examples"]