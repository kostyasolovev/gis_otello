package service

import (
	"booking/orders/entity"
	roomsEntity "booking/rooms/entity"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type Service struct {
	repo  Repository
	rooms RoomsService
}

type Repository interface {
	CreateOrder(order entity.Order) (entity.Order, error)
	MarkFailed(orderID uuid.UUID) error
}

type RoomsService interface {
	BookRoom(room roomsEntity.Room, from, to time.Time) error
}

func (s *Service) CreateOrder(order entity.Order) (entity.Order, error) {
	order.OrderID = uuid.New()

	newOrder, err := s.repo.CreateOrder(order)
	if err != nil {
		return entity.Order{}, fmt.Errorf("create order: %w", err)
	}

	room := roomsEntity.Room{
		HotelID: order.HotelID,
		RoomID:  order.RoomID,
	}

	err = s.rooms.BookRoom(room, order.From, order.To)
	if err != nil {
		// если не удалось забронировать комнаты - помечаем заказ как провальный
		markErr := s.repo.MarkFailed(newOrder.OrderID)
		if markErr != nil {
			err = errors.Join(err, markErr)
		}

		return entity.Order{}, err
	}

	return newOrder, nil
}

func NewService(repo Repository, rooms RoomsService) *Service {
	return &Service{
		repo:  repo,
		rooms: rooms,
	}
}
