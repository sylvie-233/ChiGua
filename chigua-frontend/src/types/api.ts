// API响应类型定义
export interface ApiResponse<T = any> {
  code: number
  msg: string
  data: T
}

// 分页响应类型
export interface PageResponse {
  total: number
  page: number
  pageSize: number
  totalPages: number
  hasNext: boolean
  hasPrev: boolean
}
