<script setup lang="ts">
import { ref, onMounted, computed } from "vue"
import { useRoute } from "vue-router"
import { articleApi } from "@/services"
import type { Article } from "@/types/article"
import { renderMarkdown } from "@/utils/markdown"
import { formatDateTimeFull } from "@/utils/dateFormat"

const route = useRoute()
const article = ref<Article | null>(null)
const loading = ref(true)
const error = ref("")

const renderedContent = computed(() => {
  if (!article.value?.content) return ""
  return renderMarkdown(article.value.content)
})

const displayAuthor = computed(() => {
  const author = article.value?.author
  if (!author) return { name: "未知", initial: "?" }
  return {
    name: author.nickname || author.username,
    initial: (author.nickname || author.username).charAt(0).toUpperCase()
  }
})

onMounted(async () => {
  const id = Number(route.params.id)
  if (isNaN(id)) {
    error.value = "无效的文章ID"
    loading.value = false
    return
  }

  try {
    const response = await articleApi.getArticle(id)
    if (response.code === 200) {
      article.value = response.data
    } else {
      error.value = response.msg || "获取文章失败"
    }
  } catch (err) {
    error.value = "获取文章失败"
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="bg-[#262626] text-white min-h-[calc(100vh-72px)] pb-8">
    <!-- 文章详情容器 -->
    <div class="max-w-4xl mx-auto px-4 py-8">
      <!-- 加载状态 -->
      <div v-if="loading" class="flex justify-center items-center h-64">
        <div
          class="animate-spin rounded-full h-12 w-12 border-b-2 border-purple-500"
        ></div>
      </div>

      <!-- 错误状态 -->
      <div v-else-if="error" class="text-center py-12">
        <p class="text-red-400 text-lg">{{ error }}</p>
      </div>

      <!-- 文章内容 -->
      <article v-else-if="article">
        <!-- 顶部区域：标题 + 作者信息 + 面包屑 -->
        <header class="border-b border-gray-700/80 pb-6 mb-8">
          <!-- 标题 -->
          <h1
            class="text-2xl md:text-3xl font-bold text-white leading-snug mb-5"
          >
            {{ article.title }}
          </h1>

          <!-- 作者信息行 -->
          <div class="flex items-center justify-between mb-3">
            <div class="flex items-center gap-3">
              <!-- 头像 -->
              <img
                v-if="article.author?.avatar"
                :src="article.author.avatar"
                class="w-9 h-9 rounded-full object-cover shrink-0"
                @error="$event.target.style.display = 'none'"
              />
              <div
                v-else
                class="w-9 h-9 rounded-full flex items-center justify-center shrink-0"
                style="background-color: #1e8d77"
              >
                <span class="text-sm font-bold text-white">{{ displayAuthor.initial }}</span>
              </div>
              <!-- 名称和时间 -->
              <div class="flex items-center gap-2 text-sm">
                <span class="text-gray-300 font-medium">{{ displayAuthor.name }}</span>
                <span class="text-gray-600">·</span>
                <span class="text-gray-500">{{ formatDateTimeFull(article.publishAt || article.createdAt) }}</span>
              </div>
            </div>

            <!-- 右侧：来源等信息 -->
            <div class="flex items-center gap-2 text-sm text-gray-500">
              <span>来源：{{ article.author.username }}</span>
            </div>
          </div>

          <!-- 面包屑 -->
          <nav class="flex items-center gap-2 text-sm">
            <router-link to="/" class="hover:underline transition-colors" style="color: #1e8d77">
              首页
            </router-link>
            <span class="text-gray-500">&gt;</span>
            <router-link
              v-if="article.category"
              :to="`/category/${article.category.id}`"
              class="hover:underline transition-colors"
              style="color: #1e8d77"
            >
              {{ article.category.name }}
            </router-link>
            <span v-else style="color: #1e8d77">资讯</span>
          </nav>
        </header>

        <!-- 文章内容 -->
        <div
          class="prose prose-invert max-w-none bg-white/5 backdrop-blur-sm rounded-xl p-8 markdown-body border border-white/5"
          v-html="renderedContent"
        ></div>

        <!-- 文章标签 -->
        <div class="flex flex-wrap gap-3 mt-8 pt-4 border-t border-gray-700">
          <span class="text-gray-400">标签：</span>
          <span
            v-for="tag in article.tags"
            :key="tag.id"
            class="px-3 py-1 bg-white/10 backdrop-blur-sm text-gray-300 rounded-full text-sm hover:bg-white/20 transition-colors cursor-pointer"
          >
            {{ tag.name }}
          </span>
          <span v-if="!article.tags || article.tags.length === 0" class="text-gray-500 text-sm"
            >暂无标签</span
          >
        </div>
      </article>
    </div>
  </div>
</template>
