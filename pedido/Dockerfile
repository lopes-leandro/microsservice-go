FROM golang:1.15

WORKDIR /go/src/pedido

COPY . .

RUN GOOS=linux go build

EXPOSE 9093

ENTRYPOINT ["./pedido"]