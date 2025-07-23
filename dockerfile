
FROM golang:1.24.5-alpine AS builder

RUN apk add --no-cache build-base
RUN apk add --no-cache sqlite-dev

WORKDIR /app
COPY go.mod go.sum ./
RUN ls -l && cat go.mod && cat go.sum
RUN go mod tidy
RUN go mod download
COPY . .
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init -g cmd/api/main.go
RUN go build -o /app/server ./cmd/api/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/server .
COPY --from=builder /app/docs ./docs

EXPOSE 8080

CMD ["./server"]