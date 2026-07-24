import request from './request'
import type { BaseResponse } from '@/types'

export interface MenuItem {
  id: number
  parentId: number
  title: string
  path: string
  component: string
  icon: string
  permissionCode: string
  sortOrder: number
  visible: boolean
  menuType: string
  children?: MenuItem[]
}

export const getMenuTree = async (): Promise<BaseResponse<MenuItem[]>> => {
  const response = await request.get('/admin/menu/tree')
  return response.data
}

export const getAllMenus = async (): Promise<BaseResponse<MenuItem[]>> => {
  const response = await request.get('/admin/menu')
  return response.data
}

export const createMenu = async (data: Partial<MenuItem>): Promise<BaseResponse<MenuItem>> => {
  const response = await request.post('/admin/menu', data)
  return response.data
}

export const updateMenu = async (id: number, data: Partial<MenuItem>): Promise<BaseResponse<void>> => {
  const response = await request.put(`/admin/menu/${id}`, data)
  return response.data
}

export const deleteMenu = async (id: number): Promise<BaseResponse<void>> => {
  const response = await request.delete(`/admin/menu/${id}`)
  return response.data
}
