#build stage
FROM golang:latest AS builder
WORKDIR /go/src/app
ADD . /go/src/app
RUN go get -d -v ./...
RUN go build -o /go/bin/app

#final stage
FROM gcr.io/distroless/base
LABEL Name=docker-spaces-backup Version=0.0.1 maintainer="Yash Thakkar<thakkaryash94@gmail.com>"
COPY --from=builder /go/bin/app /
CMD ["./app"]
