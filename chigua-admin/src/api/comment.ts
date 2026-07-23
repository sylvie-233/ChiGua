import request from './request'
import type { Comment, PageResponse, BaseResponse } from '@/types'

export const getCommentList = async (params: { page: number; pageSize: number; keyword?: string }): Promise<BaseResponse<PageResponse<Comment>>> => {
  const response = await request.get('/admin/comment', { params })
  return response.data
}

export const deleteComment = async (id: number): Promise<BaseResponse<void>> => {
  const response = await request.delete(`/admin/comment/${id}`)
  return response.data
}
