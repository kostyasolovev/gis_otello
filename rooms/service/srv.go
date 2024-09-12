package service

import (
	"booking/rooms/entity"
	"time"
)

type RoomsSrv struct {
	repo Repository
}

type Repository interface {
	BookRoom(hotelID, roomID string, from, to time.Time) error
}

func (srv *RoomsSrv) BookRoom(room entity.Room, from, to time.Time) error {
	return srv.repo.BookRoom(room.HotelID, room.RoomID, from, to)
}

func NewRoomsService(repo Repository) *RoomsSrv {
	return &RoomsSrv{
		repo: repo,
	}
}
