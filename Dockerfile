# syntax=docker/dockerfile:1

# Build stage
FROM golang:1.23-alpine AS build

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o /app/main ./cmd/main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Copy binary and templates from build stage
COPY --from=build /app/main .
COPY --from=build /app/web/templates ./web/templates

EXPOSE 8080

CMD ["./main"]

