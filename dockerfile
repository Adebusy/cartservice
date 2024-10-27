# Step 1: Build the Go binary
FROM golang:1.19 as build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o cartservice

# Step 2: Create a minimal Docker image and add the binary
FROM alpine:latest
WORKDIR /root/
COPY --from=build /app/cartservice .
EXPOSE 8080
ENTRYPOINT ["./cartservice"]