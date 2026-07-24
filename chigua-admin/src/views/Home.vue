<script setup lang="ts">
import { reactive, onMounted, onUnmounted, computed, ref } from 'vue'
import { useRouter } from 'vue-router'
import { Modal, message } from 'ant-design-vue'
import {
  ReloadOutlined,
  FileOutlined,
  AppstoreOutlined,
  TagsOutlined,
  MessageOutlined,
  UserOutlined,
  FileTextOutlined
} from '@ant-design/icons-vue'
import { getDashboardStats } from '@/api/stats'
import type { DashboardStats } from '@/api/stats'
import { useUserStore } from '@/stores/user'
import { formatDate } from '@/utils/date'
 
const router = useRouter()
const userStore = useUserStore()
 
const state = reactive({
  loading: false,
  stats: {
    articles: { total: 0, draft: 0, published: 0, unpublished: 0, pending: 0 },
    categories: 0,
    tags: 0,
    comments: 0,
    users: 0,
    recentArticles: [],
    articleDailyStats: []
  } as DashboardStats
})
 
const topCards = [
  { title: '分类总数', suffix: '个', icon: AppstoreOutlined, color: '#52c41a', bg: '#f6ffed', path: '/categories', key: 'categories' },
  { title: '标签总数', suffix: '个', icon: TagsOutlined, color: '#faad14', bg: '#fffbe6', path: '/tags', key: 'tags' },
  { title: '评论总数', suffix: '条', icon: MessageOutlined, color: '#722ed1', bg: '#f9f0ff', path: '/comments', key: 'comments' },
  { title: '用户总数', suffix: '人', icon: UserOutlined, color: '#eb2f96', bg: '#fff0f6', path: '/users', key: 'users' }
]
 
const articleStatusList: { label: string; key: keyof DashboardStats['articles']; color: string }[] = [
  { label: '草稿', key: 'draft', color: '#8c8c8c' },
  { label: '已发布', key: 'published', color: '#52c41a' },
  { label: '已下架', key: 'unpublished', color: '#faad14' },
  { label: '审核中', key: 'pending', color: '#1890ff' }
]
 
const months = ['1月', '2月', '3月', '4月', '5月', '6月', '7月', '8月', '9月', '10月', '11月', '12月']

// ====== 热力图响应式尺寸 ======
const activityRef = ref<HTMLElement>()
const activityWidth = ref(780)

const heatmapSizes = computed(() => {
  const width = activityWidth.value
  const weekdayWidth = 20
  const availableWidth = Math.max(width - weekdayWidth - 12, 200)
  const gapRatio = 0.2
  const cellSize = Math.max(5, Math.min(20,
    Math.round(availableWidth / (52 + 51 * gapRatio))
  ))
  const cellGap = Math.max(2, Math.round(cellSize * gapRatio))
  const weekWidth = cellSize + cellGap
  const fontSize = Math.max(8, Math.min(12, Math.round(cellSize * 0.8)))
  const borderRadius = Math.max(1, Math.round(cellSize * 0.18))
  return { cellSize, cellGap, weekWidth, fontSize, borderRadius }
})

const activityCssVars = computed(() => ({
  '--cell-size': `${heatmapSizes.value.cellSize}px`,
  '--cell-gap': `${heatmapSizes.value.cellGap}px`,
  '--week-width': `${heatmapSizes.value.weekWidth}px`,
  '--font-size': `${heatmapSizes.value.fontSize}px`,
  '--border-radius': `${heatmapSizes.value.borderRadius}px`,
}))

const monthLabelStyle = (weekIndex: number) => {
  const s = heatmapSizes.value
  // weekdays 列宽 + activity-content gap = 标签需要向右偏移的量
  const offset = s.fontSize * 2.2 + 4
  return { left: `${weekIndex * s.weekWidth + offset}px` }
}
 
const getYearAgoDate = () => {
  const now = new Date()
  const date = new Date(now.getFullYear(), now.getMonth(), now.getDate())
  date.setDate(date.getDate() - 364)
  return date
}
 
const formatDateStr = (date: Date) => {
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}
 
