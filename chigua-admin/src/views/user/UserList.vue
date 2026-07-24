<script setup lang="ts">
import { ref, reactive } from 'vue'
import { PlusOutlined, CameraOutlined } from '@ant-design/icons-vue'
import { Modal, message } from 'ant-design-vue'
import { getUserList, deleteUser, createUser, updateUser, getCurrentUser } from '@/api/user'
import { getRoles, updateUserRoles } from '@/api/rbac'
import { uploadFile } from '@/api/upload'
import { formatDate } from '@/utils/date'
import { createPagination, zebraRow, emptyText } from '@/utils/table'
import { useUserStore } from '@/stores/user'
import AvatarCropper from '@/components/AvatarCropper.vue'
import type { UserResponse, Role } from '@/types'

const userStore = useUserStore()
const searchText = ref('')
const loading = ref(false)
const users = ref<UserResponse[]>([])

const pagination = createPagination()

const roleLabelMap: Record<string, { label: string; color: string }> = {
  admin: { label: '管理员', color: 'red' },
  reviewer: { label: '审核员', color: 'orange' },
  user: { label: '普通用户', color: 'default' },
}

const columns = [
  { title: '序号', dataIndex: 'index', key: 'index', width: 80, align: 'center' as const },
  { title: '用户名', dataIndex: 'username', key: 'username', width: 120 },
  { title: '昵称', dataIndex: 'nickname', key: 'nickname', width: 120 },
  { title: '头像', dataIndex: 'avatar', key: 'avatar', width: 60, align: 'center' as const },
  { title: '创建时间', dataIndex: 'createdAt', key: 'createdAt', width: 180, align: 'center' as const },
  { title: '操作', key: 'actions', width: 180, align: 'center' as const }
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
        message.error('删除失败，请稍后重试')
      }
    }
  })
}

const formVisible = ref(false)
const isEdit = ref(false)
const editingId = ref(0)
const form = reactive({
  username: '',
  password: '',
  nickname: '',
  avatar: '',
})

const resetForm = () => {
  form.username = ''
  form.password = ''
  form.nickname = ''
  form.avatar = ''
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
  form.avatar = record.avatar || ''
  formVisible.value = true
}

// 头像裁剪
const avatarModalVisible = ref(false)
const avatarImageUrl = ref('')
const cropperRef = ref<any>()
const avatarUploading = ref(false)

const handleAvatarSelect = () => {
  const input = document.createElement('input')
  input.type = 'file'
  input.accept = 'image/*'
  input.onchange = (e: Event) => {
    const file = (e.target as HTMLInputElement).files?.[0]
    if (!file) return
    avatarImageUrl.value = URL.createObjectURL(file)
    avatarModalVisible.value = true
  }
  input.click()
}

const handleAvatarCrop = () => {
  cropperRef.value?.handleCrop()
}

const onAvatarCropped = async (blob: Blob) => {
  avatarUploading.value = true
  try {
    const file = new File([blob], 'avatar.jpg', { type: 'image/jpeg' })
    const res = await uploadFile(file)
    if (res.code === 200) {
      form.avatar = res.data.fileUrl
      avatarModalVisible.value = false
      message.success('头像已上传')
    } else {
      message.error('上传失败')
    }
  } catch { message.error('上传失败') }
  finally { avatarUploading.value = false }
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

  try {
    if (isEdit.value) {
      const response = await updateUser(editingId.value, { nickname: form.nickname, avatar: form.avatar })
      if (response.code === 200) {
        message.success('编辑成功')
        formVisible.value = false
        // 如果编辑的是当前登录用户，刷新本地 store
        if (userStore.userInfo && editingId.value === userStore.userInfo.id) {
          getCurrentUser().then(res => { if (res.code === 200) userStore.setUserInfo(res.data) })
        }
        await fetchData()
      } else {
        message.error(response.msg || '保存失败')
      }
    } else {
      const response = await createUser({ username: form.username, password: form.password, nickname: form.nickname, avatar: form.avatar })
      if (response.code === 200) {
        message.success('新增用户成功')
        formVisible.value = false
        await fetchData()
      } else {
        message.error(response.msg || '保存失败')
      }
    }
  } catch {
    message.error('保存失败，请稍后重试')
  }
}

// ====== 角色分配抽屉 ======
const roleDrawerVisible = ref(false)
const roleUser = ref<UserResponse | null>(null)
const allRoles = ref<Role[]>([])
const selectedRoleIds = ref<number[]>([])
const roleSaving = ref(false)

const openRoleDrawer = async (record: UserResponse) => {
  roleUser.value = record
  roleDrawerVisible.value = true
  try {
    const res = await getRoles({ page: 1, pageSize: 100 })
    if (res.code === 200) {
      allRoles.value = res.data.items
      selectedRoleIds.value = res.data.items
        .filter(r => record.roles?.includes(r.code))
        .map(r => r.id)
    }
  } catch {
    message.error('加载角色列表失败')
  }
}

