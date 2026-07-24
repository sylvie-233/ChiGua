import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export interface UserInfo {
  id: number
  username: string
  nickname: string
  avatar: string
  roles: string[]
  permissions: string[]
  createdAt: string
  updateAt: string
}

export const useUserStore = defineStore('user', () => {
  const token = ref(localStorage.getItem('token') || '')
  const userInfo = ref<UserInfo | null>(null)

  const isLoggedIn = computed(() => !!token.value)

  const permissions = computed(() => userInfo.value?.permissions ?? [])

  const displayName = computed(() => {
    if (!userInfo.value) return ''
    return userInfo.value.nickname || userInfo.value.username
  })

  const roleLabel = computed(() => {
    if (!userInfo.value) return ''
    const roles = userInfo.value.roles || []
    if (roles.includes('admin')) return '管理员'
    if (roles.includes('reviewer')) return '审核员'
    return '普通用户'
  })

  function hasPermission(code: string): boolean {
    return permissions.value.includes(code)
  }

  function hasAnyPermission(...codes: string[]): boolean {
    return codes.some(c => permissions.value.includes(c))
  }

  const setToken = (newToken: string) => {
    token.value = newToken
    localStorage.setItem('token', newToken)
  }

  const setUserInfo = (info: UserInfo) => {
    userInfo.value = info
    localStorage.setItem('userInfo', JSON.stringify(info))
    localStorage.setItem('permissions', JSON.stringify(info.permissions ?? []))
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
    localStorage.removeItem('permissions')
  }

  initFromStorage()

  return {
    token,
    userInfo,
    isLoggedIn,
    permissions,
    displayName,
    roleLabel,
    hasPermission,
    hasAnyPermission,
    setToken,
    setUserInfo,
    logout
  }
})
