# syntax=docker/dockerfile:1

# --- Builder Stage ---
# Use the official Golang image to create a build artifact.
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files to download dependencies.
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application's source code.
COPY . .

# Build the Go application. CGO_ENABLED=0 is important for a static binary.
RUN CGO_ENABLED=0 GOOS=linux go build -o /main .

# --- Final Stage ---
# Use a minimal, non-root base image for a small and secure container.
FROM gcr.io/distroless/static-debian11
COPY --from=builder /main /main
CMD ["/main"]