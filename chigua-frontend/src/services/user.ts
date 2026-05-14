import { typedApi } from "./client"
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
  }
}