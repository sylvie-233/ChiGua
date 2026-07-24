<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ArrowLeftOutlined, PlusOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { MdEditor } from 'md-editor-v3'
import 'md-editor-v3/lib/style.css'
import { createArticle, updateArticle, getArticleById, getArticleReviewRecords } from '@/api/article'
import { getCategoryList } from '@/api/category'
import { getTagList } from '@/api/tag'
import { uploadFile } from '@/api/upload'
import type { Category, Tag, ArticleReviewRecord } from '@/types'
import { formatDate } from '@/utils/date'

const router = useRouter()
const route = useRoute()
const loading = ref(false)
const saving = ref(false)
const isEdit = ref(false)
const articleId = ref(0)
const reviewRecords = ref<ArticleReviewRecord[]>([])

const form = reactive({
  title: '',
  content: '',
  categoryId: undefined as number | undefined,
  tagIds: [] as number[]
})

const categories = ref<Category[]>([])
const tags = ref<Tag[]>([])

// 封面图片：固定3个槽位
const coverImages = ref([
  { url: '', uploading: false },
  { url: '', uploading: false },
  { url: '', uploading: false },
])

// 用 ref 数组来触发隐藏的 file input
const fileInputRefs = ref<HTMLInputElement[]>([])

const triggerCoverUpload = (index: number) => {
  fileInputRefs.value[index]?.click()
}

const handleCoverFileChange = async (event: Event, index: number) => {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file) return

  // 设置上传中状态
  coverImages.value = coverImages.value.map((img, i) =>
    i === index ? { ...img, uploading: true } : img
  )
  try {
    const res = await uploadFile(file)
    if (res.code === 200) {
      coverImages.value = coverImages.value.map((img, i) =>
        i === index ? { ...img, url: res.data.fileUrl, uploading: false } : img
      )
      message.success(`封面图${index + 1}上传成功`)
    } else {
      coverImages.value = coverImages.value.map((img, i) =>
        i === index ? { ...img, uploading: false } : img
      )
      message.error(`封面图${index + 1}上传失败`)
    }
  } catch {
    coverImages.value = coverImages.value.map((img, i) =>
      i === index ? { ...img, uploading: false } : img
    )
    message.error(`封面图${index + 1}上传失败`)
  } finally {
    // 清空 input 以允许重复上传同一文件
    input.value = ''
  }
}

const removeCoverImage = (index: number) => {
  coverImages.value = coverImages.value.map((img, i) =>
    i === index ? { ...img, url: '' } : img
  )
}

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
        form.categoryId = article.categoryId
        form.tagIds = article.tags?.map(t => t.id) ?? []
        // 解析已有的封面图片（逗号分隔的URL）
        const urls = article.coverImage
          ? article.coverImage.split(',').filter(u => u.trim())
          : []
        // 整体替换数组以确保 Vue 响应式更新
        coverImages.value = coverImages.value.map((img, i) => ({
          ...img,
          url: i < urls.length ? urls[i].trim() : '',
        }))

        // 加载审核记录
        try {
          const reviewRes = await getArticleReviewRecords(articleId.value)
          if (reviewRes.code === 200) {
            reviewRecords.value = reviewRes.data
          }
        } catch {
          // 审核记录加载失败不阻塞
        }
      }
    } catch {
      message.error('加载文章失败')
    } finally {
      loading.value = false
    }
  }
})

