<script setup lang="ts">
import { onMounted, ref } from "vue"

import NewsCard from "@/components/NewsCard.vue"
import Pagination from "@/components/Pagination.vue"
import SectionDesc from "@/components/SectionDesc.vue"
import { getRandomImages } from "@/utils/randomImage"

// 模拟数据
const news = ref([
  {
    id: 1,
    title:
      "推特网黄CP 碳水王子 与starluck拍片捞金 被素人男友曝光渣女毫无歉意扬言追求爱情！",
    images: getRandomImages(3),
    author: "瓜小哥",
    date: new Date("2026-04-07"),
    tags: ["今日吃瓜", "网红黑料"],
    isHot: true
  },
  {
    id: 2,
    title: "Vue 3.5 发布，带来新特性和性能提升",
    images: getRandomImages(3),
    author: "Vue团队",
    date: new Date("2026-04-10"),
    tags: ["Vue", "前端", "JavaScript"]
  },
  {
    id: 3,
    title: "TypeScript 6.0 新特性详解",
    images: getRandomImages(3),
    author: "TS爱好者",
    date: new Date("2026-04-09"),
    tags: ["TypeScript", "前端"]
  }
])

// 分页相关
const currentPage = ref(2)
const totalPages = ref(1539)
const totalItems = ref(4617)

const handlePageChange = (page: number) => {
  currentPage.value = page
  // 这里可以调用API获取对应页码的数据
  console.log(`切换到第${page}页`)
}

onMounted(() => {
  // 这里可以调用API获取真实数据
  console.log("HomeView mounted")
})
</script>

<template>
  <SectionDesc />
  <div class="max-w-300 mx-auto py-8">
    <!-- 主内容区 -->
    <div class="space-y-8">
      <!-- 文章列表 -->
      <div>
        <div class="space-y-12 px-4">
          <!-- 新闻卡片 -->
          <NewsCard
            v-for="item in news"
            :key="item.id"
            :title="item.title"
            :images="item.images"
            :author="item.author"
            :date="item.date"
            :tags="item.tags"
            :is-hot="item.isHot"
          />
        </div>

        <!-- 分页 -->
        <div class="mt-8">
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
