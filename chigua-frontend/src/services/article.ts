import api from "./client"

// 文章相关API
export const createArticle = (data: any) => api.post("/article", data)

export const getArticleList = (params: { page: number; pageSize: number }) =>
  api.get("/article", { params })

export const getArticle = (id: number) => api.get(`/article/${id}`)

export const updateArticle = (id: number, data: any) =>
  api.put(`/article/${id}`, data)

export const deleteArticle = (id: number) => api.delete(`/article/${id}`)

export const publishArticle = (id: number) => api.put(`/article/${id}/publish`)
