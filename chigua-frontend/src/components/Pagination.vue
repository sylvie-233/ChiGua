<script setup lang="ts">
import { ref } from "vue"

interface Props {
  currentPage: number
  totalPages: number
  totalItems: number
}

const props = defineProps<Props>()

const emit = defineEmits<{
  (e: "pageChange", page: number): void
}>()

const goToPage = (page: number) => {
  if (page >= 1 && page <= props.totalPages) {
    emit("pageChange", page)
  }
}

const goToPrevious = () => {
  goToPage(props.currentPage - 1)
}

const goToNext = () => {
  goToPage(props.currentPage + 1)
}

const inputPage = ref(props.currentPage.toString())

const handleJump = () => {
  const page = parseInt(inputPage.value)
  if (!isNaN(page)) {
    goToPage(page)
  }
}
</script>

<template>
  <div class="flex items-center justify-between w-full py-4">
    <!-- 上一页按钮 -->
    <button
      :disabled="currentPage === 1"
      class="px-4 py-2 border rounded disabled:opacity-50 disabled:cursor-not-allowed"
      @click="goToPrevious"
    >
      上一页
    </button>

    <!-- 页码信息 -->
    <div class="flex items-center gap-4">
      <span class="text-sm">{{ currentPage }}/{{ totalPages }}</span>
      <div class="flex items-center gap-2">
        <input
          v-model="inputPage"
          min="1"
          type="number"
          :max="totalPages"
          class="w-16 px-2 py-1 border rounded text-center"
        />
        <button class="px-3 py-1 border rounded" @click="handleJump">
          跳转
        </button>
      </div>
    </div>

    <!-- 下一页按钮 -->
    <button
      :disabled="currentPage === totalPages"
      class="px-4 py-2 border rounded disabled:opacity-50 disabled:cursor-not-allowed"
      @click="goToNext"
    >
      下一页
    </button>
  </div>
</template>

<style scoped>
/* 隐藏数字输入框的上下箭头按钮 */
input[type="number"]::-webkit-inner-spin-button,
input[type="number"]::-webkit-outer-spin-button {
  -webkit-appearance: none;
  margin: 0;
}

input[type="number"] {
  appearance: textfield;
  -moz-appearance: textfield;
}
</style>
