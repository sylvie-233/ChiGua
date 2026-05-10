import api from "./client"
import type { ApiResponse, Article } from "@/types/api"
import { USE_MOCK, mockArticles, mockSuccess } from "./mock"

// 文章相关API
export const createArticle = (data: any): Promise<ApiResponse<Article>> => {
  if (USE_MOCK) {
    return Promise.resolve(
      mockSuccess({
        id: Date.now(),
        title: data.title,
        content: data.content,
        status: 0,
        author_id: 1,
        author: {
          id: 1,
          username: "zhangsan",
          nickname: "张三",
          email: "zhangsan@example.com",
          created_at: new Date().toISOString(),
          updated_at: new Date().toISOString()
        },
        category_id: data.category_id || 1,
        category: { id: 1, name: "科技", created_at: "", updated_at: "" },
        tags: [],
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString()
      })
    )
  }
  return api.post("/article", data)
}

export const getArticleList = (params: {
  page: number
  pageSize: number
}): Promise<ApiResponse<Article[]>> => {
  if (USE_MOCK) {
    const start = (params.page - 1) * params.pageSize
    const end = start + params.pageSize
    const paginatedArticles = mockArticles.slice(start, end)
    return Promise.resolve(mockSuccess(paginatedArticles))
  }
  return api.get("/article", { params })
}

export const getArticle = (id: number): Promise<ApiResponse<Article>> => {
  if (USE_MOCK) {
    const article = mockArticles.find(a => a.id === id) || mockArticles[0]
    return Promise.resolve(mockSuccess(article))
  }
  return api.get(`/article/${id}`)
}

export const updateArticle = (
  id: number,
  data: any
): Promise<ApiResponse<Article>> => {
  if (USE_MOCK) {
    const article = mockArticles.find(a => a.id === id) || mockArticles[0]
    return Promise.resolve(mockSuccess({ ...article, ...data }))
  }
  return api.put(`/article/${id}`, data)
}

export const deleteArticle = (id: number): Promise<ApiResponse<null>> => {
  if (USE_MOCK) {
    return Promise.resolve(mockSuccess(null))
  }
  return api.delete(`/article/${id}`)
}

export const publishArticle = (id: number): Promise<ApiResponse<null>> => {
  if (USE_MOCK) {
    return Promise.resolve(mockSuccess(null))
  }
  return api.put(`/article/${id}/publish`)
}
