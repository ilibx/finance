package service

import (
	"context"
	"fmt"

	"finance/internal/domain/inventory/entity"
	inventoryRepo "finance/internal/domain/inventory/repository"
	productRepo "finance/internal/domain/product/repository"
)

// InventoryAlertService handles inventory alert business logic
type InventoryAlertService struct {
	alertRepo       inventoryRepo.InventoryAlertRepository
	thresholdRepo   inventoryRepo.InventoryThresholdRepository
	productRepo     productRepo.ProductRepository
}

// NewInventoryAlertService creates a new inventory alert service
func NewInventoryAlertService(
	alertRepo inventoryRepo.InventoryAlertRepository,
	thresholdRepo inventoryRepo.InventoryThresholdRepository,
	productRepo productRepo.ProductRepository,
) *InventoryAlertService {
	return &InventoryAlertService{
		alertRepo:     alertRepo,
		thresholdRepo: thresholdRepo,
		productRepo:   productRepo,
	}
}

// CheckStockAndCreateAlert checks stock level and creates alert if needed
func (s *InventoryAlertService) CheckStockAndCreateAlert(ctx context.Context, productID int64) error {
	// Get product
	product, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return fmt.Errorf("failed to get product: %w", err)
	}

	// Get threshold for product
	threshold, err := s.thresholdRepo.GetByProductID(ctx, productID)
	if err != nil {
		// If no threshold exists, use default values
		threshold = entity.NewInventoryThreshold(productID, 10, 1000, 5, 20)
		if createErr := s.thresholdRepo.Create(ctx, threshold); createErr != nil {
			return fmt.Errorf("failed to create default threshold: %w", createErr)
		}
	}

	if !threshold.Enabled {
		return nil
	}

	// Check if alert should be triggered
	alertType := threshold.CheckStock(product.Stock)
	if alertType == nil {
		return nil
	}

	// Create alert
	alertLevel := threshold.GetAlertLevel(product.Stock)
	message := s.generateAlertMessage(*alertType, product.Name, product.SKU, product.Stock, threshold)

	alert := &entity.InventoryAlert{
		ProductID:      product.ID,
		ProductName:    product.Name,
		ProductSKU:     product.SKU,
		AlertType:      *alertType,
		AlertLevel:     alertLevel,
		CurrentStock:   product.Stock,
		ThresholdValue: threshold.MinStock,
		Message:        message,
		IsRead:         false,
		CreatedAt:      product.UpdatedAt,
		UpdatedAt:      product.UpdatedAt,
	}

	return s.alertRepo.Create(ctx, alert)
}

// generateAlertMessage generates a human-readable alert message
func (s *InventoryAlertService) generateAlertMessage(
	alertType entity.AlertType,
	productName, productSKU string,
	currentStock int,
	threshold *entity.InventoryThreshold,
) string {
	switch alertType {
	case entity.AlertTypeOutOfStock:
		return fmt.Sprintf("【缺货警告】商品 %s (%s) 当前库存为 0，需要立即补货！", productName, productSKU)
	case entity.AlertTypeLowStock:
		return fmt.Sprintf("【低库存警告】商品 %s (%s) 当前库存 %d，低于最低库存阈值 %d，请及时补货！", productName, productSKU, currentStock, threshold.MinStock)
	case entity.AlertTypeOverStock:
		return fmt.Sprintf("【高库存警告】商品 %s (%s) 当前库存 %d，超过最高库存阈值 %d，请注意库存积压！", productName, productSKU, currentStock, threshold.MaxStock)
	default:
		return fmt.Sprintf("【库存警告】商品 %s (%s) 库存异常", productName, productSKU)
	}
}

// ListAlerts lists inventory alerts with filters
func (s *InventoryAlertService) ListAlerts(ctx context.Context, productID *int64, isRead *bool, limit, offset int) ([]*entity.InventoryAlert, error) {
	return s.alertRepo.List(ctx, productID, isRead, limit, offset)
}