const normalizeDateStr = (dateStr: string) => {
  if (dateStr.includes('T')) {
    return dateStr.split('T')[0]
  }
  return dateStr
}
 
const getActivityData = computed(() => {
  const countMap = new Map<string, number>()
  state.stats.articleDailyStats.forEach(item => {
    countMap.set(normalizeDateStr(item.date), item.count)
  })
 
  const startDate = getYearAgoDate()
  const weeks: { date: string; count: number; month: number }[][] = []
  const monthLabels: { weekIndex: number; label: string }[] = []
  // 用 年-月 作为 key，避免同年月份重叠（如2025年7月和2026年7月）
  type MonthSpan = { firstWeek: number; lastWeek: number; month: number; days: number }
  const monthSpanMap = new Map<string, MonthSpan>()
  const monthOrder: string[] = []

  for (let i = 0; i < 365; i++) {
    const currentDate = new Date(startDate)
    currentDate.setDate(startDate.getDate() + i)
    const dateStr = formatDateStr(currentDate)
    const weekIndex = Math.floor(i / 7)
    const month = currentDate.getMonth()
    const key = `${currentDate.getFullYear()}-${month}`

    if (!weeks[weekIndex]) {
      weeks[weekIndex] = []
    }

    weeks[weekIndex].push({
      date: dateStr,
      count: countMap.get(dateStr) || 0,
      month
    })

    if (!monthSpanMap.has(key)) {
      monthSpanMap.set(key, { firstWeek: weekIndex, lastWeek: weekIndex, month, days: 1 })
      monthOrder.push(key)
    } else {
      const span = monthSpanMap.get(key)!
      span.lastWeek = weekIndex
      span.days++
    }
  }

  // 同月跨年去重：只保留天数更多的那个
  const seenMonths = new Map<number, string>()
  const toRemove = new Set<string>()
  for (const key of monthOrder) {
    const span = monthSpanMap.get(key)!
    const existing = seenMonths.get(span.month)
    if (existing !== undefined) {
      const existingSpan = monthSpanMap.get(existing)!
      if (span.days > existingSpan.days) {
        toRemove.add(existing)
        seenMonths.set(span.month, key)
      } else {
        toRemove.add(key)
      }
    } else {
      seenMonths.set(span.month, key)
    }
  }

  const filteredOrder = monthOrder.filter(k => !toRemove.has(k))
  const firstKey = filteredOrder[0]
  const startKey = `${startDate.getFullYear()}-${startDate.getMonth()}`
  // 只有首月正好是数据起始月且不是1号开始，才算残月 → 标签靠右
  const isFirstMonthPartial = firstKey === startKey && startDate.getDate() !== 1

  for (const key of filteredOrder) {
    const span = monthSpanMap.get(key)!
    const isFirst = key === filteredOrder[0]
    const weekIndex = (isFirst && isFirstMonthPartial) ? span.lastWeek : span.firstWeek
    monthLabels.push({ weekIndex, label: months[span.month] })
  }
  return { weeks, monthLabels }
})
 
const getColorClass = (count: number) => {
  if (count === 0) return 'level-0'
  if (count <= 1) return 'level-1'
  if (count <= 3) return 'level-2'
  if (count <= 6) return 'level-3'
  return 'level-4'
}
 
const getCardValue = (key: string) => {
  return state.stats[key as keyof DashboardStats] as number
}
 
const fetchStats = async () => {
  state.loading = true
  try {
    const res = await getDashboardStats()
    if (res.code === 200) {
      const data = { ...res.data }
      if (!data.recentArticles) {
        data.recentArticles = []
      }
      Object.assign(state.stats, data)
    } else {
      message.error(res.msg || '获取统计数据失败')
    }
  } catch (error) {
    console.error('获取统计数据失败:', error)
    message.error('获取统计数据失败')
  } finally {
    state.loading = false
  }
}
 
let resizeObserver: ResizeObserver | null = null

onMounted(() => {
  fetchStats()
  if (activityRef.value) {
    activityWidth.value = activityRef.value.clientWidth
    resizeObserver = new ResizeObserver((entries) => {
      for (const entry of entries) {
        activityWidth.value = entry.contentRect.width
      }
    })
    resizeObserver.observe(activityRef.value)
  }
})

