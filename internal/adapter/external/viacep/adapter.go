package viacep

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"rfoh/cloud-run/internal/domain/entity"
)

const viaCEPURL = "http://viacep.com.br/ws/%s/json/"

type viaCEPResponse struct {
	CEP   string `json:"cep"`
	City  string `json:"localidade"`
	State string `json:"uf"`
	Error bool   `json:"erro"`
}

type ViaCEPAdapter struct {
	client *http.Client
}

func NewViaCEPAdapter(client *http.Client) *ViaCEPAdapter {
	return &ViaCEPAdapter{client: client}
}

func (a *ViaCEPAdapter) FindLocationByCEP(cep *entity.CEP) (*entity.Location, error) {
	url := fmt.Sprintf(viaCEPURL, cep.Value())
	resp, err := a.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data viaCEPResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	if data.Error || data.City == "" {
		return nil, &entity.CEPNotFoundError{CEP: cep.Value()}
	}

	return entity.NewLocation(data.City, data.State), nil
}
