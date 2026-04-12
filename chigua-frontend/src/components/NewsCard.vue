<script setup lang="ts">
import { formatDateWithChinese } from "@/utils/dateFormat"

interface Props {
  title: string
  images: string[]
  author: string
  date: Date
  tags: string[]
  isHot?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  isHot: false
})
</script>

<template>
  <div
    class="relative h-75 rounded-lg shadow-md overflow-hidden hover:shadow-lg transition-all duration-300 transform hover:scale-105"
  >
    <!-- 图片区域 -->
    <div class="absolute inset-0">
      <div class="flex h-full w-full">
        <img
          v-for="(image, index) in images"
          :key="index"
          :src="image"
          :alt="title"
          class="w-1/3 h-full object-cover"
        />
      </div>
    </div>

    <!-- 内容区域 -->
    <div
      class="absolute inset-0 bg-linear-to-t from-black/80 to-transparent p-4 flex flex-col justify-center items-center text-center"
    >
      <h3
        class="text-4xl font-weight-[600] mb-2 line-clamp-2 hover:text-primary transition-colors text-white"
      >
        {{ title }}
      </h3>
      <div class="flex flex-col items-center text-lg text-gray-300 mb-1 gap-2">
        <div class="flex items-center">
          <span>{{ author }}</span>
          <span class="mx-2">·</span>
          <span>{{ formatDateWithChinese(date) }}</span>
          <span class="mx-2">·</span>
          <span>{{ tags.join(", ") }}</span>
        </div>
      </div>
    </div>

    <!-- 热榜标识 -->
    <div
      v-if="isHot"
      class="absolute top-0 right-0 bg-green-500 text-white px-2 py-1 text-xs font-bold transform rotate-45 translate-x-2 translate-y-[50%] z-10"
    >
      热榜 HOT
    </div>
  </div>
</template>

<style scoped>
.line-clamp-2 {
  display: -webkit-box;
  line-clamp: 2;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
