package repository

import (
	"context"

	"labs/l0/database"
	"labs/l0/database/models"

	"gorm.io/gorm"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *models.Order) error
	GetOrderByUID(ctx context.Context, orderUID string) (*models.Order, error)
	GetAllOrders(ctx context.Context) ([]*models.Order, error)
	DeleteOrder(ctx context.Context, orderUID string) error
}

type GormOrderRepository struct {
	db *gorm.DB
}

func NewGormOrderRepository() *GormOrderRepository {
	return &GormOrderRepository{db: database.DB}
}

func (r *GormOrderRepository) CreateOrder(ctx context.Context, order *models.Order) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(order).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *GormOrderRepository) GetOrderByUID(ctx context.Context, orderUID string) (*models.Order, error) {
	var order models.Order

	err := r.db.WithContext(ctx).
		Preload("Delivery"). //Используем всякие джоины для полной картинки
		Preload("Payment").
		Preload("Items").
		Where("order_uid = ?", orderUID).
		First(&order).Error

	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (r *GormOrderRepository) GetAllOrders(ctx context.Context) ([]*models.Order, error) {
	var orders []*models.Order

	err := r.db.WithContext(ctx).
		Preload("Delivery"). //Используем всякие джоины для полной картинки
		Preload("Payment").
		Preload("Items").
		Find(&orders).Error

	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *GormOrderRepository) DeleteOrder(ctx context.Context, orderUID string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Удаляем связанные items
		if err := tx.Where("order_uid = ?", orderUID).Delete(&models.Item{}).Error; err != nil {
			return err
		}

		// Удаляем связанные payment
		if err := tx.Where("order_uid = ?", orderUID).Delete(&models.Payment{}).Error; err != nil {
			return err
		}

		// Удаляем связанные delivery
		if err := tx.Where("order_uid = ?", orderUID).Delete(&models.Delivery{}).Error; err != nil {
			return err
		}

		// Удаляем основной заказ
		if err := tx.Where("order_uid = ?", orderUID).Delete(&models.Order{}).Error; err != nil {
			return err
		}

		return nil
	})
}
