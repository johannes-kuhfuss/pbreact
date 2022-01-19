# syntax=docker/dockerfile:1

##
## Build
##
FROM golang:1.17-alpine AS build

# Setup ENV
WORKDIR /app

# Download prereqs
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy sources
COPY . .

# Build
RUN go build -o /pbreact

##
## Deploy
##
FROM alpine:latest

WORKDIR /

# Copy certificate
RUN mkdir -p /cert
COPY ./cert.pem /cert/cert.pem
COPY ./cert.key /cert/cert.key
RUN chown nobody:nobody /cert/cert.pem
RUN chown nobody:nobody /cert/cert.key
RUN chmod 0400 /cert/cert.pem
RUN chmod 0400 /cert/cert.key

# Copy executable
COPY --from=build /pbreact /pbreact
RUN chmod +x /pbreact

# Expose ports
EXPOSE 8080
EXPOSE 8443

# Runtime config
USER nobody:nobody

ENTRYPOINT ["/pbreact"]
