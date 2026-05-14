import api from "./client"
import type { ApiResponse } from "@/types/api"
import type { Tag } from "@/types/tag"

export const tagApi = {
  getAllTags() {
    return api.get<ApiResponse<Tag[]>>("/tag/list")
  },

  createTag(data: { name: string }) {
    return api.post<ApiResponse<Tag>>("/tag", data)
  },

  updateTag(id: number, data: { name: string }) {
    return api.put<ApiResponse<Tag>>(`/tag/${id}`, data)
  },

  deleteTag(id: number) {
    return api.delete<ApiResponse<void>>(`/tag/${id}`)
  }
}
