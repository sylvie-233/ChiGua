import { typedApi } from "./client"
import type { Article } from "@/types/article"
import type { Category } from "@/types/category"
import type { Tag } from "@/types/tag"

export const adminApi = {
  article: {
    async list(params: { page: number; pageSize: number }) {
      return typedApi.get<{
        items: Article[]
        total: number
        page: number
        pageSize: number
      }>("/article/list", { params })
    },
    async create(data: {
      title: string
      content: string
      category_id: number
      tag_ids?: number[]
      coverImage?: string
    }) {
      return typedApi.post<Article>("/article", data)
    },
    async update(
      id: number,
      data: {
        title?: string
        content?: string
        category_id?: number
        tag_ids?: number[]
        coverImage?: string
      }
    ) {
      return typedApi.put<Article>(`/article/${id}`, data)
    },
    async delete(id: number) {
      return typedApi.delete<void>(`/article/${id}`)
    },
    async publish(id: number) {
      return typedApi.post<void>(`/article/${id}/publish`)
    },
    async updateStatus(id: number, status: number) {
      return typedApi.put<void>(`/article/${id}/status`, { status })
    }
  },
  category: {
    async getAll() {
      return typedApi.get<Category[]>("/categorie")
    },
    async create(data: { name: string }) {
      return typedApi.post<Category>("/categorie", data)
    },
    async update(id: number, data: { name: string }) {
      return typedApi.put<Category>(`/categorie/${id}`, data)
    },
    async delete(id: number) {
      return typedApi.delete<void>(`/categorie/${id}`)
    }
  },
  tag: {
    async getAll() {
      return typedApi.get<Tag[]>("/tag")
    },
    async create(data: { name: string }) {
      return typedApi.post<Tag>("/tag", data)
    },
    async update(id: number, data: { name: string }) {
      return typedApi.put<Tag>(`/tag/${id}`, data)
    },
    async delete(id: number) {
      return typedApi.delete<void>(`/tag/${id}`)
    }
  }
}