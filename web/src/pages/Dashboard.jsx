import { Card, Row, Col, Statistic, Table, Tag, Space, DatePicker, Button, Select } from 'antd'
import {
  ProjectOutlined,
  UserOutlined,
  ShoppingCartOutlined,
  WalletOutlined,
  RiseOutlined,
  FallOutlined,
  DollarOutlined,
  ShoppingOutlined,
} from '@ant-design/icons'
import { useState, useEffect } from 'react'
import api from '../services/api'
import dayjs from 'dayjs'

const { RangePicker } = DatePicker
const { Option } = Select

// 简单的柱状图组件（使用纯 CSS）
function SimpleBarChart({ data, title }) {
  if (!data || data.length === 0) return null
  
  const maxValue = Math.max(...data.map(item => item.value))
  
  return (
    <div style={{ marginTop: 20 }}>
      <h4 style={{ marginBottom: 16 }}>{title}</h4>
      <div style={{ display: 'flex', alignItems: 'flex-end', height: 200, gap: 8 }}>
        {data.map((item, index) => (
          <div key={index} style={{ flex: 1, textAlign: 'center' }}>
            <div 
              style={{
                height: `${(item.value / maxValue) * 160}px`,
                background: `linear-gradient(to top, #1890ff, #36cfc9)`,
                borderRadius: '4px 4px 0 0',
                minHeight: 4,
                transition: 'height 0.3s',
              }}
              title={`${item.name}: ${item.value}`}
            />
            <div style={{ fontSize: 12, marginTop: 8, color: '#666' }}>{item.name}</div>
            <div style={{ fontSize: 11, color: '#999' }}>¥{item.value.toLocaleString()}</div>
          </div>
        ))}
      </div>
    </div>
  )
}

// 简单的折线图组件（使用 SVG）
function SimpleLineChart({ data, title }) {
  if (!data || data.length === 0) return null
  
  const width = 500
  const height = 200
  const padding = 40
  const chartWidth = width - padding * 2
  const chartHeight = height - padding * 2
  
  const values = data.map(item => item.value)
  const maxValue = Math.max(...values)
  const minValue = Math.min(...values)
  const valueRange = maxValue - minValue || 1
  
  const points = data.map((item, index) => {
    const x = padding + (index / (data.length - 1)) * chartWidth
    const y = height - padding - ((item.value - minValue) / valueRange) * chartHeight
    return `${x},${y}`
  }).join(' ')
  
  const areaPoints = `${padding},${height - padding} ${points} ${width - padding},${height - padding}`
  
  return (
    <div style={{ marginTop: 20 }}>
      <h4 style={{ marginBottom: 16 }}>{title}</h4>
      <svg width="100%" height={height} viewBox={`0 0 ${width} ${height}`} style={{ maxWidth: '100%' }}>
        {/* 渐变填充 */}
        <defs>
          <linearGradient id="areaGradient" x1="0" y1="0" x2="0" y2="1">
            <stop offset="0%" stopColor="#1890ff" stopOpacity={0.3} />
            <stop offset="100%" stopColor="#1890ff" stopOpacity={0.05} />
          </linearGradient>
        </defs>
        
        {/* 网格线 */}
        {[0, 0.25, 0.5, 0.75, 1].map((ratio, i) => (
          <line
            key={i}
            x1={padding}
            y1={height - padding - ratio * chartHeight}
            x2={width - padding}
            y2={height - padding - ratio * chartHeight}
            stroke="#f0f0f0"
            strokeWidth={1}
          />
        ))}
        
        {/* 面积区域 */}
        <polygon points={areaPoints} fill="url(#areaGradient)" />
        
        {/* 折线 */}
        <polyline
          points={points}
          fill="none"
          stroke="#1890ff"
          strokeWidth={2}
          strokeLinecap="round"
          strokeLinejoin="round"
        />
        
        {/* 数据点 */}
        {data.map((item, index) => {
          const x = padding + (index / (data.length - 1)) * chartWidth
          const y = height - padding - ((item.value - minValue) / valueRange) * chartHeight
          return (
            <g key={index}>
              <circle cx={x} cy={y} r={4} fill="#fff" stroke="#1890ff" strokeWidth={2} />
              <text
                x={x}
                y={y - 12}
                textAnchor="middle"
                fontSize={10}
                fill="#666"
              >
                ¥{(item.value / 1000).toFixed(0)}k
              </text>
            </g>
          )
        })}
        
        {/* X 轴标签 */}
        {data.map((item, index) => {
          const x = padding + (index / (data.length - 1)) * chartWidth
          return (
            <text
              key={index}
              x={x}
              y={height - 15}
              textAnchor="middle"
              fontSize={10}
              fill="#999"
            >
              {item.name}
            </text>
          )
        })}
      </svg>
    </div>
  )
}

