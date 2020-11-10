FROM golang:1.15

WORKDIR /go/src/voucher

COPY . .

RUN GOOS=linux go build

EXPOSE 9092

ENTRYPOINT ["./voucher"]