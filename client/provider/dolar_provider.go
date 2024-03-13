package provider

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

type DolarProvider struct {
	url string
}

func NewDolarProvider(url string) *DolarProvider {
	return &DolarProvider{url: url}
}

type DolarValue struct {
	Bid string `json:"bid"`
}

func (d *DolarProvider) GetDolar() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", d.url, nil)
	if err != nil {
		return "", err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		errString := fmt.Sprintf("unexpected status from dolar provider: %s", res.Status)
		return "", errors.New(errString)
	}
	var bid DolarValue
	if err := json.NewDecoder(res.Body).Decode(&bid); err != nil {
		return "", err
	}
	log.Println("Got bid: ", bid.Bid)
	return bid.Bid, nil
}
