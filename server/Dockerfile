## Start from a Debian image with the latest version of Go installed
## and a workspace(GOPATH) configured at /go.
FROM golang:latest
RUN mkdir /goserver
ADD . /goserver
WORKDIR /goserver
RUN go build
CMD ["./goserver"]

# Expose the port used in the particular GoChat server. Default for GoChat
# is port 8081 but this can be changed by user.
EXPOSE 8081
