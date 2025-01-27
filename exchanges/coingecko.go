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


type CoinGecko struct{}

func NewCoinGeckoExchanger() *CoinGecko {
	return &CoinGecko{}
}

func (ex *CoinGecko) TrackCurrencyValue(crypto string) (*models.TrackCurrencyResponse, error) {
	// prepare the url
	url := fmt.Sprintf("%s/simple/price?ids=%s&vs_currencies=%s", constants.COINGECKO_URL, strings.ToLower(crypto), strings.ToLower("usd"))

	response, err := utils.MakeAPICall(url, http.MethodGet, nil, nil)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var coinGeckoResponse map[string]map[string]float64
	err = json.Unmarshal(body, &coinGeckoResponse)
	if err != nil {
		return nil, err
	}

	price := coinGeckoResponse[strings.ToLower(crypto)]["usd"]

	return &models.TrackCurrencyResponse{
		Value: price,
	}, nil
}