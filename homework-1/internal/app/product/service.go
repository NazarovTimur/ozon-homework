package product

import (
	"bytes"
	"encoding/json"
	"fmt"
	"homework-1/internal/pkg/retry"
	"io"
	"net/http"
)

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

func (ps *ProductService) ValidateProduct(sku uint32) (*ProductResponse, error) {
	request := ProductRequest{
		Token:      ps.token,
		SkuProduct: sku,
	}

	jsonRequest, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("Error marshalling request: %v", err)
	}

	req, err := http.NewRequest("POST", "http://route256.pavl.uk:8080/get_product", bytes.NewBuffer(jsonRequest))
	if err != nil {
		return nil, fmt.Errorf("Error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := ps.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Ошибка запроса с ретраями: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ошибка %d: %s", resp.StatusCode, string(body))
	}

	var product ProductResponse
	if err = json.NewDecoder(resp.Body).Decode(&product); err != nil {
		return nil, fmt.Errorf("некорректный ответ: %v", err)
	}

	return &product, nil
}
