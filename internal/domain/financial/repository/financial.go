package repository

import (
	"database/sql"
	"time"

	"finance/internal/domain/financial/entity"
)

// FinancialRepository 财务报表数据接口
type FinancialRepository interface {
	GenerateReport(reportType string, startDate, endDate time.Time) (*entity.FinancialReport, error)
	GetRevenueExpenseItems(startDate, endDate time.Time) ([]entity.RevenueExpenseItem, error)
	GetProfitTrend(startDate, endDate time.Time, groupBy string) ([]entity.ProfitTrend, error)
	GetCategoryStatistics(startDate, endDate time.Time, categoryType string) ([]entity.CategoryStatistics, error)
}

// postgresFinancialRepository PostgreSQL 实现
type postgresFinancialRepository struct {
	db *sql.DB
}

func NewFinancialRepository(db *sql.DB) FinancialRepository {
	return &postgresFinancialRepository{db: db}
}

func (r *postgresFinancialRepository) GenerateReport(reportType string, startDate, endDate time.Time) (*entity.FinancialReport, error) {
	report := &entity.FinancialReport{
		ReportType: reportType,
		StartDate:  startDate,
		EndDate:    endDate,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// 查询总收入（订单金额）
	orderQuery := `
		SELECT COUNT(*), COALESCE(SUM(total_amount), 0)
		FROM orders 
		WHERE created_at >= $1 AND created_at <= $2 AND status IN ('completed', 'paid')
	`
	err := r.db.QueryRow(orderQuery, startDate, endDate).Scan(&report.OrderCount, &report.TotalRevenue)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	// 查询充值总额
	rechargeQuery := `
		SELECT COALESCE(SUM(amount), 0)
		FROM recharges 
		WHERE created_at >= $1 AND created_at <= $2 AND status = 'completed'
	`
	err = r.db.QueryRow(rechargeQuery, startDate, endDate).Scan(&report.RechargeAmount)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	// 查询发票总数
	invoiceQuery := `
		SELECT COUNT(*)
		FROM invoices 
		WHERE created_at >= $1 AND created_at <= $2
	`
	err = r.db.QueryRow(invoiceQuery, startDate, endDate).Scan(&report.InvoiceCount)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	// 计算总支出（采购成本）
	expenseQuery := `
		SELECT COALESCE(SUM(total_amount), 0)
		FROM purchase_orders 
		WHERE created_at >= $1 AND created_at <= $2 AND status IN ('completed', 'approved')
	`
	err = r.db.QueryRow(expenseQuery, startDate, endDate).Scan(&report.TotalExpense)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	// 计算净利润和利润率
	report.NetProfit = report.TotalRevenue - report.TotalExpense
	if report.TotalRevenue > 0 {
		report.ProfitMargin = (report.NetProfit / report.TotalRevenue) * 100
	}

	return report, nil
}

func (r *postgresFinancialRepository) GetRevenueExpenseItems(startDate, endDate time.Time) ([]entity.RevenueExpenseItem, error) {
	var items []entity.RevenueExpenseItem

	// 查询收入项（订单）
	orderQuery := `
		SELECT created_at, 'revenue' as type, 'order' as category, 
		       total_amount as amount, '订单收入' as description, id as related_id
		FROM orders 
		WHERE created_at >= $1 AND created_at <= $2 AND status IN ('completed', 'paid')
		ORDER BY created_at
	`
	rows, err := r.db.Query(orderQuery, startDate, endDate)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item entity.RevenueExpenseItem
		err := rows.Scan(&item.Date, &item.Type, &item.Category, &item.Amount, &item.Description, &item.RelatedID)
		if err != nil {
			continue
		}
		items = append(items, item)
	}

	// 查询支出项（采购）
	purchaseQuery := `
		SELECT created_at, 'expense' as type, 'purchase' as category, 
		       total_amount as amount, '采购支出' as description, id as related_id
		FROM purchase_orders 
		WHERE created_at >= $1 AND created_at <= $2 AND status IN ('completed', 'approved')
		ORDER BY created_at
	`
	rows, err = r.db.Query(purchaseQuery, startDate, endDate)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item entity.RevenueExpenseItem
		err := rows.Scan(&item.Date, &item.Type, &item.Category, &item.Amount, &item.Description, &item.RelatedID)
		if err != nil {
			continue
		}
		items = append(items, item)
	}

	return items, nil
}

