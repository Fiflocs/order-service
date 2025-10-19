package http

import (
	"compress/gzip"
	"encoding/json"
	"net/http"
	"order-service/internal/cache"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type Handler struct {
	cache *cache.Cache
}

func NewHandler(cache *cache.Cache) *Handler {
	return &Handler{cache: cache}
}

func (h *Handler) GetOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderUID := vars["id"]

	order, exists := h.cache.Get(orderUID)
	if !exists {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	// Optimization
	w.Header().Set("Cache-Control", "public, max-age=120")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// Optimization
	if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		json.NewEncoder(gz).Encode(order)
	} else {
		json.NewEncoder(w).Encode(order)
	}
}

func (h *Handler) GetOrders(w http.ResponseWriter, r *http.Request) {
	orders := h.cache.GetAll()

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		json.NewEncoder(gz).Encode(orders)
	} else {
		json.NewEncoder(w).Encode(orders)
	}
}

func (h *Handler) ServeOrderPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/templates/order.html")
}

// Optimization
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	healthStatus := map[string]interface{}{
		"status":     "healthy",
		"service":    "order-service",
		"cache_size": len(h.cache.GetAll()),
		"timestamp":  time.Now().Format(time.RFC3339),
	}

	json.NewEncoder(w).Encode(healthStatus)
}
