package repository

import (
"context"

"erp-system/internal/domain/invoice/entity"
)

// InvoiceRepository defines the interface for invoice data access
type InvoiceRepository interface {
Create(ctx context.Context, invoice *entity.Invoice) error
GetByID(ctx context.Context, id int64) (*entity.Invoice, error)
GetByInvoiceNo(ctx context.Context, invoiceNo string) (*entity.Invoice, error)
Update(ctx context.Context, invoice *entity.Invoice) error
List(ctx context.Context, limit, offset int) ([]*entity.Invoice, error)
}
