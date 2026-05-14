import { typedApi } from "./client"
import type { Category } from "@/types/category"

export const categoryApi = {
  async getAllCategories() {
    return typedApi.get<Category[]>("/category/list")
  },

  async createCategory(data: { name: string }) {
    return typedApi.post<Category>("/category", data)
  },

  async updateCategory(id: number, data: { name: string }) {
    return typedApi.put<Category>(`/category/${id}`, data)
  },

  async deleteCategory(id: number) {
    return typedApi.delete<void>(`/category/${id}`)
  }
}