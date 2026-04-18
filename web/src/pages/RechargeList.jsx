import { useState, useEffect } from 'react'
import { Table, Button, Space, Modal, Form, Input, InputNumber, Select, message, Tag, Upload } from 'antd'
import { PlusOutlined, EditOutlined, DeleteOutlined, DownloadOutlined, UploadOutlined } from '@ant-design/icons'
import api from '../services/api'

function RechargeList() {
  const [recharges, setRecharges] = useState([])
  const [loading, setLoading] = useState(false)
  const [modalVisible, setModalVisible] = useState(false)
  const [editingRecharge, setEditingRecharge] = useState(null)
  const [form] = Form.useForm()

  const fetchRecharges = async () => {
    setLoading(true)
    try {
      const response = await api.get('/recharges')
      setRecharges(response.data || [])
    } catch (error) {
      message.error('获取充值记录失败')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchRecharges()
  }, [])

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields()
      if (editingRecharge) {
        await api.put(`/recharges/${editingRecharge.id}`, values)
        message.success('更新成功')
      } else {
        await api.post('/recharges', values)
        message.success('创建成功')
      }
      setModalVisible(false)
      fetchRecharges()
    } catch (error) {
      message.error('操作失败')
    }
  }

  const handleDelete = async (id) => {
    Modal.confirm({
      title: '确认删除',
      content: '确定要删除这条充值记录吗？',
      onOk: async () => {
        try {
          await api.delete(`/recharges/${id}`)
          message.success('删除成功')
          fetchRecharges()
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
      await api.post('/recharges/import', formData, { headers: { 'Content-Type': 'multipart/form-data' } })
      message.success('导入成功')
      fetchRecharges()
    } catch (error) {
      message.error('导入失败')
    }
    return false
  }

  const columns = [
    { title: '流水号', dataIndex: 'transactionNo', key: 'transactionNo' },
    { title: '用户', dataIndex: 'userName', key: 'userName' },
    { title: '金额', dataIndex: 'amount', key: 'amount', render: (v) => `¥${v?.toFixed(2)}` },
    { title: '充值时间', dataIndex: 'rechargeTime', key: 'rechargeTime' },
    {
      title: '类型',
      dataIndex: 'type',
      key: 'type',
      render: (type) => <Tag color={type === 'user' ? 'blue' : 'green'}>{type === 'user' ? '用户充值' : '供应商充值'}</Tag>,
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status) => {
        const map = { pending: '待处理', completed: '已完成', failed: '失败' }
        return <Tag color={status === 'completed' ? 'green' : 'orange'}>{map[status] || status}</Tag>
      },
    },
    {
      title: '操作',
      key: 'action',
      render: (_, record) => (
        <Space>
          <Button icon={<EditOutlined />} onClick={() => { setEditingRecharge(record); form.setFieldsValue(record); setModalVisible(true) }}>编辑</Button>
          <Button icon={<DeleteOutlined />} danger onClick={() => handleDelete(record.id)}>删除</Button>
        </Space>
      ),
    },
  ]

  return (
    <div>
      <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between' }}>
        <h2>充值管理</h2>
        <Space>
          <Upload showUploadList={false} beforeUpload={handleImport}>
            <Button icon={<UploadOutlined />}>导入</Button>
          </Upload>
          <Button icon={<DownloadOutlined />} onClick={() => message.info('导出功能开发中')}>导出</Button>
          <Button type="primary" icon={<PlusOutlined />} onClick={() => { setEditingRecharge(null); form.resetFields(); setModalVisible(true) }}>新建充值</Button>
        </Space>
      </div>
      <Table columns={columns} dataSource={recharges} loading={loading} rowKey="id" />
      <Modal title={editingRecharge ? '编辑充值' : '新建充值'} open={modalVisible} onOk={handleSubmit} onCancel={() => setModalVisible(false)}>
        <Form form={form} layout="vertical">
          <Form.Item name="transactionNo" label="流水号" rules={[{ required: true }]}><Input /></Form.Item>
          <Form.Item name="userName" label="用户名称" rules={[{ required: true }]}><Input /></Form.Item>
          <Form.Item name="amount" label="金额" rules={[{ required: true }]}><InputNumber min={0} step={0.01} style={{ width: '100%' }} /></Form.Item>
          <Form.Item name="type" label="类型" initialValue="user">
            <Select>
              <Select.Option value="user">用户充值</Select.Option>
              <Select.Option value="supplier">供应商充值</Select.Option>
            </Select>
          </Form.Item>
          <Form.Item name="status" label="状态" initialValue="completed">
            <Select>
              <Select.Option value="pending">待处理</Select.Option>
              <Select.Option value="completed">已完成</Select.Option>
              <Select.Option value="failed">失败</Select.Option>
            </Select>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}

export default RechargeList
