###############################################################################
#                                  BUILD
###############################################################################

FROM golang:1.25-bookworm AS builder

WORKDIR /app

COPY . /app

RUN CGO_ENABLED=1 go build -o main main.go

###############################################################################
#                                INSTALLATION
###############################################################################

FROM debian:bookworm-slim

# Add Debian mirror and install CA certs
# Using Debian slim instead of headless-shell to significantly reduce container size
RUN apt-get update && apt-get install -y \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

ENV WORKDIR=/app
ENV OMP_THREAD_LIMIT=1

ADD manifest/i18n $WORKDIR/manifest/i18n/
ADD resource $WORKDIR/resource/
ADD storage $WORKDIR/storage/

COPY --from=builder /app/main $WORKDIR/main
RUN chmod +x $WORKDIR/main

###############################################################################
#                                   START
###############################################################################
WORKDIR $WORKDIR
ENTRYPOINT ["/app/main"]
