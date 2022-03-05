# Build the manager binary
FROM golang:1.16 as builder

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY main.go main.go
COPY api/ api/
COPY controllers/ controllers/

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o manager main.go

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
# FROM gcr.io/distroless/static:nonroot
# WORKDIR /
# COPY --from=builder /workspace/manager .
# USER 65532:65532

# ENTRYPOINT ["/manager"]

FROM fedora:34

WORKDIR /

RUN dnf -y install --nodocs --setopt=install_weak_deps=0 \
           iproute iptables iptables-nft ipset procps-ng && \
    dnf -y clean all

COPY --from=builder /workspace/manager .

# Wrapper scripts to choose the appropriate iptables
# https://github.com/kubernetes-sigs/iptables-wrappers
COPY scripts/iptables-wrapper-installer.sh /
RUN /iptables-wrapper-installer.sh

ENTRYPOINT ["/manager"]

