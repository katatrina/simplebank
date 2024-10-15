# Build stage
FROM golang:1.23.2-alpine AS buidler
WORKDIR /app
COPY . .
RUN go build -o main main.go
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.1/migrate.linux-amd64.tar.gz | tar xvz

# Run stage
FROM alpine:3.20
WORKDIR /app
COPY --from=buidler /app/main .
COPY --from=buidler /app/migrate .
COPY app.env .
COPY db/migrations ./migrations
COPY start.sh .
COPY wait-for.sh .

EXPOSE 8080
