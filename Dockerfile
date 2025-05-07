# Stage: Build Bash Completion Script
FROM ruby:3.4.3-alpine as completion_builder
WORKDIR /_bash_completion
RUN apk add --no-cache build-base
RUN gem install completely --no-document

COPY completely.yaml .
RUN completely generate

# Stage: Build the Go Binaries (API and CLI)
FROM golang:1.24.2-alpine as go_builder
WORKDIR /_go
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
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o /_go/hyde ./backend/cmd/cli/main.go
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o /_go/api ./backend/cmd/api/main.go

# Stage: Final Runtime 
FROM node:22-alpine as runner
RUN apk add --no-cache bash bash-completion

ARG USER_NAME=hyde
ARG GROUP_NAME=hyde
RUN addgroup -S ${GROUP_NAME}
RUN adduser -S ${USER_NAME} -G ${GROUP_NAME}

# Stage: Install pnpm
FROM runner as builder
# Install pnpm
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN wget -qO- https://get.pnpm.io/install.sh | ENV="$HOME/.bashrc" SHELL="$(which bash)" PNPM_VERSION="10.7.1" bash -
	
# Stage: Build Frontend
FROM builder as frontend_builder
USER ${USER_NAME}
WORKDIR /_frontend

# Install dependencies
COPY frontend/package.json frontend/pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile

# Build frontend
COPY frontend/ ./
RUN pnpm build
RUN pnpm prune --production

# Stage: Final Image
FROM runner as final_runner
COPY --from=completion_builder /_bash_completion/completely.bash /etc/bash_completion.d/hyde
COPY --from=go_builder /_go/hyde /usr/local/bin/hyde
COPY --from=go_builder /_go/api /usr/local/bin/api
COPY --from=frontend_builder /_frontend/dist /_frontend/dist
COPY --from=frontend_builder /_frontend/node_modules /_frontend/node_modules
COPY --from=frontend_builder /_frontend/package.json /_frontend/package.json
COPY --from=frontend_builder /pnpm /pnpm

# Create directories for playground files
RUN chown -R ${USER_NAME}:${GROUP_NAME} /home/${USER_NAME}
COPY test-file.txt /home/${USER_NAME}/test-file.txt

COPY entrypoint.bash /usr/local/bin/entrypoint.bash
ENTRYPOINT [ "/bin/bash", "/usr/local/bin/entrypoint.bash" ]