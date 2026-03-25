FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o ingestion-api .

FROM alpine:3.20

WORKDIR /app
COPY --from=builder /app/ingestion-api .

EXPOSE 8080

CMD ["./ingestion-api"]