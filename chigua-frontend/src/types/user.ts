// 用户信息类型（对齐后端 UserResponse）
export interface User {
  id: number
  username: string
  nickname: string
  avatar: string
  role?: string
  roles?: string[]
  permissions?: string[]
  createdAt?: string
  updateAt?: string
}
