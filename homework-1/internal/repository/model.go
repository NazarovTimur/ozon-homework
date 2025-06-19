package repository

type ProductRequest struct {
	Token      string `json:"token"`
	SkuProduct uint32 `json:"sku"`
}
type ProductResponse struct {
	Name  string `json:"name"`
	Price uint32 `json:"price"`
}
