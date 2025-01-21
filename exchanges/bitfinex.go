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

type BitFinex struct{}

func NewBitFinexExchanger() *BitFinex {
	return &BitFinex{}
}

func (ex *BitFinex) TrackCurrencyValue(crypto string) (*models.TrackCurrencyResponse, error) {
	// prepare the url
	url := fmt.Sprintf("%s/ticker/t%s%s", constants.BITFINEX_URL, strings.ToUpper(crypto), strings.ToUpper("USD"))

	response, err := utils.MakeAPICall(url, http.MethodGet, nil, nil)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var bitfinexResponse struct {
		LastPrice string `json:"last_price"`
	}

	err = json.Unmarshal(body, &bitfinexResponse)
	if err != nil {
		return nil, err
	}

	price, err := strconv.ParseFloat(bitfinexResponse.LastPrice, 64)
	if err != nil {
		return nil, err
	}

	return &models.TrackCurrencyResponse{
		Value: price,
	}, nil

}
