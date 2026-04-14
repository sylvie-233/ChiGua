import api from "./client"

// 评论相关API
export const createComment = (data: { article_id: number; content: string }) =>
  api.post("/comment", data)

export const getCommentsByArticleId = (articleId: number) =>
  api.get(`/comment/article/${articleId}`)

export const deleteComment = (id: number) => api.delete(`/comment/${id}`)
