import api from "./client"
import type { ApiResponse, User } from "@/types/api"
import { USE_MOCK, mockUser, mockSuccess } from "./mock"

// 用户相关API
export const register = (data: {
  username: string
  password: string
  email: string
}): Promise<ApiResponse<User>> => {
  if (USE_MOCK) {
    return Promise.resolve(
      mockSuccess({
        id: Date.now(),
        username: data.username,
        nickname: data.username,
        email: data.email,
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString()
      })
    )
  }
  return api.post("/user/register", data)
}

export const login = (data: {
  username: string
  password: string
}): Promise<ApiResponse<{ user: User; token: string }>> => {
  if (USE_MOCK) {
    return Promise.resolve(
      mockSuccess({
        user: mockUser,
        token: "mock-jwt-token"
      })
    )
  }
  return api.post("/user/login", data)
}

export const getCurrentUser = (): Promise<ApiResponse<User>> => {
  if (USE_MOCK) {
    return Promise.resolve(mockSuccess(mockUser))
  }
  return api.get("/user/current")
}
