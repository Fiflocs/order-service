package service

import (
	"order-service/internal/models"
	"testing"
)

func TestNatsSubscriber_ValidateOrder(t *testing.T) {
	ns := &NatsSubscriber{}

	tests := []struct {
		name    string
		order   *models.Order
		wantErr bool
	}{
		{
			name: "Valid order",
			order: &models.Order{
				OrderUID:    "test-123",
				TrackNumber: "TRACK-123",
				Payment:     models.Payment{Transaction: "txn-123"},
			},
			wantErr: false,
		},
		{
			name: "Missing order UID",
			order: &models.Order{
				TrackNumber: "TRACK-123",
				Payment:     models.Payment{Transaction: "txn-123"},
			},
			wantErr: true,
		},
		{
			name: "Missing track number",
			order: &models.Order{
				OrderUID: "test-123",
				Payment:  models.Payment{Transaction: "txn-123"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ns.validateOrder(tt.order)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateOrder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
