FROM golang:1.25-alpine AS builder

WORKDIR /app

# Устанавливаем зависимости для скачивания модулей
RUN apk add --no-cache git

# Копируем файлы модулей и скачиваем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем бинарник
# CGO_ENABLED=0 для статической линковки (важно для alpine)
RUN CGO_ENABLED=0 GOOS=linux go build -o app cmd/main.go

# Финальный этап (Production Image)
FROM alpine:latest

WORKDIR /root/

# Устанавливаем сертификаты (нужны для HTTPS запросов, если будут)
RUN apk --no-cache add ca-certificates

# Копируем бинарник из этапа сборки
COPY --from=builder /app/app .

# Копируем конфиги (Viper ищет их в папке configs)
COPY --from=builder /app/configs ./configs

# Копируем папку с миграциями (если захочешь запускать их из кода, но мы будем делать это через отдельный контейнер)
COPY --from=builder /app/migrations ./migrations

# Открываем порт
EXPOSE 8080

# Запускаем приложение
CMD ["./app"]
