import api from "./client"
import type { ApiResponse } from "@/types/api"
import type { Category } from "@/types/category"

export const categoryApi = {
  getAllCategories() {
    return api.get<ApiResponse<Category[]>>("/category/list")
  },

  createCategory(data: { name: string }) {
    return api.post<ApiResponse<Category>>("/category", data)
  },

  updateCategory(id: number, data: { name: string }) {
    return api.put<ApiResponse<Category>>(`/category/${id}`, data)
  },

  deleteCategory(id: number) {
    return api.delete<ApiResponse<void>>(`/category/${id}`)
  }
}
