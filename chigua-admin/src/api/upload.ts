import request from './request'
import type { BaseResponse } from '@/types'

export interface UploadFileResult {
  fileName: string
  filePath: string
  fileUrl: string
}

/** 上传单个文件，返回文件信息 */
export const uploadFile = async (file: File): Promise<BaseResponse<UploadFileResult>> => {
  const formData = new FormData()
  formData.append('file', file)
  const response = await request.post('/upload/file', formData, {
    headers: { 'Content-Type': 'multipart/form-data' }
  })
  return response.data
}
