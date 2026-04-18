package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"finance/internal/domain/financial/service"
)

// FinancialHandler 财务报表处理器
type FinancialHandler struct {
	service service.FinancialService
}

func NewFinancialHandler(service service.FinancialService) *FinancialHandler {
	return &FinancialHandler{service: service}
}

// GenerateReportRequest 生成报表请求
type GenerateReportRequest struct {
	ReportType string `json:"report_type"` // monthly, quarterly, yearly, custom
	StartDate  string `json:"start_date"`  // YYYY-MM-DD
	EndDate    string `json:"end_date"`    // YYYY-MM-DD
	GroupBy    string `json:"group_by"`    // day, week, month
}

// @Summary 生成财务报表
// @Description 根据时间段生成财务报表
// @Tags 财务报表
// @Accept json
// @Produce json
// @Param request body GenerateReportRequest true "报表参数"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/financial/report [post]
func (h *FinancialHandler) GenerateReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req GenerateReportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 解析日期
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		startDate = time.Now().AddDate(0, -1, 0) // 默认上个月
	}
	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		endDate = time.Now() // 默认今天
	}

	report, err := h.service.GenerateReport(req.ReportType, startDate, endDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"code":    0,
		"message": "success",
		"data":    report,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// @Summary 获取收支明细
// @Description 获取指定时间段的收支明细
// @Tags 财务报表
// @Accept json
// @Produce json
// @Param start_date query string false "开始日期 YYYY-MM-DD"
// @Param end_date query string false "结束日期 YYYY-MM-DD"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/financial/items [get]
func (h *FinancialHandler) GetRevenueExpenseItems(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		startDate = time.Now().AddDate(0, -1, 0)
	}
	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		endDate = time.Now()
	}

	items, err := h.service.GetRevenueExpenseItems(startDate, endDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"code":    0,
		"message": "success",
		"data":    items,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// @Summary 获取利润趋势
// @Description 获取利润趋势数据
// @Tags 财务报表
// @Accept json
// @Produce json
// @Param start_date query string false "开始日期 YYYY-MM-DD"
// @Param end_date query string false "结束日期 YYYY-MM-DD"
// @Param group_by query string false "分组方式 day/week/month"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/financial/trend [get]
func (h *FinancialHandler) GetProfitTrend(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")
	groupBy := r.URL.Query().Get("group_by")

	if groupBy == "" {
		groupBy = "month"
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		startDate = time.Now().AddDate(0, -6, 0) // 默认最近 6 个月
	}
	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		endDate = time.Now()
	}

	trends, err := h.service.GetProfitTrend(startDate, endDate, groupBy)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"code":    0,
		"message": "success",
		"data":    trends,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// @Summary 获取类别统计
// @Description 获取按类别统计的财务数据
// @Tags 财务报表
// @Accept json
// @Produce json
// @Param start_date query string false "开始日期 YYYY-MM-DD"
// @Param end_date query string false "结束日期 YYYY-MM-DD"
// @Param category_type query string false "类别类型 product/other"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/financial/category [get]
func (h *FinancialHandler) GetCategoryStatistics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")
	categoryType := r.URL.Query().Get("category_type")

	if categoryType == "" {
		categoryType = "product"
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		startDate = time.Now().AddDate(0, -1, 0)
	}
	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		endDate = time.Now()
	}

	stats, err := h.service.GetCategoryStatistics(startDate, endDate, categoryType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"code":    0,
		"message": "success",
		"data":    stats,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// @Summary 获取仪表盘汇总
// @Description 获取财务仪表盘汇总数据
// @Tags 财务报表
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/financial/dashboard [get]
func (h *FinancialHandler) GetDashboardSummary(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	summary, err := h.service.GetDashboardSummary()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"code":    0,
		"message": "success",
		"data":    summary,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
