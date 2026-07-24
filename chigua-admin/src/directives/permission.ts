import type { Directive } from 'vue'

const getPermissions = (): string[] => {
  try {
    return JSON.parse(localStorage.getItem('permissions') || '[]')
  } catch {
    return []
  }
}

// v-permission="'article:create'"   — 无权限则隐藏元素
// v-permission:any="['a','b']"      — 有任一权限则显示
export const permission: Directive = {
  mounted(el, binding) {
    const perms = getPermissions()
    const required = binding.value
    const mode = binding.arg // undefined = all, 'any' = any

    if (typeof required === 'string') {
      if (!perms.includes(required)) {
        el.parentNode?.removeChild(el)
      }
    } else if (Array.isArray(required)) {
      const has = mode === 'any'
        ? required.some((p: string) => perms.includes(p))
        : required.every((p: string) => perms.includes(p))
      if (!has) {
        el.parentNode?.removeChild(el)
      }
    }
  }
}
