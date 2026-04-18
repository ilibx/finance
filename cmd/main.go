package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"erp-system/internal/config"
	"erp-system/internal/handler"
	"erp-system/internal/repository"
	"erp-system/internal/service"

	_ "github.com/lib/pq"
)

func main() {
	// 加载配置
	cfg := config.Load()

	// 连接数据库
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User,
		cfg.Database.Password, cfg.Database.DBName)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Printf("警告：数据库连接失败：%v，系统将以有限模式运行", err)
		// 继续运行，但部分功能不可用
	} else {
		defer db.Close()
		if err := db.Ping(); err != nil {
			log.Printf("警告：数据库 ping 失败：%v", err)
		} else {
			log.Println("数据库连接成功")
		}
	}

	// 初始化仓库层
	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	invoiceRepo := repository.NewInvoiceRepository(db)
	rechargeRepo := repository.NewRechargeRecordRepository(db)
	supplierRepo := repository.NewSupplierRepository(db)
	billRepo := repository.NewConsumptionBillRepository(db)

	// 初始化服务层
	userService := service.NewUserService(userRepo)
	productService := service.NewProductService(productRepo)
	orderService := service.NewOrderService(orderRepo, productRepo, userRepo, billRepo, invoiceRepo)
	invoiceService := service.NewInvoiceService(invoiceRepo)
	rechargeService := service.NewRechargeService(rechargeRepo, userRepo, supplierRepo)

	// 初始化处理器
	userHandler := handler.NewUserHandler(userService)
	productHandler := handler.NewProductHandler(productService)
	orderHandler := handler.NewOrderHandler(orderService)
	invoiceHandler := handler.NewInvoiceHandler(invoiceService, orderService)
	rechargeHandler := handler.NewRechargeHandler(rechargeService)
	excelImportHandler := handler.NewExcelImportHandler(rechargeService, orderService)

	// 注册路由
	http.HandleFunc("/api/users/create", userHandler.CreateUser)
	http.HandleFunc("/api/users/get", userHandler.GetUser)
	http.HandleFunc("/api/products/create", productHandler.CreateProduct)
	http.HandleFunc("/api/orders/create", orderHandler.CreateOrder)
	http.HandleFunc("/api/invoices/generate", invoiceHandler.GenerateInvoice)
	http.HandleFunc("/api/recharge/process", rechargeHandler.ProcessRecharge)
	
	// Excel 导入导出
	http.HandleFunc("/api/excel/import/consumption-bills", excelImportHandler.ImportConsumptionBills)
	http.HandleFunc("/api/excel/import/recharge-records", excelImportHandler.ImportRechargeRecords)
	http.HandleFunc("/api/excel/import/supplier-recharges", excelImportHandler.ImportSupplierRecharges)
	http.HandleFunc("/api/excel/import/supplier-invoices", excelImportHandler.ImportSupplierInvoices)
	http.HandleFunc("/api/excel/export/consumption-bills", excelImportHandler.ExportConsumptionBills)

	// 启动服务器
	addr := ":" + cfg.Server.Port
	log.Printf("ERP 系统启动在 %s", addr)
	log.Println("API 端点:")
	log.Println("  - POST /api/users/create - 创建用户")
	log.Println("  - GET  /api/users/get?id=1 - 获取用户")
	log.Println("  - POST /api/products/create - 创建产品")
	log.Println("  - POST /api/orders/create - 创建订单")
	log.Println("  - POST /api/invoices/generate - 生成发票")
	log.Println("  - POST /api/recharge/process - 处理充值")
	log.Println("  - POST /api/excel/import/consumption-bills - 导入消费账单")
	log.Println("  - POST /api/excel/import/recharge-records - 导入充值记录")
	log.Println("  - POST /api/excel/import/supplier-recharges - 导入供应商充值")
	log.Println("  - POST /api/excel/import/supplier-invoices - 导入供应商发票")
	log.Println("  - GET  /api/excel/export/consumption-bills - 导出消费账单")

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
