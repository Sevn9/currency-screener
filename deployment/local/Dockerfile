# Stage 1: Build the Go application
FROM golang:1.25-alpine AS builder

WORKDIR /build

COPY . .

ARG BUILD_TARGET

RUN go build -o /build/app $BUILD_TARGET

# Stage 2: Create a smaller image for running the binary
FROM alpine:latest
#FROM scratch

COPY --from=builder /build/app /

ENTRYPOINT ["/app"]
