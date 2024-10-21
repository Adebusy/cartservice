# Stage 1: Build stage
FROM golang:1.19 AS build

# Set the working directory inside the container
WORKDIR /go/src/app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go binary
RUN go build -o /go/src/app/bin/app

# Stage 2: Production stage
FROM alpine:3.18

# Copy the Go binary from the 'build' stage
COPY --from=build /go/src/app/bin/app /go/bin/app

# Set the binary as the entry point
ENTRYPOINT ["/go/bin/app"]