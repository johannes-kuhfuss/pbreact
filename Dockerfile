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

EXPOSE 443

USER nobody:nogroup

ENTRYPOINT ["/pbreact"]
