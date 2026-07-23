import request from './request'
import type { UserResponse, PageResponse, BaseResponse } from '@/types'

export interface LoginParams {
  username: string
  password: string
}

export interface LoginResponse {
  token: string
  user: UserResponse
}

export const login = async (data: LoginParams): Promise<BaseResponse<LoginResponse>> => {
  const response = await request.post('/user/login', data)
  return response.data
}

export const getCurrentUser = async (): Promise<BaseResponse<UserResponse>> => {
  const response = await request.get('/user/current')
  return response.data
}

export const getUserList = async (params: { page: number; pageSize: number; keyword?: string }): Promise<BaseResponse<PageResponse<UserResponse>>> => {
  const response = await request.get('/admin/user', { params })
  return response.data
}

export const deleteUser = async (id: number): Promise<BaseResponse<void>> => {
  const response = await request.delete(`/admin/user/${id}`)
  return response.data
}

export const updateUserRole = async (id: number, role: string): Promise<BaseResponse<void>> => {
  const response = await request.put(`/admin/user/${id}/role`, { role })
  return response.data
}

export interface UserCreate {
  username: string
  password: string
  nickname?: string
  role?: string
}

export interface UserUpdate {
  nickname?: string
  role?: string
}

export const createUser = async (data: UserCreate): Promise<BaseResponse<UserResponse>> => {
  const response = await request.post('/admin/user', data)
  return response.data
}

export const updateUser = async (id: number, data: UserUpdate): Promise<BaseResponse<void>> => {
  const response = await request.put(`/admin/user/${id}`, data)
  return response.data
}