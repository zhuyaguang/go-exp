# Build the manager binary
FROM golang:1.16 as builder

RUN apt-get -y update && apt-get -y install upx

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
# RUN go mod download

# Copy the go source
COPY main.go main.go
COPY api/ api/
COPY controllers/ controllers/
COPY cmd/ cmd/
COPY pkg/ pkg/

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
ENV GO111MODULE=on
ENV GOPROXY="https://goproxy.cn"

# Build
# RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o manager main.go

# Build
RUN go mod download && \
		go build -mod=readonly -o manager main.go && \
		go build -mod=readonly -o backup cmd/backup/main.go && \
		upx manager backup

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM debian:10.6
WORKDIR /
COPY --from=builder /workspace/manager .
USER 65532:65532

ENTRYPOINT ["/manager"]

FROM debian:10.6 as backup
WORKDIR /
COPY --from=builder /workspace/backup .
USER nonroot:nonroot
ENTRYPOINT ["/backup"]