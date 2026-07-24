<script setup lang="ts">
import { ref, onMounted } from "vue"
import { useRouter } from "vue-router"
import { useAuthStore } from "@/stores/auth"
import { message } from "ant-design-vue"
import { categoryApi } from "@/services/category"
import type { Category } from "@/types/category"
import AuthModal from "./AuthModal.vue"

const router = useRouter()
const authStore = useAuthStore()
const showAuthModal = ref(false)
const showLogoutMenu = ref(false)
const categories = ref<Category[]>([])

const openAuthModal = () => {
  showAuthModal.value = true
}

const closeAuthModal = () => {
  showAuthModal.value = false
}

const handleProfile = () => {
  showLogoutMenu.value = false
  router.push("/profile")
}

const handleLogout = () => {
  authStore.logout()
  showLogoutMenu.value = false
  router.push("/")
  message.success("退出登录成功")
}

const handleCategoryClick = (categoryId: number) => {
  router.push(`/category/${categoryId}`)
}

onMounted(async () => {
  try {
    const res = await categoryApi.getAllCategories()
    if (res.code === 200 && res.data) {
      categories.value = res.data
    }
  } catch (error) {
    console.error("获取分类列表失败:", error)
  }
})
</script>

<template>
  <div
    class="fixed top-0 left-0 right-0 z-50 flex items-center justify-center bg-[#343232] shadow-md"
  >
    <nav class="backdrop-blur-sm w-full max-w-7xl">
      <div class="px-8">
        <div class="flex items-center justify-between h-18">
          <!-- 左侧：Logo -->
          <div class="flex items-center">
            <router-link to="/" class="flex items-center">
              <div
                class="text-2xl font-bold bg-linear-to-tr from-purple-600 via-red-500 to-white bg-clip-text text-transparent"
              >
                吃瓜网
              </div>
            </router-link>
          </div>

          <!-- 中间：分类列表 -->
          <div class="hidden md:flex items-center justify-center gap-6">
            <router-link
              to="/"
              class="text-white hover:text-primary font-medium"
              >首页</router-link
            >
            <a
              v-for="cat in categories"
              :key="cat.id"
              class="text-white hover:text-primary font-medium cursor-pointer"
              @click="handleCategoryClick(cat.id)"
            >
              {{ cat.name }}
            </a>
          </div>

          <!-- 右侧：搜索和登录按钮 -->
          <div class="flex items-center space-x-4">
            <!-- 搜索按钮 -->
            <button
              class="w-10 h-10 rounded-full cursor-pointer bg-[#2C2A2A] flex items-center justify-center hover:bg-gray-700 transition-colors"
            >
              <svg
                class="w-5 h-5 text-gray-300"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
                xmlns="http://www.w3.org/2000/svg"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
                ></path>
              </svg>
            </button>

            <!-- 登录按钮/用户头像 -->
            <div class="relative">
              <!-- 未登录时显示登录按钮 -->
              <button
                v-if="!authStore.isLoggedIn()"
                class="w-10 h-10 rounded-full cursor-pointer bg-[#2C2A2A] flex items-center justify-center hover:bg-gray-700 transition-colors"
                @click="openAuthModal"
              >
                <svg
                  class="w-5 h-5 text-gray-300"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                  xmlns="http://www.w3.org/2000/svg"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"
                  ></path>
                </svg>
              </button>

              <!-- 已登录时显示用户头像 -->
              <div
                v-else
                class="relative"
                @mouseenter="showLogoutMenu = true"
                @mouseleave="showLogoutMenu = false"
              >
                <button class="w-10 h-10 rounded-full overflow-hidden flex items-center justify-center transition-colors">
                  <img
                    v-if="authStore.user?.avatar"
                    :src="authStore.user.avatar"
                    class="w-full h-full object-cover"
                  />
                  <svg
                    v-else
                    class="w-6 h-6 text-white"
                    fill="currentColor"
                    viewBox="0 0 24 24"
                    xmlns="http://www.w3.org/2000/svg"
                  >
                    <path
                      fill-rule="evenodd"
                      d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"
                      clip-rule="evenodd"
                    />
                  </svg>
                </button>

                <!-- 下拉菜单 -->
                <div
                  v-show="showLogoutMenu"
                  class="absolute right-0 top-full w-40 bg-gray-800 rounded-lg shadow-xl z-50"
                >
                  <!-- 用户信息 -->
                  <div class="px-4 py-3 border-b border-gray-700">
                    <p class="text-white text-sm font-medium truncate">
                      {{ authStore.user?.nickname || authStore.user?.username }}
                    </p>
                    <p class="text-gray-500 text-xs truncate mt-0.5">
                      {{ authStore.user?.username }}
                    </p>
                  </div>
                  <div class="py-1">
                  <button
                    class="w-full px-4 py-2 text-left text-gray-300 hover:bg-gray-700 hover:text-white transition-colors"
                    @click="handleProfile"
                  >
                    个人中心
                  </button>
                  <button
                    class="w-full px-4 py-2 text-left text-gray-300 hover:bg-gray-700 hover:text-white transition-colors"
                    @click="handleLogout"
                  >
                    退出登录
                  </button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </nav>
  </div>

  <!-- 占位符，防止内容被固定导航栏遮挡 -->
  <div class="h-18"></div>

  <!-- 登录/注册弹窗 -->
  <AuthModal :visible="showAuthModal" @close="closeAuthModal" />
</template>

<style scoped>
/* 导航栏样式 */
</style>
