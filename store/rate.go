package store

import (
	"fmt"
	"time"
)

const (
	tinkoffGroup = "SavingAccountTransfers"
	EUR          = "EUR"
	RUB          = "RUB"
	USD          = "USD"
)

type TinkoffRate struct {
	ResultCode string `json:"resultCode"`
	Payload    struct {
		LastUpdate struct {
			Milliseconds int64 `json:"milliseconds"`
		} `json:"lastUpdate"`
		Rates []struct {
			Category     string `json:"category"`
			FromCurrency struct {
				Code    int    `json:"code"`
				Name    string `json:"name"`
				StrCode string `json:"strCode"`
			} `json:"fromCurrency"`
			ToCurrency struct {
				Code    int    `json:"code"`
				Name    string `json:"name"`
				StrCode string `json:"strCode"`
			} `json:"toCurrency"`
			Buy  float32 `json:"buy,omitempty"`
			Sell float32 `json:"sell,omitempty"`
		} `json:"rates"`
	} `json:"payload"`
	TrackingID string `json:"trackingId"`
}

type SberRate struct {
	Base struct {
		Num840 struct {
			Num0 struct {
				IsoCur          string  `json:"isoCur"`
				CurrencyName    string  `json:"currencyName"`
				CurrencyNameEng string  `json:"currencyNameEng"`
				RateType        string  `json:"rateType"`
				CategoryCode    string  `json:"categoryCode"`
				Scale           int     `json:"scale"`
				BuyValue        float32 `json:"buyValue"`
				SellValue       float32 `json:"sellValue"`
				ActiveFrom      int64   `json:"activeFrom"`
				BuyValuePrev    float32 `json:"buyValuePrev"`
				SellValuePrev   float32 `json:"sellValuePrev"`
				AmountFrom      int     `json:"amountFrom"`
				AmountTo        float64 `json:"amountTo"`
			} `json:"0"`
			Num1000 struct {
				IsoCur          string  `json:"isoCur"`
				CurrencyName    string  `json:"currencyName"`
				CurrencyNameEng string  `json:"currencyNameEng"`
				RateType        string  `json:"rateType"`
				CategoryCode    string  `json:"categoryCode"`
				Scale           int     `json:"scale"`
				BuyValue        float32 `json:"buyValue"`
				SellValue       float32 `json:"sellValue"`
				ActiveFrom      int64   `json:"activeFrom"`
				BuyValuePrev    float32 `json:"buyValuePrev"`
				SellValuePrev   float32 `json:"sellValuePrev"`
				AmountFrom      int     `json:"amountFrom"`
				AmountTo        float64 `json:"amountTo"`
			} `json:"1000"`
		} `json:"840"`
		Num978 struct {
			Num0 struct {
				IsoCur          string  `json:"isoCur"`
				CurrencyName    string  `json:"currencyName"`
				CurrencyNameEng string  `json:"currencyNameEng"`
				RateType        string  `json:"rateType"`
				CategoryCode    string  `json:"categoryCode"`
				Scale           int     `json:"scale"`
				BuyValue        float32 `json:"buyValue"`
				SellValue       float32 `json:"sellValue"`
				ActiveFrom      int64   `json:"activeFrom"`
				BuyValuePrev    float32 `json:"buyValuePrev"`
				SellValuePrev   float32 `json:"sellValuePrev"`
				AmountFrom      int     `json:"amountFrom"`
				AmountTo        float64 `json:"amountTo"`
			} `json:"0"`
			Num1000 struct {
				IsoCur          string  `json:"isoCur"`
				CurrencyName    string  `json:"currencyName"`
				CurrencyNameEng string  `json:"currencyNameEng"`
				RateType        string  `json:"rateType"`
				CategoryCode    string  `json:"categoryCode"`
				Scale           int     `json:"scale"`
				BuyValue        float32 `json:"buyValue"`
				SellValue       float32 `json:"sellValue"`
				ActiveFrom      int64   `json:"activeFrom"`
				BuyValuePrev    float32 `json:"buyValuePrev"`
				SellValuePrev   float32 `json:"sellValuePrev"`
				AmountFrom      int     `json:"amountFrom"`
				AmountTo        float64 `json:"amountTo"`
			} `json:"1000"`
		} `json:"978"`
	} `json:"base"`
}

