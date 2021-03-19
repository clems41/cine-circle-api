FROM golang:1.16.2-alpine3.12 as build-env
RUN apk add --update git openssh-client
ENV CGO_ENABLED 0

WORKDIR /src

ADD go.mod go.sum ./
RUN go mod download

ADD . ./
RUN set -ex; \
    if grep -q '^package main *' *.go; then go install .; fi; \
    if [ -d cmd ]; then go install ./cmd/...; fi

FROM alpine:3.11
RUN apk add --no-cache curl
ENTRYPOINT ["/bin/cine-circle-api"]
#RUN echo "hosts: files dns" > /etc/nsswitch.conf
#ADD https://github.com/golang/go/raw/master/lib/time/zoneinfo.zip /zoneinfo.zip
#ENV ZONEINFO /zoneinfo.zip
COPY --from=build-env /go/bin/ /bin
