package repository

import (
	"booking/rooms/entity"
	"fmt"
	"sync"
	"time"
)

// конструктор репозитория комнат.
func New() *RoomRepo {
	repo := &RoomRepo{
		rooms: make(map[string]*roomAvailability),
	}

	type internalRoom struct {
		HotelID string    `json:"hotel_id"`
		RoomID  string    `json:"room_id"`
		Date    time.Time `json:"date"`
		Quota   int       `json:"quota"`
	}

	var Availability = []internalRoom{
		{"reddison", "lux", date(2024, 1, 1), 1},
		{"reddison", "lux", date(2024, 1, 2), 1},
		{"reddison", "lux", date(2024, 1, 3), 1},
		{"reddison", "lux", date(2024, 1, 4), 1},
		{"reddison", "lux", date(2024, 1, 5), 0},
	}

	room := roomAvailability{
		Room: entity.Room{
			HotelID: "reddison",
			RoomID:  "lux",
		},
		firstDate:   date(2024, 1, 1),
		datesQuotas: make([]int, len(Availability)),
		mu:          sync.Mutex{},
	}

	for i, av := range Availability {
		room.datesQuotas[i] = av.Quota
	}

	repo.rooms[fmt.Sprintf("%s-%s", room.HotelID, room.RoomID)] = &room

	return repo
}

func date(year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}
