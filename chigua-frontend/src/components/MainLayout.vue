<script setup lang="ts">
import NavBar from "@/components/NavBar.vue"
</script>

<template>
  <div>
    <NavBar />
    <div style="min-height: calc(100vh - 72px)">
      <slot />
    </div>
  </div>
</template>
