FROM golang:1.17.1-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main /app/main.go

FROM alpine:3.14
EXPOSE 8000
WORKDIR /app
COPY --from=builder /app/main /app
ENTRYPOINT /app/main