import { typedApi } from "./client"
import api from "./client"
import type { User } from "@/types/user"

export interface LoginResponse {
  token: string
  user: User
}

export const userApi = {
  async register(data: {
    username: string
    password: string
    nickname: string
  }) {
    return typedApi.post<User>("/user/register", data)
  },

  async login(data: {
    username: string
    password: string
  }) {
    return typedApi.post<LoginResponse>("/user/login", data)
  },

  async getProfile() {
    return typedApi.get<User>("/user/me")
  },

  async updateProfile(data: Partial<User>) {
    return typedApi.put<string>("/user/me", data)
  },

  async updateAvatar(avatar: string) {
    return typedApi.put<string>("/user/avatar", { avatar })
  },

  async uploadFile(file: File) {
    const formData = new FormData()
    formData.append("file", file)
    const response = await api.post("/upload/file", formData)
    return response as any
  }
}