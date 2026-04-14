import axios from "axios"

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
    // 从localStorage获取token
    const token = localStorage.getItem("token")
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  error => {
    return Promise.reject(error)
  }
)

// 响应拦截器
api.interceptors.response.use(
  response => {
    // 统一处理响应格式
    return response.data
  },
  error => {
    // 统一处理错误
    console.error("API Error:", error)
    return Promise.reject(error)
  }
)

export default api
