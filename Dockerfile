FROM node:22-alpine AS frontend-builder
WORKDIR /src/frontend
COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

FROM golang:1.24-bookworm AS backend-builder
WORKDIR /src/backend
RUN apt-get update && apt-get install -y --no-install-recommends gcc libc6-dev \
    && rm -rf /var/lib/apt/lists/*
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ ./
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -o /out/brain-training .

FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates tzdata \
    && rm -rf /var/lib/apt/lists/* \
    && useradd --system --uid 10001 --create-home appuser
WORKDIR /app
COPY --from=backend-builder /out/brain-training /app/brain-training
COPY --from=frontend-builder /src/frontend/dist /app/static
RUN mkdir -p /app/data /app/avatars && chown -R appuser:appuser /app

ENV PORT=8000 \
    GIN_MODE=release \
    DB_DRIVER=sqlite \
    DB_PATH=/app/data/brain-training.db \
    FRONTEND_DIR=/app/static

USER appuser
VOLUME ["/app/data", "/app/avatars"]
EXPOSE 8000
ENTRYPOINT ["/app/brain-training"]
