# Stage 1: Build the Go application
FROM golang:1.22.6 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Copy .env file to the Docker image
COPY .env .env

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -o main .

# Stage 2: Create the minimal Disrootless image
FROM gcr.io/distroless/base-debian11

# Copy the pre-built binary from the builder stage
COPY --from=builder /app/main /main

# Copy .env file to the final image if it's needed
COPY --from=builder /app/.env /.

# Copy the templates directory into the final image
COPY --from=builder /app/templates /templates

# Copy the static folder into the final image
COPY --from=builder /app/static /static

# Command to run the executable
ENTRYPOINT ["/main"]

# Expose the port the app runs on
EXPOSE 9001
