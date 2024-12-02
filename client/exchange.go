package client

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	jsoniter "github.com/json-iterator/go"
)

const (
	AppId        = "app_id"
	BaseURL      = "https://openexchangerates.org"
	LatestFile   = "api/latest.json"
	CurrencyFile = "api/currencies.json"
)

type Converter struct {
	appId string
}

func NewConverter(appId string) *Converter {
	return &Converter{
		appId: appId,
	}
}

func (c *Converter) Convert(base string, newCurrency string, value string) (float64, error) {
	latest, err := c.GetLatest(base)
	if err != nil {
		return 0, err
	}

	baseCurrencyRate, ok := latest.Rates[base]
	if !ok {
		return 0, fmt.Errorf("unknown currency %s", base)
	}

	newCurrencyRate, ok := latest.Rates[newCurrency]
	if !ok {
		return 0, fmt.Errorf("unknown currency %s", newCurrency)
	}

	amount, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, err
	}

	return (amount / baseCurrencyRate) * newCurrencyRate, nil
}

func (c *Converter) GetCurrencies() (map[string]string, error) {
	u, _ := url.ParseRequestURI(BaseURL)
	u.Path = CurrencyFile
	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	currencies := make(map[string]string)
	err = jsoniter.Unmarshal(body, &currencies)
	if err != nil {
		return nil, err
	}

	return currencies, nil
}

func (c *Converter) GetLatest(base string) (LatestResponse, error) {
	u, _ := url.ParseRequestURI(BaseURL)
	u.Path = LatestFile
	params := url.Values{}
	params.Add(AppId, c.appId)
	u.RawQuery = params.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return LatestResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return LatestResponse{}, fmt.Errorf("non 200 status code: status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return LatestResponse{}, err
	}

	var latest LatestResponse
	err = jsoniter.Unmarshal(body, &latest)
	if err != nil {
		return LatestResponse{}, err
	}

	return latest, nil
}
