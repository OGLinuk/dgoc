FROM golang:1.11.0
ADD . /go/src/queue
WORKDIR /go/src/queue
RUN go get
RUN go build -o queue .
EXPOSE 8003
CMD ["./queue"]