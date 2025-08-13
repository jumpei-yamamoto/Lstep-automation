# syntax=docker/dockerfile:1.7
##
## Build
##
FROM golang:1.22-alpine AS build
WORKDIR /src
RUN apk add --no-cache git ca-certificates tzdata

# Go modules
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Source
COPY backend/ ./

# Build (CGO無効=静的リンク). Echo + pgx/stdlib はCGO不要で動きます
ENV CGO_ENABLED=0
RUN go build -trimpath -ldflags="-s -w" -o /out/server ./cmd/server

##
## Runtime
##
FROM alpine:3.20
WORKDIR /app
RUN apk add --no-cache ca-certificates tzdata

# 非rootユーザ
RUN addgroup -S app && adduser -S app -G app
USER app

COPY --from=build /out/server /app/server
ENV PORT=8080
EXPOSE 8080

# NOTE: DockerのHEALTHCHECKはdistrolessが不要ですが、ALB/RenderのHCに任せる運用を推奨
ENTRYPOINT ["/app/server"]