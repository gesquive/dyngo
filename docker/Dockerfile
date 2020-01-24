FROM index.docker.io/gesquive/go-builder:latest AS builder

ENV APP=dyngo
ARG TARGETARCH
ARG TARGETOS
ARG TARGETVARIANT

COPY dist/ /dist/
RUN copy-release

RUN mkdir -p /etc/${APP}
COPY docker/config.yml /etc/${APP}
COPY docker/snakeoil.sh /etc/${APP}

FROM scratch
LABEL maintainer="Gus Esquivel <gesquive@gmail.com>"

# Import the user/group files from builder
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

# Import from builder
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /app/${APP} /app/
COPY --from=builder /etc/${APP} /etc/

# Use an unprivileged user
USER runner:runner
WORKDIR /config

VOLUME /config

ENTRYPOINT ["/app/dyngo"]
