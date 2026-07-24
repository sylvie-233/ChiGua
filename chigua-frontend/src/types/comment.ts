import type { User } from "./user"

// 评论类型（对齐后端 CommentResponse）
export interface Comment {
  id: number
  parentId: number
  articleId: number
  replyUserId: number
  userId: number
  content: string
  createdAt: string
  user?: User
  replyUser?: User
  children?: CommentPage
}

export interface CommentPage {
  total: number
  page: number
  pageSize: number
  totalPages: number
  hasNext: boolean
  hasPrev: boolean
  items: Comment[]
}
