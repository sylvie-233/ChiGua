import request from './request'
import type { Article, ArticleResponse, ArticleCreate, ArticleUpdate, PageResponse, BaseResponse } from '@/types'

export const getArticleList = async (params: { page: number; pageSize: number; keyword?: string }): Promise<BaseResponse<PageResponse<ArticleResponse>>> => {
  const response = await request.get('/admin/article', { params })
  return response.data
}

export const getArticleById = async (id: number): Promise<BaseResponse<ArticleResponse>> => {
  const response = await request.get(`/admin/article/${id}`)
  return response.data
}

export const createArticle = async (data: ArticleCreate): Promise<BaseResponse<Article>> => {
  const response = await request.post('/admin/article', data)
  return response.data
}

export const updateArticle = async (id: number, data: ArticleUpdate): Promise<BaseResponse<Article>> => {
  const response = await request.put(`/admin/article/${id}`, data)
  return response.data
}

export const deleteArticle = async (id: number): Promise<BaseResponse<void>> => {
  const response = await request.delete(`/admin/article/${id}`)
  return response.data
}

export const publishArticle = async (id: number): Promise<BaseResponse<void>> => {
  const response = await request.post(`/admin/article/${id}/publish`)
  return response.data
}

export const updateArticleStatus = async (id: number, status: number): Promise<BaseResponse<void>> => {
  const response = await request.put(`/admin/article/${id}/status`, { status })
  return response.data
}
