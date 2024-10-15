# Build stage
FROM golang:1.23.2-alpine AS buidler
WORKDIR /app
COPY . .
RUN go build -o main main.go

# Run stage
FROM alpine:3.20
WORKDIR /app
COPY --from=buidler /app/main .
COPY app.env .

ENTRYPOINT ["/app/main"]
EXPOSE 8080
