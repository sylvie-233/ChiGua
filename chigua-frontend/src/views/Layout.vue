<script setup lang="ts">
import MainLayout from "@/components/MainLayout.vue"
</script>

<template>
  <MainLayout>
    <router-view />
  </MainLayout>
</template>

<style scoped>
:deep(.ant-menu-item-selected) {
  background-color: #262626 !important;
  border-right: 2px solid #6b7280;
}

:deep(.ant-menu-item:hover) {
  background-color: #262626 !important;
}

:deep(.ant-menu-item .anticon) {
  color: #d1d5db !important;
  vertical-align: middle !important;
  transform: translateY(-2px);
}

:deep(.ant-card-head-title) {
  color: #d1d5db;
}

:deep(.ant-card-extra) {
  display: flex;
  align-items: center;
}
</style>
