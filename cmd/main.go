package main

import (
"database/sql"
"fmt"
"log"
"net/http"

"finance/internal/config"
"finance/internal/handler"
userRepo "finance/internal/domain/user/repository"
productRepo "finance/internal/domain/product/repository"
orderRepo "finance/internal/domain/order/repository"
invoiceRepo "finance/internal/domain/invoice/repository"
rechargeRepo "finance/internal/domain/recharge/repository"
supplierRepo "finance/internal/domain/supplier/repository"
projectRepo "finance/internal/domain/project/repository"
userService "finance/internal/domain/user/service"
productService "finance/internal/domain/product/service"
orderService "finance/internal/domain/order/service"
invoiceService "finance/internal/domain/invoice/service"
rechargeService "finance/internal/domain/recharge/service"
projectService "finance/internal/domain/project/service"
supplierService "finance/internal/domain/supplier/service"

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
projRepo := projectRepo.NewProjectRepository(db)

userSvc := userService.NewUserService(uRepo)
productSvc := productService.NewProductService(pRepo)
orderSvc := orderService.NewOrderService(oRepo, pRepo)
invoiceSvc := invoiceService.NewInvoiceService(iRepo)
rechargeSvc := rechargeService.NewRechargeService(rRepo, uRepo, sRepo)
projSvc := projectService.NewProjectService(projRepo)
supplierSvc := supplierService.NewSupplierService(sRepo)

userHandler := handler.NewUserHandler(userSvc)
productHandler := handler.NewProductHandler(productSvc)
orderHandler := handler.NewOrderHandler(orderSvc)
invoiceHandler := handler.NewInvoiceHandler(invoiceSvc, orderSvc)
rechargeHandler := handler.NewRechargeHandler(rechargeSvc)
projectHandler := handler.NewProjectHandler(projSvc)
supplierHandler := handler.NewSupplierHandler(supplierSvc)
excelImportHandler := handler.NewExcelImportHandler()

http.HandleFunc("/api/users/create", userHandler.CreateUser)
http.HandleFunc("/api/users/get", userHandler.GetUser)
http.HandleFunc("/api/products/create", productHandler.CreateProduct)
http.HandleFunc("/api/orders/create", orderHandler.CreateOrder)
http.HandleFunc("/api/invoices/generate", invoiceHandler.GenerateInvoice)
http.HandleFunc("/api/recharge/process", rechargeHandler.ProcessRecharge)

// Project management endpoints
http.HandleFunc("/api/projects/create", projectHandler.CreateProject)
http.HandleFunc("/api/projects/get", projectHandler.GetProject)
http.HandleFunc("/api/projects/update", projectHandler.UpdateProject)
http.HandleFunc("/api/projects/delete", projectHandler.DeleteProject)
http.HandleFunc("/api/projects/list", projectHandler.ListProjects)
http.HandleFunc("/api/projects/status", projectHandler.UpdateProjectStatus)
http.HandleFunc("/api/projects/track-progress", projectHandler.TrackProjectProgress)

// Supplier management endpoints
http.HandleFunc("/api/suppliers/create", supplierHandler.CreateSupplier)
http.HandleFunc("/api/suppliers/get", supplierHandler.GetSupplier)
http.HandleFunc("/api/suppliers/list", supplierHandler.ListSuppliers)
http.HandleFunc("/api/suppliers/update", supplierHandler.UpdateSupplier)
http.HandleFunc("/api/suppliers/delete", supplierHandler.DeleteSupplier)

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
log.Println("项目管理:")
log.Println("  - POST /api/projects/create - 创建项目")
log.Println("  - GET  /api/projects/get?id=1 - 获取项目")
log.Println("  - PUT  /api/projects/update - 更新项目")
log.Println("  - DELETE /api/projects/delete?id=1 - 删除项目")
log.Println("  - GET  /api/projects/list - 项目列表")
log.Println("  - PATCH /api/projects/status?id=1 - 更新项目状态")
log.Println("  - POST /api/projects/track-progress?id=1 - 跟踪项目进度")
log.Println("供应商管理:")
log.Println("  - POST /api/suppliers/create - 创建供应商")
log.Println("  - GET  /api/suppliers/get?id=1 - 获取供应商")
log.Println("  - GET  /api/suppliers/list - 供应商列表")
log.Println("  - PUT  /api/suppliers/update - 更新供应商")
log.Println("  - DELETE /api/suppliers/delete?id=1 - 删除供应商")
log.Println("Excel 导入导出:")
log.Println("  - POST /api/excel/import/* - 导入 Excel")
log.Println("  - GET  /api/excel/export/* - 导出 Excel")

if err := http.ListenAndServe(addr, nil); err != nil {
log.Fatal(err)
}
}
