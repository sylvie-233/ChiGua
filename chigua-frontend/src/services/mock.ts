import type { Article, User, Category, Tag, Comment } from "@/types/api"

// Mock数据开关
export const USE_MOCK = true

// Mock用户数据
export const mockUser: User = {
  id: 1,
  username: "zhangsan",
  nickname: "张三",
  email: "zhangsan@example.com",
  created_at: "2024-01-01 10:00:00",
  updated_at: "2024-01-01 10:00:00"
}

// Mock分类数据
export const mockCategories: Category[] = [
  { id: 1, name: "科技", created_at: "2024-01-01", updated_at: "2024-01-01" },
  { id: 2, name: "娱乐", created_at: "2024-01-01", updated_at: "2024-01-01" },
  { id: 3, name: "体育", created_at: "2024-01-01", updated_at: "2024-01-01" },
  { id: 4, name: "财经", created_at: "2024-01-01", updated_at: "2024-01-01" }
]

// Mock标签数据
export const mockTags: Tag[] = [
  { id: 1, name: "Vue", created_at: "2024-01-01", updated_at: "2024-01-01" },
  {
    id: 2,
    name: "TypeScript",
    created_at: "2024-01-01",
    updated_at: "2024-01-01"
  },
  { id: 3, name: "前端", created_at: "2024-01-01", updated_at: "2024-01-01" },
  { id: 4, name: "React", created_at: "2024-01-01", updated_at: "2024-01-01" }
]

// Mock文章数据
export const mockArticles: Article[] = [
  {
    id: 1,
    title: "Vue 3 组合式API详解",
    content: `# Vue 3 组合式API详解

## 什么是组合式API

组合式API是Vue 3中引入的一种新的代码组织方式，它允许我们使用函数来组织组件逻辑，而不是选项对象。

## 为什么使用组合式API

组合式API提供了更好的代码组织和复用能力，使得逻辑关注点可以被分组在一起，而不是分散在不同的选项中。

## ref和reactive

在组合式API中，我们使用ref来创建响应式的基本类型数据，使用reactive来创建响应式的对象。

\`\`\`typescript
import { ref, reactive } from 'vue'

const count = ref(0)
const state = reactive({
  name: '张三',
  age: 25
})
\`\`\`

## computed和watch

computed用于创建计算属性，watch用于监听响应式数据的变化。

\`\`\`typescript
import { computed, watch } from 'vue'

const doubled = computed(() => count.value * 2)

watch(count, (newVal, oldVal) => {
  console.log(\`count changed from \${oldVal} to \${newVal}\`)
})
\`\`\``,
    status: 1,
    author_id: 1,
    author: mockUser,
    category_id: 1,
    category: mockCategories[0],
    tags: [mockTags[0], mockTags[1], mockTags[2]],
    created_at: "2024-01-15 14:30:00",
    updated_at: "2024-01-15 14:30:00"
  },
  {
    id: 2,
    title: "TypeScript 高级类型技巧",
    content: `# TypeScript 高级类型技巧

## 泛型

泛型允许我们编写可重用的组件，同时保持类型安全。

\`\`\`typescript
function identity<T>(arg: T): T {
  return arg
}
\`\`\`

## 条件类型

条件类型允许我们根据类型关系来选择类型。

\`\`\`typescript
type IsString<T> = T extends string ? true : false
\`\`\`

## 映射类型

映射类型允许我们基于旧类型创建新类型。

\`\`\`typescript
type Readonly<T> = {
  readonly [P in keyof T]: T[P]
}
\`\`\``,
    status: 1,
    author_id: 1,
    author: mockUser,
    category_id: 1,
    category: mockCategories[0],
    tags: [mockTags[1], mockTags[2]],
    created_at: "2024-01-14 10:00:00",
    updated_at: "2024-01-14 10:00:00"
  },
  {
    id: 3,
    title: "前端性能优化指南",
    content: `# 前端性能优化指南

## 代码分割

使用代码分割可以减小首屏加载时间。

## 图片优化

使用适当的图片格式和大小可以显著提升加载速度。

## 缓存策略

合理使用浏览器缓存可以减少重复请求。`,
    status: 1,
    author_id: 1,
    author: mockUser,
    category_id: 1,
    category: mockCategories[0],
    tags: [mockTags[2]],
    created_at: "2024-01-13 09:00:00",
    updated_at: "2024-01-13 09:00:00"
  },
  {
    id: 4,
    title: "React vs Vue 对比分析",
    content: `# React vs Vue 对比分析

## 学习曲线

Vue的学习曲线相对平缓，而React需要更多的JavaScript知识。

## 生态系统

React拥有更庞大的生态系统和社区支持。

## 性能表现

两者在性能上都表现出色，具体取决于使用方式。`,
    status: 1,
    author_id: 1,
    author: mockUser,
    category_id: 1,
    category: mockCategories[0],
    tags: [mockTags[0], mockTags[3]],
    created_at: "2024-01-12 16:00:00",
    updated_at: "2024-01-12 16:00:00"
  }
]

// Mock评论数据
export const mockComments: Comment[] = [
  {
    id: 1,
    article_id: 1,
    user_id: 1,
    user: mockUser,
    content: "这篇文章写得非常好，学到了很多！",
    created_at: "2024-01-15 15:00:00",
    updated_at: "2024-01-15 15:00:00"
  },
  {
    id: 2,
    article_id: 1,
    user_id: 1,
    user: mockUser,
    content: "期待更多这样的文章！",
    created_at: "2024-01-15 16:00:00",
    updated_at: "2024-01-15 16:00:00"
  }
]

// Mock响应包装函数
export const mockSuccess = <T>(data: T) => ({
  code: 200,
  msg: "成功",
  data
})

export const mockError = (msg: string = "失败") => ({
  code: 500,
  msg,
  data: null
})
