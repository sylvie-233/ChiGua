import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export interface UserInfo {
  id: number
  username: string
  nickname: string
  role: string
  createdAt: string
  updateAt: string
}

export const useUserStore = defineStore('user', () => {
  const token = ref(localStorage.getItem('token') || '')
  const userInfo = ref<UserInfo | null>(null)

  const isLoggedIn = computed(() => !!token.value)
  const isAdmin = computed(() => userInfo.value?.role === 'admin')
  const isAdminOrReviewer = computed(() => {
    const role = userInfo.value?.role
    return role === 'admin' || role === 'reviewer'
  })
  const displayName = computed(() => {
    if (!userInfo.value) return ''
    return userInfo.value.nickname || userInfo.value.username
  })

  const setToken = (newToken: string) => {
    token.value = newToken
    localStorage.setItem('token', newToken)
  }

  const setUserInfo = (info: UserInfo) => {
    userInfo.value = info
    localStorage.setItem('userInfo', JSON.stringify(info))
  }

  const initFromStorage = () => {
    const stored = localStorage.getItem('userInfo')
    if (stored) {
      try {
        userInfo.value = JSON.parse(stored)
      } catch (e) {
        localStorage.removeItem('userInfo')
      }
    }
  }

  const logout = () => {
    token.value = ''
    userInfo.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('userInfo')
  }

  initFromStorage()

  return {
    token,
    userInfo,
    isLoggedIn,
    isAdmin,
    isAdminOrReviewer,
    displayName,
    setToken,
    setUserInfo,
    logout
  }
})
