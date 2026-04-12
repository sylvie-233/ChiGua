<script setup lang="ts">
// 章节描述组件

interface Props {
  title?: string
  description?: string
}

const props = withDefaults(defineProps<Props>(), {
  title: "吃瓜网 - 全网更新最快最全的吃瓜网",
  description:
    "全网更新最快最全的吃瓜网站，网红黑料、明星绯闻、吃瓜视频、瓜榜TOP10实时更新"
})
</script>

<template>
  <div class="flex justify-center items-center flex-col bg-[#343232] py-8">
    <h2 class="text-3xl font-bold mb-4">{{ title }}</h2>
    <p class="text-sm">{{ description }}</p>
  </div>
</template>

<style scoped>
/* 导航栏样式 */
</style>
