import api from "./client"
import type { ApiResponse } from "@/types/api"
import type { Article } from "@/types/article"

export const articleApi = {
  getArticleList(params: {
    page: number
    pageSize: number
    categoryId?: number
  }) {
    return api.get<
      ApiResponse<{
        items: Article[]
        total: number
        page: number
        pageSize: number
      }>
    >("/article/list", {
      params
    })
  },

  getArticleById(id: number) {
    return api.get<ApiResponse<Article>>(`/article/${id}`)
  },

  createArticle(data: {
    title: string
    content: string
    category_id: number
    tag_ids?: number[]
  }) {
    return api.post<ApiResponse<Article>>("/article", data)
  },

  updateArticle(
    id: number,
    data: {
      title?: string
      content?: string
      category_id?: number
      tag_ids?: number[]
    }
  ) {
    return api.put<ApiResponse<Article>>(`/article/${id}`, data)
  },

  deleteArticle(id: number) {
    return api.delete<ApiResponse<void>>(`/article/${id}`)
  }
}

// 兼容旧的方法名
export const getArticleList = articleApi.getArticleList
export const getArticle = articleApi.getArticleById
