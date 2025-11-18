# ---------- Etapa 1: Compilación ----------
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o deleter .

# ---------- Etapa 2: Imagen final optimizada ----------
FROM alpine:latest

LABEL maintainer="danysoftdev" \
      version="0.1" \
      description="Imagen optimizada para el microservicio de eliminación en Go"

RUN adduser -D gouser

WORKDIR /home/gouser/deleter

COPY --from=builder /app/deleter .

USER gouser

EXPOSE 3000

CMD ["./deleter"]
