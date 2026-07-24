import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import AdminLayout from '@/layout/AdminLayout.vue'
import Home from '@/views/Home.vue'
import Login from '@/views/Login.vue'
import ArticleList from '@/views/article/ArticleList.vue'
import ArticleForm from '@/views/article/ArticleForm.vue'
import PendingReview from '@/views/article/PendingReview.vue'
import ReviewRecords from '@/views/article/ReviewRecords.vue'
import CategoryList from '@/views/category/CategoryList.vue'
import TagList from '@/views/tag/TagList.vue'
import CommentList from '@/views/comment/CommentList.vue'
import UserList from '@/views/user/UserList.vue'
import RoleList from '@/views/role/RoleList.vue'
import MenuList from '@/views/role/MenuList.vue'
import Forbidden from '@/views/403.vue'

// 扩展 meta 类型
declare module 'vue-router' {
  interface RouteMeta {
    requiresAuth?: boolean
    permission?: string
  }
}

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: Login,
    meta: { requiresAuth: false }
  },
  {
    path: '/',
    component: AdminLayout,
    meta: { requiresAuth: true },
    children: [
      {
        path: '',
        name: 'Home',
        component: Home
      },
      {
        path: 'articles',
        name: 'ArticleList',
        component: ArticleList,
        meta: { permission: 'article:list' }
      },
      {
        path: 'articles/new',
        name: 'ArticleCreate',
        component: ArticleForm,
        meta: { permission: 'article:create' }
      },
      {
        path: 'articles/:id/edit',
        name: 'ArticleEdit',
        component: ArticleForm,
        meta: { permission: 'article:edit' }
      },
      {
        path: 'articles/pending',
        name: 'PendingReview',
        component: PendingReview,
        meta: { permission: 'article:review' }
      },
      {
        path: 'articles/records',
        name: 'ReviewRecords',
        component: ReviewRecords,
        meta: { permission: 'article:records' }
      },
      {
        path: 'categories',
        name: 'CategoryList',
        component: CategoryList,
        meta: { permission: 'category:list' }
      },
      {
        path: 'tags',
        name: 'TagList',
        component: TagList,
        meta: { permission: 'tag:list' }
      },
      {
        path: 'comments',
        name: 'CommentList',
        component: CommentList,
        meta: { permission: 'comment:list' }
      },
      {
        path: 'system/users',
        name: 'UserList',
        component: UserList,
        meta: { permission: 'user:list' }
      },
      {
        path: 'system/roles',
        name: 'RoleList',
        component: RoleList,
        meta: { permission: 'role:list' }
      },
      {
        path: 'system/menus',
        name: 'MenuList',
        component: MenuList,
        meta: { permission: 'role:list' }
      },
      {
        path: '403',
        name: 'Forbidden',
        component: Forbidden,
        meta: { permission: undefined }
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, _, next) => {
  const token = localStorage.getItem('token')
  const requiresAuth = to.meta.requiresAuth !== false

  if (requiresAuth && !token) {
    next('/login')
  } else if (!requiresAuth && token && to.path === '/login') {
    next('/')
  } else {
    // 权限检查：如果有 permission meta，检查用户是否有该权限
    if (to.meta.permission) {
      const permissions: string[] = JSON.parse(localStorage.getItem('permissions') || '[]')
      if (permissions.length > 0 && !permissions.includes(to.meta.permission as string)) {
        // 有权限列表但不包含所需权限 → 无权限
        next('/403')
        return
      }
    }
    next()
  }
})

export default router