func (r *postgresFinancialRepository) GetProfitTrend(startDate, endDate time.Time, groupBy string) ([]entity.ProfitTrend, error) {
	var trends []entity.ProfitTrend

	var dateFormat string
	switch groupBy {
	case "day":
		dateFormat = "YYYY-MM-DD"
	case "week":
		dateFormat = "IYYY-IW"
	case "month":
		dateFormat = "YYYY-MM"
	default:
		dateFormat = "YYYY-MM"
	}

	query := `
		WITH revenue_data AS (
			SELECT 
				to_char(created_at, $3) as period,
				COALESCE(SUM(total_amount), 0) as revenue
			FROM orders 
			WHERE created_at >= $1 AND created_at <= $2 AND status IN ('completed', 'paid')
			GROUP BY to_char(created_at, $3)
		),
		expense_data AS (
			SELECT 
				to_char(created_at, $3) as period,
				COALESCE(SUM(total_amount), 0) as expense
			FROM purchase_orders 
			WHERE created_at >= $1 AND created_at <= $2 AND status IN ('completed', 'approved')
			GROUP BY to_char(created_at, $3)
		)
		SELECT 
			COALESCE(r.period, e.period) as period,
			COALESCE(r.revenue, 0) as revenue,
			COALESCE(e.expense, 0) as expense,
			COALESCE(r.revenue, 0) - COALESCE(e.expense, 0) as profit
		FROM revenue_data r
		FULL OUTER JOIN expense_data e ON r.period = e.period
		ORDER BY period
	`

	rows, err := r.db.Query(query, startDate, endDate, dateFormat)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var trend entity.ProfitTrend
		err := rows.Scan(&trend.Period, &trend.Revenue, &trend.Expense, &trend.Profit)
		if err != nil {
			continue
		}
		if trend.Revenue > 0 {
			trend.ProfitRate = (trend.Profit / trend.Revenue) * 100
		}
		trends = append(trends, trend)
	}

	return trends, nil
}

func (r *postgresFinancialRepository) GetCategoryStatistics(startDate, endDate time.Time, categoryType string) ([]entity.CategoryStatistics, error) {
	var stats []entity.CategoryStatistics

	var query string
	if categoryType == "product" {
		query = `
			SELECT 
				p.category as category,
				COALESCE(SUM(oi.quantity * oi.unit_price), 0) as total_amount,
				COUNT(DISTINCT oi.order_id) as count
			FROM order_items oi
			JOIN products p ON oi.product_id = p.id
			JOIN orders o ON oi.order_id = o.id
			WHERE o.created_at >= $1 AND o.created_at <= $2 AND o.status IN ('completed', 'paid')
			GROUP BY p.category
			ORDER BY total_amount DESC
		`
	} else {
		query = `
			SELECT 
				'other' as category,
				COALESCE(SUM(total_amount), 0) as total_amount,
				COUNT(*) as count
			FROM orders 
			WHERE created_at >= $1 AND created_at <= $2 AND status IN ('completed', 'paid')
		`
	}

	rows, err := r.db.Query(query, startDate, endDate)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	defer rows.Close()

	// 计算总额用于百分比
	var totalAmount float64
	for _, stat := range stats {
		totalAmount += stat.TotalAmount
	}

	for rows.Next() {
		var stat entity.CategoryStatistics
		err := rows.Scan(&stat.Category, &stat.TotalAmount, &stat.Count)
		if err != nil {
			continue
		}
		if totalAmount > 0 {
			stat.Percentage = (stat.TotalAmount / totalAmount) * 100
		}
		stats = append(stats, stat)
	}

	return stats, nil
}
