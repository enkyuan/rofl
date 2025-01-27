package exchanges

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"lume/rofl/constants"
	"lume/rofl/models"
	"lume/rofl/utils"
)

type Coinbase struct{}

func NewCoinbaseExchanger() *Coinbase {
	return &Coinbase{}
}


func (ex *Coinbase) TrackCurrencyValue(crypto string) (*models.TrackCurrencyResponse, error) {
	// prepare the url
	url := fmt.Sprintf("%s/prices/%s-USD/spot", constants.COINBASE_URL, strings.ToUpper(crypto))

	response, err := utils.MakeAPICall(url, http.MethodGet, nil, nil)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var coinbaseResponse struct {
		Data struct {
			Base     string `json:"base"`
			Currency string `json:"currency"`
			Amount   string `json:"amount"`
		} `json:"data"`
	}

	err = json.Unmarshal(body, &coinbaseResponse)
	if err != nil {
		return nil, err
	}

	price, err := strconv.ParseFloat(coinbaseResponse.Data.Amount, 64)
	if err != nil {
		return nil, err
	}

	return &models.TrackCurrencyResponse{
		Value: price,
	}, nil

}