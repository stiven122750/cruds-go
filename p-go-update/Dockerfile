# ---------- Etapa 1: Compilación ----------
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Nombre del binario del microservicio UPDATE
RUN go build -o updater .

# ---------- Etapa 2: Imagen final optimizada ----------
FROM alpine:latest

LABEL maintainer="danysoftdev" \
      version="0.1" \
      description="Imagen optimizada para el microservicio de actualización en Go"

RUN adduser -D gouser

WORKDIR /home/gouser/updater

COPY --from=builder /app/updater .

USER gouser

EXPOSE 5001

CMD ["./updater"]
