FROM golang

ADD . /go/src/cities-api

ENV GIN_MODE release
WORKDIR /go/src/cities-api
RUN make build

ENTRYPOINT /go/src/cities-api

EXPOSE 8080