const doSave = async () => {
  saving.value = true
  try {
    const coverImageStr = coverImages.value.map(img => img.url).join(',')

    const data = {
      title: form.title,
      content: form.content || undefined,
      coverImage: coverImageStr || undefined,
      categoryId: form.categoryId!,
      tagIds: form.tagIds.length > 0 ? form.tagIds : undefined
    }

    if (isEdit.value) {
      const res = await updateArticle(articleId.value, data)
      if (res.code !== 200) {
        message.error('保存失败')
        return
      }
    } else {
      const res = await createArticle(data)
      if (res.code !== 200) {
        message.error('保存失败')
        return
      }
      isEdit.value = true
      articleId.value = res.data.id
    }

    message.success('保存成功')
    router.push('/articles')
  } catch {
    message.error('保存失败')
  } finally {
    saving.value = false
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
        <a-button type="primary" :loading="saving" @click="doSave">保存</a-button>
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
            <a-col :span="24">
              <div class="cover-images-section">
                <div class="cover-label">封面图片 <span class="cover-required">（必须上传3张）</span></div>
                <div class="cover-images-grid">
                  <div
                    v-for="(img, index) in coverImages"
                    :key="index"
                    class="cover-image-item"
                  >
                    <!-- 隐藏的原生 file input -->
                    <input
                      :ref="(el: any) => { if (el) fileInputRefs[index] = el }"
                      type="file"
                      accept="image/*"
                      class="cover-file-input-hidden"
                      @change="(e: Event) => handleCoverFileChange(e, index)"
                    />
                    <div v-if="img.url" class="cover-image-preview">
                      <img :src="img.url" :alt="`封面图${index + 1}`" />
                      <div class="cover-image-mask">
                        <a-button size="small" danger @click="removeCoverImage(index)">删除</a-button>
                      </div>
                    </div>
                    <div
                      v-else
                      class="cover-upload-trigger"
                      :class="{ 'is-uploading': img.uploading }"
                      @click="triggerCoverUpload(index)"
                    >
                      <a-spin v-if="img.uploading" size="small" />
                      <PlusOutlined v-else />
                      <div class="cover-upload-text">{{ img.uploading ? '上传中...' : `封面 ${index + 1}` }}</div>
                    </div>
                  </div>
                </div>
              </div>
            </a-col>
          </a-row>
        </div>

        <!-- 审核记录时间线 -->
        <a-row v-if="reviewRecords.length > 0" :gutter="16" style="margin-top: 12px;">
          <a-col :span="24">
            <a-card title="审核记录" size="small">
              <a-timeline>
                <a-timeline-item
                  v-for="record in reviewRecords"
                  :key="record.id"
                  :color="record.action === 'approve' ? 'green' : 'red'"
                >
                  <div>
                    <a-tag :color="record.action === 'approve' ? 'green' : 'red'">
                      {{ record.action === 'approve' ? '审核通过' : '驳回' }}
                    </a-tag>
                    <span style="margin-left: 8px; color: #666;">
                      {{ record.reviewer.nickname || record.reviewer.username }}
                    </span>
                    <span style="margin-left: 8px; color: #999; font-size: 12px;">
                      {{ formatDate(record.createdAt) }}
                    </span>
                  </div>
                  <div v-if="record.comment" style="margin-top: 4px; padding: 8px; background: #f5f5f5; border-radius: 4px;">
                    {{ record.comment }}
                  </div>
                </a-timeline-item>
              </a-timeline>
            </a-card>
          </a-col>
        </a-row>

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

/* 封面图片上传样式 */
.cover-images-section {
  padding: 4px 0;
}

.cover-label {
  font-size: 14px;
  font-weight: 500;
  margin-bottom: 10px;
  color: #333;
}

.cover-required {
  color: #999;
  font-weight: 400;
  font-size: 13px;
}

.cover-images-grid {
  display: flex;
  gap: 12px;
}

.cover-image-item {
  width: 180px;
  height: 120px;
  border-radius: 6px;
  overflow: hidden;
  position: relative;
}

.cover-file-input-hidden {
  position: absolute;
  width: 0;
  height: 0;
  opacity: 0;
  pointer-events: none;
}

.cover-image-preview {
  position: relative;
  width: 100%;
  height: 100%;
  border: 1px solid #e8e8e8;
  border-radius: 6px;
  overflow: hidden;
}

.cover-image-preview img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
}

.cover-image-mask {
  position: absolute;
  inset: 0;
  background: rgba(0, 0, 0, 0.45);
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0;
  transition: opacity 0.2s;
}

.cover-image-preview:hover .cover-image-mask {
  opacity: 1;
}

.cover-upload-trigger {
  width: 100%;
  height: 100%;
  min-height: 120px;
  border: 1px dashed #d9d9d9;
  border-radius: 6px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: border-color 0.2s;
  background: #fafafa;
  gap: 6px;
  color: #999;
  font-size: 22px;
}

.cover-upload-trigger:hover {
  border-color: #1677ff;
  color: #1677ff;
}

.cover-upload-trigger.is-uploading {
  cursor: not-allowed;
  border-color: #d9d9d9;
  color: #999;
}

.cover-upload-text {
  font-size: 13px;
}
</style>
