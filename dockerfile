# syntax=docker/dockerfile:1

FROM golang:1.21.8-alpine3.19

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod tidy

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux

RUN go build -o /backendapi
EXPOSE 8080
CMD ["go", "run", "main.go"] 