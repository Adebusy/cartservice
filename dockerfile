#!/bin/bash
# Stage 1: Build stage
FROM golang:alpine AS build

# Set the working directory inside the container
WORKDIR /go/src/app

# Copy the Go module files and download dependencies
COPY go.mod go.sum .env ./

# Copy the source code into the container
COPY . .

ENV GOOS=linux GOARCH=amd64 CGO_ENABLED=0
# Build the Go binary
RUN go build -o /go/bin/app

# Stage 2: Production stage
FROM alpine:3.13

ENV ENVIRONMENT = "live"
ENV DATABASE_SERVER_dev = "localhost"
ENV PASSWORD_dev = "Password1"
ENV DATABASE_dev = "DigitalCartDB"
ENV USERID_dev = "postgres"
ENV PORT_dev = "5432"
ENV LISTEN_ADDR_dev = ":8080"
ENV DATABASE_SERVER_live="my-db-postgresql-nyc3-62498-do-user-17863435-0.m.db.ondigitalocean.com"
ENV PASSWORD_live = "AVNS_4p8LzBbUn5iE6NeHLQP"
ENV DATABASE_live = "cartbackeddb"
ENV USERID_live = "cartusr"
ENV PORT_live = "25060"
ENV LISTEN_ADDR_live = ":8080"

# Copy the Go binary from the 'build' stage
COPY --from=build /go/bin/app /go/bin/app

COPY .env /go/bin/app

RUN chmod +x /go/bin/app

# Set the binary as the entry point
ENTRYPOINT ["/go/bin/app"]