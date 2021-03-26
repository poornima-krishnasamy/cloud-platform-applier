FROM golang:1.13.3 as builder
WORKDIR $GOPATH/src/github.com/poornima-krishnasamy/cloud-platform-applier
COPY . $GOPATH/src/github.com/poornima-krishnasamy/cloud-platform-applier
RUN make build

FROM ubuntu
LABEL maintainer="Poornima Krishnasamy"
WORKDIR /root/
RUN apt-get update && \
    apt-get install -y git
ADD https://storage.googleapis.com/kubernetes-release/release/v1.17.12/bin/linux/amd64/kubectl /usr/local/bin/kubectl
RUN chmod +x /usr/local/bin/kubectl
COPY --from=builder /go/src/github.com/poornima-krishnasamy/cloud-platform-applier
