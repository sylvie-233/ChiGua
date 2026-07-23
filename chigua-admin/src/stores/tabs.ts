import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export interface TabItem {
  path: string
  name: string
  title: string
}

export const useTabsStore = defineStore('tabs', () => {
  const tabs = ref<TabItem[]>([
    { path: '/', name: 'Home', title: '首页' }
  ])
  const activeKey = ref('/')

  const addTab = (tab: TabItem) => {
    const exists = tabs.value.find(t => t.path === tab.path)
    if (!exists) {
      tabs.value.push(tab)
    }
    activeKey.value = tab.path
  }

  const removeTab = (path: string) => {
    const index = tabs.value.findIndex(t => t.path === path)
    if (index === -1) return

    tabs.value.splice(index, 1)

    if (activeKey.value === path) {
      const nextTab = tabs.value[index] || tabs.value[index - 1] || tabs.value[0]
      if (nextTab) {
        activeKey.value = nextTab.path
      }
    }
  }

  const removeOtherTabs = (path: string) => {
    tabs.value = tabs.value.filter(t => t.path === '/' || t.path === path)
    activeKey.value = path
  }

  const removeAllTabs = () => {
    tabs.value = [{ path: '/', name: 'Home', title: '首页' }]
    activeKey.value = '/'
  }

  const setActiveKey = (path: string) => {
    activeKey.value = path
  }

  return {
    tabs,
    activeKey,
    addTab,
    removeTab,
    removeOtherTabs,
    removeAllTabs,
    setActiveKey
  }
})
