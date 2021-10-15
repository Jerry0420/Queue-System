FROM node:16.11.1-alpine AS builder-frontend
WORKDIR /app
COPY ./src ./src
COPY ./public ./public
COPY package.json package-lock.json ./
RUN npm install && npm run build

FROM golang:1.17.1-alpine AS builder
WORKDIR /app
COPY . .
COPY --from=builder-frontend /app/build /app/build
RUN go build -o main /app/main.go && \
    /app/scripts/install-migrate.sh

FROM alpine:3.14
EXPOSE 8000
WORKDIR /app
COPY --from=builder /app/main /app
COPY --from=builder /usr/bin/migrate /usr/bin/migrate
COPY ./scripts/migration.sh ./scripts/migration.sh
COPY ./migrations ./migrations
RUN addgroup -S appgroup && \
    adduser -S appuser -G appgroup && \
    mkdir /app/logs && \
    chown appuser:appgroup -R /app
USER appuser
ENTRYPOINT /app/scripts/migration.sh up && /app/main