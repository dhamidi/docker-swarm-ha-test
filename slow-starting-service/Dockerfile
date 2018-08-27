FROM golang:1.10-alpine

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["app", "--delay-min", "1", "--delay-max", "20"]
