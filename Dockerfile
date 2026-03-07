# Build stage — Go backend
FROM golang:1.24-alpine AS builder

WORKDIR /build

COPY backend/ .
RUN go mod tidy && CGO_ENABLED=0 GOOS=linux go build -o dashboard ./cmd/server

# Build stage — Svelte frontend
FROM node:22-alpine AS frontend

WORKDIR /frontend
COPY frontend/package*.json ./
RUN npm ci
COPY frontend/ .
RUN npm run build

# Final image
FROM alpine:3.21

RUN adduser -D -u 1000 app

WORKDIR /app

COPY --from=builder /build/dashboard .
COPY --from=frontend /frontend/dist ./frontend/dist

RUN mkdir -p /app/data/uploads && chown -R app:app /app
USER app

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --retries=2 \
  CMD wget -qO- http://localhost:8080/health || exit 1

CMD ["./dashboard"]
