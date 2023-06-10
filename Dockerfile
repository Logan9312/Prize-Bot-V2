# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang:1.19

# Enable CGO
ENV CGO_ENABLED=1

# Copy the local package files to the container's workspace.
ADD . /go/src/my/app

# Build the application inside the container.
RUN go install my/app

# Run the application by default when the container starts.
ENTRYPOINT /go/bin/my/app

# Document that the service listens on port 8080.
EXPOSE 8080
