# Build the webhook binary
FROM golang:1.16 as builder

RUN apt-get -y update && apt-get -y install upx


WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum

# Copy the go source
COPY main.go main.go
COPY pkg/ pkg/
COPY cmd/ cmd/

# Build
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
ENV GO111MODULE=on
ENV GOPROXY="https://goproxy.cn"

# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download && \
    go build -a -o admission-registry main.go && \
    go build -a -o tls-manager cmd/tls/main.go && \
    upx admission-registry tls-manager

FROM alpine:3.9.2 as manager
COPY --from=builder /workspace/admission-registry .
ENTRYPOINT ["/admission-registry"]

FROM alpine:3.9.2 as tls
COPY --from=builder /workspace/tls-manager .
ENTRYPOINT ["/tls-manager"]


# docker build --target manager -t cnych/admission-registry:v0.0.2 --target tls -t  cnych/admission-registry-tls:v0.0.2  -f Dockerfile .