package priceFeed

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/lightyeario/kelp/support/exchange"
	"github.com/lightyeario/kelp/support/exchange/api/assets"
)

// priceFeed allows you to fetch the price of a feed
type priceFeed interface {
	getPrice() (float64, error)
}

func priceFeedFactory(feedType string, url string) priceFeed {
	switch feedType {
	case "crypto":
		return newCMCFeed(url)
	case "fiat":
		return newFiatFeed(url)
	case "fixed":
		return newFixedFeed(url)
	case "exchange":
		// [0] = exchangeType, [1] = base, [2] = quote
		urlParts := strings.Split(url, "/")
		exchange := exchange.ExchangeFactory(urlParts[0])
		tradingPair := assets.TradingPair{
			Base:  exchange.GetAssetConverter().MustFromString(urlParts[1]),
			Quote: exchange.GetAssetConverter().MustFromString(urlParts[2]),
		}
		return newExchangeFeed(&exchange, &tradingPair)
	}
	return nil
}

func getJSON(client http.Client, url string, target interface{}) error {
	r, err := client.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}