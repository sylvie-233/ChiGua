import { createRouter, createWebHistory } from "vue-router"

// 路由配置
const routes = [
  {
    path: "/",
    name: "Home",
    component: () => import("@/views/HomeView.vue")
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
