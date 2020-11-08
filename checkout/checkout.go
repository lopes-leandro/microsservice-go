package main

import (
	"encoding/json"
	"html/template"
	// "io/ioutil"
	"log"
	"net/http"
	// "net/url"
	// "github.com/hashicorp/go-retryablehttp"
	"github.com/joho/godotenv"
	"github.com/wesleywillians/go-rabbitmq/queue"
)

// Order é a ordem de serviço
type Order struct {
	Coupon string
	CcNumber string
}

// Result é usado para pegar o retorno da bagaça
type Result struct {
	Status string
}

func init()  {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env")
	}
}

func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/process", process)
	http.ListenAndServe(":9898", nil)
}

func home(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprint(w, "<h1>Olá</h1>")
	t := template.Must(template.ParseFiles("templates/home.html"))
	t.Execute(w, Result{})
}

func process(w http.ResponseWriter, r *http.Request) {
	// log.Println(r.FormValue("voucher"))
	// log.Println(r.FormValue("cc-number"))
	// foi substituído por uma fila
	//result := makeHTTPCall("http://localhost:9091", r.FormValue("voucher"), r.FormValue("cc-number"))

	voucher := r.FormValue("voucher") 
	ccNumber := r.FormValue("cc-number")

	order := Order {
		Coupon: voucher,
		CcNumber: ccNumber,
	}

	jsonOrder, err := json.Marshal(order)
	if err != nil {
		log.Fatal("Error parsing to json")
	}

	rabbitMQ := queue.NewRabbitMQ()
	ch := rabbitMQ.Connect()
	defer ch.Close()

	err = rabbitMQ.Notify(string(jsonOrder), "appplication/json", "orders_ex", "")
	if err != nil {
		log.Fatal("Error sending message to the queue")
	}

	t := template.Must(template.ParseFiles("templates/process.html"))
	t.Execute(w, "")
}

// func makeHTTPCall(urlMicrosservice string, voucher string, ccNumber string) Result {
// 	values := url.Values{}
// 	values.Add("voucher", voucher)
// 	values.Add("cc-number", ccNumber)

// 	retryClient := retryablehttp.NewClient()
// 	retryClient.RetryMax = 5

// 	res, err := retryClient.PostForm(urlMicrosservice, values)
// 	if err != nil {
// 		result := Result{Status: "Servidor fora do ar!"}
// 		return result
// 	}

// 	defer res.Body.Close()

// 	data, err := ioutil.ReadAll(res.Body)
// 	if err != nil {
// 		log.Fatal("Error processing result")
// 	}

// 	result := Result{}

// 	json.Unmarshal(data, &result)

// 	return result
// }
