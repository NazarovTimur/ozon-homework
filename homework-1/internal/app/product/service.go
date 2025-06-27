package product

import (
	"bytes"
	"encoding/json"
	"fmt"
	"homework-1/internal/pkg/model"
	"homework-1/internal/pkg/retry"
	"io"
	"net/http"
)

type ProductValidator interface {
	ValidateProduct(sku uint32) (*model.ProductResponse, error)
}

type ProductService struct {
	client *retry.RetryClient
	url    string
	token  string
}

func NewProductService(client *retry.RetryClient, url, token string) *ProductService {
	return &ProductService{
		client: client,
		url:    url,
		token:  token,
	}
}

func (ps *ProductService) ValidateProduct(sku uint32) (*model.ProductResponse, error) {
	request := model.ProductRequest{
		Token:      ps.token,
		SkuProduct: sku,
	}

	jsonRequest, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("Error marshalling request: %v", err)
	}

	urlRequest := ps.url + "/get_product"
	req, err := http.NewRequest(http.MethodPost, urlRequest, bytes.NewBuffer(jsonRequest))
	if err != nil {
		return nil, fmt.Errorf("Error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := ps.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Request error with retrays: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Error %d: %s", resp.StatusCode, string(body))
	}

	var product model.ProductResponse
	if err = json.NewDecoder(resp.Body).Decode(&product); err != nil {
		return nil, fmt.Errorf("Incorrect response: %v", err)
	}

	return &product, nil
}
