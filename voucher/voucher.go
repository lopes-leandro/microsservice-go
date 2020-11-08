package main

import (
	"fmt"
	"log"
	"encoding/json"
	"net/http"
)

// Voucher é o voucher
type Voucher struct {
	Code string
}
// Vouchers é usado para pegar o retorno da bagaça
type Vouchers struct {
	Voucher []Voucher
}

// Check é um método de Check
func (v Vouchers) Check(code string) string  {
	for _, item := range v.Voucher {
		if code == item.Code {
			return "valid"
		}
	}
	return "invalid"
}

// Result é usado para preencher os resultados
type Result struct {
	Status string
}

var vouchers Vouchers

func main()  {
	voucher := Voucher {
		Code: "abc",
	}
	vouchers.Voucher = append(vouchers.Voucher, voucher)

	http.HandleFunc("/", home)
	http.ListenAndServe(":9092", nil)
}

func home(w http.ResponseWriter, r *http.Request)  {
	code := r.PostFormValue("voucher")
	valid := vouchers.Check(code)

	result := Result{Status: valid}

	jsonResult, err := json.Marshal(result)
	if err != nil {
		log.Fatal("Error converting json")
	}

	fmt.Fprintf(w, string(jsonResult))
}