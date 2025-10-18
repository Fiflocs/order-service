package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"order-service/internal/config"
	"time"

	"github.com/nats-io/stan.go"
)

func main() {
	cfg := config.Load()

	sc, err := stan.Connect(cfg.NATS.ClusterID, "publisher-100", stan.NatsURL(cfg.NATS.URL))
	if err != nil {
		log.Fatal(err)
	}
	defer sc.Close()

	// Списки для генерации случайных данных
	firstNames := []string{"Alex", "Maria", "John", "Emma", "Mike", "Sarah", "David", "Lisa", "Chris", "Anna"}
	lastNames := []string{"Johnson", "Garcia", "Smith", "Wilson", "Brown", "Davis", "Miller", "Taylor", "Anderson", "Thomas"}
	cities := []string{"New York", "Los Angeles", "Chicago", "Houston", "Phoenix", "Philadelphia", "San Antonio", "San Diego", "Dallas", "San Jose"}
	products := []struct {
		name  string
		brand string
		price int
	}{
		{"MacBook Pro 16\"", "Apple", 19999},
		{"iPhone 15 Pro", "Apple", 8999},
		{"Samsung Galaxy S24", "Samsung", 7999},
		{"Sony WH-1000XM5", "Sony", 3499},
		{"iPad Air", "Apple", 5999},
		{"PlayStation 5", "Sony", 4999},
		{"Xbox Series X", "Microsoft", 4999},
		{"Nintendo Switch", "Nintendo", 2999},
		{"AirPods Pro", "Apple", 2499},
		{"Galaxy Watch", "Samsung", 1999},
		{"Surface Laptop", "Microsoft", 12999},
		{"Kindle Paperwhite", "Amazon", 1299},
		{"GoPro Hero 12", "GoPro", 3999},
		{"DJI Mini 3", "DJI", 4699},
		{"Bose QuietComfort", "Bose", 3299},
	}

	// Новые статусы товаров
	itemStatuses := []string{"pending", "processing", "shipped", "delivered"}

	// Статусы заказов
	orderStatuses := []int{1, 2, 3} // 1=In Store, 2=In Transit, 3=Delivered

	rand.Seed(time.Now().UnixNano())

	// Создаем 100+ заказов
	for i := 1; i <= 120; i++ {
		orderUID := fmt.Sprintf("ORD-2024-%03d", i)
		firstName := firstNames[rand.Intn(len(firstNames))]
		lastName := lastNames[rand.Intn(len(lastNames))]
		customerName := fmt.Sprintf("%s %s", firstName, lastName)
		city := cities[rand.Intn(len(cities))]

		// Случайное количество товаров (1-4)
		itemCount := rand.Intn(4) + 1
		var items []map[string]interface{}
		totalAmount := 0

		// Определяем статус заказа на основе товаров
		orderStatus := orderStatuses[rand.Intn(len(orderStatuses))]

		for j := 0; j < itemCount; j++ {
			product := products[rand.Intn(len(products))]
			quantity := rand.Intn(3) + 1
			sale := rand.Intn(30)
			discountedPrice := product.price * (100 - sale) / 100

			itemTotal := discountedPrice * quantity

			// Статус товара в зависимости от статуса заказа
			var itemStatus string
			switch orderStatus {
			case 1: // In Store
				itemStatus = itemStatuses[rand.Intn(2)]
			case 2: // In Transit
				itemStatus = "shipped"
			case 3: // Delivered
				itemStatus = "delivered"
			}

			items = append(items, map[string]interface{}{
				"chrt_id":      1000 + i*10 + j,
				"track_number": fmt.Sprintf("TRK-%03d-%s", i, string('A'+j)),
				"price":        product.price,
				"rid":          fmt.Sprintf("ITEM-%s-%d", orderUID, j+1),
				"name":         product.name,
				"sale":         sale,
				"size":         "Standard",
				"total_price":  itemTotal,
				"nm_id":        10000 + i*10 + j,
				"brand":        product.brand,
				"status":       itemStatus,
				"quantity":     quantity,
			})

			totalAmount += itemTotal
		}

		// Добавляем стоимость доставки
		deliveryCost := 500
		totalAmount += deliveryCost

		order := map[string]interface{}{
			"order_uid":    orderUID,
			"track_number": fmt.Sprintf("TRK-%03d-MAIN", i),
			"entry":        "WBIL",
			"delivery": map[string]interface{}{
				"name":    customerName,
				"phone":   fmt.Sprintf("+1-555-%04d", 1000+i),
				"zip":     fmt.Sprintf("%05d", 10000+i),
				"city":    city,
				"address": fmt.Sprintf("Main St %d", i),
				"region":  "NY",
				"email":   fmt.Sprintf("%s.%s@email.com", firstName, lastName),
			},
			"payment": map[string]interface{}{
				"transaction":   fmt.Sprintf("TXN-%s", orderUID),
				"request_id":    "",
				"currency":      "USD",
				"provider":      "stripe",
				"amount":        totalAmount,
				"payment_dt":    time.Now().Add(-time.Duration(i) * time.Hour).Unix(),
				"bank":          "chase",
				"delivery_cost": deliveryCost,
				"goods_total":   totalAmount - deliveryCost,
				"custom_fee":    0,
			},
			"items":              items,
			"locale":             "en",
			"internal_signature": "",
			"customer_id":        fmt.Sprintf("CUST-%03d", i),
			"delivery_service":   "fedex",
			"shardkey":           fmt.Sprintf("%d", (i%10)+1),
			"sm_id":              100 + i,
			"date_created":       time.Now().Add(-time.Duration(i) * time.Hour).Format(time.RFC3339),
			"oof_shard":          "1",
			"status":             orderStatus,
		}

		orderData, err := json.Marshal(order)
		if err != nil {
			log.Printf("Error marshaling order %s: %v", orderUID, err)
			continue
		}

		if err := sc.Publish(cfg.NATS.Subject, orderData); err != nil {
			log.Printf("Error publishing order %s: %v", orderUID, err)
		} else {
			fmt.Printf("✅ Published order: %s (%d items, status: %d)\n", orderUID, itemCount, orderStatus)
		}

		// Небольшая задержка между сообщениями
		if i%10 == 0 {
			time.Sleep(100 * time.Millisecond)
		}
	}

	fmt.Println("🎉 Successfully published 120 orders with new structure!")
	fmt.Println("📋 Order range: ORD-2024-001 to ORD-2024-120")
	fmt.Println("🔄 Statuses: 1=In Store, 2=In Transit, 3=Delivered")
	fmt.Println("📦 Items now include quantity and proper statuses")
}
