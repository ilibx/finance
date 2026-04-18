import { useState, useEffect } from 'react'
import { Table, Button, Space, Modal, Form, Input, message, Tag, Upload } from 'antd'
import { PlusOutlined, EditOutlined, DeleteOutlined, UploadOutlined } from '@ant-design/icons'
import api from '../services/api'

function SupplierList() {
  const [suppliers, setSuppliers] = useState([])
  const [loading, setLoading] = useState(false)
  const [modalVisible, setModalVisible] = useState(false)
  const [editingSupplier, setEditingSupplier] = useState(null)
  const [form] = Form.useForm()

  const fetchSuppliers = async () => {
    setLoading(true)
    try {
      const response = await api.get('/suppliers')
      setSuppliers(response.data || [])
    } catch (error) {
      message.error('获取供应商列表失败')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchSuppliers()
  }, [])

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields()
      if (editingSupplier) {
        await api.put(`/suppliers/${editingSupplier.id}`, values)
        message.success('更新成功')
      } else {
        await api.post('/suppliers', values)
        message.success('创建成功')
      }
      setModalVisible(false)
      fetchSuppliers()
    } catch (error) {
      message.error('操作失败')
    }
  }

  const handleDelete = async (id) => {
    Modal.confirm({
      title: '确认删除',
      content: '确定要删除这个供应商吗？',
      onOk: async () => {
        try {
          await api.delete(`/suppliers/${id}`)
          message.success('删除成功')
          fetchSuppliers()
        } catch (error) {
          message.error('删除失败')
        }
      },
    })
  }

  const handleImport = async (file) => {
    const formData = new FormData()
    formData.append('file', file)
    try {
      await api.post('/suppliers/import', formData, { headers: { 'Content-Type': 'multipart/form-data' } })
      message.success('导入成功')
      fetchSuppliers()
    } catch (error) {
      message.error('导入失败')
    }
    return false
  }

  const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 80 },
    { title: '供应商名称', dataIndex: 'name', key: 'name' },
    { title: '联系人', dataIndex: 'contact', key: 'contact' },
    { title: '电话', dataIndex: 'phone', key: 'phone' },
    { title: '邮箱', dataIndex: 'email', key: 'email' },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status) => <Tag color={status === 'active' ? 'green' : 'red'}>{status === 'active' ? '合作中' : '已停用'}</Tag>,
    },
    {
      title: '操作',
      key: 'action',
      render: (_, record) => (
        <Space>
          <Button icon={<EditOutlined />} onClick={() => { setEditingSupplier(record); form.setFieldsValue(record); setModalVisible(true) }}>编辑</Button>
          <Button icon={<DeleteOutlined />} danger onClick={() => handleDelete(record.id)}>删除</Button>
        </Space>
      ),
    },
  ]

  return (
    <div>
      <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between' }}>
        <h2>供应商管理</h2>
        <Space>
          <Upload showUploadList={false} beforeUpload={handleImport}>
            <Button icon={<UploadOutlined />}>导入</Button>
          </Upload>
          <Button type="primary" icon={<PlusOutlined />} onClick={() => { setEditingSupplier(null); form.resetFields(); setModalVisible(true) }}>新建供应商</Button>
        </Space>
      </div>
      <Table columns={columns} dataSource={suppliers} loading={loading} rowKey="id" />
      <Modal title={editingSupplier ? '编辑供应商' : '新建供应商'} open={modalVisible} onOk={handleSubmit} onCancel={() => setModalVisible(false)}>
        <Form form={form} layout="vertical">
          <Form.Item name="name" label="供应商名称" rules={[{ required: true }]}><Input /></Form.Item>
          <Form.Item name="contact" label="联系人"><Input /></Form.Item>
          <Form.Item name="phone" label="电话"><Input /></Form.Item>
          <Form.Item name="email" label="邮箱"><Input /></Form.Item>
          <Form.Item name="address" label="地址"><Input.TextArea rows={3} /></Form.Item>
          <Form.Item name="status" label="状态" initialValue="active">
            <Select>
              <Select.Option value="active">合作中</Select.Option>
              <Select.Option value="inactive">已停用</Select.Option>
            </Select>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}

export default SupplierList
