# Use an official Golang image as a base
FROM golang:alpine

# Set the working directory to /app
WORKDIR /app

# Copy the main.go file
COPY main.go .
COPY go.mod .

# Install dependencies
ENV GO111MODULE=on
RUN go get -d -v

# Build the Go binary
RUN go build -o backend main.go

# Expose the port the backend will use
EXPOSE 8080

# Run the backend when the container starts
CMD ["./backend"]
