# syntax=docker/dockerfile:1

FROM golang:1.23

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# COPY **/*.go ./
COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /go-program

EXPOSE 8090

RUN ["chmod", "+x", "/go-program"]

CMD ["/go-program", "serve", "--http=0.0.0.0:8090"]
