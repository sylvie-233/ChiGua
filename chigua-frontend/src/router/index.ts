import { createRouter, createWebHistory, type RouteRecordRaw } from "vue-router"

// 路由配置
const routes: RouteRecordRaw[] = [
  {
    path: "/",
    name: "Home",
    component: () => import("@/views/HomeView.vue")
  },
  {
    path: "/article/:id",
    name: "ArticleDetail",
    component: () => import("@/views/ArticleDetailView.vue")
  },
  {
    path: "/admin",
    name: "Admin",
    component: () => import("@/views/AdminView.vue")
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
