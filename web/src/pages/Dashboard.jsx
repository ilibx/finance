import { Card, Row, Col, Statistic } from 'antd'
import {
  ProjectOutlined,
  UserOutlined,
  ShoppingCartOutlined,
  WalletOutlined,
} from '@ant-design/icons'

function Dashboard() {
  return (
    <div>
      <h2 style={{ marginBottom: 24 }}>仪表盘</h2>
      <Row gutter={16}>
        <Col span={6}>
          <Card>
            <Statistic
              title="项目总数"
              value={12}
              prefix={<ProjectOutlined />}
              suffix="个"
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title="用户总数"
              value={156}
              prefix={<UserOutlined />}
              suffix="人"
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title="订单总数"
              value={328}
              prefix={<ShoppingCartOutlined />}
              suffix="单"
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title="充值总额"
              value={89650}
              prefix={<WalletOutlined />}
              precision={2}
              suffix="元"
            />
          </Card>
        </Col>
      </Row>
    </div>
  )
}

export default Dashboard