const handleRoleSave = async () => {
  if (!roleUser.value) return
  roleSaving.value = true
  try {
    const res = await updateUserRoles(roleUser.value.id, selectedRoleIds.value)
    if (res.code === 200) {
      message.success('角色已更新')
      roleDrawerVisible.value = false
      fetchData()
    } else {
      message.error(res.msg || '保存失败')
    }
  } catch {
    message.error('保存失败')
  } finally {
    roleSaving.value = false
  }
}

fetchData()
</script>

<template>
  <div>
    <a-card style="margin-bottom: 16px;">
      <div style="display: flex; justify-content: space-between; align-items: center; flex-wrap: wrap; gap: 12px;">
        <a-input-search v-model:value="searchText" placeholder="搜索用户名或昵称..." allow-clear style="width: 240px;" @search="handleSearch" />
        <a-button type="primary" @click="openCreateDialog" v-if="userStore.hasPermission('user:create')">
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
        <template v-if="column.key === 'avatar'">
          <a-avatar :size="32" :src="record.avatar || undefined" style="background: #f0f0f0;">
            {{ record.nickname?.charAt(0) || record.username?.charAt(0) || '?' }}
          </a-avatar>
        </template>
        <template v-if="column.key === 'roles'">
          <a-space :size="4">
            <a-tag v-for="r in (record.roles?.length ? record.roles : ['user'])" :key="r" :color="roleLabelMap[r]?.color || 'default'">
              {{ roleLabelMap[r]?.label || r }}
            </a-tag>
          </a-space>
        </template>
        <template v-if="column.key === 'createdAt'">
          {{ formatDate(record.createdAt) }}
        </template>
        <template v-if="column.key === 'actions'">
          <a-space>
            <a-button size="small" @click="openRoleDrawer(record)" v-if="userStore.hasPermission('user:update')">角色</a-button>
            <a-button type="primary" size="small" @click="openEditDialog(record)" v-if="userStore.hasPermission('user:update')">编辑</a-button>
            <a-button size="small" danger @click="handleDelete(record.id)" v-if="userStore.hasPermission('user:delete')">删除</a-button>
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
        <a-form-item label="用户名">
          <a-input v-model:value="form.username" :disabled="isEdit" />
        </a-form-item>
        <a-form-item label="密码" :rules="[{ required: !isEdit, message: '请输入密码' }]">
          <a-input-password v-model:value="form.password" :placeholder="isEdit ? '不填则保持不变' : '请输入密码'" />
        </a-form-item>
        <a-form-item label="昵称">
          <a-input v-model:value="form.nickname" placeholder="请输入昵称" />
        </a-form-item>
        <a-form-item label="头像">
          <div style="display: flex; align-items: center; gap: 12px;">
            <a-avatar :size="64" :src="form.avatar || undefined" style="background: #f0f0f0;">
              <template v-if="!form.avatar"><CameraOutlined style="font-size: 24px; color: #bbb;" /></template>
            </a-avatar>
            <a-button @click="handleAvatarSelect">上传头像</a-button>
          </div>
        </a-form-item>
        <a-form-item>
          <a-space>
            <a-button type="primary" @click="handleSubmit" :loading="loading">保存</a-button>
            <a-button @click="formVisible = false">取消</a-button>
          </a-space>
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 头像裁剪弹窗 -->
    <a-modal v-model:open="avatarModalVisible" title="裁剪头像" @ok="handleAvatarCrop" :confirm-loading="avatarUploading" width="440px">
      <AvatarCropper ref="cropperRef" :image-url="avatarImageUrl" @cropped="onAvatarCropped" />
    </a-modal>

    <!-- 角色分配抽屉 -->
    <a-drawer
      v-model:open="roleDrawerVisible"
      :title="`角色配置 - ${roleUser?.nickname || roleUser?.username || ''}`"
      :width="400"
      :footer-style="{ textAlign: 'right' }"
    >
      <template #footer>
        <a-button style="margin-right: 8px;" @click="roleDrawerVisible = false">取消</a-button>
        <a-button type="primary" :loading="roleSaving" @click="handleRoleSave" v-if="userStore.hasPermission('user:update')">保存</a-button>
      </template>

      <a-checkbox-group v-model:value="selectedRoleIds" style="display: flex; flex-direction: column; gap: 12px;">
        <a-checkbox v-for="role in allRoles" :key="role.id" :value="role.id">
          <span style="font-weight: 500;">{{ role.name }}</span>
          <span style="color: #999; margin-left: 8px; font-size: 12px;">{{ role.description }}</span>
        </a-checkbox>
      </a-checkbox-group>
    </a-drawer>
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
