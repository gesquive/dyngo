FROM index.docker.io/gesquive/go-builder:latest AS builder

# This docker file relys on `make dist` being run first
ENV APP=dyngo
ARG TARGETARCH
ARG TARGETOS
ARG TARGETVARIANT

COPY dist/ /dist/
RUN copy-release

FROM scratch
LABEL maintainer="Gus Esquivel <gesquive@gmail.com>"

# Import from builder
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd

COPY --from=builder /app/${APP} /app/

# Use an unprivileged user
USER runner

ENV DYNGO_CONFIG "/config/config.yml"
ENV DYNGO_LOG_FILE "stdout"
ENV DYNGO_RUN_ONCE false

VOLUME /config

ENTRYPOINT ["/app/dyngo"]
