package service

import (
	"encoding/json"
	"fmt"
	"log"
	"order-service/internal/cache"
	"order-service/internal/models"
	"order-service/internal/repository"

	"github.com/nats-io/stan.go"
)

type NatsSubscriber struct {
	sc      stan.Conn
	repo    *repository.OrderRepository
	cache   *cache.Cache
	subject string
}

func NewNatsSubscriber(sc stan.Conn, repo *repository.OrderRepository, cache *cache.Cache, subject string) *NatsSubscriber {
	return &NatsSubscriber{
		sc:      sc,
		repo:    repo,
		cache:   cache,
		subject: subject,
	}
}

func (ns *NatsSubscriber) Subscribe() (stan.Subscription, error) {
	return ns.sc.Subscribe(ns.subject, func(msg *stan.Msg) {
		log.Printf("Received message: %s", string(msg.Data))

		var order models.Order
		if err := json.Unmarshal(msg.Data, &order); err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			return
		}

		// Валидация данных
		if err := ns.validateOrder(&order); err != nil {
			log.Printf("Invalid order data: %v", err)
			return
		}

		// Сохранение в БД
		if err := ns.repo.SaveOrder(&order); err != nil {
			log.Printf("Error saving order to DB: %v", err)
			return
		}

		// Сохранение в кэш
		ns.cache.Set(&order)

		log.Printf("Order %s processed successfully", order.OrderUID)
	}, stan.DurableName("order-service"))
}

func (ns *NatsSubscriber) validateOrder(order *models.Order) error {
	if order.OrderUID == "" {
		return fmt.Errorf("order_uid is required")
	}
	if order.TrackNumber == "" {
		return fmt.Errorf("track_number is required")
	}
	if order.Payment.Transaction == "" {
		return fmt.Errorf("payment transaction is required")
	}
	return nil
}
