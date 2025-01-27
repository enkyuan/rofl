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


type Kraken struct{}

func NewKrakenExchanger() *Kraken {
	return &Kraken{}
}

func (ex *Kraken) TrackCurrencyValue(crypto string) (*models.TrackCurrencyResponse, error) {
	// prepare the url
	url := fmt.Sprintf("%s/Ticker?pair=%s%s", constants.KRAKEN_URL, strings.ToUpper(crypto), strings.ToUpper("USD"))

	response, err := utils.MakeAPICall(url, http.MethodGet, nil, nil)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var krakenResponse models.KrakenTickerResponse
	err = json.Unmarshal(body, &krakenResponse)
	if err != nil {
		return nil, err
	}

	price, err := strconv.ParseFloat(krakenResponse.Result[0].C[0], 64)
	if err != nil {
		return nil, err
	}

	return &models.TrackCurrencyResponse{
		Value: price,
	}, nil
}