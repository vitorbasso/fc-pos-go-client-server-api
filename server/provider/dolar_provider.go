package provider

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/vitorbass93/challenge1/server/repository"
)

type DolarProvider struct {
	Endpoint        string
	DolarRepository *repository.DolarRepository
}

func NewDolarProvider(endpoint string, dolarRepository *repository.DolarRepository) *DolarProvider {
	return &DolarProvider{Endpoint: endpoint, DolarRepository: dolarRepository}
}

type Dolar struct {
	USDBRL *repository.USDBRL
}

func (d *DolarProvider) GetDolar() (*repository.USDBRL, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", d.Endpoint, nil)
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var dolar Dolar
	if err := json.NewDecoder(res.Body).Decode(&dolar); err != nil {
		return nil, err
	}
	if err := d.DolarRepository.InsertDolar(dolar.USDBRL); err != nil {
		return nil, err
	}
	return dolar.USDBRL, nil
}
