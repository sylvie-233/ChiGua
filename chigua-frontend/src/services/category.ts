import api from "./client"
import type { ApiResponse, Category } from "@/types/api"
import { USE_MOCK, mockCategories, mockSuccess } from "./mock"

// 分类相关API
export const createCategory = (data: {
  name: string
}): Promise<ApiResponse<Category>> => {
  if (USE_MOCK) {
    return Promise.resolve(
      mockSuccess({
        id: Date.now(),
        name: data.name,
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString()
      })
    )
  }
  return api.post("/category", data)
}

export const getCategoryList = (): Promise<ApiResponse<Category[]>> => {
  if (USE_MOCK) {
    return Promise.resolve(mockSuccess(mockCategories))
  }
  return api.get("/category")
}

export const deleteCategory = (id: number): Promise<ApiResponse<null>> => {
  if (USE_MOCK) {
    return Promise.resolve(mockSuccess(null))
  }
  return api.delete(`/category/${id}`)
}
