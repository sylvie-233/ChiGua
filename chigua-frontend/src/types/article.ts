import type { User } from "./user"
import type { Category } from "./category"
import type { Tag } from "./tag"

// 文章类型
export interface Article {
  id: number
  title: string
  content: string
  status: number
  author_id: number
  author: User
  category_id: number
  category: Category
  tags: Tag[]
  created_at: string
  updated_at: string
}
