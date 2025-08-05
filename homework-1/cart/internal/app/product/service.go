package product

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/time/rate"
	"homework-1/cart/internal/pkg/model"
	"homework-1/cart/internal/pkg/my_errgroup"
	"homework-1/cart/internal/pkg/retry"
	"net/http"
	"sync"
)

type ProductValidator interface {
	ValidateProduct(sku uint32) (*model.ProductResponse, error)
	ValidateProductParallel(SKUs []uint32) (map[uint32]model.ProductResponse, error)
}

type ProductService struct {
	client  *retry.RetryClient
	url     string
	token   string
	limiter *rate.Limiter
}

func NewProductService(client *retry.RetryClient, url, token string) *ProductService {
	return &ProductService{
		client:  client,
		url:     url,
		token:   token,
		limiter: rate.NewLimiter(10, 10),
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
		var body []byte
		if err = json.NewDecoder(resp.Body).Decode(&body); err != nil {
			return nil, fmt.Errorf("Error decoding response body: %v", err)
		}
		return nil, fmt.Errorf("Error %d: %s", resp.StatusCode, string(body))
	}

	var product model.ProductResponse
	if err = json.NewDecoder(resp.Body).Decode(&product); err != nil {
		return nil, fmt.Errorf("Incorrect response: %v", err)
	}

	return &product, nil
}

func (ps *ProductService) ValidateProductParallel(SKUs []uint32) (map[uint32]model.ProductResponse, error) {
	myG, ctx := my_errgroup.WithContext(context.Background())
	models := make(map[uint32]model.ProductResponse, len(SKUs))
	var mu sync.Mutex

	urlRequest := ps.url + "/get_product"

	for _, sku := range SKUs {
		sku := sku
		myG.Go(func() error {
			err := ps.limiter.Wait(ctx)
			if err != nil {
				return err
			}

			request := model.ProductRequest{
				Token:      ps.token,
				SkuProduct: sku,
			}
			jsonRequest, err := json.Marshal(request)
			if err != nil {
				return fmt.Errorf("Error marshalling request: %v", err)
			}

			req, err := http.NewRequestWithContext(ctx, http.MethodPost, urlRequest, bytes.NewBuffer(jsonRequest))
			if err != nil {
				return fmt.Errorf("Error creating request: %v", err)
			}
			req.Header.Set("Content-Type", "application/json")

			resp, err := ps.client.Do(req)
			if err != nil {
				return fmt.Errorf("Request error with retrays: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				var body []byte
				if err = json.NewDecoder(resp.Body).Decode(&body); err != nil {
					return fmt.Errorf("Error decoding response body: %v", err)
				}
				return fmt.Errorf("Error %d: %s", resp.StatusCode, string(body))
			}

			var product model.ProductResponse
			if err = json.NewDecoder(resp.Body).Decode(&product); err != nil {
				return fmt.Errorf("Incorrect response: %v", err)
			}
			mu.Lock()
			models[sku] = product
			mu.Unlock()
			return nil
		})

	}

	err := myG.Wait()
	if err != nil {
		return nil, err
	}

	return models, nil
}
