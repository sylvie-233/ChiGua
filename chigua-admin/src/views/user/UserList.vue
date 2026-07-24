<script setup lang="ts">
import { ref, reactive } from 'vue'
import { PlusOutlined } from '@ant-design/icons-vue'
import { Modal, message } from 'ant-design-vue'
import { getUserList, deleteUser, updateUserRole, createUser, updateUser } from '@/api/user'
import { formatDate } from '@/utils/date'
import { createPagination, zebraRow, emptyText } from '@/utils/table'
import type { UserCreate, UserUpdate } from '@/api/user'
import type { UserResponse } from '@/types'

const searchText = ref('')
const loading = ref(false)
const users = ref<UserResponse[]>([])

const pagination = createPagination()

const roleMap: Record<string, { label: string; color: string }> = {
  admin: { label: '管理员', color: 'red' },
  user: { label: '普通用户', color: 'default' }
}

const columns = [
  { title: '序号', dataIndex: 'index', key: 'index', width: 80, align: 'center' as const },
  { title: '用户名', dataIndex: 'username', key: 'username', width: 150 },
  { title: '昵称', dataIndex: 'nickname', key: 'nickname', width: 120 },
  { title: '角色', dataIndex: 'role', key: 'role', width: 100, align: 'center' as const },
  { title: '创建时间', dataIndex: 'createdAt', key: 'createdAt', width: 180, align: 'center' as const },
  { title: '操作', key: 'actions', width: 200, align: 'center' as const }
]

