import { useState, useEffect } from 'react'
import { Table, Button, Space, Modal, Form, Input, InputNumber, message } from 'antd'
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons'
import api from '../services/api'

function ProductList() {
  const [products, setProducts] = useState([])
  const [loading, setLoading] = useState(false)
  const [modalVisible, setModalVisible] = useState(false)
  const [editingProduct, setEditingProduct] = useState(null)
  const [form] = Form.useForm()

  const fetchProducts = async () => {
    setLoading(true)
    try {
      const response = await api.get('/products')
      setProducts(response.data || [])
    } catch (error) {
      message.error('获取产品列表失败')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchProducts()
  }, [])

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields()
      if (editingProduct) {
        await api.put(`/products/${editingProduct.id}`, values)
        message.success('更新成功')
      } else {
        await api.post('/products', values)
        message.success('创建成功')
      }
      setModalVisible(false)
      fetchProducts()
    } catch (error) {
      message.error('操作失败')
    }
  }

  const handleDelete = async (id) => {
    Modal.confirm({
      title: '确认删除',
      content: '确定要删除这个产品吗？',
      onOk: async () => {
        try {
          await api.delete(`/products/${id}`)
          message.success('删除成功')
          fetchProducts()
        } catch (error) {
          message.error('删除失败')
        }
      },
    })
  }

  const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 80 },
    { title: '产品名称', dataIndex: 'name', key: 'name' },
    { title: '编码', dataIndex: 'code', key: 'code' },
    { title: '价格', dataIndex: 'price', key: 'price', render: (v) => `¥${v?.toFixed(2)}` },
    { title: '库存', dataIndex: 'stock', key: 'stock' },
    {
      title: '操作',
      key: 'action',
      render: (_, record) => (
        <Space>
          <Button icon={<EditOutlined />} onClick={() => { setEditingProduct(record); form.setFieldsValue(record); setModalVisible(true) }}>编辑</Button>
          <Button icon={<DeleteOutlined />} danger onClick={() => handleDelete(record.id)}>删除</Button>
        </Space>
      ),
    },
  ]

  return (
    <div>
      <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between' }}>
        <h2>产品管理</h2>
        <Button type="primary" icon={<PlusOutlined />} onClick={() => { setEditingProduct(null); form.resetFields(); setModalVisible(true) }}>新建产品</Button>
      </div>
      <Table columns={columns} dataSource={products} loading={loading} rowKey="id" />
      <Modal title={editingProduct ? '编辑产品' : '新建产品'} open={modalVisible} onOk={handleSubmit} onCancel={() => setModalVisible(false)}>
        <Form form={form} layout="vertical">
          <Form.Item name="name" label="产品名称" rules={[{ required: true }]}><Input /></Form.Item>
          <Form.Item name="code" label="编码" rules={[{ required: true }]}><Input /></Form.Item>
          <Form.Item name="price" label="价格" rules={[{ required: true }]}><InputNumber min={0} step={0.01} style={{ width: '100%' }} /></Form.Item>
          <Form.Item name="stock" label="库存" initialValue={0}><InputNumber min={0} style={{ width: '100%' }} /></Form.Item>
        </Form>
      </Modal>
    </div>
  )
}

export default ProductList
