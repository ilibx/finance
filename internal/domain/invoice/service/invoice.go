package service

import (
"context"
"fmt"
"time"

"erp-system/internal/common/errors"
"erp-system/internal/domain/invoice/entity"
"erp-system/internal/domain/invoice/repository"
)

// InvoiceService handles invoice business logic
type InvoiceService struct {
invoiceRepo repository.InvoiceRepository
}

// NewInvoiceService creates a new invoice service
func NewInvoiceService(invoiceRepo repository.InvoiceRepository) *InvoiceService {
return &InvoiceService{
invoiceRepo: invoiceRepo,
}
}

// CreateInvoice creates a new invoice
func (s *InvoiceService) CreateInvoice(ctx context.Context, orderID, userID int64, amount, taxRate float64) (*entity.Invoice, error) {
invoiceNo := generateInvoiceNo()
taxAmount := amount * taxRate
invoice := entity.NewInvoice(invoiceNo, orderID, userID, amount, taxAmount, entity.InvoiceTypeSales)

if err := s.invoiceRepo.Create(ctx, invoice); err != nil {
return nil, err
}

return invoice, nil
}

// GetInvoiceByID gets an invoice by ID
func (s *InvoiceService) GetInvoiceByID(ctx context.Context, id int64) (*entity.Invoice, error) {
invoice, err := s.invoiceRepo.GetByID(ctx, id)
if err != nil {
return nil, errors.ErrInvoiceNotFound
}
return invoice, nil
}

// IssueInvoice issues an invoice
func (s *InvoiceService) IssueInvoice(ctx context.Context, invoiceID int64) error {
invoice, err := s.invoiceRepo.GetByID(ctx, invoiceID)
if err != nil {
return errors.ErrInvoiceNotFound
}

if err := invoice.Issue(); err != nil {
return err
}

return s.invoiceRepo.Update(ctx, invoice)
}

// ListInvoices lists invoices with pagination
func (s *InvoiceService) ListInvoices(ctx context.Context, limit, offset int) ([]*entity.Invoice, error) {
return s.invoiceRepo.List(ctx, limit, offset)
}

func generateInvoiceNo() string {
return fmt.Sprintf("INV%s", time.Now().Format("20060102150405"))
}