type AlfaRate struct {
	Usd []struct {
		Type  string  `json:"type"`
		Date  string  `json:"date"`
		Value float32 `json:"value"`
		Order string  `json:"order"`
	} `json:"usd"`
	Eur []struct {
		Type  string  `json:"type"`
		Date  string  `json:"date"`
		Value float32 `json:"value"`
		Order string  `json:"order"`
	} `json:"eur"`
	Chf []struct {
		Type  string  `json:"type"`
		Date  string  `json:"date"`
		Value float32 `json:"value"`
		Order string  `json:"order"`
	} `json:"chf"`
	Gbp []struct {
		Type  string  `json:"type"`
		Date  string  `json:"date"`
		Value float32 `json:"value"`
		Order string  `json:"order"`
	} `json:"gbp"`
}

type Rate struct {
	Owner        string  `json:"owner"`
	Buy          float32 `json:"buy,omitempty"`
	Sell         float32 `json:"sell,omitempty"`
	ToCurrency   string  `json:"toCurrency"`
	FromCurrency string  `json:"fromCurrency"`
	LastUpdate   int64   `json:"lastUpdate"`
}

func MakeFromTinkoff(tinkoffRate *TinkoffRate) *Rate {
	rates := tinkoffRate.Payload.Rates
	for _, r := range rates {
		if r.Category == tinkoffGroup {
			return &Rate{
				Owner:        "TINKOFF",
				Buy:          r.Buy,
				Sell:         r.Sell,
				ToCurrency:   r.ToCurrency.Name,
				FromCurrency: r.FromCurrency.Name,
				LastUpdate:   tinkoffRate.Payload.LastUpdate.Milliseconds,
			}
		}
	}
	return nil
}

func MakeFromSber(sberRate *SberRate, amount int, currency string) *Rate {
	rate := &Rate{Owner: "SBER"}
	rate.ToCurrency = RUB
	rate.FromCurrency = currency
	switch currency {
	case EUR:
		if amount >= 1000 {
			rate.Sell = sberRate.Base.Num978.Num1000.SellValue
			rate.Buy = sberRate.Base.Num978.Num1000.BuyValue
			rate.LastUpdate = sberRate.Base.Num978.Num1000.ActiveFrom
		} else {
			rate.Sell = sberRate.Base.Num978.Num0.SellValue
			rate.Buy = sberRate.Base.Num978.Num0.BuyValue
			rate.LastUpdate = sberRate.Base.Num978.Num0.ActiveFrom
		}
	case USD:
		if amount >= 1000 {
			rate.Sell = sberRate.Base.Num840.Num1000.SellValue
			rate.Buy = sberRate.Base.Num840.Num1000.BuyValue
			rate.LastUpdate = sberRate.Base.Num840.Num1000.ActiveFrom
		} else {
			rate.Sell = sberRate.Base.Num840.Num0.SellValue
			rate.Buy = sberRate.Base.Num840.Num0.BuyValue
			rate.LastUpdate = sberRate.Base.Num840.Num0.ActiveFrom
		}
	}
	return rate
}

func MakeFromAlfa(alfaRate *AlfaRate, currency string) *Rate {
	rate := &Rate{Owner: "ALFA"}
	rate.ToCurrency = RUB
	rate.FromCurrency = currency
	switch currency {
	case USD:
		rate.Buy = alfaRate.Usd[0].Value
		rate.Sell = alfaRate.Usd[1].Value
	case EUR:
		rate.Buy = alfaRate.Eur[0].Value
		rate.Sell = alfaRate.Eur[1].Value
	}
	dateString := ""
	s := alfaRate.Usd[0].Date
	for i, r := range s {
		if s[i] != ' ' {
			dateString += string(r)
		}
		if len(dateString) == 10 {
			dateString += "T"
		}
	}
	dateString += ".03Z"
	fmt.Println("TIME:", dateString)
	dateTime, err := time.Parse(time.RFC3339, dateString)
	if err != nil {
		fmt.Println("err:", err)
		rate.LastUpdate = 0
		return rate
	}
	fmt.Println("dateTime:", dateTime)
	rate.LastUpdate = dateTime.UnixNano() / 1000000

	return rate
}
