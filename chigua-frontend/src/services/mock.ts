import type { Article } from "@/types/article"
import type { Category } from "@/types/category"
import type { Tag } from "@/types/tag"
import type { Comment } from "@/types/comment"

export const mockArticles: Article[] = [
  {
    id: 1,
    title: "探索人工智能的未来发展趋势",
    content: "人工智能正在改变我们的生活方式...",
    status: 1,
    author_id: 1,
    author: {
      id: 1,
      username: "admin",
      nickname: "管理员",
      email: "admin@example.com",
      created_at: "2024-01-01T00:00:00Z",
      updated_at: "2024-01-01T00:00:00Z"
    },
    category_id: 1,
    category: {
      id: 1,
      name: "科技",
      created_at: "2024-01-01T00:00:00Z",
      updated_at: "2024-01-01T00:00:00Z"
    },
    tags: [
      {
        id: 1,
        name: "AI",
        created_at: "2024-01-01T00:00:00Z",
        updated_at: "2024-01-01T00:00:00Z"
      },
      {
        id: 2,
        name: "技术",
        created_at: "2024-01-01T00:00:00Z",
        updated_at: "2024-01-01T00:00:00Z"
      }
    ],
    created_at: "2024-01-15T10:00:00Z",
    updated_at: "2024-01-15T10:00:00Z"
  },
  {
    id: 2,
    title: "前端开发最佳实践",
    content: "现代前端开发需要掌握多种技术栈...",
    status: 1,
    author_id: 1,
    author: {
      id: 1,
      username: "admin",
      nickname: "管理员",
      email: "admin@example.com",
      created_at: "2024-01-01T00:00:00Z",
      updated_at: "2024-01-01T00:00:00Z"
    },
    category_id: 2,
    category: {
      id: 2,
      name: "前端",
      created_at: "2024-01-01T00:00:00Z",
      updated_at: "2024-01-01T00:00:00Z"
    },
    tags: [
      {
        id: 3,
        name: "Vue",
        created_at: "2024-01-01T00:00:00Z",
        updated_at: "2024-01-01T00:00:00Z"
      },
      {
        id: 4,
        name: "TypeScript",
        created_at: "2024-01-01T00:00:00Z",
        updated_at: "2024-01-01T00:00:00Z"
      }
    ],
    created_at: "2024-01-14T08:00:00Z",
    updated_at: "2024-01-14T08:00:00Z"
  },
  {
    id: 3,
    title: "云计算架构设计",
    content: "云原生架构正在成为企业级应用的新标准...",
    status: 1,
    author_id: 1,
    author: {
      id: 1,
      username: "admin",
      nickname: "管理员",
      email: "admin@example.com",
      created_at: "2024-01-01T00:00:00Z",
      updated_at: "2024-01-01T00:00:00Z"
    },
    category_id: 3,
    category: {
      id: 3,
      name: "后端",
      created_at: "2024-01-01T00:00:00Z",
      updated_at: "2024-01-01T00:00:00Z"
    },
    tags: [
      {
        id: 5,
        name: "Cloud",
        created_at: "2024-01-01T00:00:00Z",
        updated_at: "2024-01-01T00:00:00Z"
      },
      {
        id: 6,
        name: "Go",
        created_at: "2024-01-01T00:00:00Z",
        updated_at: "2024-01-01T00:00:00Z"
      }
    ],
    created_at: "2024-01-13T14:00:00Z",
    updated_at: "2024-01-13T14:00:00Z"
  },
  {
    id: 4,
    title: "数据库性能优化技巧",
    content: "数据库优化是提升应用性能的关键...",
    status: 1,
    author_id: 1,
    author: {
      id: 1,
      username: "admin",
      nickname: "管理员",
      email: "admin@example.com",
      created_at: "2024-01-01T00:00:00Z",
      updated_at: "2024-01-01T00:00:00Z"
    },
    category_id: 3,
    category: {
      id: 3,
      name: "后端",
      created_at: "2024-01-01T00:00:00Z",
      updated_at: "2024-01-01T00:00:00Z"
    },
    tags: [
      {
        id: 7,
        name: "PostgreSQL",
        created_at: "2024-01-01T00:00:00Z",
        updated_at: "2024-01-01T00:00:00Z"
      }
    ],
    created_at: "2024-01-12T09:00:00Z",
    updated_at: "2024-01-12T09:00:00Z"
  },
  {
    id: 5,
    title: "微服务架构实践",
    content: "微服务架构带来了更高的可扩展性...",
    status: 1,
    author_id: 1,
    author: {
      id: 1,
      username: "admin",
      nickname: "管理员",
      email: "admin@example.com",
      created_at: "2024-01-01T00:00:00Z",
      updated_at: "2024-01-01T00:00:00Z"
    },
    category_id: 3,
    category: {
      id: 3,
      name: "后端",
      created_at: "2024-01-01T00:00:00Z",
      updated_at: "2024-01-01T00:00:00Z"
    },
    tags: [
      {
        id: 8,
        name: "Microservices",
        created_at: "2024-01-01T00:00:00Z",
        updated_at: "2024-01-01T00:00:00Z"
      }
    ],
    created_at: "2024-01-11T11:00:00Z",
    updated_at: "2024-01-11T11:00:00Z"
  },
  {
    id: 6,
    title: "React vs Vue：选择哪个框架",
    content: "两大主流前端框架的对比分析...",
    status: 1,
    author_id: 1,
    author: {
      id: 1,
      username: "admin",
      nickname: "管理员",
      email: "admin@example.com",
      created_at: "2024-01-01T00:00:00Z",
      updated_at: "2024-01-01T00:00:00Z"
    },
    category_id: 2,
    category: {
      id: 2,
      name: "前端",
      created_at: "2024-01-01T00:00:00Z",
      updated_at: "2024-01-01T00:00:00Z"
    },
    tags: [
      {
        id: 3,
        name: "Vue",
        created_at: "2024-01-01T00:00:00Z",
        updated_at: "2024-01-01T00:00:00Z"
      },
      {
        id: 9,
        name: "React",
        created_at: "2024-01-01T00:00:00Z",
        updated_at: "2024-01-01T00:00:00Z"
      }
    ],
    created_at: "2024-01-10T16:00:00Z",
    updated_at: "2024-01-10T16:00:00Z"
  }
]

