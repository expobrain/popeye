# -----------------------------------------------------------------------------
# Base...
FROM alpine:3.11.2 AS base

RUN apk add --no-cache python3 py3-crcmod ca-certificates

# -----------------------------------------------------------------------------
# Build...
FROM golang:1.13.5-alpine AS popeye

RUN apk upgrade && apk add git ca-certificates

ENV VERSION=v0.8.0 GO111MODULE=on PACKAGE=github.com/derailed/popeye

WORKDIR /go/src/$PACKAGE

COPY go.mod go.sum main.go ./
COPY internal internal
COPY types types
COPY pkg pkg
COPY cmd cmd

RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/popeye \
  -trimpath -ldflags="-w -s -X $PACKAGE/cmd.version=$VERSION" *.go

# -----------------------------------------------------------------------------
# gcloud...
FROM base AS gcloud

RUN apk add curl

ENV CLOUD_SDK_VERSION=288.0.0
ENV CLOUDSDK_PYTHON=python3
ENV CLOUDSDK_URL=https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-sdk-${CLOUD_SDK_VERSION}-linux-x86_64.tar.gz

ENV PATH /google-cloud-sdk/bin:$PATH

RUN curl -L $CLOUDSDK_URL | tar -zx
RUN gcloud config set core/disable_usage_reporting true
RUN gcloud config set component_manager/disable_update_check true
RUN gcloud config set metrics/environment github_docker_image
RUN gcloud --version

# -----------------------------------------------------------------------------
# Image...
FROM base

COPY --from=gcloud /google-cloud-sdk /usr/local/Caskroom/google-cloud-sdk/latest/google-cloud-sdk
COPY --from=popeye /go/bin/popeye /bin/popeye

ENV PATH /usr/local/Caskroom/google-cloud-sdk/latest/google-cloud-sdk/bin:$PATH

RUN gcloud --version
RUN popeye -h

ENTRYPOINT [ "/bin/popeye" ]
