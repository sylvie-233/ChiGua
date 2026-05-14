import { createRouter, createWebHistory, type RouteRecordRaw } from "vue-router"
import { useAuthStore } from "@/stores/auth"
import { message } from "ant-design-vue"

// 路由配置
const routes: RouteRecordRaw[] = [
  {
    path: "/",
    component: () => import("@/views/Layout.vue"),
    children: [
      {
        path: "",
        name: "Home",
        component: () => import("@/views/HomeView.vue")
      },
      {
        path: "article/:id",
        name: "ArticleDetail",
        component: () => import("@/views/ArticleDetailView.vue")
      },
      {
        path: "profile",
        name: "Profile",
        component: () => import("@/views/ProfileView.vue"),
        meta: { requiresAuth: true }
      }
    ]
  },
  {
    path: "/admin",
    component: () => import("@/components/AdminLayout.vue"),
    meta: { requiresAuth: true, requiresAdmin: true },
    children: [
      {
        path: "article",
        name: "AdminArticle",
        component: () => import("@/views/admin/AdminArticle.vue")
      },
      {
        path: "category",
        name: "AdminCategory",
        component: () => import("@/views/admin/AdminCategory.vue")
      },
      {
        path: "tag",
        name: "AdminTag",
        component: () => import("@/views/admin/AdminTag.vue")
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 导航守卫
router.beforeEach((to, _from, next) => {
  const authStore = useAuthStore()

  // 如果路由需要认证但用户未登录
  if (to.meta.requiresAuth && !authStore.isLoggedIn()) {
    message.error("请登录后访问")
    next("/")
    return
  }

  // 如果路由需要管理员权限
  if (to.meta.requiresAdmin && authStore.user?.role !== "admin") {
    message.error("无权限访问")
    next("/")
    return
  }

  next()
})

export default router
