import { useState, useEffect } from 'react'
import { Table, Button, Space, Modal, Form, Input, InputNumber, Select, message, Tag, Upload } from 'antd'
import { PlusOutlined, EditOutlined, DeleteOutlined, DownloadOutlined, UploadOutlined } from '@ant-design/icons'
import api from '../services/api'

function InvoiceList() {
  const [invoices, setInvoices] = useState([])
  const [loading, setLoading] = useState(false)
  const [modalVisible, setModalVisible] = useState(false)
  const [editingInvoice, setEditingInvoice] = useState(null)
  const [form] = Form.useForm()

  const fetchInvoices = async () => {
    setLoading(true)
    try {
      const response = await api.get('/invoices')
      setInvoices(response.data || [])
    } catch (error) {
      message.error('获取发票列表失败')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchInvoices()
  }, [])

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields()
      if (editingInvoice) {
        await api.put(`/invoices/${editingInvoice.id}`, values)
        message.success('更新成功')
      } else {
        await api.post('/invoices', values)
        message.success('创建成功')
      }
      setModalVisible(false)
      fetchInvoices()
    } catch (error) {
      message.error('操作失败')
    }
  }

  const handleDelete = async (id) => {
    Modal.confirm({
      title: '确认删除',
      content: '确定要删除这张发票吗？',
      onOk: async () => {
        try {
          await api.delete(`/invoices/${id}`)
          message.success('删除成功')
          fetchInvoices()
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
      await api.post('/invoices/import', formData, { headers: { 'Content-Type': 'multipart/form-data' } })
      message.success('导入成功')
      fetchInvoices()
    } catch (error) {
      message.error('导入失败')
    }
    return false
  }

  const columns = [
    { title: '发票号', dataIndex: 'invoiceNo', key: 'invoiceNo' },
    { title: '金额', dataIndex: 'amount', key: 'amount', render: (v) => `¥${v?.toFixed(2)}` },
    { title: '开票日期', dataIndex: 'issueDate', key: 'issueDate' },
    {
      title: '类型',
      dataIndex: 'type',
      key: 'type',
      render: (type) => <Tag color={type === 'sales' ? 'blue' : 'green'}>{type === 'sales' ? '销售发票' : '采购发票'}</Tag>,
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status) => {
        const map = { draft: '草稿', issued: '已开具', paid: '已付款', cancelled: '已作废' }
        return <Tag>{map[status] || status}</Tag>
      },
    },
    {
      title: '操作',
      key: 'action',
      render: (_, record) => (
        <Space>
          <Button icon={<EditOutlined />} onClick={() => { setEditingInvoice(record); form.setFieldsValue(record); setModalVisible(true) }}>编辑</Button>
          <Button icon={<DeleteOutlined />} danger onClick={() => handleDelete(record.id)}>删除</Button>
        </Space>
      ),
    },
  ]

  return (
    <div>
      <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between' }}>
        <h2>发票管理</h2>
        <Space>
          <Upload showUploadList={false} beforeUpload={handleImport}>
            <Button icon={<UploadOutlined />}>导入</Button>
          </Upload>
          <Button icon={<DownloadOutlined />} onClick={() => message.info('导出功能开发中')}>导出</Button>
          <Button type="primary" icon={<PlusOutlined />} onClick={() => { setEditingInvoice(null); form.resetFields(); setModalVisible(true) }}>新建发票</Button>
        </Space>
      </div>
      <Table columns={columns} dataSource={invoices} loading={loading} rowKey="id" />
      <Modal title={editingInvoice ? '编辑发票' : '新建发票'} open={modalVisible} onOk={handleSubmit} onCancel={() => setModalVisible(false)}>
        <Form form={form} layout="vertical">
          <Form.Item name="invoiceNo" label="发票号" rules={[{ required: true }]}><Input /></Form.Item>
          <Form.Item name="amount" label="金额" rules={[{ required: true }]}><InputNumber min={0} step={0.01} style={{ width: '100%' }} /></Form.Item>
          <Form.Item name="issueDate" label="开票日期"><Input type="date" /></Form.Item>
          <Form.Item name="type" label="类型" initialValue="sales">
            <Select>
              <Select.Option value="sales">销售发票</Select.Option>
              <Select.Option value="purchase">采购发票</Select.Option>
            </Select>
          </Form.Item>
          <Form.Item name="status" label="状态" initialValue="draft">
            <Select>
              <Select.Option value="draft">草稿</Select.Option>
              <Select.Option value="issued">已开具</Select.Option>
              <Select.Option value="paid">已付款</Select.Option>
              <Select.Option value="cancelled">已作废</Select.Option>
            </Select>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}

export default InvoiceList
