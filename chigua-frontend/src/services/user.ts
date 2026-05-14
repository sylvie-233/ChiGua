import api from "./client"
import type { ApiResponse } from "@/types/api"
import type { User } from "@/types/user"

export interface LoginResponse {
  token: string
  user: User
}

export const userApi = {
  register(data: {
    username: string
    password: string
    nickname: string
  }): Promise<ApiResponse<User>> {
    return api.post("/user/register", data)
  },

  login(data: {
    username: string
    password: string
  }): Promise<ApiResponse<LoginResponse>> {
    return api.post("/user/login", data)
  },

  getProfile(): Promise<ApiResponse<User>> {
    return api.get("/user/me")
  },

  updateProfile(data: Partial<User>): Promise<ApiResponse<string>> {
    return api.put("/user/me", data)
  }
}
