# Stage 1: Build Bash Completion Script
FROM ruby:3.4.3-alpine3.21 as completion_builder
WORKDIR /app
RUN apk add --no-cache build-base
RUN gem install completely --no-document
COPY completely.yaml .
RUN completely generate

# Stage 2: Build the Go Binary
FROM golang:1.24.2-alpine3.21 as builder
WORKDIR /app
RUN apk add --no-cache git ca-certificates build-base gcc musl-dev

# Copy only dependency files first to leverage Docker cache
# Download dependencies (this layer is cached unless go.mod/go.sum change)
COPY go.mod go.sum ./
RUN go mod download

# Copy source code *after* dependencies are downloaded
COPY backend/ ./backend/

# Build the Go application statically linked (if possible) and stripped
# -ldflags="-w -s" reduces binary size by removing debug info
# CGO_ENABLED=0 makes it statically linked against musl libc on Alpine if no C code is used.
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o /app/hyde ./backend/cmd/cli/main.go

# Stage 3: Final Runtime Image
FROM alpine:3.21
WORKDIR /app

# # Install runtime dependencies: bash for the shell and bash-completion package
RUN apk add --no-cache bash bash-completion

# # Copy the bash completion script to the standard system-wide location
# # This makes it available automatically for bash sessions if bash-completion is sourced.
COPY --from=completion_builder /app/completely.bash /etc/bash_completion.d/hyde
COPY --from=builder /app/hyde /usr/local/bin/hyde

# Create a non-root user and group named 'hyde-user'
RUN addgroup -S hyde-user && adduser -S hyde-user -G hyde-user
USER hyde-user

# Copy the test file to be able to test the encryption/decryption
COPY test-file.txt .

ENTRYPOINT [ "sleep", "infinity" ]