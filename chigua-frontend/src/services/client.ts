import axios from "axios"
import type { AxiosRequestConfig } from "axios"
import NProgress from "nprogress"
import type { ApiResponse } from "@/types/api"

const api = axios.create({
  baseURL:
    import.meta.env.VITE_API_BASE_URL ||
    (import.meta.env.PROD ? "/api" : "http://localhost:8080/api"),
  timeout: 10000,
  headers: {
    "Content-Type": "application/json"
  }
})

api.interceptors.request.use(
  config => {
    NProgress.start()
    // FormData 需要浏览器自动设置 Content-Type（含 boundary）
    if (config.data instanceof FormData) {
      config.headers.set("Content-Type", null)
    }
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
    NProgress.done()
    return Promise.reject(error)
  }
)

api.interceptors.response.use(
  response => {
    NProgress.done()
    return response.data
  },
  error => {
    NProgress.done()
    console.error("API Error:", error)
    if (error.response && error.response.data) {
      return Promise.resolve(error.response.data)
    }
    return Promise.reject(error)
  }
)

interface TypedApi {
  get<T>(url: string, config?: AxiosRequestConfig): Promise<ApiResponse<T>>
  post<T>(url: string, data?: unknown, config?: AxiosRequestConfig): Promise<ApiResponse<T>>
  put<T>(url: string, data?: unknown, config?: AxiosRequestConfig): Promise<ApiResponse<T>>
  delete<T>(url: string, config?: AxiosRequestConfig): Promise<ApiResponse<T>>
}

const typedApi: TypedApi = {
  get<T>(url: string, config?: AxiosRequestConfig) {
    return api.get<ApiResponse<T>>(url, config) as unknown as Promise<ApiResponse<T>>
  },
  post<T>(url: string, data?: unknown, config?: AxiosRequestConfig) {
    return api.post<ApiResponse<T>>(url, data, config) as unknown as Promise<ApiResponse<T>>
  },
  put<T>(url: string, data?: unknown, config?: AxiosRequestConfig) {
    return api.put<ApiResponse<T>>(url, data, config) as unknown as Promise<ApiResponse<T>>
  },
  delete<T>(url: string, config?: AxiosRequestConfig) {
    return api.delete<ApiResponse<T>>(url, config) as unknown as Promise<ApiResponse<T>>
  }
}

export { typedApi }
export default api