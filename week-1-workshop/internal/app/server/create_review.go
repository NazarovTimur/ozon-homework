package server

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"gitlab.ozon.dev/14/week-1-workshop/internal/pkg/reviews/model"
	"io"
	"log"
	"net/http"
	"strconv"
)

type CreateReviewRequest struct {
	SKU     int64     `json:"sku"`
	Comment string    `json:"comment"`
	UserID  uuid.UUID `json:"user_id"`
}

type CreateReviewResponse struct {
	SKU     int64     `json:"sku"`
	Comment string    `json:"comment"`
	UserID  uuid.UUID `json:"user_id"`
}

func (s *Server) CreateReview(w http.ResponseWriter, r *http.Request) {
	rawID := r.PathValue("id")
	sku, err := strconv.ParseInt(rawID, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		_, errOut := fmt.Fprintf(w, "{\"message\":\"%s\"}", err)
		if errOut != nil {
			log.Printf("POST /products/{id}/reviews out failed: %s", errOut.Error())
			return
		}

		return
	}

	if sku < 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		_, errOut := fmt.Fprintf(w, "{\"message\":\"%s\"}", "sku must be more than 0")
		if errOut != nil {
			log.Printf("POST /products/{id}/reviews out failed: %s", errOut.Error())
			return
		}

		return
	}

	body, err := io.ReadAll(r.Body)

	var createRequest CreateReviewRequest

	err = json.Unmarshal(body, &createRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		_, errOut := fmt.Fprintf(w, "{\"message\":\"%s\"}", err)
		if errOut != nil {
			log.Printf("POST /products/{id}/reviews out failed: %s", errOut.Error())
			return
		}

		return
	}

	if createRequest.UserID == uuid.Nil || createRequest.SKU < 1 ||
		len(createRequest.Comment) == 0 || createRequest.SKU != sku {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		_, errOut := fmt.Fprintf(w, "{\"message\":\"%s\"}", "invalid arguments")
		if errOut != nil {
			log.Printf("POST /products/{id}/reviews out failed: %s", errOut.Error())
			return
		}

		return
	}

	inputReview := model.Review{
		SKU:     createRequest.SKU,
		Comment: createRequest.Comment,
		UserID:  createRequest.UserID,
	}

	reviewOutput, err := s.reviewService.AddReview(r.Context(), inputReview)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		_, errOut := fmt.Fprintf(w, "{\"message\":\"%s\"}", err)
		if errOut != nil {
			log.Printf("POST /products/{id}/reviews out failed: %s", errOut.Error())
			return
		}

		return
	}

	rawResponse, err := json.Marshal(&CreateReviewResponse{
		SKU:     reviewOutput.SKU,
		Comment: reviewOutput.Comment,
		UserID:  reviewOutput.UserID,
	})

	fmt.Fprint(w, string(rawResponse))
}
