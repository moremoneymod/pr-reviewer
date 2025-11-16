FROM alpine:3.18

RUN apk update && \
    apk add --no-cache postgresql-client

RUN wget -O /bin/goose https://github.com/pressly/goose/releases/download/v3.14.0/goose_linux_x86_64 && \
    chmod +x /bin/goose

WORKDIR /app
COPY migrations/ ./migrations/

CMD ["sh", "-c", "until pg_isready -h $DB_HOST -p $DB_PORT -U $DB_USER; do sleep 2; done && goose -dir /app/migrations postgres \"user=$DB_USER password=$DB_PASSWORD dbname=$DB_NAME host=$DB_HOST port=$DB_PORT sslmode=disable\" up"]