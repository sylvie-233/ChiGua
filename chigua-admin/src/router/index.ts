import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import AdminLayout from '@/layout/AdminLayout.vue'
import Home from '@/views/Home.vue'
import Login from '@/views/Login.vue'
import ArticleList from '@/views/article/ArticleList.vue'
import ArticleForm from '@/views/article/ArticleForm.vue'
import CategoryList from '@/views/category/CategoryList.vue'
import TagList from '@/views/tag/TagList.vue'
import CommentList from '@/views/comment/CommentList.vue'
import UserList from '@/views/user/UserList.vue'

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
        component: ArticleList
      },
      {
        path: 'articles/new',
        name: 'ArticleCreate',
        component: ArticleForm
      },
      {
        path: 'articles/:id/edit',
        name: 'ArticleEdit',
        component: ArticleForm
      },
      {
        path: 'categories',
        name: 'CategoryList',
        component: CategoryList
      },
      {
        path: 'tags',
        name: 'TagList',
        component: TagList
      },
      {
        path: 'comments',
        name: 'CommentList',
        component: CommentList
      },
      {
        path: 'users',
        name: 'UserList',
        component: UserList
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
    next()
  }
})

export default router
