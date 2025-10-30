FROM composer:2 AS vendor
WORKDIR /app
COPY composer.json composer.lock ./
RUN composer install --no-scripts --prefer-dist --no-interaction --no-progress

FROM php:8.3-cli-alpine AS app
WORKDIR /app
RUN apk add --no-cache bash
RUN mkdir -p /app/data /app/var
COPY . /app
COPY --from=vendor /app/vendor /app/vendor
RUN chmod +x /app/docker-entrypoint-cron.sh
CMD ["php", "-v"]
