FROM golang:1.15

WORKDIR /go/src/checkout

COPY . .

RUN GOOS=linux go build

EXPOSE 9898

ENTRYPOINT ["./checkout"]