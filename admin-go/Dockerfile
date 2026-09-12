# syntax=docker/dockerfile:1

# --- build stage: static Go binary ---
FROM golang:1.25-alpine AS build
WORKDIR /src
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/portfolio-admin .

# --- runtime stage: distroless-lite alpine ---
FROM alpine:3.20
RUN apk add --no-cache ca-certificates \
    && adduser -D -H -u 1000 admin
WORKDIR /app
COPY --from=build /out/portfolio-admin /app/portfolio-admin
# Content volume (Markdown source of truth, bind-mounted RW) + data dir;
# created as root so the runtime user (uid 1000) can write into them.
RUN mkdir -p /content/projects /data
USER admin
EXPOSE 8000
ENTRYPOINT ["/app/portfolio-admin"]