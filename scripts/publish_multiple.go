package main

import (
	"encoding/json"
	"fmt"
	"log"
	"order-service/internal/config"
	"time"

	"github.com/nats-io/stan.go"
)

func main() {
	cfg := config.Load()

	sc, err := stan.Connect(cfg.NATS.ClusterID, "publisher-multi", stan.NatsURL(cfg.NATS.URL))
	if err != nil {
		log.Fatal(err)
	}
	defer sc.Close()

	// Красивые примеры заказов с понятными ID
	orders := []map[string]interface{}{
		createOrder("ORD-2024-001", "Alex Johnson", "MacBook Pro 16\"", 19999, "TRK-001-ABC"),
		createOrder("ORD-2024-002", "Maria Garcia", "iPhone 15 Pro", 8999, "TRK-002-DEF"),
		createOrder("ORD-2024-003", "John Smith", "Samsung Galaxy S24", 7999, "TRK-003-GHI"),
		createOrder("ORD-2024-004", "Emma Wilson", "Sony Headphones", 2999, "TRK-004-JKL"),
		createOrder("ORD-2024-005", "Mike Brown", "iPad Air", 5999, "TRK-005-MNO"),
	}

	// Публикуем заказы
	for _, order := range orders {
		orderData, err := json.Marshal(order)
		if err != nil {
			log.Printf("Error marshaling order: %v", err)
			continue
		}

		if err := sc.Publish(cfg.NATS.Subject, orderData); err != nil {
			log.Printf("Error publishing order: %v", err)
		} else {
			fmt.Printf("✅ Published order: %s\n", order["order_uid"])
		}

		time.Sleep(300 * time.Millisecond)
	}

	fmt.Println("🎉 All sample orders published successfully!")
	fmt.Println("📋 Sample Order IDs: ORD-2024-001, ORD-2024-002, ORD-2024-003, ORD-2024-004, ORD-2024-005")
}

func createOrder(orderUID, customerName, productName string, amount int, trackNumber string) map[string]interface{} {
	return map[string]interface{}{
		"order_uid":    orderUID,
		"track_number": trackNumber,
		"entry":        "WBIL",
		"delivery": map[string]interface{}{
			"name":    customerName,
			"phone":   "+1-555-0101",
			"zip":     "10001",
			"city":    "New York",
			"address": "Main Street 123",
			"region":  "NY",
			"email":   fmt.Sprintf("%s@email.com", customerName),
		},
		"payment": map[string]interface{}{
			"transaction":   fmt.Sprintf("TXN-%s", orderUID),
			"request_id":    "",
			"currency":      "USD",
			"provider":      "stripe",
			"amount":        amount,
			"payment_dt":    time.Now().Unix(),
			"bank":          "chase",
			"delivery_cost": 500,
			"goods_total":   amount - 500,
			"custom_fee":    0,
		},
		"items": []map[string]interface{}{
			{
				"chrt_id":      1001,
				"track_number": trackNumber,
				"price":        amount - 500,
				"rid":          fmt.Sprintf("ITEM-%s", orderUID),
				"name":         productName,
				"sale":         0,
				"size":         "Standard",
				"total_price":  amount - 500,
				"nm_id":        10001,
				"brand":        "Apple",
				"status":       200,
			},
		},
		"locale":             "en",
		"internal_signature": "",
		"customer_id":        fmt.Sprintf("cust-%s", orderUID),
		"delivery_service":   "fedex",
		"shardkey":           "1",
		"sm_id":              101,
		"date_created":       time.Now().Format(time.RFC3339),
		"oof_shard":          "1",
	}
}
