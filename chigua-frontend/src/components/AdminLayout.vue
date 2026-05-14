<script setup lang="ts">
import { ref, h } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { Layout, Menu, Breadcrumb, Dropdown, Avatar, Button } from 'ant-design-vue'
import type { ItemType } from 'ant-design-vue'
import {
  FileTextOutlined,
  FolderOutlined,
  TagOutlined,
  UserOutlined,
  HomeOutlined,
  MenuFoldOutlined,
  MenuUnfoldOutlined
} from '@ant-design/icons-vue'
import { useAuthStore } from '@/stores/auth'

const { Header, Sider, Content } = Layout

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const collapsed = ref(false)

const menuItems: ItemType[] = [
  {
    key: '/admin/article',
    icon: () => h(FileTextOutlined),
    label: '文章管理'
  },
  {
    key: '/admin/category',
    icon: () => h(FolderOutlined),
    label: '分类管理'
  },
  {
    key: '/admin/tag',
    icon: () => h(TagOutlined),
    label: '标签管理'
  }
]

const selectedKeys = ref([route.path])

const handleMenuClick = (info: { key: string | number }) => {
  router.push(String(info.key))
}

const toggleCollapsed = () => {
  collapsed.value = !collapsed.value
}

const handleBack = () => {
  router.push('/')
}

const getBreadcrumb = () => {
  const name = route.name as string
  if (name === 'AdminArticle') return '文章管理'
  if (name === 'AdminCategory') return '分类管理'
  if (name === 'AdminTag') return '标签管理'
  return ''
}
</script>

<template>
  <Layout class="admin-layout">
    <Sider
      v-model:collapsed="collapsed"
      :trigger="null"
      :width="200"
      :collapsed-width="64"
      class="bg-[#001529]!"
    >
      <div class="logo text-white text-center py-4 font-bold text-lg">
        <template v-if="!collapsed">后台管理</template>
        <template v-else>管理</template>
      </div>
      <Menu
        v-model:selectedKeys="selectedKeys"
        theme="dark"
        mode="inline"
        :items="menuItems"
        @click="handleMenuClick"
      />
    </Sider>

    <Layout>
      <Header class="bg-white! px-4! flex! items-center! justify-between! shadow-md!">
        <div class="flex items-center gap-4">
          <Button
            type="text"
            @click="toggleCollapsed"
          >
            <template #icon>
              <MenuUnfoldOutlined v-if="collapsed" />
              <MenuFoldOutlined v-else />
            </template>
          </Button>
          <Breadcrumb>
            <Breadcrumb.Item href="" @click="handleBack">
              <HomeOutlined />
            </Breadcrumb.Item>
            <Breadcrumb.Item>{{ getBreadcrumb() }}</Breadcrumb.Item>
          </Breadcrumb>
        </div>

        <Dropdown>
          <div class="flex items-center gap-2 cursor-pointer">
            <Avatar>
              <template #icon><UserOutlined /></template>
            </Avatar>
            <span>{{ authStore.user?.nickname || authStore.user?.username }}</span>
          </div>
          <template #overlay>
            <Menu>
              <Menu.Item key="back" @click="handleBack">
                <HomeOutlined /> 返回首页
              </Menu.Item>
              <Menu.Divider />
            </Menu>
          </template>
        </Dropdown>
      </Header>

      <Content class="bg-[#f0f2f5]! p-6!">
        <div class="bg-white! rounded-lg! min-h-[calc(100vh-120px)]! shadow-sm!">
          <slot />
        </div>
      </Content>
    </Layout>
  </Layout>
</template>

<style scoped>
.admin-layout {
  min-height: 100vh;
}

.logo {
  transition: all 0.3s;
}
</style>
