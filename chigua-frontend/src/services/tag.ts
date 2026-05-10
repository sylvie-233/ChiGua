import api from "./client"
import type { ApiResponse, Tag } from "@/types/api"
import { USE_MOCK, mockTags, mockSuccess } from "./mock"

// 标签相关API
export const createTag = (data: {
  name: string
}): Promise<ApiResponse<Tag>> => {
  if (USE_MOCK) {
    return Promise.resolve(
      mockSuccess({
        id: Date.now(),
        name: data.name,
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString()
      })
    )
  }
  return api.post("/tag", data)
}

export const getTagList = (): Promise<ApiResponse<Tag[]>> => {
  if (USE_MOCK) {
    return Promise.resolve(mockSuccess(mockTags))
  }
  return api.get("/tag")
}

export const deleteTag = (id: number): Promise<ApiResponse<null>> => {
  if (USE_MOCK) {
    return Promise.resolve(mockSuccess(null))
  }
  return api.delete(`/tag/${id}`)
}
