FROM golang:1.17

WORKDIR /go/src/app

COPY ./ .
RUN ls 
RUN go get -d -v ./...
RUN go install -v ./...

CMD ["generator"]