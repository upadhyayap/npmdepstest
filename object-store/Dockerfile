FROM golang:1.18

COPY . /go/src/object-store

WORKDIR /go/src/object-store

RUN make build

CMD ["./bin/object-store"]