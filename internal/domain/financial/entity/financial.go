package entity

import "time"

// FinancialReport 财务报表实体
type FinancialReport struct {
	ID              int64       `json:"id"`
	ReportType      string      `json:"report_type"` // monthly: 月报，quarterly: 季报，yearly: 年报，custom: 自定义
	StartDate       time.Time   `json:"start_date"`
	EndDate         time.Time   `json:"end_date"`
	TotalRevenue    float64     `json:"total_revenue"`    // 总收入
	TotalExpense    float64     `json:"total_expense"`    // 总支出
	NetProfit       float64     `json:"net_profit"`       // 净利润
	ProfitMargin    float64     `json:"profit_margin"`    // 利润率
	OrderCount      int         `json:"order_count"`      // 订单数量
	InvoiceCount    int         `json:"invoice_count"`    // 发票数量
	RechargeAmount  float64     `json:"recharge_amount"`  // 充值金额
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
}

// RevenueExpenseItem 收支明细项
type RevenueExpenseItem struct {
	Date        time.Time `json:"date"`
	Type        string    `json:"type"` // revenue: 收入，expense: 支出
	Category    string    `json:"category"`
	Amount      float64   `json:"amount"`
	Description string    `json:"description"`
	RelatedID   int64     `json:"related_id"` // 关联的订单ID、发票ID等
}

// ProfitTrend 利润趋势数据
type ProfitTrend struct {
	Period     string  `json:"period"`
	Revenue    float64 `json:"revenue"`
	Expense    float64 `json:"expense"`
	Profit     float64 `json:"profit"`
	ProfitRate float64 `json:"profit_rate"`
}

// CategoryStatistics 类别统计
type CategoryStatistics struct {
	Category    string  `json:"category"`
	TotalAmount float64 `json:"total_amount"`
	Count       int     `json:"count"`
	Percentage  float64 `json:"percentage"`
}
