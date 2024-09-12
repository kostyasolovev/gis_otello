package http

import (
	"booking/orders/entity"
	"net/http"
)

type Handler struct {
	orders Orders
	logger Logger
	*http.ServeMux
}

func New(srv Orders, logger Logger) *Handler {
	h := new(Handler)
	mux := http.NewServeMux()
	mux.HandleFunc("/orders/create", h.CreateOrder)

	h.ServeMux = mux
	h.logger = logger
	h.orders = srv

	return h
}

type Orders interface {
	CreateOrder(order entity.Order) (entity.Order, error)
}

type Logger interface {
	LogInfo(format string, v ...any)
	LogErrorf(format string, v ...any)
}
