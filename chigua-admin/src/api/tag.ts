import request from './request'
import type { Tag, PageResponse, BaseResponse } from '@/types'

export const getTagList = async (params: { page: number; pageSize: number; keyword?: string }): Promise<BaseResponse<PageResponse<Tag>>> => {
  const response = await request.get('/admin/tag', { params })
  return response.data
}

export const createTag = async (data: { name: string }): Promise<BaseResponse<Tag>> => {
  const response = await request.post('/admin/tag', data)
  return response.data
}

export const updateTag = async (id: number, data: { name: string }): Promise<BaseResponse<Tag>> => {
  const response = await request.put(`/admin/tag/${id}`, data)
  return response.data
}

export const deleteTag = async (id: number): Promise<BaseResponse<void>> => {
  const response = await request.delete(`/admin/tag/${id}`)
  return response.data
}
