<script setup lang="ts">
import { ref, onMounted } from "vue"
import { useRouter } from "vue-router"
import { categoryApi } from "@/services/category"
import { getArticleList } from "@/services/article"
import type { Category } from "@/types/category"

const router = useRouter()
const categories = ref<Category[]>([])
const articleCount = ref(0)

const navTo = (path: string) => {
  router.push(path)
}

onMounted(async () => {
  try {
    const [catRes, articleRes] = await Promise.all([
      categoryApi.getAllCategories(),
      getArticleList({ page: 1, pageSize: 1 })
    ])
    if (catRes.code === 200) categories.value = catRes.data || []
    if (articleRes.code === 200) articleCount.value = articleRes.data?.total || 0
  } catch { /* ignore */ }
})
</script>

<template>
  <footer class="bg-[#1a1a1a] border-t border-gray-800 mt-16">
    <div class="max-w-280 mx-auto px-8 py-12">
      <!-- 上部：多列信息 -->
      <div class="grid grid-cols-1 md:grid-cols-4 gap-8 mb-10">
        <!-- 关于 -->
        <div>
          <h4 class="text-white text-lg font-bold mb-4">🍉 吃瓜网</h4>
          <p class="text-gray-400 text-sm leading-relaxed">
            全网更新最快最全的吃瓜网站，网红黑料、明星绯闻、吃瓜视频、瓜榜TOP10实时更新。
          </p>
        </div>

        <!-- 快捷导航 -->
        <div>
          <h4 class="text-white text-base font-semibold mb-4">快捷导航</h4>
          <ul class="space-y-2">
            <li>
              <a class="text-gray-400 hover:text-primary text-sm cursor-pointer transition-colors" @click="navTo('/')">首页</a>
            </li>
            <li v-for="cat in categories.slice(0, 5)" :key="cat.id">
              <a class="text-gray-400 hover:text-primary text-sm cursor-pointer transition-colors" @click="navTo(`/category/${cat.id}`)">{{ cat.name }}</a>
            </li>
          </ul>
        </div>

        <!-- 友情链接 -->
        <div>
          <h4 class="text-white text-base font-semibold mb-4">友情链接</h4>
          <ul class="space-y-2">
            <li><a class="text-gray-400 hover:text-primary text-sm transition-colors" href="https://github.com" target="_blank">GitHub</a></li>
            <li><a class="text-gray-400 hover:text-primary text-sm transition-colors" href="https://vuejs.org" target="_blank">Vue.js</a></li>
            <li><a class="text-gray-400 hover:text-primary text-sm transition-colors" href="https://go.dev" target="_blank">Go</a></li>
          </ul>
        </div>

        <!-- 站点统计 -->
        <div>
          <h4 class="text-white text-base font-semibold mb-4">站点统计</h4>
          <div class="grid grid-cols-2 gap-4">
            <div class="bg-[#262626] rounded-lg p-3 text-center">
              <div class="text-primary text-xl font-bold">{{ articleCount }}</div>
              <div class="text-gray-500 text-xs mt-1">文章总数</div>
            </div>
            <div class="bg-[#262626] rounded-lg p-3 text-center">
              <div class="text-primary text-xl font-bold">{{ categories.length }}</div>
              <div class="text-gray-500 text-xs mt-1">分类数</div>
            </div>
          </div>
        </div>
      </div>

      <!-- 下部：版权 -->
      <div class="border-t border-gray-800 pt-6 flex flex-col md:flex-row items-center justify-between gap-4">
        <p class="text-gray-500 text-xs">
          © {{ new Date().getFullYear() }} 吃瓜网 · 版权所有 · 仅供学习交流使用
        </p>
        <p class="text-gray-600 text-xs">
          Powered by Vue 3 + Go + PostgreSQL
        </p>
      </div>
    </div>
  </footer>
</template>
