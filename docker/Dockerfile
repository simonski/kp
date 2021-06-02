FROM alpine
RUN apk update
RUN apk add openssh-keygen
RUN apk add go
ENV GO111MODULE=on
ENV GOPATH=/root/go
ENV GOBIN=/root/go/bin
ENV PATH=$PATH:$GOBIN


