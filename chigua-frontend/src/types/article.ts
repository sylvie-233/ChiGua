import type { User } from "./user"
import type { Category } from "./category"
import type { Tag } from "./tag"

// 文章类型（对齐后端 ArticleResponse）
export interface Article {
  id: number
  authorId: number
  categoryId: number
  title: string
  content: string
  coverImage: string
  status: number
  reviewerId: number
  reviewComment: string
  submittedAt: string
  publishAt: string
  createdAt: string
  updateAt: string
  tags: Tag[]
  category: Category
  author: User
  reviewer?: User | null
}
