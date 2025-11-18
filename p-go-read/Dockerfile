# ---------- Etapa 1: Compilación ----------
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Nombre del binario del microservicio READ
RUN go build -o reader .

# ---------- Etapa 2: Imagen final optimizada ----------
FROM alpine:latest

LABEL maintainer="danysoftdev" \
      version="0.1" \
      description="Imagen optimizada para el microservicio de lectura en Go"

RUN adduser -D gouser

WORKDIR /home/gouser/reader

COPY --from=builder /app/reader .

USER gouser

EXPOSE 3001

CMD ["./reader"]
