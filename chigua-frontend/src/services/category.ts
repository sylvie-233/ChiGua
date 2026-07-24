import { typedApi } from "./client"
import type { Category } from "@/types/category"

export const categoryApi = {
  async getAllCategories() {
    return typedApi.get<Category[]>("/categorie")
  },

  async createCategory(data: { name: string }) {
    return typedApi.post<Category>("/categorie", data)
  },

  async updateCategory(id: number, data: { name: string }) {
    return typedApi.put<Category>(`/categorie/${id}`, data)
  },

  async deleteCategory(id: number) {
    return typedApi.delete<void>(`/categorie/${id}`)
  }
}
