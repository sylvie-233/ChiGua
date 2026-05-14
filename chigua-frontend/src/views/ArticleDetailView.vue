<script setup lang="ts">
import { ref, onMounted, computed } from "vue"
import { useRoute } from "vue-router"
import { articleApi } from "@/services"
import type { Article } from "@/types/article"
import { renderMarkdown } from "@/utils/markdown"

const route = useRoute()
const article = ref<Article | null>(null)
const loading = ref(true)
const error = ref("")

const renderedContent = computed(() => {
  if (!article.value?.content) return ""
  return renderMarkdown(article.value.content)
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
  <div class="min-h-screen bg-gray-900 text-white">
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
      <article v-else-if="article" class="space-y-8">
        <!-- 标题 -->
        <header class="text-center border-b border-gray-700 pb-6">
          <h1
            class="text-3xl md:text-4xl font-bold bg-gradient-to-r from-purple-400 via-pink-400 to-red-400 bg-clip-text text-transparent"
          >
            {{ article.title }}
          </h1>
        </header>

        <!-- 作者信息 -->
        <div
          class="flex items-center justify-center gap-4 py-4 bg-gray-800/50 rounded-lg px-6"
        >
          <div
            class="w-12 h-12 bg-gradient-to-br from-purple-500 to-pink-500 rounded-full flex items-center justify-center"
          >
            <span class="text-xl font-bold">{{
              article.author.username.charAt(0).toUpperCase()
            }}</span>
          </div>
          <div class="text-left">
            <p class="font-semibold">{{ article.author.username }}</p>
            <p class="text-gray-400 text-sm">{{ article.author.email }}</p>
          </div>
          <div class="text-gray-500 text-sm ml-auto">
            <span>{{ article.created_at }}</span>
          </div>
        </div>

        <!-- 文章分类 -->
        <div class="flex items-center gap-2">
          <span class="text-gray-400">分类：</span>
          <span
            class="px-4 py-1 bg-purple-600/30 text-purple-300 rounded-full text-sm"
          >
            {{ article.category?.name || "未分类" }}
          </span>
        </div>

        <!-- 文章内容 -->
        <div
          class="prose prose-invert max-w-none bg-gray-800/30 rounded-xl p-8 markdown-body"
          v-html="renderedContent"
        ></div>

        <!-- 文章标签 -->
        <div class="flex flex-wrap gap-3 pt-4 border-t border-gray-700">
          <span class="text-gray-400">标签：</span>
          <span
            v-for="tag in article.tags"
            :key="tag.id"
            class="px-3 py-1 bg-gray-700 text-gray-300 rounded-full text-sm hover:bg-gray-600 transition-colors cursor-pointer"
          >
            #{{ tag.name }}
          </span>
          <span v-if="article.tags.length === 0" class="text-gray-500 text-sm"
            >暂无标签</span
          >
        </div>
      </article>
    </div>
  </div>
</template>
