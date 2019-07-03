#build stage
FROM golang:1.11 AS build-env

WORKDIR /go/src/github.com/redefik/notificationmanagement

COPY . .

RUN go get -d -v ./...

RUN cd cmd && CGO_ENABLED=0 GOOS=linux go build -installsuffix cgo -o /go/bin/notificationmanagement

#production stage
FROM alpine:latest

WORKDIR /root/

COPY --from=build-env /go/bin/notificationmanagement .
COPY --from=build-env /go/src/github.com/redefik/notificationmanagement/config/config.json .

EXPOSE 80

RUN apk add ca-certificates

CMD ["./notificationmanagement", "-config=config.json"]
