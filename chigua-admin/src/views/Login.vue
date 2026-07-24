<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { UserOutlined, LockOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
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

const form = reactive({
  username: '',
  password: ''
})

const handleLogin = async () => {
  if (!form.username.trim()) {
    message.warning('请输入用户名')
    return
  }
  if (!form.password) {
    message.warning('请输入密码')
    return
  }

  loading.value = true
  try {
    const response = await adminLogin(form)
    if (response.code === 200) {
      userStore.setToken(response.data.token)
      userStore.setUserInfo(response.data.user)
      message.success('登录成功')
      router.push('/')
    } else if (response.code === 403) {
      message.error('无管理员权限，请联系管理员')
    } else {
      message.error(response.msg || '登录失败，请重试')
    }
  } catch (error: any) {
    console.error('登录失败:', error)
    if (error.response) {
      if (error.response.data?.code === 403) {
        message.error('无管理员权限，请联系管理员')
      } else {
        message.error(error.response.data?.msg || '登录失败，请重试')
      }
    } else {
      message.error('网络连接失败，请稍后重试')
    }
  } finally {
    loading.value = false
  }
}

const onKeydown = (e: KeyboardEvent) => {
  if (e.key === 'Enter') {
    handleLogin()
  }
}
</script>

<template>
  <div class="login-container" @keydown="onKeydown">
    <div class="login-content">
      <div class="login-header">
        <h1>🍉 吃瓜网</h1>
        <p>后台管理系统</p>
      </div>
      <a-form :model="form" class="login-form" @submit.prevent="handleLogin">
        <a-form-item name="username" :rules="[{ required: true, message: '请输入用户名' }]">
          <a-input v-model:value="form.username" placeholder="用户名" size="large" autocomplete="username">
            <template #prefix><component :is="UserOutlined" /></template>
          </a-input>
        </a-form-item>
        <a-form-item name="password" :rules="[{ required: true, message: '请输入密码' }]">
          <a-input-password v-model:value="form.password" placeholder="密码" size="large" autocomplete="current-password">
            <template #prefix><component :is="LockOutlined" /></template>
          </a-input-password>
        </a-form-item>

        <a-form-item>
          <a-button type="primary" size="large" block html-type="submit" :loading="loading">
            登录
          </a-button>
        </a-form-item>
      </a-form>
    </div>
  </div>
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
