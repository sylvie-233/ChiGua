import axios from "axios"
import NProgress from "nprogress"

// 创建axios实例
const api = axios.create({
  baseURL:
    import.meta.env.VITE_API_BASE_URL ||
    (import.meta.env.PROD ? "/api" : "http://localhost:8080/api"),
  timeout: 10000,
  headers: {
    "Content-Type": "application/json"
  }
})

// 请求拦截器
api.interceptors.request.use(
  config => {
    // 显示加载进度条
    NProgress.start()

    // 从Pinia持久化的存储中获取token
    const authStore = localStorage.getItem("auth")
    if (authStore) {
      try {
        const auth = JSON.parse(authStore)
        if (auth.token) {
          config.headers.Authorization = `Bearer ${auth.token}`
        }
      } catch (e) {
        console.error("Failed to parse auth store:", e)
      }
    }
    return config
  },
  error => {
    // 隐藏加载进度条
    NProgress.done()
    return Promise.reject(error)
  }
)

// 响应拦截器
api.interceptors.response.use(
  response => {
    // 隐藏加载进度条
    NProgress.done()
    // 统一处理响应格式
    return response.data
  },
  error => {
    // 隐藏加载进度条
    NProgress.done()
    // 统一处理错误
    console.error("API Error:", error)
    // 如果后端返回了业务错误信息，也返回它
    if (error.response && error.response.data) {
      return Promise.resolve(error.response.data)
    }
    return Promise.reject(error)
  }
)

export default api
