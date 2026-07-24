<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import * as Icons from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { useTabsStore } from '@/stores/tabs'
import { useUserStore } from '@/stores/user'
import { getCurrentUser } from '@/api/user'
import { getMenuTree } from '@/api/menu'
import type { MenuItem } from '@/api/menu'

const collapsed = ref(false)
const route = useRoute()
const router = useRouter()
const tabsStore = useTabsStore()
const userStore = useUserStore()

const openKeys = ref<string[]>([])
const menuTree = ref<MenuItem[]>([])

// 动态图标：从字符串名解析为组件
const resolveIcon = (iconName: string) => {
  if (!iconName) return null
  const icon = (Icons as Record<string, any>)[iconName]
  return icon || null
}

// 加载菜单树
const loadMenu = async () => {
  try {
    const res = await getMenuTree()
    if (res.code === 200) {
      menuTree.value = res.data
    }
  } catch {
    // 菜单加载失败不阻塞
  }
}

// 构建 menuMap（路径 → 标题）
const menuMap = computed(() => {
  const map: Record<string, string> = { '/': '首页' }
  const walk = (items: MenuItem[]) => {
    for (const item of items) {
      if (item.path) map[item.path] = item.title
      if (item.children) walk(item.children)
    }
  }
  walk(menuTree.value)
  return map
})

const selectedKeys = computed(() => {
  const path = route.path
  if (path === '/articles/new' || (path.startsWith('/articles/') && path.endsWith('/edit'))) {
    return ['/articles']
  }
  return [path]
})

// 自动展开当前路径所在子菜单
watch(() => route.path, (path) => {
  const walk = (items: MenuItem[]) => {
    for (const item of items) {
      if (item.children?.length) {
        const childPaths = item.children.map(c => c.path).filter(Boolean)
        if (childPaths.some(p => path.startsWith(p!))) {
          if (!openKeys.value.includes(String(item.id))) {
            openKeys.value = [...openKeys.value, String(item.id)]
          }
        }
        walk(item.children)
      }
    }
  }
  walk(menuTree.value)
}, { immediate: true })

const getPageTitle = (path: string): string | undefined => {
  if (path.startsWith('/articles/') && path.endsWith('/edit')) return '编辑文章'
  return menuMap.value[path]
}

watch(() => route.path, (path) => {
  const title = getPageTitle(path)
  if (title) {
    tabsStore.addTab({ path, name: (route.name as string) || path, title })
  }
}, { immediate: true })

onMounted(async () => {
  if (userStore.token && !userStore.userInfo) {
    try {
      const res = await getCurrentUser()
      if (res.code === 200) userStore.setUserInfo(res.data)
    } catch { /* ignore */ }
  }
  await loadMenu()
})

const handleMenuClick = (e: { key: string }) => {
  if (e.key !== route.path) router.push(e.key)
}

const handleBreadcrumbClick = (path: string) => {
  if (path !== route.path) router.push(path)
}

const handleTabClick = (path: string) => {
  if (path !== route.path) router.push(path)
}

const handleTabEdit = (targetKey: string | MouseEvent, action: 'add' | 'remove') => {
  if (action === 'remove') {
    tabsStore.removeTab(targetKey as string)
    if (tabsStore.activeKey !== router.currentRoute.value.path) {
      router.push(tabsStore.activeKey)
    }
  }
}

const handleLogout = () => {
  userStore.logout()
  tabsStore.removeAllTabs()
  message.success('已退出登录')
  router.push('/login')
}

const breadcrumbs = computed(() => {
  const items = [{ path: '/', title: '首页' }]
  if (route.path !== '/') {
    // 找到当前路径的父菜单链
    const findParentChain = (items: MenuItem[], target: string, chain: { path: string; title: string }[]): boolean => {
      for (const item of items) {
        if (item.path === target) {
          if (item.children?.length) chain.push({ path: item.path, title: item.title })
          return true
        }
        if (item.children?.length) {
          chain.push({ path: item.children[0]?.path || '', title: item.title })
          if (findParentChain(item.children, target, chain)) return true
          chain.pop()
        }
      }
      return false
    }

    const chain: { path: string; title: string }[] = []
    findParentChain(menuTree.value, route.path, chain)
    // 去重
    const seen = new Set<string>()
    for (const c of chain) {
      if (!seen.has(c.title)) { seen.add(c.title); items.push(c) }
    }
    const title = getPageTitle(route.path)
    if (title) items.push({ path: route.path, title })
  }
  return items
})

