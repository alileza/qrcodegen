# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .
RUN go build -o qrcodegen

# Minimal runtime image
FROM scratch

WORKDIR /app
COPY --from=builder /app/qrcodegen /app/qrcodegen

EXPOSE 8080

ENTRYPOINT ["/app/qrcodegen", "server", "--addr", ":8080"] 