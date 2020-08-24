package main

import (
	"./store"
	"encoding/json"
	"fmt"

	"net/http"
	"time"
)

const (
	tinkoffApi = "https://api.tinkoff.ru/v1/currency_rates"
	sberApi    = "https://www.sberbank.ru/portalserver/proxy/?pipe=shortCachePipe&url=http%3A%2F%2Flocalhost%2Frates-web%2FrateService%2Frate%2Fcurrent%3FregionId%3D77%26rateCategory%3Dbase%26currencyCode%3D978%26currencyCode%3D840"
	alfaApi    = "https://alfabank.ru/ext-json/0.2/exchange/cash?offset=0&limit=1&mode=rest"
)

var client = &http.Client{Timeout: 10 * time.Second}

func getTinkoffRate(from, to string) (*store.Rate, error) {
	r, err := client.Get(tinkoffApi + "?from=" + from + "&to=" + to)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	tr := &store.TinkoffRate{}
	err = json.NewDecoder(r.Body).Decode(tr)
	if err != nil {
		return nil, err
	}
	return store.MakeFromTinkoff(tr), nil
}

func getSberRate(currency string, amount int) (*store.Rate, error) {
	r, err := client.Get(sberApi)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	sr := &store.SberRate{}
	err = json.NewDecoder(r.Body).Decode(sr)
	if err != nil {
		return nil, err
	}
	return store.MakeFromSber(sr, amount, currency), nil
}

func getAlfaRate(currency string) (*store.Rate, error) {
	r, err := client.Get(alfaApi)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	ar := &store.AlfaRate{}
	err = json.NewDecoder(r.Body).Decode(ar)
	if err != nil {
		return nil, err
	}
	return store.MakeFromAlfa(ar, currency), nil
}

func main() {
	rate, err := getTinkoffRate(store.USD, store.RUB)
	if err != nil {
		panic(err)
	}
	fmt.Printf("TINKOFF: %s -> %s BUY: %.2f  SELL: %.2f\n", rate.FromCurrency, rate.ToCurrency, rate.Buy, rate.Sell)

	rateSber, err := getSberRate(store.USD, 100)
	if err != nil {
		panic(err)
	}
	fmt.Printf("SBER: %s -> %s BUY: %.2f  SELL: %.2f\n", rateSber.FromCurrency, rateSber.ToCurrency, rateSber.Buy, rateSber.Sell)

	rateAlfa, err := getAlfaRate(store.USD)
	if err != nil {
		panic(err)
	}
	fmt.Printf("ALFA: %s -> %s BUY: %.2f  SELL: %.2f\n", rateAlfa.FromCurrency, rateAlfa.ToCurrency, rateAlfa.Buy, rateAlfa.Sell)
}
