package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"order-service/internal/cache"
	"order-service/internal/models"
	"testing"
	"time"

	"github.com/gorilla/mux"
)

func TestHandler_GetOrder(t *testing.T) {
	cache := cache.New()
	handler := NewHandler(cache)

	// Добавляем тестовый заказ в кэш
	order := &models.Order{
		OrderUID:    "test-123",
		TrackNumber: "TRACK-123",
		DateCreated: time.Now(),
	}
	cache.Set(order)

	// Создаем запрос
	req, err := http.NewRequest("GET", "/orders/test-123", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Добавляем параметры маршрута
	req = mux.SetURLVars(req, map[string]string{"id": "test-123"})

	// Создаем ResponseRecorder
	rr := httptest.NewRecorder()

	// Вызываем хендлер
	handler.GetOrder(rr, req)

	// Проверяем статус код
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status 200, got %d", status)
	}

	// Проверяем тело ответа
	var responseOrder models.Order
	if err := json.Unmarshal(rr.Body.Bytes(), &responseOrder); err != nil {
		t.Errorf("Invalid JSON response: %v", err)
	}

	if responseOrder.OrderUID != "test-123" {
		t.Errorf("Expected order ID test-123, got %s", responseOrder.OrderUID)
	}
}

func TestHandler_GetOrderNotFound(t *testing.T) {
	cache := cache.New()
	handler := NewHandler(cache)

	req, err := http.NewRequest("GET", "/orders/non-existent", nil)
	if err != nil {
		t.Fatal(err)
	}

	req = mux.SetURLVars(req, map[string]string{"id": "non-existent"})
	rr := httptest.NewRecorder()

	handler.GetOrder(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", status)
	}
}
