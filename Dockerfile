FROM golang:1.16.9-alpine as build-env
RUN apk add --no-cache ca-certificates git
RUN apk --no-cache add tzdata

# Install depedencies if needed (swag, easytags, etc...)
# RUN go get -u -v github.com/swaggo/swag/cmd/swag
# RUN go get -u -v github.com/betacraft/easytags

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

ADD . ./

RUN go build -o /go/bin ./...

FROM alpine:3.11
RUN apk add --no-cache curl

# Useful when hostAliases are specified in values chart
#RUN echo "hosts: files dns" > /etc/nsswitch.conf

# Useful to use location with datetime
COPY --from=build-env /usr/share/zoneinfo /usr/share/zoneinfo

# Copy all go binaries previously built
COPY --from=build-env /go/bin/ /bin

# Useful to get error codes from admin webService. Can be removed if not used.
COPY --from=build-env /app/internal/domain /app/internal/domain

ENTRYPOINT ["/bin/cine-circle-api"]
