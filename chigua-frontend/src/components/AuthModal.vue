<script setup lang="ts">
import { ref, reactive } from "vue"
import { useAuthStore } from "@/stores/auth"

defineProps<{
  visible: boolean
}>()

const emit = defineEmits<{
  (e: "close"): void
}>()

const authStore = useAuthStore()

const activeTab = ref<'login' | 'register'>('login')
const loginForm = reactive({
  username: '',
  password: ''
})
const registerForm = reactive({
  username: '',
  password: '',
  email: ''
})
const loginError = ref('')
const registerError = ref('')
const isLoading = ref(false)

const handleLogin = async () => {
  if (!loginForm.username || !loginForm.password) {
    loginError.value = '请填写账号和密码'
    return
  }
  
  isLoading.value = true
  loginError.value = ''
  
  const success = await authStore.login(loginForm.username, loginForm.password)
  
  isLoading.value = false
  
  if (success) {
    emit('close')
    // 清空表单
    loginForm.username = ''
    loginForm.password = ''
  } else {
    loginError.value = '登录失败，请检查账号密码'
  }
}

const handleRegister = async () => {
  if (!registerForm.username || !registerForm.password || !registerForm.email) {
    registerError.value = '请填写完整信息'
    return
  }
  if (registerForm.username.length < 6) {
    registerError.value = '账号至少6位'
    return
  }
  
  isLoading.value = true
  registerError.value = ''
  
  const success = await authStore.register(registerForm.username, registerForm.password, registerForm.email)
  
  isLoading.value = false
  
  if (success) {
    // 注册成功后切换到登录页
    activeTab.value = 'login'
    // 清空表单
    registerForm.username = ''
    registerForm.password = ''
    registerForm.email = ''
  } else {
    registerError.value = '注册失败'
  }
}

const closeModal = () => {
  emit('close')
}
</script>

<template>
  <div
    v-if="visible"
    class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm"
    @click.self="closeModal"
  >
    <div class="w-full max-w-md mx-4 bg-gray-800 rounded-xl shadow-2xl overflow-hidden">
      <!-- 头部 -->
      <div class="text-center py-6">
        <div
          class="text-3xl font-bold bg-gradient-to-tr from-purple-600 via-red-500 to-white bg-clip-text text-transparent"
        >
          51吃瓜网
        </div>
      </div>

      <!-- Tab切换 -->
      <div class="flex border-b border-gray-700">
        <button
          class="flex-1 py-3 text-center font-medium transition-colors"
          :class="activeTab === 'login' ? 'text-white border-b-2 border-green-500' : 'text-gray-400 hover:text-white'"
          @click="activeTab = 'login'; loginError = ''"
        >
          登录
        </button>
        <button
          class="flex-1 py-3 text-center font-medium transition-colors"
          :class="activeTab === 'register' ? 'text-white border-b-2 border-green-500' : 'text-gray-400 hover:text-white'"
          @click="activeTab = 'register'; registerError = ''"
        >
          注册
        </button>
      </div>

      <!-- 表单区域 -->
      <div class="p-6">
        <!-- 登录表单 -->
        <div v-if="activeTab === 'login'">
          <div class="space-y-4">
            <div>
              <input
                v-model="loginForm.username"
                type="text"
                placeholder="请输入账号(最小6位数)或手机号"
                class="w-full px-4 py-3 bg-gray-700 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-green-500"
              />
            </div>
            <div>
              <input
                v-model="loginForm.password"
                type="password"
                placeholder="请输入登录密码"
                class="w-full px-4 py-3 bg-gray-700 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-green-500"
              />
            </div>
            <div v-if="loginError" class="text-red-400 text-sm">
              {{ loginError }}
            </div>
            <button
              class="w-full py-3 bg-green-500 hover:bg-green-600 text-white font-medium rounded-lg transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
              :disabled="isLoading"
              @click="handleLogin"
            >
              {{ isLoading ? '登录中...' : '登录' }}
            </button>
          </div>
        </div>

        <!-- 注册表单 -->
        <div v-else>
          <div class="space-y-4">
            <div>
              <input
                v-model="registerForm.username"
                type="text"
                placeholder="请输入账号(最小6位数)"
                class="w-full px-4 py-3 bg-gray-700 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-green-500"
              />
            </div>
            <div>
              <input
                v-model="registerForm.email"
                type="email"
                placeholder="请输入邮箱"
                class="w-full px-4 py-3 bg-gray-700 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-green-500"
              />
            </div>
            <div>
              <input
                v-model="registerForm.password"
                type="password"
                placeholder="请输入密码"
                class="w-full px-4 py-3 bg-gray-700 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-green-500"
              />
            </div>
            <div v-if="registerError" class="text-red-400 text-sm">
              {{ registerError }}
            </div>
            <button
              class="w-full py-3 bg-green-500 hover:bg-green-600 text-white font-medium rounded-lg transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
              :disabled="isLoading"
              @click="handleRegister"
            >
              {{ isLoading ? '注册中...' : '注册' }}
            </button>
          </div>
        </div>

        <!-- 温馨提示 -->
        <div class="mt-6 pt-4 border-t border-gray-700">
          <div class="text-gray-400 text-sm font-medium mb-2">温馨提示</div>
          <ul class="text-gray-500 text-xs space-y-1">
            <li>1.请务必记住自己的账号密码;</li>
            <li>2.不提供密码修改;</li>
            <li>3.不提供密码找回功能;</li>
          </ul>
        </div>
      </div>

      <!-- 关闭按钮 -->
      <button
        class="absolute top-4 right-4 w-8 h-8 rounded-full bg-gray-700 hover:bg-gray-600 flex items-center justify-center transition-colors"
        @click="closeModal"
      >
        <svg
          class="w-4 h-4 text-gray-300"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M6 18L18 6M6 6l12 12"
          ></path>
        </svg>
      </button>
    </div>
  </div>
</template>
