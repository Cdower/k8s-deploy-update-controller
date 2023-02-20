FROM --platform=$BUILDPLATFORM golang:1.20.1-alpine AS builder
# https://github.com/docker/buildx#building-multi-platform-images
# https://github.com/tonistiigi/xx
COPY --from=tonistiigi/xx / /
# if CGO_ENABLED=1
# RUN apk add clang lld
ARG TARGETPLATFORM
WORKDIR /app
COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY cmd ./cmd
COPY pkg ./pkg

ENV CGO_ENABLED=0
RUN xx-go build -o /deploy-update-controller /app/cmd/deployUpdater && \
    xx-verify /deploy-update-controller
# TODO: Not root
FROM --platform=$BUILDPLATFORM alpine:latest AS app
ARG TARGETPLATFORM
WORKDIR /root/
COPY --from=builder /deploy-update-controller ./
CMD [ "./deploy-update-controller" ]
