import request from './request'
import type { BaseResponse, PageResponse, Role } from '@/types'

export const getRoles = async (params?: { page: number; pageSize: number; keyword?: string }): Promise<BaseResponse<PageResponse<Role>>> => {
  const response = await request.get('/admin/role', { params })
  return response.data
}

export const createRole = async (data: { code: string; name: string; description?: string }): Promise<BaseResponse<Role>> => {
  const response = await request.post('/admin/role', data)
  return response.data
}

export const updateRole = async (id: number, data: { name: string; description?: string }): Promise<BaseResponse<void>> => {
  const response = await request.put(`/admin/role/${id}`, data)
  return response.data
}

export const deleteRole = async (id: number): Promise<BaseResponse<void>> => {
  const response = await request.delete(`/admin/role/${id}`)
  return response.data
}

export const getRoleMenuIds = async (roleId: number): Promise<BaseResponse<number[]>> => {
  const response = await request.get(`/admin/role/${roleId}/menus`)
  return response.data
}

export const updateRoleMenus = async (roleId: number, menuIds: number[]): Promise<BaseResponse<void>> => {
  const response = await request.put(`/admin/role/${roleId}/menus`, { menuIds })
  return response.data
}

export const updateUserRoles = async (userId: number, roleIds: number[]): Promise<BaseResponse<void>> => {
  const response = await request.put(`/admin/user/${userId}/roles`, { roleIds })
  return response.data
}
