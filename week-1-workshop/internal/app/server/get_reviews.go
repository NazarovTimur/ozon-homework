package server

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net/http"
	"strconv"
)

type GetReviewsReviewResponse struct {
	SKU     int64     `json:"sku"`
	Comment string    `json:"comment"`
	UserID  uuid.UUID `json:"user_id"`
}

type GetReviewsResponse struct {
	Reviews []GetReviewsReviewResponse `json:"reviews"`
}

func (s *Server) GetReviews(w http.ResponseWriter, r *http.Request) {
	rawID := r.PathValue("id")
	sku, err := strconv.ParseInt(rawID, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		_, errOut := fmt.Fprintf(w, "{\"message\":\"%s\"}", err)
		if errOut != nil {
			log.Printf("GET /products/{id}/reviews out failed: %s", errOut.Error())
			return
		}

		return
	}

	reviews, err := s.reviewService.GetReviews(r.Context(), sku)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		_, errOut := fmt.Fprintf(w, "{\"message\":\"%s\"}", err)
		if errOut != nil {
			log.Printf("GET /products/{id}/reviews out failed: %s", errOut.Error())
			return
		}

		return
	}

	reviewsResponse := make([]GetReviewsReviewResponse, 0, len(reviews))
	for _, reviewItem := range reviews {
		reviewsResponse = append(reviewsResponse, GetReviewsReviewResponse{
			SKU:     reviewItem.SKU,
			Comment: reviewItem.Comment,
			UserID:  reviewItem.UserID,
		})
	}

	response := GetReviewsResponse{Reviews: reviewsResponse}

	rawResponse, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		_, errOut := fmt.Fprintf(w, "{\"message\":\"%s\"}", err)
		if errOut != nil {
			log.Printf("GET /products/{id}/reviews out failed: %s", errOut.Error())
			return
		}

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	w.Write(rawResponse)
}
