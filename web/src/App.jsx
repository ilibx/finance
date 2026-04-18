import { Routes, Route, Navigate } from 'react-router-dom'
import { Layout } from 'antd'
import Dashboard from './pages/Dashboard'
import ProjectList from './pages/ProjectList'
import UserList from './pages/UserList'
import ProductList from './pages/ProductList'
import OrderList from './pages/OrderList'
import InvoiceList from './pages/InvoiceList'
import RechargeList from './pages/RechargeList'
import SupplierList from './pages/SupplierList'
import Sidebar from './components/Sidebar'
import Header from './components/Header'

const { Content } = Layout

function App() {
  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Sidebar />
      <Layout>
        <Header />
        <Content style={{ margin: '24px 16px', padding: 24, background: '#fff' }}>
          <Routes>
            <Route path="/" element={<Navigate to="/dashboard" replace />} />
            <Route path="/dashboard" element={<Dashboard />} />
            <Route path="/projects" element={<ProjectList />} />
            <Route path="/users" element={<UserList />} />
            <Route path="/products" element={<ProductList />} />
            <Route path="/orders" element={<OrderList />} />
            <Route path="/invoices" element={<InvoiceList />} />
            <Route path="/recharges" element={<RechargeList />} />
            <Route path="/suppliers" element={<SupplierList />} />
          </Routes>
        </Content>
      </Layout>
    </Layout>
  )
}

export default App
