## Start from a Debian image with the latest version of Go installed
## and a workspace(GOPATH) configured at /go.
FROM golang:latest
RUN go get github.com/monnand/dhkx
RUN mkdir /goclient
ADD . /goclient
WORKDIR /goclient
RUN go build
CMD ["./goclient"]
