import { marked } from "marked"

// 配置marked选项
marked.setOptions({
  gfm: true,
  breaks: true
})

export const renderMarkdown = (content: string): string => {
  return marked(content) as string
}
