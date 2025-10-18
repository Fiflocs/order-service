package cache

import (
    "sync"
    "order-service/internal/models"
)

type Cache struct {
    mu     sync.RWMutex
    orders map[string]*models.Order
}

func New() *Cache {
    return &Cache{
        orders: make(map[string]*models.Order),
    }
}

func (c *Cache) Set(order *models.Order) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.orders[order.OrderUID] = order
}

func (c *Cache) Get(uid string) (*models.Order, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    order, exists := c.orders[uid]
    return order, exists
}

func (c *Cache) GetAll() map[string]*models.Order {
    c.mu.RLock()
    defer c.mu.RUnlock()
    return c.orders
}

func (c *Cache) Restore(orders []*models.Order) {
    c.mu.Lock()
    defer c.mu.Unlock()
    for _, order := range orders {
        c.orders[order.OrderUID] = order
    }
}