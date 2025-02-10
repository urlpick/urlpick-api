FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
COPY configs/ configs/
EXPOSE 8080
CMD ["./main"]