onUnmounted(() => {
  if (resizeObserver) {
    resizeObserver.disconnect()
    resizeObserver = null
  }
})
 
const handleRefresh = () => {
  fetchStats()
}
 
const navTo = (path: string) => {
  router.push(path)
}
 
const handleArticleClick = (article: { id: number; title: string }) => {
  Modal.info({
    title: '文章详情',
    content: `标题：${article.title}\nID：${article.id}\n\n如需编辑请前往文章管理。`,
    okText: '前往文章管理',
    onOk: () => {
      router.push('/articles')
    }
  })
}
</script>
 
<template>
  <div>
    <a-page-header
      :title="`欢迎回来，${userStore.displayName || '管理员'}`"
      sub-title="系统概览"
    >
      <template #extra>
        <a-space>
          <a-button type="primary" :loading="state.loading" @click="handleRefresh">
            <template #icon><component :is="ReloadOutlined" /></template>
            刷新数据
          </a-button>
        </a-space>
      </template>
    </a-page-header>
 
    <a-spin :spinning="state.loading">
      <a-row :gutter="16" style="margin-top: 24px;">
        <a-col
          v-for="card in topCards"
          :key="card.key"
          :span="6"
        >
          <a-card hoverable @click="navTo(card.path)" style="cursor: pointer;">
            <div style="display: flex; align-items: center; gap: 16px;">
              <a-avatar
                :size="48"
                :style="{ background: card.bg, color: card.color }"
              >
                <component :is="card.icon" style="font-size: 24px;" />
              </a-avatar>
              <div>
                <div style="font-size: 14px; color: #8c8c8c;">{{ card.title }}</div>
                <div style="font-size: 28px; font-weight: 600; color: #262626;">
                  {{ getCardValue(card.key) }}<span style="font-size: 14px; font-weight: normal; margin-left: 4px; color: #8c8c8c;">{{ card.suffix }}</span>
                </div>
              </div>
            </div>
          </a-card>
        </a-col>
      </a-row>
 
      <a-row :gutter="16" style="margin-top: 16px;">
        <a-col :span="24">
          <a-card hoverable @click="navTo('/articles')" style="cursor: pointer;">
            <div style="display: flex; align-items: center; justify-content: space-between; flex-wrap: wrap; gap: 24px;">
              <div style="display: flex; align-items: center; gap: 16px;">
                <a-avatar
                  :size="56"
                  style="background: #e6f7ff; color: #1890ff;"
                >
                  <component :is="FileOutlined" style="font-size: 28px;" />
                </a-avatar>
                <div>
                  <div style="font-size: 14px; color: #8c8c8c;">文章总数</div>
                  <div style="font-size: 32px; font-weight: 600; color: #262626;">
                    {{ state.stats.articles.total }}<span style="font-size: 16px; font-weight: normal; margin-left: 4px; color: #8c8c8c;">篇</span>
                  </div>
                </div>
              </div>
              <div style="display: flex; gap: 32px;">
                <div
                  v-for="status in articleStatusList"
                  :key="status.key"
                  style="text-align: center;"
                >
                  <div style="font-size: 13px; color: #8c8c8c; margin-bottom: 4px;">{{ status.label }}</div>
                  <div style="font-size: 22px; font-weight: 600; color: #262626;">
                    {{ state.stats.articles[status.key] }}
                  </div>
                  <div :style="{ width: '24px', height: '3px', borderRadius: '2px', background: status.color, margin: '6px auto 0' }"></div>
                </div>
              </div>
            </div>
          </a-card>
        </a-col>
      </a-row>
 
      <a-row :gutter="16" style="margin-top: 16px;">
        <a-col :span="24">
          <a-card title="文章发布活跃度">
            <div ref="activityRef" class="activity-wrapper" :style="activityCssVars">
              <div class="months-header">
                <div
                  v-for="label in getActivityData.monthLabels"
                  :key="label.weekIndex"
                  class="month-label"
                  :style="monthLabelStyle(label.weekIndex)"
                >
                  {{ label.label }}
                </div>
              </div>
              <div class="activity-content">
                <div class="weekdays">
                  <div class="weekday">一</div>
                  <div class="weekday">三</div>
                  <div class="weekday">五</div>
                </div>
                <div class="weeks-container">
                  <div
                    v-for="(week, weekIndex) in getActivityData.weeks"
                    :key="weekIndex"
                    class="week-column"
                  >
                    <template
                      v-for="(day, _dayIndex) in week"
                      :key="_dayIndex"
                    >
                      <a-tooltip :title="day.date + ': ' + day.count + '篇文章'">
                        <div
                          class="day-cell"
                          :class="getColorClass(day.count)"
                        ></div>
                      </a-tooltip>
                    </template>
                  </div>
                </div>
              </div>
              <div class="legend">
                <span>少</span>
                <span class="legend-cell level-0"></span>
                <span class="legend-cell level-1"></span>
                <span class="legend-cell level-2"></span>
                <span class="legend-cell level-3"></span>
                <span class="legend-cell level-4"></span>
                <span>多</span>
              </div>
            </div>
          </a-card>
        </a-col>
      </a-row>
 
      <a-row :gutter="16" style="margin-top: 16px;">
        <a-col :span="24">
          <a-card title="最近发布的文章">
            <template #extra>
              <a-button type="link" @click="navTo('/articles')">查看全部</a-button>
            </template>
            <a-empty v-if="!state.loading && state.stats.recentArticles.length === 0" description="暂无文章" />
            <a-list
              v-else
              :data-source="state.stats.recentArticles"
              size="large"
              :loading="state.loading"
            >
              <template #renderItem="{ item }">
                <a-list-item>
                  <a-list-item-meta>
                    <template #avatar>
                      <a-avatar style="background: #1890ff;">
                        <template #icon><component :is="FileTextOutlined" /></template>
                      </a-avatar>
                    </template>
                    <template #title>
                      <a @click="handleArticleClick(item)">{{ item.title }}</a>
                    </template>
                    <template #description>
                      {{ item.authorName }} · {{ formatDate(item.createdAt) }}
                    </template>
                  </a-list-item-meta>
                </a-list-item>
              </template>
            </a-list>
          </a-card>
        </a-col>
      </a-row>
    </a-spin>
  </div>
