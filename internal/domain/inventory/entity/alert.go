package entity

import (
	"time"
)

// AlertLevel represents the severity level of inventory alert
type AlertLevel string

const (
	AlertLevelLow      AlertLevel = "low"
	AlertLevelMedium   AlertLevel = "medium"
	AlertLevelCritical AlertLevel = "critical"
)

// AlertType represents the type of inventory alert
type AlertType string

const (
	AlertTypeLowStock       AlertType = "low_stock"
	AlertTypeOutOfStock     AlertType = "out_of_stock"
	AlertTypeOverStock      AlertType = "over_stock"
	AlertTypeExpiringSoon   AlertType = "expiring_soon"
	AlertTypeExpired        AlertType = "expired"
)

// InventoryAlert represents an inventory warning alert
type InventoryAlert struct {
	ID             int64
	ProductID      int64
	ProductName    string
	ProductSKU     string
	AlertType      AlertType
	AlertLevel     AlertLevel
	CurrentStock   int
	ThresholdValue int
	Message        string
	IsRead         bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// InventoryThreshold represents the threshold settings for a product
type InventoryThreshold struct {
	ID            int64
	ProductID     int64
	MinStock      int // Minimum stock level (triggers low stock alert)
	MaxStock      int // Maximum stock level (triggers over stock alert)
	SafetyStock   int // Safety stock buffer
	ReorderPoint  int // Reorder point for purchasing
	Enabled       bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// NewInventoryThreshold creates a new inventory threshold
func NewInventoryThreshold(productID int64, minStock, maxStock, safetyStock, reorderPoint int) *InventoryThreshold {
	now := time.Now()
	return &InventoryThreshold{
		ProductID:    productID,
		MinStock:     minStock,
		MaxStock:     maxStock,
		SafetyStock:  safetyStock,
		ReorderPoint: reorderPoint,
		Enabled:      true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// CheckStock checks if current stock triggers any alert
func (t *InventoryThreshold) CheckStock(currentStock int) *AlertType {
	if currentStock <= 0 {
		alertType := AlertTypeOutOfStock
		return &alertType
	}
	if currentStock <= t.MinStock {
		alertType := AlertTypeLowStock
		return &alertType
	}
	if currentStock > t.MaxStock {
		alertType := AlertTypeOverStock
		return &alertType
	}
	return nil
}

// GetAlertLevel determines the alert level based on stock and threshold
func (t *InventoryThreshold) GetAlertLevel(currentStock int) AlertLevel {
	if currentStock <= 0 {
		return AlertLevelCritical
	}
	if currentStock <= t.SafetyStock {
		return AlertLevelCritical
	}
	if currentStock <= t.MinStock {
		return AlertLevelLow
	}
	if currentStock <= t.ReorderPoint {
		return AlertLevelMedium
	}
	return "" // No alert
}
