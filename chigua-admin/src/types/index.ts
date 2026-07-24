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
  updateAt: string
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
  reviewerId: number
  reviewComment: string
  submittedAt: string | null
  publishAt: string
  createdAt: string
  updateAt: string
}

export interface ArticleResponse extends Article {
  tags: Tag[]
  category: Category
  author: UserResponse
  reviewer?: UserResponse
}

export interface ArticleReviewRecord {
  id: number
  articleId: number
  action: 'approve' | 'reject'
  comment: string
  createdAt: string
  reviewer: UserResponse
}

/** 审核记录列表项（含文章标题） */
export interface ArticleReviewRecordWithTitle extends ArticleReviewRecord {
  articleTitle: string
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
