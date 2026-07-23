<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ArrowLeftOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { MdEditor } from 'md-editor-v3'
import 'md-editor-v3/lib/style.css'
import { createArticle, updateArticle, getArticleById, updateArticleStatus } from '@/api/article'
import { getCategoryList } from '@/api/category'
import { getTagList } from '@/api/tag'
import { uploadFile } from '@/api/upload'
import type { Category, Tag } from '@/types'

const router = useRouter()
const route = useRoute()
const loading = ref(false)
const savingDraft = ref(false)
const publishing = ref(false)
const isEdit = ref(false)
const articleId = ref(0)

const form = reactive({
  title: '',
  content: '',
  coverImage: '',
  categoryId: undefined as number | undefined,
  tagIds: [] as number[]
})

const categories = ref<Category[]>([])
const tags = ref<Tag[]>([])

// 编辑器图片上传
const onUploadImg = async (files: File[], callback: (urls: string[]) => void) => {
  const urls: string[] = []
  for (const file of files) {
    try {
      const res = await uploadFile(file)
      if (res.code === 200) {
        urls.push(res.data.fileUrl)
      }
    } catch {
      message.error(`上传 ${file.name} 失败`)
    }
  }
  callback(urls)
}

onMounted(async () => {
  const [catRes, tagRes] = await Promise.all([
    getCategoryList({ page: 1, pageSize: 100 }),
    getTagList({ page: 1, pageSize: 100 })
  ])
  if (catRes.code === 200) categories.value = catRes.data.items
  if (tagRes.code === 200) tags.value = tagRes.data.items

  const id = route.params.id
  if (id) {
    isEdit.value = true
    articleId.value = Number(id)
    loading.value = true
    try {
      const res = await getArticleById(articleId.value)
      if (res.code === 200) {
        const article = res.data
        form.title = article.title
        form.content = article.content
        form.coverImage = article.coverImage
        form.categoryId = article.categoryId
        form.tagIds = article.tags.map(t => t.id)
      }
    } catch {
      message.error('加载文章失败')
    } finally {
      loading.value = false
    }
  }
})

const validate = (): boolean => {
  if (!form.title.trim()) {
    message.warning('请输入文章标题')
    return false
  }
  if (!form.categoryId) {
    message.warning('请选择分类')
    return false
  }
  return true
}

const doSave = async (andPublish: boolean) => {
  if (!validate()) return

  if (andPublish) {
    publishing.value = true
  } else {
    savingDraft.value = true
  }

  try {
    if (isEdit.value) {
      const data = {
        title: form.title,
        content: form.content,
        coverImage: form.coverImage || undefined,
        categoryId: form.categoryId,
        tagIds: form.tagIds.length > 0 ? form.tagIds : undefined
      }
      const res = await updateArticle(articleId.value, data)
      if (res.code === 200) {
        if (andPublish) {
          await updateArticleStatus(articleId.value, 1)
        }
        message.success(andPublish ? '发布成功' : '草稿已保存')
        router.push('/articles')
      }
    } else {
      const data = {
        title: form.title,
        content: form.content || undefined,
        coverImage: form.coverImage || undefined,
        categoryId: form.categoryId!,
        tagIds: form.tagIds.length > 0 ? form.tagIds : undefined
      }
      const res = await createArticle(data)
      if (res.code === 200) {
        if (andPublish) {
          await updateArticleStatus(res.data.id, 1)
        }
        message.success(andPublish ? '发布成功' : '草稿已保存')
        router.push('/articles')
      }
    }
  } catch {
    message.error('保存失败')
  } finally {
    savingDraft.value = false
    publishing.value = false
  }
}

const handleBack = () => {
  router.push('/articles')
}
</script>

<template>
  <div class="article-editor">
    <!-- 顶部工具栏 -->
    <div class="editor-header">
      <a-space>
        <a-button @click="handleBack">
          <ArrowLeftOutlined /> 返回
        </a-button>
        <span class="editor-title">{{ isEdit ? '编辑文章' : '新建文章' }}</span>
      </a-space>
      <a-space>
        <a-button :loading="savingDraft" @click="doSave(false)">保存草稿</a-button>
        <a-button type="primary" :loading="publishing" @click="doSave(true)">发布文章</a-button>
      </a-space>
    </div>

    <a-spin :spinning="loading">
      <div class="editor-body">
        <!-- 文章信息区 -->
        <div class="editor-meta">
          <a-row :gutter="16">
            <a-col :span="16">
              <a-input
                v-model:value="form.title"
                placeholder="请输入文章标题..."
                size="large"
                class="title-input"
                :bordered="false"
              />
            </a-col>
            <a-col :span="4">
              <a-select
                v-model:value="form.categoryId"
                placeholder="选择分类"
                style="width: 100%"
                size="large"
              >
                <a-select-option v-for="cat in categories" :key="cat.id" :value="cat.id">
                  {{ cat.name }}
                </a-select-option>
              </a-select>
            </a-col>
            <a-col :span="4">
              <a-select
                v-model:value="form.tagIds"
                mode="multiple"
                placeholder="选择标签"
                style="width: 100%"
                size="large"
              >
                <a-select-option v-for="tag in tags" :key="tag.id" :value="tag.id">
                  {{ tag.name }}
                </a-select-option>
              </a-select>
            </a-col>
          </a-row>
          <a-row :gutter="16" style="margin-top: 12px;">
            <a-col :span="16">
              <a-input v-model:value="form.coverImage" placeholder="封面图片 URL（可选）" />
            </a-col>
          </a-row>
        </div>

        <!-- Markdown 编辑器 -->
        <div class="editor-main">
          <MdEditor
            v-model="form.content"
            :toolbars="[
              'bold', 'italic', 'strikeThrough', 'title', '-',
              'unorderedList', 'orderedList', 'task', '-',
              'code', 'quote', 'link', 'image', 'table', '-',
              'revoke', 'next', 'save', 'pageFullscreen', 'fullscreen', '-',
              'preview', 'htmlPreview', 'catalog'
            ]"
            :footers="['markdownTotal']"
            placeholder="开始写作..."
            style="height: calc(100vh - 220px);"
            @on-upload-img="onUploadImg"
          />
        </div>
      </div>
    </a-spin>
  </div>
</template>

<style scoped>
.article-editor {
  background: #fff;
  min-height: 100%;
}

.editor-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 16px;
  height: 48px;
  border-bottom: 1px solid #f0f0f0;
}

.editor-title {
  font-size: 15px;
  font-weight: 500;
}

.editor-body {
  padding: 16px;
}

.editor-meta {
  margin-bottom: 12px;
}

.title-input {
  font-size: 20px;
  font-weight: 600;
  padding: 4px 0;
}

.editor-main {
  border: 1px solid #e8e8e8;
  border-radius: 6px;
  overflow: hidden;
}

/* 覆盖 md-editor 默认样式，使容器自适应 */
.editor-main :deep(.md-editor) {
  height: 100% !important;
}

.editor-main :deep(.md-editor-toolbar) {
  border-radius: 0;
}
</style>
