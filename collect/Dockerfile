FROM golang:1.10.0
ADD . /go/src/collect
WORKDIR /go/src/collect
RUN go get
RUN go build -o collect .
EXPOSE 9001
CMD ["./collect"]