FROM golang
COPY . /go/src/github.com/go-web-docker
RUN go install github.com/go-web-docker
ENTRYPOINT /go/bin/go-web-docker
