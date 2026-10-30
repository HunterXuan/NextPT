###############################################################################
#                                  BUILD
###############################################################################

FROM --platform=$BUILDPLATFORM golang:1.25-bookworm AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

# Source edits should not invalidate the downloaded dependency layer.
COPY go.mod go.sum ./
RUN go mod download

COPY . /app

RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -o main main.go

###############################################################################
#                                INSTALLATION
###############################################################################

FROM debian:bookworm-slim

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENV WORKDIR=/app
ENV OMP_THREAD_LIMIT=1

ADD manifest/i18n $WORKDIR/manifest/i18n/
ADD resource $WORKDIR/resource/
ADD storage $WORKDIR/storage/

COPY --from=builder /app/main $WORKDIR/main

###############################################################################
#                                   START
###############################################################################
WORKDIR $WORKDIR
ENTRYPOINT ["/app/main"]
