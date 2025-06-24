package e2e

import (
	"encoding/json"
	"github.com/gojuno/minimock/v3"
	"homework-1/internal/app/product/mock"
	"homework-1/internal/app/server"
	"homework-1/internal/http/handler"
	"homework-1/internal/pkg/model"
	"homework-1/internal/repository"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDeleteItemE2E(t *testing.T) {
	ctrl := minimock.NewController(t)
	mockProduct := mock.NewProductValidatorMock(ctrl)
	mockProduct.ValidateProductMock.Expect(52).Return(&repository.ProductResponse{Name: "TestProduct", Price: 100}, nil)
	repo := repository.New()
	repo.AddCart(144, 52, 6)
	cartService := server.New(repo, mockProduct)
	handler := handler.New(cartService, mockProduct)

	router := handler.InitRoutes()
	srv := httptest.NewServer(router)
	defer srv.Close()

	req, err := http.NewRequest("DELETE", srv.URL+"/user/144/cart/52", nil)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("want status 200; got %v", resp.StatusCode)
	}
	cart, _ := repo.GetCart(144)
	if _, exists := cart[52]; exists {
		t.Errorf("item not deleted from cart")
	}
}

func TestGetCartE2E(t *testing.T) {
	ctrl := minimock.NewController(t)
	mockProduct := mock.NewProductValidatorMock(ctrl)
	mockProduct.ValidateProductMock.Expect(52).Return(&repository.ProductResponse{Name: "TestProduct", Price: 100}, nil)
	repo := repository.New()
	repo.AddCart(144, 52, 6)
	cartService := server.New(repo, mockProduct)
	handler := handler.New(cartService, mockProduct)

	router := handler.InitRoutes()
	srv := httptest.NewServer(router)
	defer srv.Close()

	req, err := http.NewRequest("GET", srv.URL+"/user/144/cart", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("want status 200; got %v", resp.StatusCode)
	}

	var cartResponse model.CartResponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}
	if err = json.Unmarshal(body, &cartResponse); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	expectedName := "TestProduct"
	if cartResponse.Items[0].Name != expectedName {
		t.Errorf("want %v; got %v", expectedName, cartResponse.Items[0].Name)
	}
	expectedCount := uint16(6)
	if cartResponse.Items[0].Count != expectedCount {
		t.Errorf("Expected quantity %d, got %d", expectedCount, cartResponse.Items[0].Count)
	}
}
