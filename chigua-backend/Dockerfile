# 构建阶段
FROM golang:1.26.2 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o server ./cmd/server

# 运行阶段
FROM debian:stable-slim

WORKDIR /app

COPY --from=builder /app/server .
COPY --from=builder /app/config.yaml .
COPY --from=builder /app/resources /app/resources

EXPOSE 8080

CMD ["./server"]