package cache

import (
	"order-service/internal/models"
	"testing"
	"time"
)

func TestCache_SetAndGet(t *testing.T) {
	cache := New()

	order := &models.Order{
		OrderUID:    "test-order-123",
		TrackNumber: "TRACK-123",
		DateCreated: time.Now(),
	}

	// Тестируем запись
	cache.Set(order)

	// Тестируем чтение
	got, exists := cache.Get("test-order-123")
	if !exists {
		t.Error("Order should exist in cache")
	}
	if got.OrderUID != order.OrderUID {
		t.Errorf("Expected %s, got %s", order.OrderUID, got.OrderUID)
	}
}

func TestCache_GetNonExistent(t *testing.T) {
	cache := New()

	_, exists := cache.Get("non-existent")
	if exists {
		t.Error("Non-existent order should not be found")
	}
}

func TestCache_Restore(t *testing.T) {
	cache := New()

	orders := []*models.Order{
		{OrderUID: "order-1", TrackNumber: "TRACK-1"},
		{OrderUID: "order-2", TrackNumber: "TRACK-2"},
	}

	cache.Restore(orders)

	if len(cache.GetAll()) != 2 {
		t.Errorf("Expected 2 orders, got %d", len(cache.GetAll()))
	}
}
