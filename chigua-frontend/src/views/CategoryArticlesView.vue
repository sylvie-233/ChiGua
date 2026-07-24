<script setup lang="ts">
import { onMounted, ref, watch, computed } from "vue"
import { useRoute, useRouter } from "vue-router"
import NewsCard from "@/components/NewsCard.vue"
import Pagination from "@/components/Pagination.vue"
import { getRandomImages } from "@/utils/randomImage"
import { articleApi } from "@/services/article"
import { categoryApi } from "@/services/category"
import type { Article } from "@/types/article"
import type { Category } from "@/types/category"

const route = useRoute()
const router = useRouter()

const articles = ref<Article[]>([])
const loading = ref(false)
const currentPage = ref(1)
const totalPages = ref(0)
const totalItems = ref(0)
const pageSize = 10
const category = ref<Category | null>(null)

const categoryId = computed(() => Number(route.params.id))

const parseCoverImages = (coverImage: string): string[] => {
  if (!coverImage) return getRandomImages(3)
  const urls = coverImage.split(",").map(u => u.trim()).filter(Boolean)
  return urls.length > 0 ? urls : getRandomImages(3)
}

const fetchCategory = async () => {
  if (isNaN(categoryId.value)) return
  try {
    const res = await categoryApi.getAllCategories()
    if (res.code === 200 && res.data) {
      category.value = res.data.find((c: Category) => c.id === categoryId.value) || null
    }
  } catch (error) {
    console.error("获取分类信息失败:", error)
  }
}

const fetchArticles = async (page: number) => {
  if (isNaN(categoryId.value)) return
  loading.value = true
  try {
    const res = await articleApi.getArticleList({
      page,
      pageSize,
      categoryId: categoryId.value
    })
    if (res.code === 200 && res.data) {
      articles.value = res.data.items || []
      totalPages.value = res.data.totalPages || 0
      totalItems.value = res.data.total || 0
      currentPage.value = res.data.page || page
    }
  } catch (error) {
    console.error("获取文章列表失败:", error)
  } finally {
    loading.value = false
  }
}

const handleCardClick = (id: number) => {
  router.push(`/article/${id}`)
}

const handlePageChange = (page: number) => {
  currentPage.value = page
  fetchArticles(page)
  window.scrollTo({ top: 0, behavior: "smooth" })
}

watch(categoryId, () => {
  fetchCategory()
  fetchArticles(1)
})

onMounted(() => {
  fetchCategory()
  fetchArticles(1)
})
</script>

<template>
  <div class="bg-[#262626] text-white min-h-[calc(100vh-72px)]">
    <!-- 分类信息头部（与首页 SectionDesc 风格一致） -->
    <div class="flex justify-center items-center flex-col bg-[#343232] py-8">
      <h2 class="text-3xl font-bold mb-4">
        {{ category?.name || "分类" }}
      </h2>
      <p class="text-sm text-gray-300">
        {{ category?.description || "浏览该分类下的所有文章" }}
      </p>
    </div>

    <!-- 文章列表区域 -->
    <div class="max-w-200 mx-auto py-8">
      <div class="space-y-8">
        <!-- 加载状态 -->
        <div v-if="loading" class="flex justify-center items-center h-64">
          <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-purple-500"></div>
        </div>

        <!-- 文章列表 -->
        <div v-else class="space-y-12 px-4">
          <NewsCard
            v-for="item in articles"
            :key="item.id"
            :id="item.id"
            :title="item.title"
            :images="parseCoverImages(item.coverImage)"
            :author="item.author?.nickname || item.author?.username || '未知'"
            :date="new Date(item.publishAt || item.createdAt)"
            :tags="item.tags?.map(t => t.name) || []"
            :is-hot="false"
            @click="handleCardClick"
          />
        </div>

        <!-- 空状态 -->
        <div v-if="!loading && articles.length === 0" class="text-center py-16">
          <p class="text-gray-400 text-lg">该分类下暂无文章</p>
        </div>

        <!-- 分页 -->
        <div v-if="totalPages > 1" class="mt-8">
          <Pagination
            :current-page="currentPage"
            :total-pages="totalPages"
            :total-items="totalItems"
            @page-change="handlePageChange"
          />
        </div>
      </div>
    </div>
  </div>
</template>
