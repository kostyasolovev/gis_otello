package dto

import (
	"booking/orders/entity"
	"fmt"
	"github.com/google/uuid"
	"time"
)

const timeLayout = "2006-01-02"

type CreateOrderRequest struct {
	HotelID   string `json:"hotel_id"`
	RoomID    string `json:"room_id"`
	UserEmail string `json:"email"`
	From      string `json:"from"`
	To        string `json:"to"`
}

func (r CreateOrderRequest) Validate() error {
	if r.HotelID == "" {
		return fmt.Errorf("%w: HotelID is required", entity.ErrValidation)
	}

	if r.RoomID == "" {
		return fmt.Errorf("%w: RoomID is required", entity.ErrValidation)
	}

	if r.UserEmail == "" {
		return fmt.Errorf("%w: UserEmail is required", entity.ErrValidation)
	}

	if r.To == "" || r.From == "" {
		return fmt.Errorf("%w: from-to params is required", entity.ErrValidation)
	}

	return nil
}

func (r CreateOrderRequest) ToEntity() (entity.Order, error) {
	from, err := time.Parse(timeLayout, r.From)
	if err != nil {
		return entity.Order{}, fmt.Errorf("%w: from time format is invalid", entity.ErrValidation)
	}

	to, err := time.Parse(timeLayout, r.To)
	if err != nil {
		return entity.Order{}, fmt.Errorf("%w: to time format is invalid", entity.ErrValidation)
	}

	if from.After(to) {
		return entity.Order{}, fmt.Errorf("%w: from time is after to", entity.ErrValidation)
	}

	return entity.Order{
		HotelID:   r.HotelID,
		RoomID:    r.RoomID,
		UserEmail: r.UserEmail,
		From:      from,
		To:        to,
	}, nil
}

type CreateOrderResponse struct {
	OrderID   uuid.UUID `json:"order_id"`
	HotelID   string    `json:"hotel_id"`
	RoomID    string    `json:"room_id"`
	UserEmail string    `json:"email"`
	From      time.Time `json:"from"`
	To        time.Time `json:"to"`
}

func NewCreateOrderResponse(order entity.Order) CreateOrderResponse {
	return CreateOrderResponse{
		OrderID:   order.OrderID,
		HotelID:   order.HotelID,
		RoomID:    order.RoomID,
		From:      order.From,
		To:        order.To,
		UserEmail: order.UserEmail,
	}
}
