package services

import (
	"context"
	"encoding/json"
	"fmt"
	"labs/l0/cache"
	"labs/l0/database/models"
	"labs/l0/database/repository"
	"log"
)

type OrderService struct {
	repo  repository.OrderRepository
	cache cache.Cache
}

func NewOrderServiceWithCache(repo repository.OrderRepository, cache cache.Cache) *OrderService {
	return &OrderService{
		repo:  repo,
		cache: cache,
	}
}

func NewOrderService(repo repository.OrderRepository) *OrderService {
	return &OrderService{
		repo:  repo,
		cache: nil,
	}
}

func (s *OrderService) ProcessMessage(ctx context.Context, msgData []byte) error {
	var order models.Order

	if err := json.Unmarshal(msgData, &order); err != nil {
		return fmt.Errorf("failed to unmarshal order: %w", err)
	}

	if err := s.validateOrder(&order); err != nil {
		return fmt.Errorf("order validation failed: %w", err)
	}

	if err := s.repo.CreateOrder(ctx, &order); err != nil {
		return fmt.Errorf("failed to create ordew %w", err)
	}

	s.cache.Set(order.OrderUid, &order)

	log.Printf("Order %s processed successfully", order.OrderUid)

	return nil
}

func (s *OrderService) validateOrder(order *models.Order) error {
	if order.OrderUid == "" {
		return fmt.Errorf("order_uid is required")
	}
	if order.TrackNumber == "" {
		return fmt.Errorf("track_number is required")
	}
	if order.Payment.Transaction == "" {
		return fmt.Errorf("payment transaction is required")
	}
	if len(order.Items) == 0 {
		return fmt.Errorf("order must have at least one item")
	}

	return nil
}

func (s *OrderService) GetOrder(ctx context.Context, orderUid string) (*models.Order, error) {
	if val, ok := s.cache.Get(orderUid); ok {
		order, ok := val.(*models.Order)
		if ok {
			return order, nil
		}
	}

	order, err := s.repo.GetOrderByUID(ctx, orderUid)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (s *OrderService) GetAllOrders(ctx context.Context) ([]*models.Order, error) {
	orders, err := s.repo.GetAllOrders(ctx)

	if err != nil {
		return nil, err
	}

	// обновляем кэш, если нет в нём заказа
	for _, order := range orders {
		if _, ok := s.cache.Get(order.OrderUid); !ok {
			s.cache.Set(order.OrderUid, order)
		}
	}

	return orders, nil
}
