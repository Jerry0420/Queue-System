FROM golang:1.17.1-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main /app/main.go

FROM alpine:3.14
EXPOSE 8000
WORKDIR /app
COPY --from=builder /app/main /app
COPY --from=builder /app/migrations /app/migrations
COPY --from=builder /app/scripts/up_db.sh /app/scripts/up_db.sh
RUN addgroup -S appgroup && \
    adduser -S appuser -G appgroup && \
    mkdir /app/logs && \
    chown appuser:appgroup -R /app && \
    /app/scripts/up_db.sh
USER appuser
ENTRYPOINT /app/main