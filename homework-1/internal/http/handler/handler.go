package handler

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"homework-1/internal/app/product"
	"homework-1/internal/app/server"
	"homework-1/internal/pkg/response"
	"io"
	"net/http"
	"strconv"
)

type Handler struct {
	cartService    *server.Server
	productService *product.ProductService
}

func New(cartService *server.Server, productService *product.ProductService) *Handler {
	return &Handler{
		cartService:    cartService,
		productService: productService,
	}
}

type CreateReviewRequest struct {
	Count uint16 `json:"count" validate:"required,gt=0"`
}

func parseIDFromPath(r *http.Request, key string) (int64, error) {
	idStr := r.PathValue(key)
	return strconv.ParseInt(idStr, 10, 64)
}

func (h *Handler) AddItemToCart(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Ожидается application/json")
		return
	}

	userId, err := parseIDFromPath(r, "user_id")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	skuId := r.PathValue("sku_id")
	sku, err := strconv.ParseInt(skuId, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var validate = validator.New()
	var createRequest CreateReviewRequest
	err = json.Unmarshal(body, &createRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if sku < 1 || userId < 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = validate.Struct(createRequest)
	if err != nil {
		response.WriteError(w, http.StatusPreconditionFailed, "invalid sku")
		return
	}
	total, existed := h.cartService.Add(userId, uint32(sku), createRequest.Count)
	if total == 0 {
		w.WriteHeader(http.StatusBadRequest)
	}

	if !existed {
		fmt.Fprintf(w, "must add %d item", total)
	} else {
		fmt.Fprintf(w, "must add %d more item, %d - must be %d items", createRequest.Count, sku, total)
	}
}

func (h *Handler) DeleteItemFromCart(w http.ResponseWriter, r *http.Request) {
	userId, err := parseIDFromPath(r, "user_id")
	if err != nil || userId < 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	skuId := r.PathValue("sku_id")
	sku, err := strconv.ParseInt(skuId, 10, 64)
	if err != nil || sku < 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = h.productService.ValidateProduct(uint32(sku))
	if err != nil {
		response.WriteError(w, http.StatusPreconditionFailed, "invalid sku")
		return
	}
	h.cartService.Remove(userId, uint32(sku))
	fmt.Fprint(w, "must delete item from cart")
}

func (h *Handler) ClearCart(w http.ResponseWriter, r *http.Request) {
	userId, err := parseIDFromPath(r, "user_id")
	if err != nil || userId < 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	h.cartService.Clear(userId)
	fmt.Fprint(w, "must delete cart")
}

func (h *Handler) GetCart(w http.ResponseWriter, r *http.Request) {
	userId, err := parseIDFromPath(r, "user_id")
	if err != nil || userId < 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	respon, err := h.cartService.Get(userId)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(respon)
}
