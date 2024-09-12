package repository

import (
	commonEntity "booking/orders/entity"
	"booking/rooms/entity"
	"errors"
	"fmt"
	"sync"
	"time"
)

type RoomRepo struct {
	// ключ - сочетание hotelID-roomID
	rooms map[string]*roomAvailability
}

type roomAvailability struct {
	entity.Room
	// дата начала расписания
	firstDate time.Time
	// расписание на каждый день
	datesQuotas []int
	// индивидуальная лочка для каждой комнаты, чтобы не блокировать весь сервис
	mu sync.Mutex
}

func (repo *RoomRepo) BookRoom(hotelID, roomID string, from, to time.Time) error {
	target, ok := repo.rooms[fmt.Sprintf("%s-%s", hotelID, roomID)]
	if !ok {
		return commonEntity.ErrNotFound
	}

	target.mu.Lock()
	defer target.mu.Unlock()

	if !isRoomAvailable(target, from, to) {
		return errors.New("room not available")
	}

	startPos, daysCount := calculateStartPos(target, from), calculateDaysCount(from, to)
	// бронируем комнату - забирая квоты
	for i := startPos; i < startPos+daysCount; i++ {
		target.datesQuotas[i]--
	}

	return nil
}

const hoursInDay = 24

func isRoomAvailable(room *roomAvailability, from, to time.Time) bool {
	if from.Before(room.firstDate) {
		return false
	}

	from = from.Truncate(24 * time.Hour)
	to = to.Truncate(24 * time.Hour)

	daysCount := calculateDaysCount(from, to)

	if daysCount <= 0 {
		return false
	}
	// если даты бронирования за пределами расписания
	startPos := calculateStartPos(room, from)
	if startPos < 0 || startPos+daysCount > len(room.datesQuotas) {
		return false
	}
	// проверяем есть ли доступные квоты в указанный период
	for i := startPos; i < startPos+daysCount; i++ {
		if room.datesQuotas[i] <= 0 {
			return false
		}
	}

	return true
}

func calculateStartPos(room *roomAvailability, from time.Time) int {
	return int(from.Sub(room.firstDate).Hours() / hoursInDay)
}

func calculateDaysCount(from, to time.Time) int {
	return int(to.Sub(from).Hours()/hoursInDay + 1)
}
