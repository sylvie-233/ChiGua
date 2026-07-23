<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { UserOutlined, LockOutlined } from '@ant-design/icons-vue'
import request from '@/api/request'
import type { BaseResponse, UserResponse } from '@/types'
import { useUserStore } from '@/stores/user'

interface LoginParams {
  username: string
  password: string
}

interface LoginResponse {
  token: string
  user: UserResponse
}

const adminLogin = async (data: LoginParams): Promise<BaseResponse<LoginResponse>> => {
  const response = await request.post('/admin/login', data)
  return response.data
}

const router = useRouter()
const userStore = useUserStore()
const loading = ref(false)
const modalVisible = ref(false)
const modalTitle = ref('')
const modalContent = ref('')

const form = reactive({
  username: '',
  password: ''
})

const showModal = (title: string, content: string) => {
  modalTitle.value = title
  modalContent.value = content
  modalVisible.value = true
}

const handleLogin = async () => {
  if (!form.username || !form.password) {
    showModal('提示', '请输入用户名和密码')
    return
  }

  loading.value = true
  try {
    const response = await adminLogin(form)
    if (response.code === 200) {
      userStore.setToken(response.data.token)
      userStore.setUserInfo(response.data.user)
      router.push('/')
    } else if (response.code === 403) {
      showModal('登录失败', '无管理员权限')
    } else {
      showModal('登录失败', response.msg || '登录失败')
    }
  } catch (error: any) {
    console.error('登录失败:', error)
    if (error.response) {
      if (error.response.data?.code === 403) {
        showModal('登录失败', '无管理员权限')
      } else {
        showModal('登录失败', error.response.data?.msg || '登录失败')
      }
    } else {
      showModal('网络错误', '网络连接失败，请稍后重试')
    }
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="login-container">
    <div class="login-content">
      <div class="login-header">
        <h1>🍉 吃瓜网</h1>
        <p>后台管理系统</p>
      </div>
      <a-form :model="form" class="login-form">
        <a-form-item name="username" :rules="[{ required: true, message: '请输入用户名' }]">
          <a-input v-model:value="form.username" placeholder="用户名" size="large">
            <template #prefix><component :is="UserOutlined" /></template>
          </a-input>
        </a-form-item>
        <a-form-item name="password" :rules="[{ required: true, message: '请输入密码' }]">
          <a-input-password v-model:value="form.password" placeholder="密码" size="large">
            <template #prefix><component :is="LockOutlined" /></template>
          </a-input-password>
        </a-form-item>

        <a-form-item>
          <a-button type="primary" size="large" block @click="handleLogin" :loading="loading">
            登录
          </a-button>
        </a-form-item>
      </a-form>
    </div>
  </div>

  <a-modal
    v-model:visible="modalVisible"
    :title="modalTitle"
    :footer="null"
    centered
  >
    <p>{{ modalContent }}</p>
  </a-modal>
</template>

<style scoped>
.login-container {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.login-content {
  width: 400px;
  background: #fff;
  border-radius: 8px;
  padding: 40px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
}

.login-header {
  text-align: center;
  margin-bottom: 32px;
}

.login-header h1 {
  font-size: 32px;
  margin-bottom: 8px;
}

.login-header p {
  color: #999;
  font-size: 14px;
}

.login-form {
  margin-top: 16px;
}
</style>
