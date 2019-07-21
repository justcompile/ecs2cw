FROM golang:alpine as deps


ENV GO111MODULE on
ENV GOLINT_VERSION 1.17.1

WORKDIR $GOPATH/src/github.com/justcompile/ecs2cw

RUN mkdir -p /build && \
    apk --no-cache add curl git bash gcc libc-dev ca-certificates && \
    update-ca-certificates && \
    curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s v$GOLINT_VERSION

COPY ./go.mod $GOPATH/src/github.com/justcompile/ecs2cw
COPY ./go.sum $GOPATH/src/github.com/justcompile/ecs2cw

RUN go mod download

ADD . $GOPATH/src/github.com/justcompile/ecs2cw

FROM deps as build

WORKDIR $GOPATH/src/github.com/justcompile/ecs2cw
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /build/ecs2cw . 

FROM scratch
WORKDIR /root/
COPY --from=build /build/ .
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
CMD ["./ecs2cw", "--interval", "60s"]