export const mockCategories: Category[] = [
  {
    id: 1,
    name: "科技",
    created_at: "2024-01-01T00:00:00Z",
    updated_at: "2024-01-01T00:00:00Z"
  },
  {
    id: 2,
    name: "前端",
    created_at: "2024-01-01T00:00:00Z",
    updated_at: "2024-01-01T00:00:00Z"
  },
  {
    id: 3,
    name: "后端",
    created_at: "2024-01-01T00:00:00Z",
    updated_at: "2024-01-01T00:00:00Z"
  },
  {
    id: 4,
    name: "生活",
    created_at: "2024-01-01T00:00:00Z",
    updated_at: "2024-01-01T00:00:00Z"
  }
]

export const mockTags: Tag[] = [
  {
    id: 1,
    name: "AI",
    created_at: "2024-01-01T00:00:00Z",
    updated_at: "2024-01-01T00:00:00Z"
  },
  {
    id: 2,
    name: "技术",
    created_at: "2024-01-01T00:00:00Z",
    updated_at: "2024-01-01T00:00:00Z"
  },
  {
    id: 3,
    name: "Vue",
    created_at: "2024-01-01T00:00:00Z",
    updated_at: "2024-01-01T00:00:00Z"
  },
  {
    id: 4,
    name: "TypeScript",
    created_at: "2024-01-01T00:00:00Z",
    updated_at: "2024-01-01T00:00:00Z"
  },
  {
    id: 5,
    name: "Cloud",
    created_at: "2024-01-01T00:00:00Z",
    updated_at: "2024-01-01T00:00:00Z"
  },
  {
    id: 6,
    name: "Go",
    created_at: "2024-01-01T00:00:00Z",
    updated_at: "2024-01-01T00:00:00Z"
  },
  {
    id: 7,
    name: "PostgreSQL",
    created_at: "2024-01-01T00:00:00Z",
    updated_at: "2024-01-01T00:00:00Z"
  },
  {
    id: 8,
    name: "Microservices",
    created_at: "2024-01-01T00:00:00Z",
    updated_at: "2024-01-01T00:00:00Z"
  },
  {
    id: 9,
    name: "React",
    created_at: "2024-01-01T00:00:00Z",
    updated_at: "2024-01-01T00:00:00Z"
  }
]

export const mockComments: Comment[] = [
  {
    id: 1,
    article_id: 1,
    user_id: 1,
    user: {
      id: 1,
      username: "admin",
      nickname: "管理员",
      email: "admin@example.com",
      created_at: "2024-01-01T00:00:00Z",
      updated_at: "2024-01-01T00:00:00Z"
    },
    content: "这篇文章很有见解！",
    created_at: "2024-01-15T12:00:00Z",
    updated_at: "2024-01-15T12:00:00Z"
  },
  {
    id: 2,
    article_id: 1,
    user_id: 1,
    user: {
      id: 1,
      username: "admin",
      nickname: "管理员",
      email: "admin@example.com",
      created_at: "2024-01-01T00:00:00Z",
      updated_at: "2024-01-01T00:00:00Z"
    },
    content: "期待更多相关内容",
    created_at: "2024-01-15T13:00:00Z",
    updated_at: "2024-01-15T13:00:00Z"
  }
]
