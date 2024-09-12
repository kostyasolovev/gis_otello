package entity

import (
	"github.com/google/uuid"
	"time"
)

type Room struct {
	HotelID string    `json:"hotel_id"`
	RoomID  string    `json:"room_id"`
	Date    time.Time `json:"date"`
	Quota   int       `json:"quota"`
}

type Status int

const (
	InProcessing Status = iota
	Successful
	Failed
)

type Order struct {
	OrderID   uuid.UUID `json:"order_id"`
	HotelID   string    `json:"hotel_id"`
	RoomID    string    `json:"room_id"`
	UserEmail string    `json:"email"`
	From      time.Time `json:"from"`
	To        time.Time `json:"to"`
	CreatedAt time.Time `json:"created_at"`
	Status    Status    `json:"status"`
}
