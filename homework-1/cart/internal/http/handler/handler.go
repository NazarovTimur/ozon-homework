package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"homework-1/cart/internal/app/cart"
	"homework-1/cart/internal/app/product"
	"homework-1/cart/internal/pkg/errorx"
	"homework-1/cart/internal/pkg/response"
	"log"
	"net/http"
	"strconv"
)

type Handler struct {
	cartService    cart.ServiceMethods
	productService product.ProductValidator
}

func New(cartService cart.ServiceMethods, productService product.ProductValidator) *Handler {
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
		fmt.Fprintln(w, "expected application/json")
		return
	}

	userId, err := parseIDFromPath(r, "user_id")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Error parsing user id:", err)
		return
	}
	skuId := r.PathValue("sku_id")
	sku, err := strconv.ParseInt(skuId, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Error parsing sku id:", err)
		return
	}

	var createRequest CreateReviewRequest
	err = json.NewDecoder(r.Body).Decode(&createRequest)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Error parsing request:", err)
		return
	}

	var validate = validator.New()
	if sku < 1 || userId < 1 {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Error parsing sku id:", err)
		return
	}
	err = validate.Struct(createRequest)
	if err != nil {
		response.WriteError(w, http.StatusPreconditionFailed, "invalid sku")
		return
	}
	total, existed, err := h.cartService.Add(r.Context(), userId, uint32(sku), createRequest.Count)
	if errors.Is(err, errorx.ErrInsufficientStock) {
		response.WriteError(w, http.StatusPreconditionFailed, "insufficient stock")
		return
	}
	if total == 0 || err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Error adding item to cart:", err)
		return
	}

	w.WriteHeader(http.StatusOK)
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
	err = h.cartService.Remove(r.Context(), userId, uint32(sku))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
	fmt.Fprint(w, "must delete item from cart")
}

func (h *Handler) ClearCart(w http.ResponseWriter, r *http.Request) {
	userId, err := parseIDFromPath(r, "user_id")
	if err != nil || userId < 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = h.cartService.Clear(r.Context(), userId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
	fmt.Fprint(w, "must delete cart")
}

func (h *Handler) GetCart(w http.ResponseWriter, r *http.Request) {
	userId, err := parseIDFromPath(r, "user_id")
	if err != nil || userId < 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	respon, err := h.cartService.Get(r.Context(), userId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(respon)
}

func (h *Handler) CheckoutAll(w http.ResponseWriter, r *http.Request) {
	userId, err := parseIDFromPath(r, "user_id")
	if err != nil || userId < 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	respon, err := h.cartService.Checkout(r.Context(), userId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(respon)
}

func (h *Handler) InitRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("DELETE /user/{user_id}/cart/{sku_id}", h.DeleteItemFromCart)
	mux.HandleFunc("GET /user/{user_id}/cart", h.GetCart)

	return mux
}
