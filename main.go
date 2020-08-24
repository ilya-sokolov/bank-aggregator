package main

import (
	"./store"
	"encoding/json"
	"fmt"

	"net/http"
	"time"
)

const (
	TinkoffApi = "https://api.tinkoff.ru/v1/currency_rates"
)

var client = &http.Client{Timeout: 10 * time.Second}

func getTinkoffRate(from, to string) (*store.Rate, error) {

	r, err := client.Get(TinkoffApi)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	tr := &store.TinkoffRate{}
	err = json.NewDecoder(r.Body).Decode(tr)
	if err != nil {
		return nil, err
	}
	return store.Rate{}.MakeFromTinkoff(tr), nil
}

func main() {
	rate, err := getTinkoffRate("USD", "RUB")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s -> %s BUY: %.2f  SELL: %.2f\n", rate.FromCurrency, rate.ToCurrency, rate.Buy, rate.Sell)

}
