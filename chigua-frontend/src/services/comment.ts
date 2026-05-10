import api from "./client"
import type { ApiResponse, Comment } from "@/types/api"
import { USE_MOCK, mockComments, mockSuccess } from "./mock"

// 评论相关API
export const createComment = (data: {
  article_id: number
  content: string
}): Promise<ApiResponse<Comment>> => {
  if (USE_MOCK) {
    return Promise.resolve(
      mockSuccess({
        id: Date.now(),
        article_id: data.article_id,
        user_id: 1,
        user: {
          id: 1,
          username: "zhangsan",
          nickname: "张三",
          email: "zhangsan@example.com",
          created_at: "",
          updated_at: ""
        },
        content: data.content,
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString()
      })
    )
  }
  return api.post("/comment", data)
}

export const getCommentsByArticleId = (
  articleId: number
): Promise<ApiResponse<Comment[]>> => {
  if (USE_MOCK) {
    const articleComments = mockComments.filter(c => c.article_id === articleId)
    return Promise.resolve(
      mockSuccess(articleComments.length > 0 ? articleComments : mockComments)
    )
  }
  return api.get(`/comment/article/${articleId}`)
}

export const deleteComment = (id: number): Promise<ApiResponse<null>> => {
  if (USE_MOCK) {
    return Promise.resolve(mockSuccess(null))
  }
  return api.delete(`/comment/${id}`)
}
