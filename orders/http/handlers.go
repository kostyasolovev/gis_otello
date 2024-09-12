package http

import (
	"booking/orders/dto"
	"booking/orders/entity"
	"encoding/json"
	"errors"
	"net/http"
)

func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		h.logger.LogErrorf("CreateOrder: json decode error: %s", err.Error())

		return
	}

	if err := req.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		h.logger.LogErrorf("CreateOrder: validation error: %s", err.Error())

		return
	}

	orderEnt, err := req.ToEntity()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		h.logger.LogErrorf("CreateOrder: convert to entity order: %s", err.Error())

		return
	}

	newOrder, err := h.orders.CreateOrder(orderEnt)
	if err != nil {
		h.logger.LogErrorf("CreateOrder: create order error: %s", err.Error())

		if errors.Is(err, entity.ErrNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		if errors.Is(err, entity.ErrValidation) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.logger.LogInfo("new order created: %v", newOrder)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	resp := dto.NewCreateOrderResponse(newOrder)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.logger.LogErrorf("failed to write response: %v", err)
	}
}
