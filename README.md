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

### 后端
- **语言**: Golang 1.19+
- **数据库**: PostgreSQL 12+
- **Web 框架**: Gin Web Framework
- **ORM**: Gorm v2
- **Excel 处理**: excelize/v2

### 前端
- **框架**: React 18
- **构建工具**: Vite 5
- **UI 组件库**: Ant Design 5
- **路由**: React Router 6
- **HTTP 客户端**: Axios

## 项目结构

```
erp-system/
├── cmd/                    # 应用程序入口
│   └── main.go            # 主程序
├── internal/              # 内部包
│   ├── config/           # 配置管理
│   ├── handler/          # HTTP 处理器
│   ├── domain/           # 领域层 (DDD)
│   │   ├── user/        # 用户领域
│   │   ├── product/     # 产品领域
│   │   ├── order/       # 订单领域
│   │   ├── invoice/     # 发票领域
│   │   ├── recharge/    # 充值领域
│   │   ├── supplier/    # 供应商领域
│   │   └── project/     # 项目领域
│   └── common/          # 通用组件
├── infrastructure/       # 基础设施层
│   ├── database/        # 数据库连接
│   └── excel/           # Excel 处理
├── interfaces/           # 接口层
│   └── graphql/         # GraphQL API
├── web/                  # 前端界面 (React + Vite)
│   ├── src/
│   │   ├── components/  # 公共组件
│   │   ├── pages/       # 页面组件
│   │   └── services/    # API 服务
│   └── package.json
├── sql/                  # 数据库脚本
│   └── init.sql         # 数据库初始化脚本
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

### 4. 启动服务

#### 后端服务

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

后端服务将在 http://localhost:8080 启动

#### 前端服务

```bash
cd web

# 安装依赖
npm install

# 启动开发服务器
npm run dev
```

前端界面将在 http://localhost:3000 启动

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

## 开发状态

### ✅ 已完成功能

#### 核心模块
- [x] **项目管理**: 项目 CRUD、状态更新、进度跟踪、成本核算（完整 API 已实现）
- [x] **用户管理**: 用户创建、信息查询、余额管理（完整 API 已实现）
- [x] **产品管理**: 产品创建、SKU 管理、价格与库存（API 已实现）
- [x] **订单管理**: 订单创建、订单项管理、状态流转（API 已实现）
- [x] **发票管理**: 销售/采购发票生成、发票状态跟踪（API 已实现）
- [x] **充值管理**: 用户/供应商充值记录、充值处理（API 已实现）
- [x] **供应商管理**: 供应商实体/服务层已实现，Excel 导入功能已完成

#### Excel 导入导出
- [x] 消费账单导入/导出
- [x] 用户充值记录导入
- [x] 供应商充值记录导入
- [x] 供应商发票导入

#### 前端界面
- [x] 项目架构：React + Vite + Ant Design
- [x] 仪表盘页面：数据概览展示
- [x] 项目管理页面：CRUD 操作
- [x] 用户管理页面：用户增删改查
- [x] 产品管理页面：产品信息管理
- [x] 订单管理页面：订单处理与导出
- [x] 发票管理页面：发票导入导出
- [x] 充值管理页面：充值记录管理
- [x] 供应商管理页面：供应商信息与导入

#### 技术架构
- [x] 领域驱动设计 (DDD) 项目结构
- [x] PostgreSQL 数据库设计与初始化脚本
- [x] RESTful API 接口实现
- [x] GraphQL Schema 定义

---

### 🔧 待完善功能

#### 核心业务功能
- [ ] **供应商管理 API**: 供应商实体/服务已实现，但缺少独立的 HTTP Handler 和 API 端点
- [ ] **完整采购流程**: 采购申请、审批、入库流程（当前仅有供应商基础数据）
- [ ] **库存预警**: 库存阈值设置与自动预警功能
- [ ] **财务报表分析**: 收支统计、利润分析、可视化报表（当前仅有基础对账）

#### 技术基础设施
- [x] **前端界面**: React + Vite + Ant Design 基础框架已完成，包含 8 个核心页面
- [ ] **用户认证授权**: JWT/OAuth2 登录与权限控制（当前无任何认证机制）
- [ ] **GraphQL Resolver**: 完成 GraphQL API 实现（Schema 已定义在 `interfaces/graphql/schema.graphql`，但无 Resolver 实现）
- [ ] **API 文档**: OpenAPI/Swagger 文档生成
- [ ] **数据备份恢复**: 定时备份、一键恢复功能

#### 测试与质量保障
- [ ] **单元测试**: 各模块测试用例覆盖（当前无任何 `_test.go` 文件）
- [ ] **集成测试**: API 端到端测试
- [ ] **代码覆盖率检查**: CI/CD 集成

---

### 📋 计划中功能

- [ ] 多语言支持 (i18n)
- [ ] 操作日志审计
- [ ] 消息通知系统（邮件、短信）
- [ ] 移动端适配
- [ ] 第三方系统集成（ERP、CRM）
- [ ] 高级搜索与筛选
- [ ] 数据导出为 PDF/CSV 格式
- [ ] 工作流引擎支持

## 许可证

MIT License
