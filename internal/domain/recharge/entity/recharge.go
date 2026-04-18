package entity

import (
"time"

"erp-system/internal/common/valueobject"
)

// RechargeType represents recharge type
type RechargeType string

const (
RechargeTypeUser     RechargeType = "user"
RechargeTypeSupplier RechargeType = "supplier"
)

// RechargeMethod represents recharge method
type RechargeMethod string

const (
RechargeMethodBankTransfer RechargeMethod = "bank_transfer"
RechargeMethodOnline       RechargeMethod = "online"
RechargeMethodCash         RechargeMethod = "cash"
)

// RechargeStatus represents recharge status
type RechargeStatus string

const (
RechargeStatusPending   RechargeStatus = "pending"
RechargeStatusCompleted RechargeStatus = "completed"
RechargeStatusFailed    RechargeStatus = "failed"
)

// RechargeRecord represents a recharge record aggregate root
type RechargeRecord struct {
ID        int64
UserID    int64
Amount    valueobject.Money
Type      RechargeType
Method    RechargeMethod
Status    RechargeStatus
Remark    string
CreatedAt time.Time
UpdatedAt time.Time
}

// NewRechargeRecord creates a new recharge record
func NewRechargeRecord(userID int64, amount float64, rechargeType RechargeType, method RechargeMethod, remark string) *RechargeRecord {
now := time.Now()
return &RechargeRecord{
UserID: userID,
Amount: valueobject.NewMoney(amount, "CNY"),
Type:   rechargeType,
Method: method,
Status: RechargeStatusPending,
Remark: remark,
CreatedAt: now,
UpdatedAt: now,
}
}

// Complete marks the recharge as completed
func (r *RechargeRecord) Complete() {
r.Status = RechargeStatusCompleted
r.UpdatedAt = time.Now()
}

// Fail marks the recharge as failed
func (r *RechargeRecord) Fail() {
r.Status = RechargeStatusFailed
r.UpdatedAt = time.Now()
}
