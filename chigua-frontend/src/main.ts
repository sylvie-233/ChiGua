import "@/style.css"

import { createApp } from "vue"

import App from "@/App.vue"
import router from "@/router"
import pinia from "@/stores"
import { useAuthStore } from "@/stores/auth"

const app = createApp(App)

app.use(pinia)
app.use(router)

// 初始化认证状态
const authStore = useAuthStore()
authStore.initAuth()

app.mount("#app")
