package excel

import (
"fmt"
"io"
"time"

"github.com/xuri/excelize/v2"
)

// ConsumptionBillImport represents imported consumption bill data
type ConsumptionBillImport struct {
UserID      int64   `json:"user_id"`
OrderNo     string  `json:"order_no"`
Amount      float64 `json:"amount"`
Description string  `json:"description"`
Date        string  `json:"date"`
}

// RechargeRecordImport represents imported recharge record data
type RechargeRecordImport struct {
UserID int64   `json:"user_id"`
Amount float64 `json:"amount"`
Method string  `json:"method"`
Remark string  `json:"remark"`
Date   string  `json:"date"`
}

// SupplierRechargeImport represents imported supplier recharge data
type SupplierRechargeImport struct {
SupplierID int64   `json:"supplier_id"`
Amount     float64 `json:"amount"`
Method     string  `json:"method"`
Remark     string  `json:"remark"`
Date       string  `json:"date"`
}

// SupplierInvoiceImport represents imported supplier invoice data
type SupplierInvoiceImport struct {
SupplierID int64   `json:"supplier_id"`
InvoiceNo  string  `json:"invoice_no"`
Amount     float64 `json:"amount"`
TaxAmount  float64 `json:"tax_amount"`
Date       string  `json:"date"`
}

// ParseConsumptionBills parses consumption bills from Excel file
func ParseConsumptionBills(reader io.Reader) ([]*ConsumptionBillImport, error) {
f, err := excelize.OpenReader(reader)
if err != nil {
return nil, fmt.Errorf("failed to open Excel file: %w", err)
}
defer f.Close()

sheetName := f.GetSheetName(0)
rows, err := f.GetRows(sheetName)
if err != nil {
return nil, fmt.Errorf("failed to read Excel rows: %w", err)
}

var bills []*ConsumptionBillImport
for i, row := range rows {
if i == 0 {
continue
}
if len(row) < 5 {
continue
}

bill := &ConsumptionBillImport{
UserID:      parseInt64(row[0]),
OrderNo:     row[1],
Amount:      parseFloat64(row[2]),
Description: row[3],
Date:        row[4],
}
bills = append(bills, bill)
}

return bills, nil
}

// ParseRechargeRecords parses recharge records from Excel file
func ParseRechargeRecords(reader io.Reader) ([]*RechargeRecordImport, error) {
f, err := excelize.OpenReader(reader)
if err != nil {
return nil, fmt.Errorf("failed to open Excel file: %w", err)
}
defer f.Close()

sheetName := f.GetSheetName(0)
rows, err := f.GetRows(sheetName)
if err != nil {
return nil, fmt.Errorf("failed to read Excel rows: %w", err)
}

var records []*RechargeRecordImport
for i, row := range rows {
if i == 0 {
continue
}
if len(row) < 4 {
continue
}

record := &RechargeRecordImport{
UserID: parseInt64(row[0]),
Amount: parseFloat64(row[1]),
Method: row[2],
Remark: row[3],
}
if len(row) > 4 {
record.Date = row[4]
}
records = append(records, record)
}

return records, nil
}

// ParseSupplierRecharges parses supplier recharges from Excel file
func ParseSupplierRecharges(reader io.Reader) ([]*SupplierRechargeImport, error) {
f, err := excelize.OpenReader(reader)
if err != nil {
return nil, fmt.Errorf("failed to open Excel file: %w", err)
}
defer f.Close()

sheetName := f.GetSheetName(0)
rows, err := f.GetRows(sheetName)
if err != nil {
return nil, fmt.Errorf("failed to read Excel rows: %w", err)
}

var recharges []*SupplierRechargeImport
for i, row := range rows {
if i == 0 {
continue
}
if len(row) < 4 {
continue
}

recharge := &SupplierRechargeImport{
SupplierID: parseInt64(row[0]),
Amount:     parseFloat64(row[1]),
Method:     row[2],
Remark:     row[3],
}
if len(row) > 4 {
recharge.Date = row[4]
}
recharges = append(recharges, recharge)
}

return recharges, nil
}

// ParseSupplierInvoices parses supplier invoices from Excel file
func ParseSupplierInvoices(reader io.Reader) ([]*SupplierInvoiceImport, error) {
f, err := excelize.OpenReader(reader)
if err != nil {
return nil, fmt.Errorf("failed to open Excel file: %w", err)
}
defer f.Close()

sheetName := f.GetSheetName(0)
rows, err := f.GetRows(sheetName)
if err != nil {
return nil, fmt.Errorf("failed to read Excel rows: %w", err)
}

var invoices []*SupplierInvoiceImport
for i, row := range rows {
if i == 0 {
continue
}
if len(row) < 4 {
continue
}

invoice := &SupplierInvoiceImport{
SupplierID: parseInt64(row[0]),
InvoiceNo:  row[1],
Amount:     parseFloat64(row[2]),
TaxAmount:  parseFloat64(row[3]),
}
if len(row) > 4 {
invoice.Date = row[4]
}
invoices = append(invoices, invoice)
}

return invoices, nil
}

// GenerateConsumptionBillExport generates Excel export for consumption bills
func GenerateConsumptionBillExport(bills []map[string]interface{}) ([]byte, error) {
f := excelize.NewFile()
defer f.Close()

sheetName := "消费账单"
f.SetSheetName("Sheet1", sheetName)

headers := []string{"ID", "用户 ID", "订单 ID", "金额", "描述", "状态", "支付时间", "创建时间"}
for i, header := range headers {
cell, _ := excelize.CoordinatesToCellName(i+1, 1)
f.SetCellValue(sheetName, cell, header)
}

for i, bill := range bills {
row := i + 2
f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), bill["id"])
f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), bill["user_id"])
f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), bill["order_id"])
f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), bill["amount"])
f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), bill["description"])
f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), bill["status"])
if paidAt, ok := bill["paid_at"].(time.Time); ok {
f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), paidAt.Format(time.RFC3339))
}
f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), bill["created_at"])
}

buffer, err := f.WriteToBuffer()
if err != nil {
return nil, fmt.Errorf("failed to write Excel: %w", err)
}

return buffer.Bytes(), nil
}

// Helper functions
func parseInt64(s string) int64 {
var result int64
fmt.Sscanf(s, "%d", &result)
return result
}

func parseFloat64(s string) float64 {
var result float64
fmt.Sscanf(s, "%f", &result)
return result
}
