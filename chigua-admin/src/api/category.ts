import request from './request'
import type { Category, PageResponse, BaseResponse } from '@/types'

export const getCategoryList = async (params: { page: number; pageSize: number; keyword?: string }): Promise<BaseResponse<PageResponse<Category>>> => {
  const response = await request.get('/admin/category', { params })
  return response.data
}

export const createCategory = async (data: { name: string }): Promise<BaseResponse<Category>> => {
  const response = await request.post('/admin/category', data)
  return response.data
}

export const updateCategory = async (id: number, data: { name: string }): Promise<BaseResponse<Category>> => {
  const response = await request.put(`/admin/category/${id}`, data)
  return response.data
}

export const deleteCategory = async (id: number): Promise<BaseResponse<void>> => {
  const response = await request.delete(`/admin/category/${id}`)
  return response.data
}
