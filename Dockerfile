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

COPY --from=build /pbreact /pbreact

RUN chmod +x /pbreact

EXPOSE 80
EXPOSE 8080
EXPOSE 443
EXPOSE 8443

USER root:root

ENTRYPOINT ["/pbreact"]
