import api from "./client"

// 用户相关API
export const register = (data: {
  username: string
  password: string
  email: string
}) => api.post("/user/register", data)

export const login = (data: { username: string; password: string }) =>
  api.post("/user/login", data)

export const getCurrentUser = () => api.get("/user/current")
