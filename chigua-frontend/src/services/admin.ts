import api from "./client"
import type { ApiResponse } from "@/types/api"
import type { Article } from "@/types/article"
import type { Category } from "@/types/category"
import type { Tag } from "@/types/tag"

export const adminApi = {
  article: {
    async list(params: { page: number; pageSize: number }) {
      const res = await api.get<
        ApiResponse<{
          items: Article[]
          total: number
          page: number
          pageSize: number
        }>
      >("/article/list", { params })
      return res as unknown as ApiResponse<{
        items: Article[]
        total: number
        page: number
        pageSize: number
      }>
    },
    async create(data: {
      title: string
      content: string
      category_id: number
      tag_ids?: number[]
      coverImage?: string
    }) {
      const res = await api.post<ApiResponse<Article>>("/article", data)
      return res as unknown as ApiResponse<Article>
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
      const res = await api.put<ApiResponse<Article>>(`/article/${id}`, data)
      return res as unknown as ApiResponse<Article>
    },
    async delete(id: number) {
      const res = await api.delete<ApiResponse<void>>(`/article/${id}`)
      return res as unknown as ApiResponse<void>
    },
    async publish(id: number) {
      const res = await api.post<ApiResponse<void>>(`/article/${id}/publish`)
      return res as unknown as ApiResponse<void>
    },
    async updateStatus(id: number, status: number) {
      const res = await api.put<ApiResponse<void>>(`/article/${id}/status`, {
        status
      })
      return res as unknown as ApiResponse<void>
    }
  },
  category: {
    async getAll() {
      const res = await api.get<ApiResponse<Category[]>>("/categorie")
      return res as unknown as ApiResponse<Category[]>
    },
    async create(data: { name: string }) {
      const res = await api.post<ApiResponse<Category>>("/categorie", data)
      return res as unknown as ApiResponse<Category>
    },
    async update(id: number, data: { name: string }) {
      const res = await api.put<ApiResponse<Category>>(`/categorie/${id}`, data)
      return res as unknown as ApiResponse<Category>
    },
    async delete(id: number) {
      const res = await api.delete<ApiResponse<void>>(`/categorie/${id}`)
      return res as unknown as ApiResponse<void>
    }
  },
  tag: {
    async getAll() {
      const res = await api.get<ApiResponse<Tag[]>>("/tag")
      return res as unknown as ApiResponse<Tag[]>
    },
    async create(data: { name: string }) {
      const res = await api.post<ApiResponse<Tag>>("/tag", data)
      return res as unknown as ApiResponse<Tag>
    },
    async update(id: number, data: { name: string }) {
      const res = await api.put<ApiResponse<Tag>>(`/tag/${id}`, data)
      return res as unknown as ApiResponse<Tag>
    },
    async delete(id: number) {
      const res = await api.delete<ApiResponse<void>>(`/tag/${id}`)
      return res as unknown as ApiResponse<void>
    }
  }
}
