FROM golang:1.25.4 AS builder
WORKDIR /src
COPY . ./
RUN CGO_ENABLED=0 go build -o cfs && mv cfs /usr/bin

# STEP 2 build a small image
FROM alpine:3.22.2
LABEL maintainer="Mohammad Nasr <mohammadne.dev@gmail.com>"
RUN apk add --no-cache bind-tools busybox-extras
COPY --from=builder /usr/bin/cfs /usr/bin/cfs
ENV USER=root
ENTRYPOINT ["/usr/bin/cfs"]