// 饼图组件（使用 SVG）
function SimplePieChart({ data, title }) {
  if (!data || data.length === 0) return null
  
  const total = data.reduce((sum, item) => sum + item.value, 0)
  const colors = ['#1890ff', '#52c41a', '#faad14', '#f5222d', '#722ed1', '#13c2c2', '#eb2f96', '#fa8c16']
  
  let currentAngle = 0
  const segments = data.map((item, index) => {
    const angle = (item.value / total) * 360
    const startAngle = currentAngle
    currentAngle += angle
    
    const startRad = (startAngle - 90) * Math.PI / 180
    const endRad = (startAngle + angle - 90) * Math.PI / 180
    
    const x1 = 100 + 80 * Math.cos(startRad)
    const y1 = 100 + 80 * Math.sin(startRad)
    const x2 = 100 + 80 * Math.cos(endRad)
    const y2 = 100 + 80 * Math.sin(endRad)
    
    const largeArc = angle > 180 ? 1 : 0
    
    const pathData = `M 100 100 L ${x1} ${y1} A 80 80 0 ${largeArc} 1 ${x2} ${y2} Z`
    
    return {
      path: pathData,
      color: colors[index % colors.length],
      name: item.name,
      value: item.value,
      percentage: ((item.value / total) * 100).toFixed(1)
    }
  })
  
  return (
    <div style={{ marginTop: 20 }}>
      <h4 style={{ marginBottom: 16 }}>{title}</h4>
      <div style={{ display: 'flex', alignItems: 'center', gap: 32 }}>
        <svg width="200" height="200" viewBox="0 0 200 200">
          {segments.map((seg, index) => (
            <path
              key={index}
              d={seg.path}
              fill={seg.color}
              stroke="#fff"
              strokeWidth={2}
              style={{ cursor: 'pointer' }}
            >
              <title>{`${seg.name}: ¥${seg.value.toLocaleString()} (${seg.percentage}%)`}</title>
            </path>
          ))}
          <text x="100" y="95" textAnchor="middle" fontSize={14} fill="#666">总计</text>
          <text x="100" y="115" textAnchor="middle" fontSize={16} fontWeight="bold" fill="#333">
            ¥{(total / 1000).toFixed(1)}k
          </text>
        </svg>
        
        <div style={{ flex: 1 }}>
          {segments.map((seg, index) => (
            <div key={index} style={{ display: 'flex', alignItems: 'center', marginBottom: 8 }}>
              <div style={{ width: 12, height: 12, borderRadius: 2, backgroundColor: seg.color, marginRight: 8 }} />
              <span style={{ flex: 1, fontSize: 13 }}>{seg.name}</span>
              <span style={{ fontSize: 13, color: '#666' }}>¥{seg.value.toLocaleString()}</span>
              <span style={{ fontSize: 12, color: '#999', marginLeft: 8 }}>({seg.percentage}%)</span>
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}

function Dashboard() {
  const [loading, setLoading] = useState(false)
  const [dashboardData, setDashboardData] = useState(null)
  const [trendData, setTrendData] = useState([])
  const [categoryData, setCategoryData] = useState([])
  const [dateRange, setDateRange] = useState([dayjs().subtract(1, 'month'), dayjs()])
  
  useEffect(() => {
    fetchDashboardData()
  }, [])
  
  const fetchDashboardData = async () => {
    setLoading(true)
    try {
      // 获取仪表盘汇总
      const summaryRes = await api.get('/financial/dashboard')
      if (summaryRes.code === 0) {
        setDashboardData(summaryRes.data)
      }
      
      // 获取利润趋势
      const startDate = dateRange[0].format('YYYY-MM-DD')
      const endDate = dateRange[1].format('YYYY-MM-DD')
      const trendRes = await api.get(`/financial/trend?start_date=${startDate}&end_date=${endDate}&group_by=month`)
      if (trendRes.code === 0 && trendRes.data && trendRes.data.trends) {
        setTrendData(trendRes.data.trends.map(item => ({
          name: item.period,
          value: item.profit
        })))
      }
      
      // 获取类别统计
      const categoryRes = await api.get(`/financial/category?start_date=${startDate}&end_date=${endDate}`)
      if (categoryRes.code === 0 && categoryRes.data && categoryRes.data.categories) {
        setCategoryData(categoryRes.data.categories.map(item => ({
          name: item.category,
          value: item.total
        })))
      }
    } catch (error) {
      console.error('获取数据失败:', error)
      // 使用模拟数据作为后备
      setDashboardData({
        total_revenue: 89650,
        total_expense: 45230,
        net_profit: 44420,
        profit_rate: 49.5,
        order_count: 328,
        product_count: 156
      })
      setTrendData([
        { name: '1 月', value: 32000 },
        { name: '2 月', value: 38000 },
        { name: '3 月', value: 42000 },
        { name: '4 月', value: 45000 },
        { name: '5 月', value: 48000 },
        { name: '6 月', value: 52000 }
      ])
      setCategoryData([
        { name: '电子产品', value: 35000 },
        { name: '办公用品', value: 22000 },
        { name: '原材料', value: 18000 },
        { name: '其他', value: 8000 }
      ])
    } finally {
      setLoading(false)
    }
  }
  
  const handleDateChange = (dates) => {
    if (dates && dates.length === 2) {
      setDateRange(dates)
    }
  }
  
  const revenueExpenseData = dashboardData ? [
    { name: '收入', value: dashboardData.total_revenue },
    { name: '支出', value: dashboardData.total_expense }
  ] : []
  
  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 24 }}>
        <h2 style={{ margin: 0 }}>财务仪表盘</h2>
        <Space>
          <RangePicker 
            value={dateRange}
            onChange={handleDateChange}
            onOk={fetchDashboardData}
          />
          <Button type="primary" onClick={fetchDashboardData} loading={loading}>
            刷新数据
          </Button>
        </Space>
      </div>
      
      {/* 关键指标卡片 */}
      <Row gutter={16} style={{ marginBottom: 24 }}>
        <Col span={6}>
          <Card>
            <Statistic
              title="总收入"
              value={dashboardData?.total_revenue || 0}
              prefix={<DollarOutlined />}
              precision={2}
              suffix="元"
              valueStyle={{ color: '#3f8600' }}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title="总支出"
              value={dashboardData?.total_expense || 0}
              prefix={<ShoppingOutlined />}
              precision={2}
              suffix="元"
              valueStyle={{ color: '#cf1322' }}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title="净利润"
              value={dashboardData?.net_profit || 0}
              prefix={(dashboardData?.net_profit >= 0 ? <RiseOutlined /> : <FallOutlined />)}
              precision={2}
              suffix="元"
              valueStyle={{ color: dashboardData?.net_profit >= 0 ? '#3f8600' : '#cf1322' }}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title="利润率"
              value={dashboardData?.profit_rate || 0}
              prefix={<WalletOutlined />}
              precision={1}
              suffix="%"
              valueStyle={{ color: '#1890ff' }}
            />
          </Card>
        </Col>
      </Row>
      
      {/* 图表区域 */}
      <Row gutter={16}>
        <Col span={12}>
          <Card title="收支对比">
            <SimpleBarChart data={revenueExpenseData} title="" />
          </Card>
        </Col>
        <Col span={12}>
          <Card title="支出分类">
            <SimplePieChart data={categoryData} title="" />
          </Card>
        </Col>
      </Row>
      
      <Row gutter={16} style={{ marginTop: 16 }}>
        <Col span={24}>
          <Card title="利润趋势">
            <SimpleLineChart data={trendData} title="" />
          </Card>
        </Col>
      </Row>
      
      {/* 额外统计信息 */}
      <Row gutter={16} style={{ marginTop: 16 }}>
        <Col span={8}>
          <Card>
            <Statistic
              title="订单总数"
              value={dashboardData?.order_count || 0}
              prefix={<ShoppingCartOutlined />}
              suffix="单"
            />
          </Card>
        </Col>
        <Col span={8}>
          <Card>
            <Statistic
              title="产品总数"
              value={dashboardData?.product_count || 0}
              prefix={<ProjectOutlined />}
              suffix="个"
            />
          </Card>
        </Col>
        <Col span={8}>
          <Card>
            <Statistic
              title="用户总数"
              value={156}
              prefix={<UserOutlined />}
              suffix="人"
            />
          </Card>
        </Col>
      </Row>
    </div>
  )
}

export default Dashboard
