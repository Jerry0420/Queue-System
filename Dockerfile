FROM node:16.11.1-alpine AS builder-frontend
WORKDIR /app
COPY package.json package-lock.json ./
COPY ./src ./src
COPY ./public ./public
RUN npm install && npm run build

FROM golang:1.17.1-alpine AS builder
WORKDIR /app
COPY . .
COPY --from=builder-frontend /app/build /app/build
RUN go build -o main /app/main.go

FROM alpine:3.14
EXPOSE 8000
WORKDIR /app
COPY --from=builder /app/main /app
RUN addgroup -S appgroup && \
    adduser -S appuser -G appgroup && \
    mkdir /app/logs && \
    chown appuser:appgroup -R /app
USER appuser
ENTRYPOINT /app/main