FROM alpine:latest

# Add the commands needed to put your compiled go binary in the container and
# run it when the container starts.
#
# See https://docs.docker.com/engine/reference/builder/ for a reference of all
# the commands you can use in this file.
#
# In order to use this file together with the docker-compose.yml file in the
# same directory, you need to ensure the image you build gets the name
# "kadlab", which you do by using the following command:
#
# $ docker build . -t kadlab

# Use an official Go runtime as a parent image
FROM golang:latest
# Download latest listing of available packages:
RUN apt-get -y update
# Upgrade already installed packages:
#RUN apt-get -y upgrade
# Install a new package:
RUN apt-get -y install netcat-traditional

# Set the working directory in"side the container
WORKDIR /app

# Copy the local code to the container
COPY . /Go/ .

# Build the Go application
RUN go build -o main .

# Expose a port (if needed)
EXPOSE 3000

# Define the command to run the executable
CMD ["./main"]