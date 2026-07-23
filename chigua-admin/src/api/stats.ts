import request from './request'
import type { BaseResponse } from '@/types'

// 文章状态统计
export interface ArticleStats {
  total: number
  draft: number
  published: number
  unpublished: number
  pending: number
}

// 最近文章
export interface RecentArticle {
  id: number
  title: string
  authorName: string
  createdAt: string
}

// 每日文章统计
export interface ArticleDailyCount {
  date: string
  count: number
}

// 仪表盘统计数据
export interface DashboardStats {
  articles: ArticleStats
  categories: number
  tags: number
  comments: number
  users: number
  recentArticles: RecentArticle[]
  articleDailyStats: ArticleDailyCount[]
}

export const getDashboardStats = async (): Promise<BaseResponse<DashboardStats>> => {
  const response = await request.get('/admin/stats')
  return response.data
}
