import api from "./client"
import type { ApiResponse } from "@/types/api"
import type { Comment } from "@/types/comment"

export const commentApi = {
  getComments(articleId: number, page: number, pageSize: number) {
    return api.get<ApiResponse<Comment[]>>(`/comment/article/${articleId}`, {
      params: { page, pageSize }
    })
  },

  createComment(data: {
    article_id: number
    content: string
    parent_id?: number
    reply_user_id?: number
  }) {
    return api.post<ApiResponse<Comment>>("/comment", data)
  },

  deleteComment(id: number) {
    return api.delete<ApiResponse<void>>(`/comment/${id}`)
  }
}
