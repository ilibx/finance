<script setup>
import { ref, onMounted, h } from 'vue'
import { message, Modal } from 'ant-design-vue'
import {
  UserOutlined,
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  ReloadOutlined
} from '@ant-design/icons-vue'
import { graphQLClient } from '../graphql/client'
import { GET_USERS, CREATE_USER, UPDATE_USER, DELETE_USER } from '../graphql/queries'

// State
const users = ref([])
const loading = ref(false)
const visible = ref(false)
const editingUser = ref(null)
const formModel = ref({
  name: '',
  email: ''
})

// Table columns
const columns = [
  {
    title: 'ID',
    dataIndex: 'id',
    key: 'id',
    width: 80
  },
  {
    title: 'Name',
    dataIndex: 'name',
    key: 'name'
  },
  {
    title: 'Email',
    dataIndex: 'email',
    key: 'email'
  },
  {
    title: 'Created At',
    dataIndex: 'createdAt',
    key: 'createdAt',
    width: 180,
    customRender: ({ text }) => {
      return text ? new Date(text).toLocaleString() : '-'
    }
  },
  {
    title: 'Actions',
    key: 'action',
    width: 150
  }
]

// Fetch users
const fetchUsers = async () => {
  loading.value = true
  try {
    const data = await graphQLClient.request(GET_USERS)
    users.value = data.users || []
  } catch (error) {
    console.error('Error fetching users:', error)
    message.error('Failed to load users: ' + (error.response?.errors?.[0]?.message || error.message))
    // For demo purposes, show sample data if backend is not available
    users.value = [
      { id: '1', name: 'John Doe', email: 'john@example.com', createdAt: new Date().toISOString() },
      { id: '2', name: 'Jane Smith', email: 'jane@example.com', createdAt: new Date().toISOString() },
      { id: '3', name: 'Bob Wilson', email: 'bob@example.com', createdAt: new Date().toISOString() }
    ]
    message.info('Showing demo data (backend not connected)')
  } finally {
    loading.value = false
  }
}

// Open modal for creating user
const showCreateModal = () => {
  editingUser.value = null
  formModel.value = { name: '', email: '' }
  visible.value = true
}

// Open modal for editing user
const showEditModal = (user) => {
  editingUser.value = user
  formModel.value = { ...user }
  visible.value = true
}

// Handle form submission
const handleSubmit = async () => {
  try {
    if (editingUser.value) {
      // Update existing user
      const variables = {
        id: editingUser.value.id,
        input: {
          name: formModel.value.name,
          email: formModel.value.email
        }
      }
      await graphQLClient.request(UPDATE_USER, variables)
      message.success('User updated successfully')
    } else {
      // Create new user
      const variables = {
        input: {
          name: formModel.value.name,
          email: formModel.value.email
        }
      }
      await graphQLClient.request(CREATE_USER, variables)
      message.success('User created successfully')
    }
    visible.value = false
    fetchUsers()
  } catch (error) {
    console.error('Error saving user:', error)
    message.error('Failed to save user: ' + (error.response?.errors?.[0]?.message || error.message))
  }
}

// Handle delete user
const handleDelete = (user) => {
  Modal.confirm({
    title: 'Delete User',
    content: `Are you sure you want to delete ${user.name}?`,
    okText: 'Delete',
    okType: 'danger',
    cancelText: 'Cancel',
    onOk: async () => {
      try {
        const variables = { id: user.id }
        await graphQLClient.request(DELETE_USER, variables)
        message.success('User deleted successfully')
        fetchUsers()
      } catch (error) {
        console.error('Error deleting user:', error)
        message.error('Failed to delete user: ' + (error.response?.errors?.[0]?.message || error.message))
      }
    }
  })
}

// Handle modal cancel
const handleCancel = () => {
  visible.value = false
  editingUser.value = null
  formModel.value = { name: '', email: '' }
}

// Lifecycle
onMounted(() => {
  fetchUsers()
})
</script>

<template>
  <div class="user-management">
    <a-card title="User Management" :bordered="false">
      <template #extra>
        <a-space>
          <a-button @click="fetchUsers" :loading="loading">
            <template #icon><ReloadOutlined /></template>
            Refresh
          </a-button>
          <a-button type="primary" @click="showCreateModal">
            <template #icon><PlusOutlined /></template>
            Add User
          </a-button>
        </a-space>
      </template>

      <a-table 
        :columns="columns" 
        :data-source="users" 
        :loading="loading"
        row-key="id"
        :pagination="{ pageSize: 10 }"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'action'">
            <a-space>
              <a-button type="link" size="small" @click="showEditModal(record)">
                <template #icon><EditOutlined /></template>
                Edit
              </a-button>
              <a-popconfirm
                title="Delete this user?"
                ok-text="Yes"
                cancel-text="No"
                @confirm="handleDelete(record)"
              >
                <a-button type="link" danger size="small">
                  <template #icon><DeleteOutlined /></template>
                  Delete
                </a-button>
              </a-popconfirm>
            </a-space>
          </template>
        </template>
      </a-table>
    </a-card>

    <!-- Create/Edit Modal -->
    <a-modal
      v-model:open="visible"
      :title="editingUser ? 'Edit User' : 'Create User'"
      @ok="handleSubmit"
      @cancel="handleCancel"
      :confirm-loading="loading"
    >
      <a-form :model="formModel" layout="vertical">
        <a-form-item
          label="Name"
          name="name"
          :rules="[{ required: true, message: 'Please enter name' }]"
        >
          <a-input 
            v-model:value="formModel.name" 
            placeholder="Enter user name"
            :prefix="h(UserOutlined)"
          >
            <template #prefix>
              <UserOutlined />
            </template>
          </a-input>
        </a-form-item>
        
        <a-form-item
          label="Email"
          name="email"
          :rules="[
            { required: true, message: 'Please enter email' },
            { type: 'email', message: 'Please enter valid email' }
          ]"
        >
          <a-input 
            v-model:value="formModel.email" 
            type="email"
            placeholder="Enter email address"
          />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<style scoped>
.user-management {
  max-width: 1200px;
  margin: 0 auto;
  padding: 24px;
}

:deep(.ant-card) {
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

:deep(.ant-table) {
  font-size: 14px;
}
</style>
