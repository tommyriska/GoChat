## Start from a Debian image with the latest version of Go installed
## and a workspace(GOPATH) configured at /go.
FROM golang:latest
RUN go get github.com/monnand/dhkx
RUN mkdir /goclient
ADD . /goclient
WORKDIR /goclient
RUN go build
CMD ["./client"]

# Expose the port used in the particular GoChat server. Default for GoChat
# is port 8081 but this can be changed by user.
EXPOSE 8081
