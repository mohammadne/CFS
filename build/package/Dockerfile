# Start from the latest golang base image
FROM golang:alpine AS builder

# Install git command required for building with go1.18
RUN apk add git

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o /cfs

FROM alpine:latest

# Add Maintainer Info
LABEL maintainer="Mohammad Nasr <mohammadne.dev@gmail.com>"

# download and extract alpine-minirootfs
RUN cd /tmp && rm -rf /tmp/alpine-rootfs && mkdir /tmp/alpine-rootfs && \
    curl -O https://dl-cdn.alpinelinux.org/alpine/v3.17/releases/x86_64/alpine-minirootfs-3.17.2-x86_64.tar.gz && \
    tar xzvf /tmp/alpine-minirootfs-3.17.2-x86_64.tar.gz -C /tmp/alpine-rootfs && \
    touch /tmp/alpine-rootfs/CONTAINER_ROOT_FS

WORKDIR /app/

COPY --from=builder /cfs .

ENTRYPOINT ["./cfs"]
