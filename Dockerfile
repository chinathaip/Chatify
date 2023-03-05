# Use the official golang image with version 1.19-alpine as the base image for building the go application
FROM golang:1.19-alpine as build-base

# Set the working directory in the container
WORKDIR /app

# Copy the go module file to the container's working directory
COPY go.mod .

# Download dependencies specified in go.mod
RUN go mod download

# Copy the rest of the application code to the container's working directory
COPY . .

# Build the Go application and output the binary to /app/out/go-app
RUN go build -o ./out/go-app .

EXPOSE 1111

# Use the official Alpine Linux runtime as the base image for running the go application
FROM alpine:3.16.2

# Copy the binary built in the previous stage from build-base to the /app/go-app on the new container
# copy line 17 to the /app/go-app
COPY --from=build-base /app/out/go-app /app/go-app

# Run the go application when the container starts
CMD ["/app/go-app"]
