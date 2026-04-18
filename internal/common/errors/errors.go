package errors

import "errors"

var (
	// User errors
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidUserInput   = errors.New("invalid user input")
	
	// Product errors
	ErrProductNotFound    = errors.New("product not found")
	ErrProductAlreadyExists = errors.New("product already exists")
	ErrInsufficientStock  = errors.New("insufficient stock")
	
	// Order errors
	ErrOrderNotFound      = errors.New("order not found")
	ErrInvalidOrderStatus = errors.New("invalid order status")
	ErrOrderAlreadyPaid   = errors.New("order already paid")
	
	// Invoice errors
	ErrInvoiceNotFound    = errors.New("invoice not found")
	ErrInvoiceAlreadyIssued = errors.New("invoice already issued")
	
	// Recharge errors
	ErrRechargeNotFound   = errors.New("recharge record not found")
	ErrInvalidRechargeAmount = errors.New("invalid recharge amount")
	
	// Supplier errors
	ErrSupplierNotFound   = errors.New("supplier not found")
	ErrSupplierAlreadyExists = errors.New("supplier already exists")
	
	// Project errors
	ErrProjectNotFound    = errors.New("project not found")
	ErrProjectAlreadyExists = errors.New("project already exists")
	
	// Common errors
	ErrDatabaseOperation  = errors.New("database operation failed")
	ErrInvalidID          = errors.New("invalid ID")
	ErrPermissionDenied   = errors.New("permission denied")
)

// DomainError represents a domain-specific error
type DomainError struct {
	Code    string
	Message string
	Err     error
}

func (e *DomainError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

func (e *DomainError) Unwrap() error {
	return e.Err
}

// NewDomainError creates a new domain error
func NewDomainError(code, message string, err error) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}
