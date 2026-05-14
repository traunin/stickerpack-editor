# Build executable
FROM golang:1.26.3-alpine AS builder

RUN apk add --no-cache git ffmpeg

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o api ./cmd/main.go

# Run
FROM alpine:3.23.4

RUN apk add --no-cache ffmpeg

WORKDIR /app

COPY --from=builder /app/api .

ENV PORT=8080

EXPOSE 8080

CMD ["./api"]
