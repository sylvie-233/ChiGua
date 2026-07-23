export interface BaseResponse<T = any> {
  code: number
  msg: string
  data: T
}

export interface PageResponse<T = any> {
  total: number
  items: T[]
  page: number
  pageSize: number
  totalPages: number
  hasNext: boolean
  hasPrev: boolean
}

export interface UserResponse {
  id: number
  username: string
  nickname: string
  role: string
  createdAt: string
}

export interface Category {
  id: number
  name: string
  createdAt: string
  updateAt: string
}

export interface Tag {
  id: number
  name: string
  createdAt: string
  updateAt: string
}

export interface Article {
  id: number
  authorId: number
  categoryId: number
  title: string
  content: string
  coverImage: string
  status: number
  publishAt: string
  createdAt: string
  updateAt: string
}

export interface ArticleResponse extends Article {
  tags: Tag[]
  category: Category
  author: UserResponse
}

export interface ArticleCreate {
  title: string
  content?: string
  coverImage?: string
  categoryId: number
  tagIds?: number[]
}

export interface ArticleUpdate {
  title?: string
  content?: string
  coverImage?: string
  categoryId?: number
  tagIds?: number[]
}

export interface ArticleStatusUpdate {
  status: number
}

export interface Comment {
  id: number
  parentId: number
  articleId: number
  replyUserId: number
  userId: number
  content: string
  createdAt: string
  user?: UserResponse
  replyUser?: UserResponse
}
