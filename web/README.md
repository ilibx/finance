# ERP Web 前端

基于 React + Vite + Ant Design 的 ERP 系统前端界面。

## 技术栈

- **框架**: React 18
- **构建工具**: Vite 5
- **UI 组件库**: Ant Design 5
- **路由**: React Router 6
- **HTTP 客户端**: Axios

## 功能模块

- ✅ 仪表盘 - 数据概览
- ✅ 项目管理 - CRUD 操作
- ✅ 用户管理 - 用户增删改查
- ✅ 产品管理 - 产品信息管理
- ✅ 订单管理 - 订单处理与导出
- ✅ 发票管理 - 发票导入导出
- ✅ 充值管理 - 充值记录管理
- ✅ 供应商管理 - 供应商信息与导入

## 快速开始

### 安装依赖

```bash
cd web
npm install
```

### 启动开发服务器

```bash
npm run dev
```

访问 http://localhost:3000

### 构建生产版本

```bash
npm run build
```

## 项目结构

```
web/
├── src/
│   ├── components/      # 公共组件
│   │   ├── Sidebar.jsx  # 侧边栏
│   │   └── Header.jsx   # 顶部导航
│   ├── pages/           # 页面组件
│   │   ├── Dashboard.jsx
│   │   ├── ProjectList.jsx
│   │   ├── UserList.jsx
│   │   ├── ProductList.jsx
│   │   ├── OrderList.jsx
│   │   ├── InvoiceList.jsx
│   │   ├── RechargeList.jsx
│   │   └── SupplierList.jsx
│   ├── services/        # API 服务
│   │   └── api.js       # Axios 配置
│   ├── styles/          # 样式文件
│   │   └── index.css
│   ├── App.jsx          # 应用入口
│   └── main.jsx         # 主入口
├── index.html
├── package.json
└── vite.config.js
```

## API 配置

默认代理配置为后端 API 地址 `http://localhost:8080`，可在 `vite.config.js` 中修改。

## 待完善功能

- [ ] 用户登录/认证
- [ ] 权限控制
- [ ] 数据可视化图表
- [ ] 高级搜索筛选
- [ ] 批量操作
- [ ] 响应式布局优化
