import { typedApi } from "./client"
import type { Comment } from "@/types/comment"

export const commentApi = {
  async getComments(articleId: number, page: number, pageSize: number) {
    return typedApi.get<Comment[]>(`/comment/article/${articleId}`, {
      params: { page, pageSize }
    })
  },

  async createComment(data: {
    article_id: number
    content: string
    parent_id?: number
    reply_user_id?: number
  }) {
    return typedApi.post<Comment>("/comment", data)
  },

  async deleteComment(id: number) {
    return typedApi.delete<void>(`/comment/${id}`)
  }
}