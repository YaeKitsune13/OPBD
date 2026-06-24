# Используем твою версию
FROM golang:1.26.3 AS builder

WORKDIR /app

# Копируем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь код
COPY . .

# ВАЖНО: Добавляем CGO_ENABLED=0 и GOOS=linux
# Это сделает бинарник статическим, и он будет работать в Alpine
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api/main.go

# Финальный образ
FROM alpine:latest

WORKDIR /root/

# Копируем бинарник
COPY --from=builder /app/main .

# На всякий случай даем права (хотя CGO_ENABLED=0 решает 99% проблем)
RUN chmod +x ./main

EXPOSE 3000

CMD ["./main"]
