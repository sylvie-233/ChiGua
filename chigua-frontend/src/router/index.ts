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
        path: "category/:id",
        name: "CategoryArticles",
        component: () => import("@/views/CategoryArticlesView.vue")
      },
      {
        path: "profile",
        name: "Profile",
        component: () => import("@/views/ProfileView.vue"),
        meta: { requiresAuth: true }
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior() {
    return { top: 0 }
  }
})

// 导航守卫
router.beforeEach((to, _from, next) => {
  const authStore = useAuthStore()

  if (to.meta.requiresAuth && !authStore.isLoggedIn()) {
    message.error("请登录后访问")
    next("/")
    return
  }

  next()
})

export default router
