# for some reason writing to `/srv/VERSION` in the final image doesn't
# seem to work. i suspect this is because `/srv` is mounted as a volume in
# `abiosoft/caddy:0.11.0-no-stats`. writing to `/srv/VERSION` seems to
# work in a builder though. bare this in mind when trying to refactor
# this code or if you're wondering why we don't just do it in the final
# container
FROM alpine/git:1.0.4 AS builder
ARG GIT_REV
RUN if [ -z "$GIT_REV" ]; then \
        echo -e "\033[91m  ---> compliance-suite-server.ga: GIT_REV is not set \033[0m" 1>&2; \
        exit 1; \
    else \
        echo -e "\033[92m  ---> compliance-suite-server.ga: GIT_REV=${GIT_REV} \033[0m"; \
    fi

RUN echo -e "\033[92m  ---> compliance-suite-server.ga: creating https://compliance-suite-server.ga/VERSION \033[0m"
RUN echo "compliance-suite-server.ga=${GIT_REV}" >> /srv/VERSION
# print date in ISO 8601 format:
# $ date --iso-8601=seconds
RUN echo "build_time=$(date -Iseconds)" >> /srv/VERSION

FROM abiosoft/caddy:0.11.0-no-stats
# auto-agree to Let's Encrypt Subscriber Agreement
ENV ACME_AGREE true
ENV CADDYFILE "/etc/Caddyfile"

COPY ./Caddyfile /etc/Caddyfile
COPY ./Caddyfile_staging /etc/Caddyfile_staging
COPY --from=builder /srv/VERSION /srv/VERSION

RUN cat /srv/VERSION
RUN ls -lah /srv

# ENTRYPOINT ["/bin/parent", "caddy"]
CMD ["--conf", "$CADDYFILE", "--log", "stdout", "--agree=$ACME_AGREE"]
