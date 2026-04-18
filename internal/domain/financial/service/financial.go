package service

import (
	"time"

	"finance/internal/domain/financial/entity"
	"finance/internal/domain/financial/repository"
)

// FinancialService 财务报表业务服务
type FinancialService interface {
	GenerateReport(reportType string, startDate, endDate time.Time) (*entity.FinancialReport, error)
	GetRevenueExpenseItems(startDate, endDate time.Time) ([]entity.RevenueExpenseItem, error)
	GetProfitTrend(startDate, endDate time.Time, groupBy string) ([]entity.ProfitTrend, error)
	GetCategoryStatistics(startDate, endDate time.Time, categoryType string) ([]entity.CategoryStatistics, error)
	GetDashboardSummary() (*DashboardSummary, error)
}

// DashboardSummary 仪表盘汇总数据
type DashboardSummary struct {
	TotalRevenue   float64 `json:"total_revenue"`
	TotalExpense   float64 `json:"total_expense"`
	NetProfit      float64 `json:"net_profit"`
	ProfitMargin   float64 `json:"profit_margin"`
	MonthRevenue   float64 `json:"month_revenue"`
	MonthExpense   float64 `json:"month_expense"`
	OrderCount     int     `json:"order_count"`
	PendingOrders  int     `json:"pending_orders"`
}

type financialService struct {
	repo repository.FinancialRepository
}

func NewFinancialService(repo repository.FinancialRepository) FinancialService {
	return &financialService{repo: repo}
}

func (s *financialService) GenerateReport(reportType string, startDate, endDate time.Time) (*entity.FinancialReport, error) {
	return s.repo.GenerateReport(reportType, startDate, endDate)
}

func (s *financialService) GetRevenueExpenseItems(startDate, endDate time.Time) ([]entity.RevenueExpenseItem, error) {
	return s.repo.GetRevenueExpenseItems(startDate, endDate)
}

func (s *financialService) GetProfitTrend(startDate, endDate time.Time, groupBy string) ([]entity.ProfitTrend, error) {
	return s.repo.GetProfitTrend(startDate, endDate, groupBy)
}

func (s *financialService) GetCategoryStatistics(startDate, endDate time.Time, categoryType string) ([]entity.CategoryStatistics, error) {
	return s.repo.GetCategoryStatistics(startDate, endDate, categoryType)
}

func (s *financialService) GetDashboardSummary() (*DashboardSummary, error) {
	now := time.Now()
	
	// 本月起止时间
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	monthEnd := monthStart.AddDate(0, 1, 0).Add(-time.Second)
	
	// 获取本月报表
	monthReport, err := s.repo.GenerateReport("monthly", monthStart, monthEnd)
	if err != nil {
		return nil, err
	}
	
	// 获取历史累计数据（从系统开始到现在）
	allTimeStart := time.Date(2020, 1, 1, 0, 0, 0, 0, now.Location())
	allTimeReport, err := s.repo.GenerateReport("all", allTimeStart, now)
	if err != nil {
		return nil, err
	}
	
	summary := &DashboardSummary{
		TotalRevenue:  allTimeReport.TotalRevenue,
		TotalExpense:  allTimeReport.TotalExpense,
		NetProfit:     allTimeReport.NetProfit,
		ProfitMargin:  allTimeReport.ProfitMargin,
		MonthRevenue:  monthReport.TotalRevenue,
		MonthExpense:  monthReport.TotalExpense,
		OrderCount:    monthReport.OrderCount,
		PendingOrders: 0, // TODO: 待实现待处理订单数查询
	}
	
	return summary, nil
}
