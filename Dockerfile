# Use an official Go runtime as a parent image
FROM golang:latest

# Download latest listing of available packages:
RUN apt-get -y update

# Upgrade already installed packages:
# RUN apt-get -y upgrade

# Install a new package:
RUN apt-get -y install netcat-traditional

# Set the working directory inside the container
WORKDIR /app

# Copy the local code from the Go directory to the container
COPY Go/ .

# Build the Go application
RUN go build -o main .

# Expose a port (if needed)
EXPOSE 3000

# Define the command to run the executable
CMD ["./main"]
