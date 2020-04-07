# -----------------------------------------------------------------------------
# Build...
FROM golang:1.13.5-alpine AS build

ENV VERSION=v0.8.0 GO111MODULE=on PACKAGE=github.com/derailed/popeye

WORKDIR /go/src/$PACKAGE

COPY go.mod go.sum main.go ./
COPY internal internal
COPY types types
COPY pkg pkg
COPY cmd cmd

RUN apk upgrade && apk add git ca-certificates
RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/popeye \
  -trimpath -ldflags="-w -s -X $PACKAGE/cmd.version=$VERSION" *.go

# -----------------------------------------------------------------------------
# Image...
FROM alpine:3.11.2
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /go/bin/popeye /bin/popeye
ENTRYPOINT [ "/bin/popeye" ]