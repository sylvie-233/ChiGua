import request from './request'
import type { Article, ArticleResponse, ArticleCreate, ArticleUpdate, ArticleReviewRecord, ArticleReviewRecordWithTitle, PageResponse, BaseResponse } from '@/types'

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

export const submitForReview = async (id: number): Promise<BaseResponse<void>> => {
  const response = await request.post(`/article/${id}/submit`)
  return response.data
}

export const approveArticle = async (id: number): Promise<BaseResponse<void>> => {
  const response = await request.post(`/admin/article/${id}/approve`)
  return response.data
}

export const rejectArticle = async (id: number, comment: string): Promise<BaseResponse<void>> => {
  const response = await request.post(`/admin/article/${id}/reject`, { comment })
  return response.data
}

export const unpublishArticle = async (id: number, comment: string): Promise<BaseResponse<void>> => {
  const response = await request.post(`/admin/article/${id}/unpublish`, { comment })
  return response.data
}

export const getArticleReviewRecords = async (id: number): Promise<BaseResponse<ArticleReviewRecord[]>> => {
  const response = await request.get(`/article/${id}/reviews`)
  return response.data
}

export const getPendingReviewArticles = async (params: { page: number; pageSize: number }): Promise<BaseResponse<PageResponse<ArticleResponse>>> => {
  const response = await request.get('/admin/article/pending', { params })
  return response.data
}

export const getAllReviewRecords = async (params: { page: number; pageSize: number }): Promise<BaseResponse<PageResponse<ArticleReviewRecordWithTitle>>> => {
  const response = await request.get('/admin/article/reviews', { params })
  return response.data
}
