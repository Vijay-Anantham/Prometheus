# Use an official Go runtime as a parent image
FROM golang:1.17 as builder

# Set the working directory inside the container
WORKDIR /app

# Copy the current directory contents into the container at /app
COPY ./main .
COPY ./go.mod .
COPY ./go.sum .

# Build the Go application
RUN go mod tidy
RUN go build -o main .

FROM alphine:v1

# Expose a port for the application (if necessary)
EXPOSE 8080

# Command to run the executable
CMD ["./main"]