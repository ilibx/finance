package main

import (
"database/sql"
"fmt"
"log"
"net/http"

"erp-system/internal/config"
"erp-system/internal/handler"
userRepo "erp-system/internal/domain/user/repository"
productRepo "erp-system/internal/domain/product/repository"
orderRepo "erp-system/internal/domain/order/repository"
invoiceRepo "erp-system/internal/domain/invoice/repository"
rechargeRepo "erp-system/internal/domain/recharge/repository"
supplierRepo "erp-system/internal/domain/supplier/repository"
userService "erp-system/internal/domain/user/service"
productService "erp-system/internal/domain/product/service"
orderService "erp-system/internal/domain/order/service"
invoiceService "erp-system/internal/domain/invoice/service"
rechargeService "erp-system/internal/domain/recharge/service"

_ "github.com/lib/pq"
)

func main() {
cfg := config.Load()

dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
cfg.Database.Host, cfg.Database.Port, cfg.Database.User,
cfg.Database.Password, cfg.Database.DBName)

db, err := sql.Open("postgres", dsn)
if err != nil {
log.Printf("警告：数据库连接失败：%v，系统将以有限模式运行", err)
} else {
defer db.Close()
if err := db.Ping(); err != nil {
log.Printf("警告：数据库 ping 失败：%v", err)
} else {
log.Println("数据库连接成功")
}
}

uRepo := userRepo.NewUserRepository(db)
pRepo := productRepo.NewProductRepository(db)
oRepo := orderRepo.NewOrderRepository(db)
iRepo := invoiceRepo.NewInvoiceRepository(db)
rRepo := rechargeRepo.NewRechargeRepository(db)
sRepo := supplierRepo.NewSupplierRepository(db)

userSvc := userService.NewUserService(uRepo)
productSvc := productService.NewProductService(pRepo)
orderSvc := orderService.NewOrderService(oRepo, pRepo)
invoiceSvc := invoiceService.NewInvoiceService(iRepo)
rechargeSvc := rechargeService.NewRechargeService(rRepo, uRepo, sRepo)

userHandler := handler.NewUserHandler(userSvc)
productHandler := handler.NewProductHandler(productSvc)
orderHandler := handler.NewOrderHandler(orderSvc)
invoiceHandler := handler.NewInvoiceHandler(invoiceSvc, orderSvc)
rechargeHandler := handler.NewRechargeHandler(rechargeSvc)
excelImportHandler := handler.NewExcelImportHandler()

http.HandleFunc("/api/users/create", userHandler.CreateUser)
http.HandleFunc("/api/users/get", userHandler.GetUser)
http.HandleFunc("/api/products/create", productHandler.CreateProduct)
http.HandleFunc("/api/orders/create", orderHandler.CreateOrder)
http.HandleFunc("/api/invoices/generate", invoiceHandler.GenerateInvoice)
http.HandleFunc("/api/recharge/process", rechargeHandler.ProcessRecharge)

http.HandleFunc("/api/excel/import/consumption-bills", excelImportHandler.ImportConsumptionBills)
http.HandleFunc("/api/excel/import/recharge-records", excelImportHandler.ImportRechargeRecords)
http.HandleFunc("/api/excel/import/supplier-recharges", excelImportHandler.ImportSupplierRecharges)
http.HandleFunc("/api/excel/import/supplier-invoices", excelImportHandler.ImportSupplierInvoices)
http.HandleFunc("/api/excel/export/consumption-bills", excelImportHandler.ExportConsumptionBills)

addr := ":" + cfg.Server.Port
log.Printf("ERP 系统启动在 %s", addr)
log.Println("API 端点:")
log.Println("  - POST /api/users/create - 创建用户")
log.Println("  - GET  /api/users/get?id=1 - 获取用户")
log.Println("  - POST /api/products/create - 创建产品")
log.Println("  - POST /api/orders/create - 创建订单")
log.Println("  - POST /api/invoices/generate - 生成发票")
log.Println("  - POST /api/recharge/process - 处理充值")
log.Println("  - POST /api/excel/import/* - 导入 Excel")
log.Println("  - GET  /api/excel/export/* - 导出 Excel")

if err := http.ListenAndServe(addr, nil); err != nil {
log.Fatal(err)
}
}
