FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the Go application source code into the container
COPY . .

# Expose the port your application will run on
EXPOSE 9000

# Command to run your application
CMD ["go", "run","./cmd/web/"]