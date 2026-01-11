FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o main .

# Etapa 2: Ejecución (Imagen ligera)
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
EXPOSE 9090
CMD ["./main"]