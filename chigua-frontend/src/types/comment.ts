import type { User } from "./user"

// 评论类型
export interface Comment {
  id: number
  article_id: number
  user_id: number
  user: User
  content: string
  created_at: string
  updated_at: string
}
