# Base image with common dependencies
FROM golang:1.24-alpine3.21 AS golang-base

RUN apk add --no-cache ca-certificates yq

#============================================================
# Development Stage
FROM golang-base AS dev

ARG UID=1000

RUN apk add --no-cache git python3 py3-pip openssh make curl

RUN adduser -D -u $UID app

ENV GOPATH /home/app/go
ENV PATH $PATH:/home/app/.local/bin:/home/app/go/bin

RUN mkdir /app && chown app:app /app

RUN pip install --break-system-packages wait-for-it

ENV GOMPLATE_VERSION=v3.11.5

RUN curl -sSL https://github.com/hairyhenderson/gomplate/releases/download/${GOMPLATE_VERSION}/gomplate_linux-amd64 \
    -o /usr/local/bin/gomplate && \
    chmod +x /usr/local/bin/gomplate

USER app

RUN go env -w GOPRIVATE=codehub.digitalarsenal.net

RUN git config --global url."git@codehub.digitalarsenal.net:".insteadOf "https://codehub.digitalarsenal.net"

WORKDIR /app

#============================================================
# Build Stage
FROM golang-base AS build

RUN apk add --no-cache git upx curl

WORKDIR /go/app
COPY . /go/app

# Setup go env
RUN go env -w GOOS=linux CGO_ENABLED=0

# Build app
RUN go build -ldflags="-w -s" -trimpath -o bin/app ./cmd/main.go
RUN upx bin/*

ENV GOMPLATE_VERSION=v3.11.5

RUN curl -sSL https://github.com/hairyhenderson/gomplate/releases/download/${GOMPLATE_VERSION}/gomplate_linux-amd64 \
    -o /usr/local/bin/gomplate && \
    chmod +x /usr/local/bin/gomplate

#============================================================
# Final Release Stage
FROM golang-base AS release

EXPOSE 8080

RUN addgroup -S app && adduser -S app -G app

COPY --from=build --chown=app:app --chmod=755 /go/app/config config
COPY --from=build --chown=0:0 --chmod=755 /go/app/scripts/start.sh /start.sh
COPY --from=build --chown=0:0 --chmod=755 /go/app/bin/app app
COPY --from=build --chown=0:0 --chmod=755 /usr/local/bin/gomplate /usr/local/bin/gomplate

ENTRYPOINT ["/start.sh"]

USER app
