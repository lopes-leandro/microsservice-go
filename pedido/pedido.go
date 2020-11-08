package main

import (
	"github.com/streadway/amqp"
	// "fmt"
	"encoding/json"
	"io/ioutil"
	"log"
	// "net/http"
	"net/url"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/satori/go.uuid"
	"github.com/joho/godotenv"
	"github.com/wesleywillians/go-rabbitmq/queue"
)

// Result é usado para pegar o retorno da bagaça
type Result struct {
	Status string
}

// Order é a ordem de serviço
type Order struct {
	ID uuid.UUID
	Coupon string
	CcNumber string
}

// NewOrder ...
func NewOrder() Order  {
	return Order{ID: uuid.NewV4()}
}

const (
	// InvalidVoucher ....
	InvalidVoucher = "invalid"
	// ValidVoucher ...
	ValidVoucher = "valid"
	// ConnectionError ... 
	ConnectionError = "connection error"
)


func init()  {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env")
	}
}

func main() {
	// http.HandleFunc("/", process)
	// http.ListenAndServe(":9091", nil)
	messageChannel := make(chan amqp.Delivery)
	
	rabbitMQ := queue.NewRabbitMQ()
	ch := rabbitMQ.Connect()
	defer ch.Close()

	rabbitMQ.Consume(messageChannel)

	for msg := range messageChannel {
		process(msg)
	}
}

func process(msg amqp.Delivery) {
	// voucher := r.PostFormValue("voucher")
	// ccNumber := r.PostFormValue("cc-number")

	order := NewOrder()

	json.Unmarshal(msg.Body, &order)

	resultVoucher := makeHTTPCall("http://localhost:9092", order.Coupon)

	// result := Result{Status: "declined"}

	// if ccNumber == "1" {
	// 	result.Status = "approved"
	// }

	// if resultVoucher.Status == "invalid" {
	// 	result.Status = "invalid voucher"
	// }

	switch resultVoucher.Status {
	case InvalidVoucher:
		log.Println("Order: ", order.ID, ": invalid voucher!")
	case ConnectionError:
		msg.Reject(false)
		log.Println("Order: ", order.ID, ": could not process!")
	case ValidVoucher:
		log.Println("Order: ", order.ID, ": Processed!")
	}

	// jsonData, err := json.Marshal(result)
	// if err != nil {
	// 	log.Fatal("Error processing json")
	// }

	// fmt.Fprintf(w, string(jsonData))
}

func makeHTTPCall(urlMicrosservice string, voucher string) Result {
	values := url.Values{}
	values.Add("voucher", voucher)

	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 5

	res, err := retryClient.PostForm(urlMicrosservice, values)
	if err != nil {
		result := Result{Status: ConnectionError}
		return result
	}

	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal("Error processing result")
	}

	result := Result{}

	json.Unmarshal(data, &result)

	return result
}
