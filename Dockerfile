# Stage 1: Build go binary
FROM golang:1.24.1-bookworm AS builder

# Copy the source code
COPY ./backend /opt/backend

# Set the working directory
WORKDIR /opt/backend

# Install dependencies and build
RUN export PATH="$PATH:$(go env GOPATH)/bin" && \
    go build -trimpath -o ./.local/strawhousebackd .

# Stage 2: Create the final image
FROM alpine:3

# Install dependencies
RUN apk add --no-cache gcompat libstdc++

# Copy binary
COPY --from=builder /opt/backend/.local/strawhousebackd /usr/local/bin/strawhousebackd

# Entrypoint
WORKDIR /opt
CMD ["strawhousebackd"]