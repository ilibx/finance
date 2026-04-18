package handler

import (
	"net/http"

	"finance/infrastructure/excel"
)

// ExcelImportHandler handles Excel import/export requests
type ExcelImportHandler struct{}

// NewExcelImportHandler creates a new Excel import handler
func NewExcelImportHandler() *ExcelImportHandler {
return &ExcelImportHandler{}
}

// ImportConsumptionBills handles consumption bills import
func (h *ExcelImportHandler) ImportConsumptionBills(w http.ResponseWriter, r *http.Request) {
if r.Method != http.MethodPost {
http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
return
}

file, _, err := r.FormFile("file")
if err != nil {
writeJSON(w, http.StatusBadRequest, APIResponse{
Success: false,
Error:   "Failed to read file: " + err.Error(),
})
return
}
defer file.Close()

bills, err := excel.ParseConsumptionBills(file)
if err != nil {
writeJSON(w, http.StatusBadRequest, APIResponse{
Success: false,
Error:   "Failed to parse Excel file: " + err.Error(),
})
return
}

writeJSON(w, http.StatusOK, APIResponse{
Success: true,
Data:    map[string]int{"imported_count": len(bills)},
})
}

// ImportRechargeRecords handles recharge records import
func (h *ExcelImportHandler) ImportRechargeRecords(w http.ResponseWriter, r *http.Request) {
if r.Method != http.MethodPost {
http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
return
}

file, _, err := r.FormFile("file")
if err != nil {
writeJSON(w, http.StatusBadRequest, APIResponse{
Success: false,
Error:   "Failed to read file: " + err.Error(),
})
return
}
defer file.Close()

records, err := excel.ParseRechargeRecords(file)
if err != nil {
writeJSON(w, http.StatusBadRequest, APIResponse{
Success: false,
Error:   "Failed to parse Excel file: " + err.Error(),
})
return
}

writeJSON(w, http.StatusOK, APIResponse{
Success: true,
Data:    map[string]int{"imported_count": len(records)},
})
}

// ImportSupplierRecharges handles supplier recharges import
func (h *ExcelImportHandler) ImportSupplierRecharges(w http.ResponseWriter, r *http.Request) {
if r.Method != http.MethodPost {
http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
return
}

file, _, err := r.FormFile("file")
if err != nil {
writeJSON(w, http.StatusBadRequest, APIResponse{
Success: false,
Error:   "Failed to read file: " + err.Error(),
})
return
}
defer file.Close()

recharges, err := excel.ParseSupplierRecharges(file)
if err != nil {
writeJSON(w, http.StatusBadRequest, APIResponse{
Success: false,
Error:   "Failed to parse Excel file: " + err.Error(),
})
return
}

writeJSON(w, http.StatusOK, APIResponse{
Success: true,
Data:    map[string]int{"imported_count": len(recharges)},
})
}

// ImportSupplierInvoices handles supplier invoices import
func (h *ExcelImportHandler) ImportSupplierInvoices(w http.ResponseWriter, r *http.Request) {
if r.Method != http.MethodPost {
http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
return
}

file, _, err := r.FormFile("file")
if err != nil {
writeJSON(w, http.StatusBadRequest, APIResponse{
Success: false,
Error:   "Failed to read file: " + err.Error(),
})
return
}
defer file.Close()

invoices, err := excel.ParseSupplierInvoices(file)
if err != nil {
writeJSON(w, http.StatusBadRequest, APIResponse{
Success: false,
Error:   "Failed to parse Excel file: " + err.Error(),
})
return
}

writeJSON(w, http.StatusOK, APIResponse{
Success: true,
Data:    map[string]int{"imported_count": len(invoices)},
})
}

// ExportConsumptionBills handles consumption bills export
func (h *ExcelImportHandler) ExportConsumptionBills(w http.ResponseWriter, r *http.Request) {
if r.Method != http.MethodGet {
http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
return
}

bills := []map[string]interface{}{}

data, err := excel.GenerateConsumptionBillExport(bills)
if err != nil {
writeJSON(w, http.StatusInternalServerError, APIResponse{
Success: false,
Error:   "Failed to generate export: " + err.Error(),
})
return
}

w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
w.Header().Set("Content-Disposition", "attachment; filename=consumption_bills.xlsx")
w.Write(data)
}
