# Vue + GraphQL + Ant Design 项目

这是一个基于 Vue 3 + Vite + Ant Design Vue + GraphQL 的前端交互界面项目。

## 技术栈

- **Vue 3** - 渐进式 JavaScript 框架 (Composition API)
- **Vite** - 下一代前端构建工具
- **Ant Design Vue 4** - 企业级 UI 组件库
- **GraphQL** - 数据查询语言
- **graphql-request** - 轻量级 GraphQL 客户端

## 功能特性

### 核心功能
- ✅ 用户列表展示（表格形式，支持分页）
- ✅ 创建新用户（模态框表单）
- ✅ 编辑现有用户
- ✅ 删除用户（带确认提示）
- ✅ 刷新数据
- ✅ 响应式布局
- ✅ 加载状态处理
- ✅ 错误处理与消息提示
- ✅ GraphQL 查询和突变
- ✅ 演示数据回退（当后端不可用时）

### UI 组件
- 📊 Table - 数据表格展示
- 📝 Form - 表单输入验证
- 🔘 Button - 操作按钮
- 🗑️ Popconfirm - 删除确认
- 📋 Modal - 创建/编辑对话框
- 🔔 Message - 消息提示
- 🎨 Layout - 页面布局

## 快速开始

### 1. 安装依赖

```bash
npm install
```

### 2. 配置环境变量（可选）

复制 `.env.example` 为 `.env` 并修改 GraphQL API 地址：

```bash
cp .env.example .env
```

编辑 `.env` 文件：

```env
VITE_GRAPHQL_ENDPOINT=http://your-graphql-api.com/graphql
```

如果不配置，默认使用 `http://localhost:4000/graphql`

### 3. 启动开发服务器

```bash
npm run dev
```

访问 `http://localhost:5173` 查看应用

## 项目结构

```
vue-graphql-app/
├── src/
│   ├── components/
│   │   └── UserManagement.vue    # 用户管理主组件（完整 CRUD 功能）
│   ├── graphql/
│   │   ├── client.js             # GraphQL 客户端配置
│   │   └── queries.js            # GraphQL 查询和突变定义
│   ├── assets/                   # 静态资源
│   ├── App.vue                   # 根组件（带布局）
│   ├── main.js                   # 入口文件（集成 Antd）
│   └── style.css                 # 全局样式
├── public/                       # 公共静态文件
├── .env.example                  # 环境变量示例
├── index.html                    # HTML 模板
├── package.json                  # 项目依赖
├── vite.config.js                # Vite 配置
└── README.md                     # 项目文档
```

## GraphQL Schema 示例

后端需要支持以下 GraphQL schema：

```graphql
type User {
  id: ID!
  name: String!
  email: String!
  createdAt: String!
  updatedAt: String
}

type DeleteResponse {
  success: Boolean!
  message: String!
}

input CreateUserInput {
  name: String!
  email: String!
}

input UpdateUserInput {
  name: String
  email: String
}

type Query {
  users: [User!]!
  user(id: ID!): User
}

type Mutation {
  createUser(input: CreateUserInput!): User!
  updateUser(id: ID!, input: UpdateUserInput!): User!
  deleteUser(id: ID!): DeleteResponse!
}
```

## GraphQL 操作示例

### 查询用户列表
```graphql
query GetUsers {
  users {
    id
    name
    email
    createdAt
  }
}
```

### 创建用户
```graphql
mutation CreateUser($input: CreateUserInput!) {
  createUser(input: $input) {
    id
    name
    email
    createdAt
  }
}
```

### 更新用户
```graphql
mutation UpdateUser($id: ID!, $input: UpdateUserInput!) {
  updateUser(id: $id, input: $input) {
    id
    name
    email
    updatedAt
  }
}
```

### 删除用户
```graphql
mutation DeleteUser($id: ID!) {
  deleteUser(id: $id) {
    success
    message
  }
}
```

## 自定义组件

你可以根据需要修改 `src/components/UserManagement.vue` 来适配你的业务需求：

1. **修改 GraphQL 查询** - 编辑 `src/graphql/queries.js`
2. **调整表格列** - 修改 `columns` 数组定义
3. **添加更多字段** - 扩展表单模型和 GraphQL input
4. **集成认证** - 在 `graphql/client.js` 中添加 token
5. **自定义样式** - 修改组件 scoped styles 或全局样式

## 构建生产版本

```bash
npm run build
```

构建产物将在 `dist/` 目录中生成。

## 预览生产版本

```bash
npm run preview
```

## 注意事项

- **演示模式**: 如果后端 GraphQL API 不可用，应用会自动显示演示数据
- **CORS**: 确保后端支持跨域资源共享
- **认证授权**: 根据需要在 `graphql/client.js` 中添加身份验证逻辑
- **错误处理**: 已实现完整的错误处理和用户提示

## 浏览器支持

- Chrome (最新版)
- Firefox (最新版)
- Safari (最新版)
- Edge (最新版)

## 相关资源

- [Vue 3 文档](https://vuejs.org/)
- [Ant Design Vue](https://antdv.com/)
- [GraphQL 官方文档](https://graphql.org/)
- [Vite 文档](https://vitejs.dev/)
- [graphql-request](https://github.com/graffle-js/graffle)

## License

MIT
