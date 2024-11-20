# Используем официальный базовый образ Go
FROM golang:1.23.2 as builder

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем файлы проекта
COPY . .

# Загружаем зависимости
RUN go mod tidy

# Собираем бинарный файл
RUN go build -o auth-service ./cmd/sso

# Минимизируем образ
FROM debian:bullseye-slim

# Устанавливаем минимальные зависимости
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем бинарный файл из стадии сборки
COPY --from=builder /app/auth-service .

# Копируем конфиги и другие необходимые файлы
COPY ./config ./config

# Указываем порт, который использует приложение
EXPOSE 8080

# Команда запуска
CMD ["./auth-service"]