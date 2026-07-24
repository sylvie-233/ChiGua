<script setup lang="ts">
import { useWindowScroll } from "@vueuse/core"
import NavBar from "@/components/NavBar.vue"
import FooterBar from "@/components/FooterBar.vue"

const { y } = useWindowScroll()

const scrollToTop = () => {
  window.scrollTo({ top: 0, behavior: "smooth" })
}
</script>

<template>
  <div>
    <NavBar />
    <div style="min-height: calc(100vh - 72px)">
      <slot />
    </div>
    <FooterBar />

    <!-- 回到顶部 -->
    <Transition name="fade">
      <button
        v-show="y > 400"
        @click="scrollToTop"
        class="fixed right-6 bottom-6 w-10 h-10 rounded-full bg-[#2a2a2a] border border-gray-600 flex items-center justify-center hover:border-gray-400 transition-colors shadow-lg z-40"
      >
        <svg
          class="w-4 h-4 text-gray-400"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          viewBox="0 0 24 24"
        >
          <path stroke-linecap="round" stroke-linejoin="round" d="M18 15l-6-6-6 6" />
        </svg>
      </button>
    </Transition>
  </div>
</template>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
