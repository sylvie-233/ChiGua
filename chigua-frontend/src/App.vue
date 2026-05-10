<script setup lang="ts">
// App组件
import NavBar from "@/components/NavBar.vue"
</script>

<template>
  <div>
    <NavBar />
    <div style="min-height: calc(100vh - 72px);">
      <router-view />
    </div>
  </div>
</template>

<style></style>
