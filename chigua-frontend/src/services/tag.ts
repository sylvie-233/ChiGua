import { typedApi } from "./client"
import type { Tag } from "@/types/tag"

export const tagApi = {
  async getAllTags() {
    return typedApi.get<Tag[]>("/tag/list")
  },

  async createTag(data: { name: string }) {
    return typedApi.post<Tag>("/tag", data)
  },

  async updateTag(id: number, data: { name: string }) {
    return typedApi.put<Tag>(`/tag/${id}`, data)
  },

  async deleteTag(id: number) {
    return typedApi.delete<void>(`/tag/${id}`)
  }
}