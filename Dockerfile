# Build stage
FROM golang:1.23.4-alpine AS builder
RUN apk add --no-cache git

WORKDIR /app

# Copy go.mod and go.sum first for dependency download
COPY go.mod go.sum ./

# Copy the rest of your source code
COPY . ./

# Build the binary from the cmd/main.go
RUN CGO_ENABLED=0 go build -o /go/bin/app ./cmd/main.go

# Final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates

# Copy the compiled binary
COPY --from=builder /go/bin/app /app/app

# Copy your view directory (where index.html and layouts are stored)
COPY --from=builder /app/view /app/view

# If your app references 'view/index.html', it will now exist in /app/view
WORKDIR /app

# Run the app
ENTRYPOINT ["./app"]
EXPOSE 80

LABEL Name=popliop/gopage Version=0.0.1
