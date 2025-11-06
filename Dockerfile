FROM composer:2 AS vendor
WORKDIR /app
COPY composer.json composer.lock ./
RUN composer install --no-scripts --prefer-dist --no-interaction --no-progress

FROM php:8.3-cli-alpine AS app
WORKDIR /app
RUN apk add --no-cache bash curl

# Install supercronic for non-root cron execution
ENV SUPERCRONIC_URL=https://github.com/aptible/supercronic/releases/download/v0.2.29/supercronic-linux-amd64 \
    SUPERCRONIC=supercronic-linux-amd64 \
    SUPERCRONIC_SHA1SUM=cd48d45c4b10f3f0bfdd3a57d054cd05ac96812b

RUN curl -fsSLO "$SUPERCRONIC_URL" \
    && echo "${SUPERCRONIC_SHA1SUM}  ${SUPERCRONIC}" | sha1sum -c - \
    && chmod +x "$SUPERCRONIC" \
    && mv "$SUPERCRONIC" "/usr/local/bin/${SUPERCRONIC}" \
    && ln -s "/usr/local/bin/${SUPERCRONIC}" /usr/local/bin/supercronic

# Create non-root user
RUN addgroup -g 1001 appuser && \
    adduser -D -u 1000 -G appuser appuser

RUN mkdir -p /app/data /app/var && \
    chown -R appuser:appuser /app

COPY --chown=appuser:appuser . /app
COPY --from=vendor --chown=appuser:appuser /app/vendor /app/vendor

# Switch to non-root user
USER appuser

CMD ["php", "-v"]
