import { typedApi } from "./client"
import type { Article } from "@/types/article"

export const articleApi = {
  async getArticleList(params: {
    page: number
    pageSize: number
    categoryId?: number
  }) {
    return typedApi.get<{
      items: Article[]
      total: number
      page: number
      pageSize: number
      totalPages: number
      hasNext: boolean
      hasPrev: boolean
    }>("/article", { params })
  },

  async getArticleById(id: number) {
    return typedApi.get<Article>(`/article/${id}`)
  },

  async createArticle(data: {
    title: string
    content?: string
    coverImage?: string
    categoryId: number
    tagIds?: number[]
  }) {
    return typedApi.post<Article>("/article", data)
  },

  async updateArticle(
    id: number,
    data: {
      title?: string
      content?: string
      coverImage?: string
      categoryId?: number
      tagIds?: number[]
    }
  ) {
    return typedApi.put<Article>(`/article/${id}`, data)
  },

  async deleteArticle(id: number) {
    return typedApi.delete<void>(`/article/${id}`)
  }
}

export const getArticleList = articleApi.getArticleList
export const getArticle = articleApi.getArticleById
