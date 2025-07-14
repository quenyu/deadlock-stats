# Development stage
FROM golang:1.21-alpine as development

WORKDIR /app

# Install development tools
RUN apk add --no-cache git

# Install air for hot reload
RUN go install github.com/cosmtrek/air@latest

# Copy go mod files
COPY go.mod go.sum* ./
RUN go mod download

# Copy source code
COPY . .

# Run with hot reload
EXPOSE 8080
CMD ["air", "-c", ".air.toml"]

# Build stage
FROM golang:1.21-alpine as build

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum* ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/app ./cmd/main.go

# Production stage
FROM alpine:latest as production

# Add necessary runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy binary from build stage
COPY --from=build /bin/app /app/

# Copy config
COPY config/config.example.yaml /app/config/config.yaml

EXPOSE 8080

CMD ["/app/app"] 