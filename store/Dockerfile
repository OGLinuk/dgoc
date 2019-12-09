FROM golang:1.11.0
ADD . /go/src/store
WORKDIR /go/src/store
RUN go get
RUN go build -o store .
EXPOSE 8002
CMD ["./store"]