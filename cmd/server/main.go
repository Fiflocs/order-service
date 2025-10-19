package main

import (
	"log"
	"net/http"
	"order-service/internal/cache"
	"order-service/internal/config"
	"order-service/internal/repository"
	"order-service/internal/service"

	httphandler "order-service/internal/delivery/http"

	"github.com/gorilla/mux"
	"github.com/nats-io/stan.go"
)

func main() {
	log.Println("Starting Order Service...")

	cfg := config.Load()

	db, err := repository.NewPostgresDB(repository.DBConfig{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.DBName,
		SSLMode:  cfg.Database.SSLMode,
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	repo := repository.NewOrderRepository(db)
	cache := cache.New()

	// Optimization
	orders, err := repo.GetAllOrders()
	if err != nil {
		log.Printf("Error restoring cache from DB: %v", err)
	} else {
		cache.Restore(orders)
		log.Printf("Cache restored with %d orders", len(orders))
	}

	// Optimization
	log.Printf("Connecting to NATS: %s", cfg.NATS.URL)
	sc, err := stan.Connect(cfg.NATS.ClusterID, cfg.NATS.ClientID, stan.NatsURL(cfg.NATS.URL))
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer sc.Close()
	log.Println("Connected to NATS successfully")

	// Optimization
	subscriber := service.NewNatsSubscriber(sc, repo, cache, cfg.NATS.Subject)
	sub, err := subscriber.Subscribe()
	if err != nil {
		log.Fatalf("Failed to subscribe: %v", err)
	}
	defer sub.Unsubscribe()
	log.Printf("Subscribed to subject: %s", cfg.NATS.Subject)

	// Optimization
	handler := httphandler.NewHandler(cache)
	router := mux.NewRouter()

	// Optimization
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/",
		http.FileServer(http.Dir("web/static/"))))

	// Optimization
	router.HandleFunc("/orders/{id}", handler.GetOrder).Methods("GET")
	router.HandleFunc("/orders", handler.GetOrders).Methods("GET")
	router.HandleFunc("/health", handler.HealthCheck).Methods("GET")
	router.HandleFunc("/", handler.ServeOrderPage)

	log.Printf("HTTP server starting on %s", cfg.HTTP.Address)
	log.Fatal(http.ListenAndServe(cfg.HTTP.Address, router))
}
