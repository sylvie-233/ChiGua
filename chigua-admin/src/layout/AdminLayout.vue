<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { HomeOutlined, FileOutlined, FileTextOutlined, AuditOutlined, HistoryOutlined, AppstoreOutlined, TagsOutlined, MessageOutlined, UserOutlined, DownOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { useTabsStore } from '@/stores/tabs'
import { useUserStore } from '@/stores/user'
import { getCurrentUser } from '@/api/user'

const collapsed = ref(false)
const route = useRoute()
const router = useRouter()
const tabsStore = useTabsStore()
const userStore = useUserStore()

const openKeys = ref<string[]>([])

const selectedKeys = computed(() => {
  const path = route.path
  if (path === '/articles/new' || (path.startsWith('/articles/') && path.endsWith('/edit'))) {
    // 新建/编辑文章 → 高亮"文章列表"
    return ['/articles']
  }
  if (path.startsWith('/articles')) {
    // 子菜单项精确匹配
    return [path]
  }
  return [path]
})

// 当在文章相关页面时，自动展开"文章管理"子菜单
watch(
  () => route.path,
  (path) => {
    if (path.startsWith('/articles')) {
      if (!openKeys.value.includes('/articles')) {
        openKeys.value = ['/articles']
      }
    }
  },
  { immediate: true }
)

const menuMap: Record<string, { title: string }> = {
  '/': { title: '首页' },
  '/articles': { title: '文章列表' },
  '/articles/new': { title: '新建文章' },
  '/articles/edit': { title: '编辑文章' },
  '/articles/pending': { title: '审核管理' },
  '/articles/records': { title: '审核记录' },
  '/categories': { title: '分类管理' },
  '/tags': { title: '标签管理' },
  '/comments': { title: '评论管理' },
  '/users': { title: '用户管理' }
}

/** 从路径获取菜单/标签页标题 */
const getPageTitle = (path: string): string | undefined => {
  if (path.startsWith('/articles/') && path.endsWith('/edit')) {
    return menuMap['/articles/edit']?.title
  }
  return menuMap[path]?.title
}

watch(
  () => route.path,
  (path) => {
    const title = getPageTitle(path)
    if (title) {
      tabsStore.addTab({
        path,
        name: (route.name as string) || path,
        title
      })
    }
  },
  { immediate: true }
)

onMounted(async () => {
  if (userStore.token && !userStore.userInfo) {
    try {
      const res = await getCurrentUser()
      if (res.code === 200) {
        userStore.setUserInfo(res.data)
      }
    } catch (error) {
      console.error('获取当前用户信息失败:', error)
    }
  }
})

const handleMenuClick = (e: { key: string }) => {
  if (e.key !== route.path) {
    router.push(e.key)
  }
}

const handleBreadcrumbClick = (path: string) => {
  if (path !== route.path) {
    router.push(path)
  }
}

const handleTabClick = (path: string) => {
  if (path !== route.path) {
    router.push(path)
  }
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
    if (route.path.startsWith('/articles')) {
      items.push({ path: '/articles', title: '文章管理' })
      const title = getPageTitle(route.path)
      if (title && title !== '文章列表') {
        items.push({ path: route.path, title })
      }
    } else {
      const title = menuMap[route.path]?.title
      if (title) {
        items.push({ path: route.path, title })
      }
    }
  }
  return items
})

const roleLabel = computed(() => {
  return userStore.isAdmin ? '管理员' : '普通用户'
})
</script>

<template>
  <a-layout class="layout" style="min-height: 100vh">
    <a-layout-sider v-model:collapsed="collapsed" collapsible>
      <div class="logo">
        <span>{{ collapsed ? '🍉' : '🍉 吃瓜网' }}</span>
      </div>
      <a-menu theme="dark" mode="inline" :selected-keys="selectedKeys" v-model:open-keys="openKeys" @click="handleMenuClick">
        <a-menu-item key="/">
          <component :is="HomeOutlined" />
          <span>首页</span>
        </a-menu-item>
        <a-sub-menu key="/articles">
          <template #icon><component :is="FileOutlined" /></template>
          <template #title>文章管理</template>
          <a-menu-item key="/articles">
            <component :is="FileTextOutlined" />
            <span>文章列表</span>
          </a-menu-item>
          <a-menu-item v-if="userStore.isAdminOrReviewer" key="/articles/pending">
            <component :is="AuditOutlined" />
            <span>审核管理</span>
          </a-menu-item>
          <a-menu-item v-if="userStore.isAdminOrReviewer" key="/articles/records">
            <component :is="HistoryOutlined" />
            <span>审核记录</span>
          </a-menu-item>
        </a-sub-menu>
        <a-menu-item key="/categories">
          <component :is="AppstoreOutlined" />
          <span>分类管理</span>
        </a-menu-item>
        <a-menu-item key="/tags">
          <component :is="TagsOutlined" />
          <span>标签管理</span>
        </a-menu-item>
        <a-menu-item key="/comments">
          <component :is="MessageOutlined" />
          <span>评论管理</span>
        </a-menu-item>
        <a-menu-item key="/users">
          <component :is="UserOutlined" />
          <span>用户管理</span>
        </a-menu-item>
      </a-menu>
    </a-layout-sider>
    <a-layout>
      <a-layout-header style="background: #fff; padding: 0 24px; display: flex; align-items: center; justify-content: space-between; height: 48px; line-height: 48px;">
        <div style="font-size: 16px; font-weight: bold;">后台管理系统</div>
        <a-dropdown>
          <span style="cursor: pointer; display: flex; align-items: center; gap: 8px;">
            <a-avatar style="background: #1890ff;">
              {{ (userStore.displayName || 'U').charAt(0).toUpperCase() }}
            </a-avatar>
            <span>{{ userStore.displayName || '未登录' }}</span>
            <a-tag v-if="userStore.userInfo" :color="userStore.isAdmin ? 'red' : 'blue'">{{ roleLabel }}</a-tag>
            <component :is="DownOutlined" />
          </span>
          <template #overlay>
            <a-menu>
              <a-menu-item @click="handleLogout">退出登录</a-menu-item>
            </a-menu>
          </template>
        </a-dropdown>
      </a-layout-header>
      <div class="tabs-container">
        <a-tabs
          v-model:active-key="tabsStore.activeKey"
          type="editable-card"
          hide-add
          @tab-click="handleTabClick"
          @edit="handleTabEdit"
        >
          <a-tab-pane
            v-for="tab in tabsStore.tabs"
            :key="tab.path"
            :tab="tab.title"
            :closable="tab.path !== '/'"
          />
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
      <a-layout-content style="margin: 8px 16px 16px; padding: 24px; background: #fff; min-height: 280px;">
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

.layout {
  min-height: 100vh;
}

.tabs-container {
  background: #fff;
  padding: 0 16px;
  border-bottom: 1px solid #f0f0f0;
}

.breadcrumb-area {
  margin: 0 16px;
  padding-top: 12px;
}

.breadcrumb-link {
  color: #1890ff;
  cursor: pointer;
  transition: color 0.2s;
}

.breadcrumb-link:hover {
  color: #40a9ff;
}

:deep(.ant-tabs-nav) {
  margin: 0 !important;
}
</style>