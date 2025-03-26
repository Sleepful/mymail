# syntax=docker/dockerfile:1

# FROM arm64v8/golang:1.23
# FROM golang:1.23-alpine AS builder
# FROM --platform=$BUILDPLATFORM golang:1.23
# FROM --platform=linux/amd64 golang:1.23
# FROM --platform=darwin/arm64v8 golang:1.23
# FROM --platform=linux/arm64 golang:1.23
FROM golang:1.23

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# COPY **/*.go ./
COPY . ./

# RUN CGO_ENABLED=0 GOARCH=arm64 GOOS=linux go build -o /docker-cmd
RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-cmd

EXPOSE 8090

RUN ["chmod", "+x", "/docker-cmd"]

# RUN go run .

CMD ["/docker-cmd"]
