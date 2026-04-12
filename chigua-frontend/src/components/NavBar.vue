<script setup lang="ts">
// 导航栏组件
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

          <!-- 中间：分类列表（仅在中等及以上屏幕显示） -->
          <div class="hidden md:flex items-center justify-center gap-6">
            <router-link
              to="/"
              class="text-white hover:text-primary font-medium"
              >首页</router-link
            >
            <router-link
              to="/category/1"
              class="text-white hover:text-primary font-medium"
              >前端技术</router-link
            >
            <router-link
              to="/category/2"
              class="text-white hover:text-primary font-medium"
              >后端技术</router-link
            >
            <router-link
              to="/category/3"
              class="text-white hover:text-primary font-medium"
              >移动开发</router-link
            >
            <router-link
              to="/category/4"
              class="text-white hover:text-primary font-medium"
              >人工智能</router-link
            >
            <router-link
              to="/category/5"
              class="text-white hover:text-primary font-medium"
              >区块链</router-link
            >
            <router-link
              to="/about"
              class="text-white hover:text-primary font-medium"
              >关于我们</router-link
            >
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

            <!-- 登录按钮 -->
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
                  d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"
                ></path>
              </svg>
            </button>
          </div>
        </div>
      </div>
    </nav>
  </div>

  <!-- 占位符，防止内容被固定导航栏遮挡 -->
  <div class="h-18"></div>
</template>

<style scoped>
/* 导航栏样式 */
</style>
