# Stage 1: Сборка
FROM golang:1.24 AS builder

WORKDIR /app

# Шаг 1: Копируем go.mod и go.sum → кэшируем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Шаг 2: Копируем исходный код
COPY . .

# Шаг 3: Собираем бинарник
RUN CGO_ENABLED=0 GOOS=linux go build -o server cmd/server/main.go

# Stage 2: Финальный образ
FROM alpine:latest

# Устанавливаем ca-certificates для HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Копируем бинарник из builder-стадии
COPY --from=builder /app/server .

# Копируем Swagger-документацию (если используется)
COPY --from=builder /app/docs ./docs

EXPOSE 8080

CMD ["./server"]