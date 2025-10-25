FROM golang:1.24-alpine

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy source code
COPY . .

# Build the application
RUN go mod tidy && go mod download && \
    CGO_ENABLED=0 GOOS=linux go build -o /app/main ./cmd/main.go

EXPOSE 8080

CMD ["/app/main"] 