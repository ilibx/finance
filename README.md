# ERP 系统

基于 Golang 实现的 ERP（企业资源计划）管理系统，类似 Odoo 的功能架构。

## 功能模块

### 1. 项目管理
- 项目创建与管理
- 项目进度跟踪
- 项目成本核算

### 2. 财务对账
- 用户充值记录管理
- 供应商充值记录管理
- 收支对账
- 财务报表

### 3. 发票管理
- 销售发票生成
- 采购发票管理
- 发票导入导出
- 发票状态跟踪

### 4. 供应链管理
- 供应商管理
- 采购订单管理
- 库存管理
- 供应商充值与发票

### 5. 销售管理
- 客户/用户管理
- 产品销售订单
- 消费账单管理
- 销售发票生成

## Excel 导入导出功能

支持以下 Excel 文件的导入导出：

- **导入消费账单**: 批量导入用户消费记录
- **导入用户充值记录**: 批量导入用户充值数据
- **导入供应商充值记录**: 批量导入供应商充值数据
- **导入供应商发票**: 批量导入供应商发票信息
- **导出用户消费账单**: 导出用户消费明细
- **导出销售发票**: 导出生成的销售发票

## 技术栈

- **语言**: Golang 1.19+
- **数据库**: PostgreSQL
- **Web 框架**: 标准库 net/http
- **Excel 处理**: excelize/v2
- **数据库驱动**: lib/pq

## 项目结构

```
erp-system/
├── cmd/                    # 应用程序入口
│   └── main.go            # 主程序
├── internal/              # 内部包
│   ├── config/           # 配置管理
│   ├── handler/          # HTTP 处理器
│   ├── model/            # 数据模型
│   ├── repository/       # 数据访问层
│   └── service/          # 业务逻辑层
├── pkg/                  # 公共包
│   └── excel/           # Excel 处理工具
├── sql/                  # 数据库脚本
│   └── init.sql         # 数据库初始化脚本
├── web/                  # 前端资源
│   ├── static/          # 静态文件
│   └── templates/       # HTML 模板
├── go.mod               # Go 模块定义
├── go.sum               # 依赖校验
└── README.md            # 项目说明
```

## 快速开始

### 1. 环境要求

- Go 1.19+
- PostgreSQL 12+

### 2. 数据库初始化

```bash
# 创建数据库
createdb erp_db

# 执行初始化脚本
psql -d erp_db -f sql/init.sql
```

### 3. 配置环境变量

```bash
export SERVER_PORT=8080
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=postgres
export DB_NAME=erp_db
```

### 4. 编译运行

```bash
# 编译
go build -o erp-system ./cmd/main.go

# 运行
./erp-system
```

或使用 Go 直接运行：

```bash
go run ./cmd/main.go
```

## API 端点

### 用户管理
- `POST /api/users/create` - 创建用户
- `GET /api/users/get?id=1` - 获取用户信息

### 产品管理
- `POST /api/products/create` - 创建产品

### 订单管理
- `POST /api/orders/create` - 创建订单

### 发票管理
- `POST /api/invoices/generate` - 生成发票

### 充值管理
- `POST /api/recharge/process` - 处理充值

### Excel 导入导出
- `POST /api/excel/import/consumption-bills` - 导入消费账单
- `POST /api/excel/import/recharge-records` - 导入充值记录
- `POST /api/excel/import/supplier-recharges` - 导入供应商充值
- `POST /api/excel/import/supplier-invoices` - 导入供应商发票
- `GET /api/excel/export/consumption-bills` - 导出消费账单

## Excel 文件格式

### 消费账单导入格式
| 用户 ID | 订单号 | 金额 | 描述 | 日期 |
|--------|--------|------|------|------|
| 1 | ORD20240101001 | 100.00 | 产品购买 | 2024-01-01 |

### 用户充值记录导入格式
| 用户 ID | 金额 | 支付方式 | 备注 | 日期 |
|--------|------|----------|------|------|
| 1 | 1000.00 | bank_transfer | 银行转账 | 2024-01-01 |

### 供应商充值记录导入格式
| 供应商 ID | 金额 | 支付方式 | 备注 | 日期 |
|----------|------|----------|------|------|
| 1 | 5000.00 | bank_transfer | 预付款 | 2024-01-01 |

### 供应商发票导入格式
| 供应商 ID | 发票号 | 金额 | 税额 | 日期 |
|----------|--------|------|------|------|
| 1 | INV001 | 10000.00 | 1300.00 | 2024-01-01 |

## 开发计划

- [ ] 完善前端界面
- [ ] 添加用户认证授权
- [ ] 实现完整的采购流程
- [ ] 添加库存预警功能
- [ ] 实现财务报表分析
- [ ] 添加数据备份恢复功能
- [ ] 支持更多 Excel 模板
- [ ] 添加 API 文档

## 许可证

MIT License