</template>
 
<style scoped>
  .activity-wrapper {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 8px 0;
  }
  .months-header {
    position: relative;
    height: calc(var(--font-size) + 6px);
    width: 100%;
    max-width: calc(52 * var(--cell-size) + 51 * var(--cell-gap) + var(--font-size) * 2.2 + 4px);
    margin-bottom: 4px;
  }
  .month-label {
    position: absolute;
    font-size: var(--font-size);
    color: #666;
    white-space: nowrap;
  }
  .activity-content {
    display: flex;
    gap: 4px;
  }
  .weekdays {
    display: flex;
    flex-direction: column;
    justify-content: space-around;
    width: calc(var(--font-size) * 2.2);
    padding-top: 2px;
    padding-bottom: 2px;
  }
  .weekday {
    font-size: var(--font-size);
    color: #666;
    height: var(--cell-size);
    line-height: var(--cell-size);
  }
  .weeks-container {
    display: flex;
    gap: var(--cell-gap);
  }
  .week-column {
    display: flex;
    flex-direction: column;
    gap: var(--cell-gap);
  }
  .day-cell {
    width: var(--cell-size);
    height: var(--cell-size);
    border-radius: var(--border-radius);
    cursor: pointer;
  }
  .level-0 { background-color: #ebedf0; }
  .level-1 { background-color: #9be9a8; }
  .level-2 { background-color: #40c463; }
  .level-3 { background-color: #30a14e; }
  .level-4 { background-color: #216e39; }
  .legend {
    display: flex;
    align-items: center;
    gap: 4px;
    margin-top: 8px;
    font-size: var(--font-size);
    color: #666;
    align-self: flex-end;
  }
  .legend-cell {
    width: var(--cell-size);
    height: var(--cell-size);
    border-radius: var(--border-radius);
  }
</style>