FROM --platform=amd64 golang:alpine AS builder

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Move to working directory /build
WORKDIR /build

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY . .

# Build the application
RUN go build -o frontend cmd/frontend/main.go

# Move to /dist directory as the place for resulting binary folder
WORKDIR /dist

# Copy binary from build to main folder
RUN cp /build/frontend .
RUN cp -R /build/config .

# Build a small image
FROM alpine

RUN mkdir /app

WORKDIR /app

COPY --from=builder /dist/frontend .
COPY --from=builder /dist/config ./config

# log dir
RUN ls -la

# Command to run
CMD ["/bin/sh", "-c", "./frontend"]