import { useState, useEffect } from 'react'
import { Table, Button, Space, Modal, Form, Input, InputNumber, Select, message, Tag } from 'antd'
import { PlusOutlined, EditOutlined, DeleteOutlined, DownloadOutlined } from '@ant-design/icons'
import api from '../services/api'

function OrderList() {
  const [orders, setOrders] = useState([])
  const [loading, setLoading] = useState(false)
  const [modalVisible, setModalVisible] = useState(false)
  const [editingOrder, setEditingOrder] = useState(null)
  const [form] = Form.useForm()

  const fetchOrders = async () => {
    setLoading(true)
    try {
      const response = await api.get('/orders')
      setOrders(response.data || [])
    } catch (error) {
      message.error('获取订单列表失败')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchOrders()
  }, [])

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields()
      if (editingOrder) {
        await api.put(`/orders/${editingOrder.id}`, values)
        message.success('更新成功')
      } else {
        await api.post('/orders', values)
        message.success('创建成功')
      }
      setModalVisible(false)
      fetchOrders()
    } catch (error) {
      message.error('操作失败')
    }
  }

  const handleDelete = async (id) => {
    Modal.confirm({
      title: '确认删除',
      content: '确定要删除这个订单吗？',
      onOk: async () => {
        try {
          await api.delete(`/orders/${id}`)
          message.success('删除成功')
          fetchOrders()
        } catch (error) {
          message.error('删除失败')
        }
      },
    })
  }

  const handleExport = async () => {
    try {
      const response = await api.get('/orders/export', { responseType: 'blob' })
      const url = window.URL.createObjectURL(new Blob([response]))
      const link = document.createElement('a')
      link.href = url
      link.setAttribute('download', 'orders.xlsx')
      document.body.appendChild(link)
      link.click()
      document.body.removeChild(link)
      message.success('导出成功')
    } catch (error) {
      message.error('导出失败')
    }
  }

  const columns = [
    { title: '订单号', dataIndex: 'orderNo', key: 'orderNo' },
    { title: '客户', dataIndex: 'customerName', key: 'customerName' },
    { title: '金额', dataIndex: 'amount', key: 'amount', render: (v) => `¥${v?.toFixed(2)}` },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status) => {
        const map = { pending: '待处理', processing: '处理中', completed: '已完成', cancelled: '已取消' }
        const colors = { pending: 'orange', processing: 'blue', completed: 'green', cancelled: 'red' }
        return <Tag color={colors[status]}>{map[status] || status}</Tag>
      },
    },
    {
      title: '操作',
      key: 'action',
      render: (_, record) => (
        <Space>
          <Button icon={<EditOutlined />} onClick={() => { setEditingOrder(record); form.setFieldsValue(record); setModalVisible(true) }}>编辑</Button>
          <Button icon={<DeleteOutlined />} danger onClick={() => handleDelete(record.id)}>删除</Button>
        </Space>
      ),
    },
  ]

  return (
    <div>
      <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between' }}>
        <h2>订单管理</h2>
        <Space>
          <Button icon={<DownloadOutlined />} onClick={handleExport}>导出</Button>
          <Button type="primary" icon={<PlusOutlined />} onClick={() => { setEditingOrder(null); form.resetFields(); setModalVisible(true) }}>新建订单</Button>
        </Space>
      </div>
      <Table columns={columns} dataSource={orders} loading={loading} rowKey="id" />
      <Modal title={editingOrder ? '编辑订单' : '新建订单'} open={modalVisible} onOk={handleSubmit} onCancel={() => setModalVisible(false)}>
        <Form form={form} layout="vertical">
          <Form.Item name="orderNo" label="订单号" rules={[{ required: true }]}><Input /></Form.Item>
          <Form.Item name="customerName" label="客户名称" rules={[{ required: true }]}><Input /></Form.Item>
          <Form.Item name="amount" label="金额" rules={[{ required: true }]}><InputNumber min={0} step={0.01} style={{ width: '100%' }} /></Form.Item>
          <Form.Item name="status" label="状态" initialValue="pending">
            <Select>
              <Select.Option value="pending">待处理</Select.Option>
              <Select.Option value="processing">处理中</Select.Option>
              <Select.Option value="completed">已完成</Select.Option>
              <Select.Option value="cancelled">已取消</Select.Option>
            </Select>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}

export default OrderList