const fetchData = async () => {
  loading.value = true
  try {
    const params: { page: number; pageSize: number; keyword?: string } = { page: pagination.current, pageSize: pagination.pageSize }
    if (searchText.value.trim()) {
      params.keyword = searchText.value.trim()
    }
    const response = await getUserList(params)
    if (response.code === 200) {
      users.value = response.data.items
      pagination.total = response.data.total
    }
  } catch (error) {
    console.error('获取用户列表失败:', error)
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.current = 1
  fetchData()
}

const handleTableChange = (paginationInfo: { current: number; pageSize: number }) => {
  pagination.current = paginationInfo.current
  pagination.pageSize = paginationInfo.pageSize
  fetchData()
}

const handleDelete = async (id: number) => {
  Modal.confirm({
    title: '确认删除',
    content: '确定要删除这个用户吗？此操作不可撤销。',
    okText: '删除',
    okButtonProps: { danger: true },
    cancelText: '取消',
    async onOk() {
      try {
        const response = await deleteUser(id)
        if (response.code === 200) {
          message.success('删除成功')
          fetchData()
        } else {
          message.error(response.msg || '删除失败')
        }
      } catch (error) {
        console.error('删除用户失败:', error)
        message.error('删除失败，请稍后重试')
      }
    }
  })
}

const handleRoleChange = async (id: number, currentRole: string) => {
  const newRole = currentRole === 'admin' ? 'user' : 'admin'
  Modal.confirm({
    title: '确认修改角色',
    content: `确定要将用户角色改为${roleMap[newRole].label}吗？`,
    okText: '确定',
    okType: 'primary',
    cancelText: '取消',
    async onOk() {
      try {
        const response = await updateUserRole(id, newRole)
        if (response.code === 200) {
          message.success('角色修改成功')
          fetchData()
        } else {
          message.error(response.msg || '角色修改失败')
        }
      } catch (error) {
        console.error('修改用户角色失败:', error)
        message.error('角色修改失败，请稍后重试')
      }
    }
  })
}

const formVisible = ref(false)
const isEdit = ref(false)
const editingId = ref(0)
const form = reactive<UserCreate>({
  username: '',
  password: '',
  nickname: '',
  role: 'user'
})

const resetForm = () => {
  form.username = ''
  form.password = ''
  form.nickname = ''
  form.role = 'user'
}

const openCreateDialog = () => {
  isEdit.value = false
  editingId.value = 0
  resetForm()
  formVisible.value = true
}

const openEditDialog = (record: UserResponse) => {
  isEdit.value = true
  editingId.value = record.id
  form.username = record.username
  form.password = ''
  form.nickname = record.nickname
  form.role = record.role
  formVisible.value = true
}

const handleSubmit = async () => {
  if (!form.username.trim()) {
    message.warning('请输入用户名')
    return
  }
  if (!isEdit.value && !form.password) {
    message.warning('请输入密码')
    return
  }

  loading.value = true
  try {
    if (isEdit.value) {
      const updateData: UserUpdate = {
        nickname: form.nickname,
        role: form.role
      }
      const response = await updateUser(editingId.value, updateData)
      if (response.code === 200) {
        message.success('编辑成功')
        formVisible.value = false
        fetchData()
      } else {
        message.error(response.msg || '保存失败')
      }
    } else {
      const response = await createUser(form)
      if (response.code === 200) {
        message.success('新增用户成功')
        formVisible.value = false
        fetchData()
      } else {
        message.error(response.msg || '保存失败')
      }
    }
  } catch (error: any) {
    console.error('保存用户失败:', error)
    if (error.response?.data?.code === 400 && error.response?.data?.msg === '用户已存在') {
      message.error('用户名已存在')
    } else {
      message.error('保存失败，请稍后重试')
    }
  } finally {
    loading.value = false
  }
}

fetchData()
</script>

<template>
  <div>
    <a-card style="margin-bottom: 16px;">
      <div style="display: flex; justify-content: space-between; align-items: center; flex-wrap: wrap; gap: 12px;">
        <a-input-search v-model:value="searchText" placeholder="搜索用户名或昵称..." allow-clear style="width: 240px;" @search="handleSearch" />
        <a-button type="primary" @click="openCreateDialog">
          <PlusOutlined /> 新增用户
        </a-button>
      </div>
    </a-card>

    <a-card :body-style="{ padding: 0 }">
      <a-table
        bordered
        size="middle"
        :columns="columns"
        :data-source="users"
        :pagination="pagination"
        :loading="loading"
        :row-class-name="zebraRow"
        :locale="{ emptyText }"
        row-key="id"
        @change="handleTableChange"
      >
      <template #bodyCell="{ column, record, index }">
        <template v-if="column.key === 'index'">
          {{ (pagination.current - 1) * pagination.pageSize + index + 1 }}
        </template>
        <template v-if="column.key === 'role'">
          <a-tag :color="roleMap[record.role]?.color" style="cursor: pointer;" @click="handleRoleChange(record.id, record.role)">
            {{ roleMap[record.role]?.label }}
          </a-tag>
        </template>
        <template v-if="column.key === 'createdAt'">
          {{ formatDate(record.createdAt) }}
        </template>
        <template v-if="column.key === 'actions'">
          <a-space>
            <a-button type="primary" size="small" @click="openEditDialog(record)">编辑</a-button>
            <a-button size="small" danger @click="handleDelete(record.id)">删除</a-button>
          </a-space>
        </template>
      </template>
    </a-table>
    </a-card>

    <a-modal
      v-model:visible="formVisible"
      :title="isEdit ? '编辑用户' : '新增用户'"
      :footer="null"
      width="500px"
    >
      <a-form :model="form" layout="vertical">
        <a-form-item label="用户名" :rules="[{ required: true, message: '请输入用户名' }]">
          <a-input v-model:value="form.username" :disabled="isEdit" />
        </a-form-item>
        <a-form-item label="密码" :rules="[{ required: !isEdit, message: '请输入密码' }]">
          <a-input-password v-model:value="form.password" :placeholder="isEdit ? '不填则保持不变' : '请输入密码'" />
        </a-form-item>
        <a-form-item label="昵称">
          <a-input v-model:value="form.nickname" placeholder="请输入昵称" />
        </a-form-item>
        <a-form-item label="角色">
          <a-select v-model:value="form.role" style="width: 100%">
            <a-select-option value="user">普通用户</a-select-option>
            <a-select-option value="admin">管理员</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item>
          <a-space>
            <a-button type="primary" @click="handleSubmit" :loading="loading">保存</a-button>
            <a-button @click="formVisible = false">取消</a-button>
          </a-space>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<style scoped>
:deep(.ant-table-thead > tr > th) {
  background: #fafafa;
  font-weight: 600;
}
.row-striped > td {
  background: #fafafa;
}
</style>