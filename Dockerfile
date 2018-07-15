# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

RUN apt-get update && apt-get install -y sox && apt-get install -y libsox-fmt-mp3
# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/aphpbonn/myRecognize

# Build the app inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
RUN go install github.com/aphpbonn/myRecognize

# Run the outyet command by default when the container starts.
ENTRYPOINT /go/bin/myRecognize

# Document that the service listens on port 8080.
EXPOSE 8000