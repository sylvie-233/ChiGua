import api from "./client"

// 分类相关API
export const createCategory = (data: { name: string }) =>
  api.post("/category", data)

export const getCategoryList = () => api.get("/category")

export const deleteCategory = (id: number) => api.delete(`/category/${id}`)
