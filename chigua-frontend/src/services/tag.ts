import api from "./client"

// 标签相关API
export const createTag = (data: { name: string }) => api.post("/tag", data)

export const getTagList = () => api.get("/tag")

export const deleteTag = (id: number) => api.delete(`/tag/${id}`)
