package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"time"
)

type AssetEvent struct {
	PaymentToken struct {
		Symbol string `json:"symbol"`
	} `json:"payment_token"`
	Asset struct {
		Collection struct {
			Slug string `json:"slug"`
		} `json:"collection"`
	} `json:"asset"`
	TotalPrice string `json:"total_price"`
	Permalink  string `json:"permalink"`
}

type ArrayStruct struct {
	AssetEvents []AssetEvent `json:"asset_events"`
}

type CollectionStats struct {
	Stats struct {
		OneDayVolume    float64 `json:"one_day_volume"`
		SevenDayVolume  float64 `json:"seven_day_volume"`
		ThirtyDayVolume float64 `json:"thirty_day_volume"`
		AllTimeVolume   float64 `json:"total_volume"`
		TotalSupply     float64 `json:"total_supply"`
		FloorPrice      float64 `json:"floor_price"`
	} `json:"stats"`
}

func FetchAssets() ArrayStruct {
	url := "https://api.opensea.io/api/v1/events?only_opensea=true&event_type=successful"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-API-KEY", "11671121b01f4beb9317229a88785834")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var assets ArrayStruct
	json.Unmarshal(body, &assets)
	return assets
}

func FetchCollectionStats(slug string) CollectionStats {
	url := "https://api.opensea.io/api/v1/collection/" + slug + "/stats"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-API-KEY", "11671121b01f4beb9317229a88785834")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	var collectionStats CollectionStats
	json.Unmarshal(body, &collectionStats)

	return collectionStats
}

func scientificNotationToUInt(scientificNotation string) (uint, error) {
	flt, _, err := big.ParseFloat(scientificNotation, 10, 0, big.ToNearestEven)
	if err != nil {
		return 0, err
	}
	fltVal := fmt.Sprintf("%.0f", flt)
	intVal, err := strconv.ParseInt(fltVal, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(intVal), nil
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func main() {
	collections := []string{}

	min := os.Args[1]
	max := os.Args[2]

	i := 0
	f, _ := os.Create("potential.txt")
	for i > -1 {
		assets := FetchAssets()

		for i := 0; i < len(assets.AssetEvents); i++ {
			salePrice, _ := strconv.ParseFloat(assets.AssetEvents[i].TotalPrice, 32)

			fuck, _ := strconv.ParseFloat(min, 64)
			fuckie, _ := strconv.ParseFloat(max, 64)

			if salePrice/1000000000000000000 >= fuck && salePrice/1000000000000000000 <= fuckie {

				checker := contains(collections, assets.AssetEvents[i].Asset.Collection.Slug)

				if checker != true {

					collections = append(collections, assets.AssetEvents[i].Asset.Collection.Slug)

					stats := FetchCollectionStats(assets.AssetEvents[i].Asset.Collection.Slug)

					f.WriteString(string(assets.AssetEvents[i].Asset.Collection.Slug))
					f.WriteString(", 1d V: " + fmt.Sprintf("%f", stats.Stats.OneDayVolume))
					f.WriteString(", 7d V: " + fmt.Sprintf("%f", stats.Stats.SevenDayVolume))
					f.WriteString(", 30d V: " + fmt.Sprintf("%f", stats.Stats.ThirtyDayVolume))
					f.WriteString(", allTime V: " + fmt.Sprintf("%f", stats.Stats.AllTimeVolume))
					f.WriteString(", totalSupply: " + fmt.Sprintf("%f ", stats.Stats.TotalSupply))
					f.WriteString(", floorPrice: " + fmt.Sprintf("%f", stats.Stats.FloorPrice) + "\n\n")
				}
			}
		}
		time.Sleep(946 * time.Millisecond)
	}
}
