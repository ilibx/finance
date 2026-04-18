import { Menu } from 'antd'
import {
  DashboardOutlined,
  ProjectOutlined,
  UserOutlined,
  ShopOutlined,
  ShoppingCartOutlined,
  FileTextOutlined,
  WalletOutlined,
  TeamOutlined,
} from '@ant-design/icons'
import { useNavigate, useLocation } from 'react-router-dom'

const menuItems = [
  {
    key: '/dashboard',
    icon: <DashboardOutlined />,
    label: '仪表盘',
  },
  {
    key: '/projects',
    icon: <ProjectOutlined />,
    label: '项目管理',
  },
  {
    key: '/users',
    icon: <UserOutlined />,
    label: '用户管理',
  },
  {
    key: '/products',
    icon: <ShopOutlined />,
    label: '产品管理',
  },
  {
    key: '/orders',
    icon: <ShoppingCartOutlined />,
    label: '订单管理',
  },
  {
    key: '/invoices',
    icon: <FileTextOutlined />,
    label: '发票管理',
  },
  {
    key: '/recharges',
    icon: <WalletOutlined />,
    label: '充值管理',
  },
  {
    key: '/suppliers',
    icon: <TeamOutlined />,
    label: '供应商管理',
  },
]

function Sidebar() {
  const navigate = useNavigate()
  const location = useLocation()

  return (
    <div style={{ width: 200, background: '#001529' }}>
      <div style={{ padding: '16px', textAlign: 'center', color: '#fff', fontSize: '18px', fontWeight: 'bold' }}>
        ERP 系统
      </div>
      <Menu
        theme="dark"
        mode="inline"
        selectedKeys={[location.pathname]}
        items={menuItems}
        onClick={({ key }) => navigate(key)}
      />
    </div>
  )
}

export default Sidebar
