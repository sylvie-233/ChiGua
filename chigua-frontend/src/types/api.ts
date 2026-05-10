// API响应类型定义
export interface ApiResponse<T = any> {
  code: number
  msg: string
  data: T
}

// 用户信息类型
export interface User {
  id: number
  username: string
  nickname: string
  email: string
  created_at: string
  updated_at: string
}

// 分类类型
export interface Category {
  id: number
  name: string
  created_at: string
  updated_at: string
}

// 标签类型
export interface Tag {
  id: number
  name: string
  created_at: string
  updated_at: string
}

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
