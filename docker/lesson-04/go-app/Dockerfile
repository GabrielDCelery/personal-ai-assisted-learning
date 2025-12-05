# Build stage
FROM golang:1.24.11-bookworm AS builder
WORKDIR /app
COPY main.go .
RUN go build -o server main.go

# Runtime stage
FROM alpine:3.23.0

COPY --from=builder /app/server .

EXPOSE 8080
CMD ["./server"]
