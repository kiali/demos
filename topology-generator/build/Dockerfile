FROM golang:1.17 AS build

WORKDIR /go/src/app

COPY ./ .
RUN go get -d -v ./...
RUN make build-binary-image

FROM registry.access.redhat.com/ubi8/ubi-minimal:8.5

WORKDIR /

COPY --from=build /go/src/app/topogen /topogen

EXPOSE 8080

ENTRYPOINT ["/topogen"]
