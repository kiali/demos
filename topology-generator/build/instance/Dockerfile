FROM golang:1.17 AS build

WORKDIR /go/src/app

COPY ./ .
RUN go get -d -v ./...
RUN go build cmd/instance/instance.go

FROM registry.access.redhat.com/ubi8/ubi-minimal:8.5

WORKDIR /

COPY --from=build /go/src/app/instance /instance

EXPOSE 8080

ENTRYPOINT ["/instance"]