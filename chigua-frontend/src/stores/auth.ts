import { defineStore } from "pinia"
import { ref } from "vue"
import type { User } from "@/types/user"
import { userApi } from "@/services"

export const useAuthStore = defineStore(
  "auth",
  () => {
    const user = ref<User | null>(null)
    const token = ref<string | null>(null)

    const login = async (
      username: string,
      password: string
    ): Promise<boolean> => {
      try {
        const response = await userApi.userApi.login({ username, password })
        if (response.code === 200) {
          user.value = response.data.user
          token.value = response.data.token
          return true
        }
        return false
      } catch (error) {
        return false
      }
    }

    const register = async (
      username: string,
      password: string,
      nickname: string
    ): Promise<{ success: boolean; message?: string }> => {
      try {
        const response = await userApi.userApi.register({
          username,
          password,
          nickname
        })
        if (response.code === 200) {
          return { success: true }
        }
        return { success: false, message: response.msg }
      } catch (error) {
        return { success: false, message: "网络错误" }
      }
    }

    const logout = () => {
      user.value = null
      token.value = null
    }

    const isLoggedIn = () => {
      return !!token.value
    }

    return {
      user,
      token,
      login,
      register,
      logout,
      isLoggedIn
    }
  },
  {
    persist: {
      storage: localStorage,
      pick: ["user", "token"]
    }
  }
)