// MarkAlertAsRead marks an alert as read
func (s *InventoryAlertService) MarkAlertAsRead(ctx context.Context, alertID int64) error {
	return s.alertRepo.MarkAsRead(ctx, alertID)
}

// MarkAllAlertsAsRead marks all alerts as read
func (s *InventoryAlertService) MarkAllAlertsAsRead(ctx context.Context, productID *int64) error {
	alerts, err := s.alertRepo.List(ctx, productID, func() *bool { b := false; return &b }(), 1000, 0)
	if err != nil {
		return err
	}
	for _, alert := range alerts {
		if err := s.alertRepo.MarkAsRead(ctx, alert.ID); err != nil {
			return err
		}
	}
	return nil
}

// GetUnreadAlertCount gets the count of unread alerts
func (s *InventoryAlertService) GetUnreadAlertCount(ctx context.Context, productID *int64) (int, error) {
	return s.alertRepo.GetUnreadCount(ctx, productID)
}

// SetThreshold sets or updates inventory threshold for a product
func (s *InventoryAlertService) SetThreshold(ctx context.Context, productID int64, minStock, maxStock, safetyStock, reorderPoint int) (*entity.InventoryThreshold, error) {
	threshold, err := s.thresholdRepo.GetByProductID(ctx, productID)
	if err != nil {
		// Create new threshold
		threshold = entity.NewInventoryThreshold(productID, minStock, maxStock, safetyStock, reorderPoint)
		if createErr := s.thresholdRepo.Create(ctx, threshold); createErr != nil {
			return nil, fmt.Errorf("failed to create threshold: %w", createErr)
		}
	} else {
		// Update existing threshold
		threshold.MinStock = minStock
		threshold.MaxStock = maxStock
		threshold.SafetyStock = safetyStock
		threshold.ReorderPoint = reorderPoint
		if updateErr := s.thresholdRepo.Update(ctx, threshold); updateErr != nil {
			return nil, fmt.Errorf("failed to update threshold: %w", updateErr)
		}
	}

	// After setting threshold, check if current stock triggers alert
	if checkErr := s.CheckStockAndCreateAlert(ctx, productID); checkErr != nil {
		return nil, fmt.Errorf("failed to check stock after setting threshold: %w", checkErr)
	}

	return threshold, nil
}

// GetThreshold gets the threshold for a product
func (s *InventoryAlertService) GetThreshold(ctx context.Context, productID int64) (*entity.InventoryThreshold, error) {
	return s.thresholdRepo.GetByProductID(ctx, productID)
}

// ListThresholds lists all inventory thresholds
func (s *InventoryAlertService) ListThresholds(ctx context.Context, limit, offset int) ([]*entity.InventoryThreshold, error) {
	return s.thresholdRepo.List(ctx, limit, offset)
}

// EnableThreshold enables inventory alert for a product
func (s *InventoryAlertService) EnableThreshold(ctx context.Context, thresholdID int64) error {
	threshold, err := s.thresholdRepo.GetByID(ctx, thresholdID)
	if err != nil {
		return err
	}
	threshold.Enabled = true
	return s.thresholdRepo.Update(ctx, threshold)
}

// DisableThreshold disables inventory alert for a product
func (s *InventoryAlertService) DisableThreshold(ctx context.Context, thresholdID int64) error {
	threshold, err := s.thresholdRepo.GetByID(ctx, thresholdID)
	if err != nil {
		return err
	}
	threshold.Enabled = false
	return s.thresholdRepo.Update(ctx, threshold)
}

// CheckAllProducts checks all products and generates alerts for those below threshold
func (s *InventoryAlertService) CheckAllProducts(ctx context.Context) error {
	// Get all products (pagination with large limit)
	products, err := s.productRepo.List(ctx, 10000, 0)
	if err != nil {
		return fmt.Errorf("failed to list products: %w", err)
	}

	for _, product := range products {
		if err := s.CheckStockAndCreateAlert(ctx, product.ID); err != nil {
			// Log error but continue checking other products
			fmt.Printf("Failed to check stock for product %d: %v\n", product.ID, err)
		}
	}

	return nil
}
