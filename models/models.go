package models

type BinanceTickerResponse struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

type CoinbaseResponse struct {
	Data struct {
		Rates map[string]string `json:"rates"`
	} `json:"data"`
}

type TrackCurrencyRequest struct {
	Currencies []string `json:"currencies" validate:"required"`
}

type TrackCurrencyResponse struct {
	Value float64 `json:"value"`
}
