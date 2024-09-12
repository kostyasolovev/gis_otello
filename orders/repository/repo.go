package repository

import (
	"booking/orders/entity"
	"fmt"
	"github.com/google/uuid"
	"sync"
)

type OrdersRepo struct {
	mu    sync.Mutex
	data  []entity.Order
	index map[uuid.UUID]int
}

func (repo *OrdersRepo) CreateOrder(order entity.Order) (entity.Order, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	repo.data = append(repo.data, order)
	repo.index[order.OrderID] = len(repo.data) - 1

	return order, nil
}

func (repo *OrdersRepo) MarkFailed(orderID uuid.UUID) error {
	return repo.setStatus(orderID, entity.Failed)
}

func (repo *OrdersRepo) setStatus(orderID uuid.UUID, status entity.Status) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	ind, ok := repo.index[orderID]
	if !ok {
		return fmt.Errorf("setting status to order <%s>: %w", orderID.String(), entity.ErrNotFound)
	}

	if ind >= len(repo.data) {
		return fmt.Errorf("db consistency check failed, order index %d out of range (data len %d)", ind, len(repo.data))
	}

	repo.data[ind].Status = status

	return nil
}

func New() *OrdersRepo {
	return &OrdersRepo{
		index: make(map[uuid.UUID]int),
	}
}