</script>

<template>
  <a-layout class="layout" style="height: 100vh">
    <a-layout-sider v-model:collapsed="collapsed" collapsible>
      <div class="logo">
        <span>{{ collapsed ? '🍉' : '🍉 吃瓜网' }}</span>
      </div>
      <a-menu theme="dark" mode="inline" :selected-keys="selectedKeys" v-model:open-keys="openKeys" @click="handleMenuClick">
        <template v-for="item in menuTree" :key="item.id">
          <!-- 子菜单 -->
          <a-sub-menu v-if="item.children?.length" :key="String(item.id)">
            <template #icon><component :is="resolveIcon(item.icon)" /></template>
            <template #title>{{ item.title }}</template>
            <template v-for="child in item.children" :key="child.id">
              <a-sub-menu v-if="child.children?.length" :key="String(child.id)">
                <template #icon><component :is="resolveIcon(child.icon)" /></template>
                <template #title>{{ child.title }}</template>
                <a-menu-item v-for="sub in child.children" :key="sub.path">{{ sub.title }}</a-menu-item>
              </a-sub-menu>
              <a-menu-item v-else :key="child.path">
                <template #icon><component :is="resolveIcon(child.icon)" /></template>
                <span>{{ child.title }}</span>
              </a-menu-item>
            </template>
          </a-sub-menu>
          <!-- 普通菜单项 -->
          <a-menu-item v-else :key="item.path">
            <template #icon><component :is="resolveIcon(item.icon)" /></template>
            <span>{{ item.title }}</span>
          </a-menu-item>
        </template>
      </a-menu>
    </a-layout-sider>
    <a-layout>
      <a-layout-header style="background: #fff; padding: 0 24px; display: flex; align-items: center; justify-content: space-between; height: 48px; line-height: 48px;">
        <div style="font-size: 16px; font-weight: bold;">后台管理系统</div>
        <a-dropdown>
          <span style="cursor: pointer; display: flex; align-items: center; gap: 8px;">
            <a-avatar style="background: #1890ff;" :src="userStore.userInfo?.avatar || undefined">
              {{ (userStore.displayName || 'U').charAt(0).toUpperCase() }}
            </a-avatar>
            <span>{{ userStore.displayName || '未登录' }}</span>
            <a-tag v-if="userStore.userInfo" :color="userStore.hasPermission('article:review') ? 'red' : 'blue'">{{ userStore.roleLabel }}</a-tag>
            <component :is="Icons.DownOutlined" />
          </span>
          <template #overlay>
            <a-menu>
              <a-menu-item @click="handleLogout">退出登录</a-menu-item>
            </a-menu>
          </template>
        </a-dropdown>
      </a-layout-header>
      <div class="tabs-container">
        <a-tabs v-model:active-key="tabsStore.activeKey" type="editable-card" hide-add @tab-click="handleTabClick" @edit="handleTabEdit">
          <a-tab-pane v-for="tab in tabsStore.tabs" :key="tab.path" :tab="tab.title" :closable="tab.path !== '/'" />
        </a-tabs>
      </div>
      <div class="breadcrumb-area">
        <a-breadcrumb>
          <a-breadcrumb-item v-for="(item, index) in breadcrumbs" :key="item.path">
            <span v-if="index < breadcrumbs.length - 1" class="breadcrumb-link" @click="handleBreadcrumbClick(item.path)">{{ item.title }}</span>
            <span v-else>{{ item.title }}</span>
          </a-breadcrumb-item>
        </a-breadcrumb>
      </div>
      <a-layout-content style="margin: 8px 16px 16px; padding: 24px; background: #fff; overflow: auto;">
        <router-view :key="route.fullPath" />
      </a-layout-content>
    </a-layout>

  </a-layout>
</template>

<style scoped>
.logo {
  height: 32px;
  margin: 16px;
  font-size: 18px;
  text-align: center;
  color: #fff;
}
.layout { min-height: 100vh; }
.tabs-container { background: #fff; padding: 0 16px; border-bottom: 1px solid #f0f0f0; }
.breadcrumb-area { margin: 0 16px; padding-top: 12px; }
.breadcrumb-link { color: #1890ff; cursor: pointer; transition: color 0.2s; }
.breadcrumb-link:hover { color: #40a9ff; }
:deep(.ant-tabs-nav) { margin: 0 !important; }
</style>
