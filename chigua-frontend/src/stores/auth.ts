import { defineStore } from "pinia"
import { ref } from "vue"
import type { User } from "@/types/api"
import { userApi } from "@/services"

export const useAuthStore = defineStore("auth", () => {
  // 用户信息
  const user = ref<User | null>(null)
  // 登录token
  const token = ref<string | null>(null)
  // 是否登录
  const isLoggedIn = ref(false)

  // 登录
  const login = async (username: string, password: string) => {
    try {
      const response = await userApi.login({ username, password })
      if (response.code === 200 && response.data) {
        user.value = response.data.user
        token.value = response.data.token
        isLoggedIn.value = true
        // 保存到localStorage
        localStorage.setItem("token", response.data.token)
        localStorage.setItem("user", JSON.stringify(response.data.user))
        return true
      }
      return false
    } catch (error) {
      console.error("登录失败:", error)
      return false
    }
  }

  // 注册
  const register = async (username: string, password: string, email: string) => {
    try {
      const response = await userApi.register({ username, password, email })
      if (response.code === 200 && response.data) {
        return true
      }
      return false
    } catch (error) {
      console.error("注册失败:", error)
      return false
    }
  }

  // 登出
  const logout = () => {
    user.value = null
    token.value = null
    isLoggedIn.value = false
    localStorage.removeItem("token")
    localStorage.removeItem("user")
  }

  // 初始化检查登录状态
  const initAuth = () => {
    const savedToken = localStorage.getItem("token")
    const savedUser = localStorage.getItem("user")
    if (savedToken && savedUser) {
      token.value = savedToken
      try {
        user.value = JSON.parse(savedUser)
        isLoggedIn.value = true
      } catch {
        logout()
      }
    }
  }

  return {
    user,
    token,
    isLoggedIn,
    login,
    register,
    logout,
    initAuth
  }
})
