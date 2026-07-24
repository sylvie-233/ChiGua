<script setup lang="ts">
import { ref } from 'vue'
import { message } from 'ant-design-vue'
import { EyeOutlined } from '@ant-design/icons-vue'
import { MdPreview } from 'md-editor-v3'
import 'md-editor-v3/lib/preview.css'
import { getPendingReviewArticles, getArticleById, approveArticle, rejectArticle } from '@/api/article'
import { formatDate } from '@/utils/date'
import { createPagination, zebraRow, emptyText } from '@/utils/table'
import type { ArticleResponse } from '@/types'

const loading = ref(false)
const articles = ref<ArticleResponse[]>([])
const pagination = createPagination()

const columns = [
  { title: '序号', dataIndex: 'index', key: 'index', width: 70, align: 'center' as const },
  { title: '标题', dataIndex: 'title', key: 'title', ellipsis: true, minWidth: 200 },
  { title: '作者', dataIndex: 'author', key: 'author', width: 120, align: 'center' as const },
  { title: '分类', dataIndex: 'category', key: 'category', width: 120, align: 'center' as const },
  { title: '提交时间', dataIndex: 'submittedAt', key: 'submittedAt', width: 170, align: 'center' as const },
  { title: '操作', key: 'actions', width: 200, align: 'center' as const }
]

const fetchData = async () => {
  loading.value = true
  try {
    const res = await getPendingReviewArticles({
      page: pagination.current,
      pageSize: pagination.pageSize
    })
    if (res.code === 200) {
      articles.value = res.data.items
      pagination.total = res.data.total
    }
  } catch {
    message.error('获取待审核列表失败')
  } finally {
    loading.value = false
  }
}

const handleTableChange = (paginationInfo: { current: number; pageSize: number }) => {
  pagination.current = paginationInfo.current
  pagination.pageSize = paginationInfo.pageSize
  fetchData()
}

// ====== 预览抽屉 ======
const previewVisible = ref(false)
const previewArticle = ref<ArticleResponse | null>(null)
const previewLoading = ref(false)

const handlePreview = async (record: ArticleResponse) => {
  previewVisible.value = true
  previewLoading.value = true
  try {
    const res = await getArticleById(record.id)
    if (res.code === 200) {
      previewArticle.value = res.data
    }
  } catch {
    message.error('加载文章详情失败')
  } finally {
    previewLoading.value = false
  }
}

// ====== 审核操作 ======
const rejectModalVisible = ref(false)
const rejectComment = ref('')
const currentRecord = ref<ArticleResponse | null>(null)
const approving = ref(false)
const rejecting = ref(false)

// 前端校验弹窗
const validateModalVisible = ref(false)
const validateErrors = ref<string[]>([])

const validateArticleFields = (article: ArticleResponse): string[] => {
  const errors: string[] = []
  if (!article.title.trim() || article.title === '无标题') {
    errors.push('标题不能为空')
  }
  if (!article.content.trim()) {
    errors.push('内容不能为空')
  }
  if (!article.categoryId) {
    errors.push('请选择分类')
  }
  const covers = article.coverImage ? article.coverImage.split(',').filter(u => u.trim()) : []
  if (covers.length !== 3) {
    errors.push('请上传3张封面图片')
  }
  return errors
}

const handleApprove = async (record: ArticleResponse) => {
  approving.value = true
  try {
    // 前端校验：先获取文章详情检查必填项
    const articleRes = await getArticleById(record.id)
    if (articleRes.code !== 200) {
      message.error('获取文章信息失败')
      return
    }
    const errors = validateArticleFields(articleRes.data)
    if (errors.length > 0) {
      validateErrors.value = errors
      validateModalVisible.value = true
      return
    }

    const res = await approveArticle(record.id)
    if (res.code === 200) {
      message.success('审核通过')
      if (previewArticle.value?.id === record.id) {
        previewVisible.value = false
      }
      fetchData()
    } else {
      message.error(res.msg || '操作失败')
    }
  } catch {
    message.error('操作失败，请稍后重试')
  } finally {
    approving.value = false
  }
}

const handleReject = (record: ArticleResponse) => {
  currentRecord.value = record
  rejectComment.value = ''
  rejectModalVisible.value = true
}

const handleRejectConfirm = async () => {
  if (!rejectComment.value.trim()) {
    message.warning('请输入驳回理由')
    return
  }
  if (!currentRecord.value) return

  rejecting.value = true
  try {
    const res = await rejectArticle(currentRecord.value.id, rejectComment.value)
    if (res.code === 200) {
      message.success('已驳回')
      rejectModalVisible.value = false
      if (previewArticle.value?.id === currentRecord.value.id) {
        previewVisible.value = false
      }
      fetchData()
    } else {
      message.error(res.msg || '操作失败')
    }
  } catch {
    message.error('操作失败，请稍后重试')
  } finally {
    rejecting.value = false
  }
}

fetchData()
</script>

<template>
  <div>
    <a-card :body-style="{ padding: 0 }">
      <a-table
        bordered
        size="middle"
        :columns="columns"
        :data-source="articles"
        :pagination="pagination"
        :loading="loading"
        :row-class-name="zebraRow"
        :locale="{ emptyText: emptyText || '暂无待审核文章' }"
        row-key="id"
        @change="handleTableChange"
      >
        <template #bodyCell="{ column, record, index }">
          <template v-if="column.key === 'index'">
            {{ (pagination.current - 1) * pagination.pageSize + index + 1 }}
          </template>
          <template v-if="column.key === 'title'">
            <a @click="handlePreview(record)">{{ record.title }}</a>
          </template>
          <template v-if="column.key === 'author'">
            {{ record.author?.nickname || '-' }}
          </template>
          <template v-if="column.key === 'category'">
            {{ record.category?.name || '-' }}
          </template>
          <template v-if="column.key === 'submittedAt'">
            {{ record.submittedAt ? formatDate(record.submittedAt) : '-' }}
          </template>
          <template v-if="column.key === 'actions'">
            <a-space>
              <a-button size="small" @click="handlePreview(record)">
                <EyeOutlined /> 预览
              </a-button>
              <a-button type="primary" size="small" :loading="approving" @click="handleApprove(record)" v-permission="'article:approve'">
                通过
              </a-button>
              <a-button size="small" danger :loading="rejecting" @click="handleReject(record)" v-permission="'article:reject'">
                驳回
              </a-button>
            </a-space>
          </template>
        </template>
      </a-table>
    </a-card>

    <!-- 驳回弹窗 -->
    <a-modal v-model:open="rejectModalVisible" title="驳回意见" @ok="handleRejectConfirm">
      <a-textarea v-model:value="rejectComment" placeholder="请输入驳回原因..." :rows="4" />
    </a-modal>

    <!-- 校验失败弹窗 -->
    <a-modal v-model:open="validateModalVisible" title="无法通过审核" :footer="null">
      <p style="margin-bottom: 8px; color: #666;">以下必填项未完善：</p>
      <ul style="padding-left: 20px; color: #ff4d4f;">
        <li v-for="err in validateErrors" :key="err">{{ err }}</li>
      </ul>
      <div style="margin-top: 16px; text-align: right;">
        <a-button type="primary" @click="validateModalVisible = false">知道了</a-button>
      </div>
    </a-modal>

    <!-- 文章预览抽屉 -->
    <a-drawer
      v-model:open="previewVisible"
      title="文章预览"
      :width="720"
      :footer-style="{ textAlign: 'right' }"
    >
      <template #footer>
        <a-space>
          <a-button size="large" danger :loading="rejecting" @click="handleReject(previewArticle!)" v-if="previewArticle" v-permission="'article:reject'">驳回</a-button>
          <a-button size="large" type="primary" :loading="approving" @click="handleApprove(previewArticle!)" v-if="previewArticle" v-permission="'article:approve'">审核通过</a-button>
        </a-space>
      </template>

      <a-spin :spinning="previewLoading">
        <template v-if="previewArticle">
          <h2>{{ previewArticle.title }}</h2>
          <a-descriptions :column="2" size="small" bordered style="margin-bottom: 16px;">
            <a-descriptions-item label="作者">{{ previewArticle.author?.nickname || previewArticle.author?.username }}</a-descriptions-item>
            <a-descriptions-item label="分类">{{ previewArticle.category?.name }}</a-descriptions-item>
            <a-descriptions-item label="标签">
              <a-tag v-for="tag in previewArticle.tags" :key="tag.id">{{ tag.name }}</a-tag>
              <span v-if="!previewArticle.tags?.length">-</span>
            </a-descriptions-item>
            <a-descriptions-item label="提交时间">{{ previewArticle.submittedAt ? formatDate(previewArticle.submittedAt) : '-' }}</a-descriptions-item>
          </a-descriptions>

          <!-- 封面图 -->
          <div v-if="previewArticle.coverImage" style="display: flex; gap: 8px; margin-bottom: 16px;">
            <img
              v-for="(url, i) in previewArticle.coverImage.split(',').filter(u => u.trim())"
              :key="i"
              :src="url.trim()"
              style="width: 200px; height: 133px; object-fit: cover; border-radius: 4px; border: 1px solid #e8e8e8;"
            />
          </div>

          <!-- 文章内容 -->
          <MdPreview :modelValue="previewArticle.content" />
        </template>
      </a-spin>
    </a-drawer>
  </div>
</template>

<style scoped>
:deep(.ant-table-thead > tr > th) {
  background: #fafafa;
  font-weight: 600;
}
.row-striped > td {
  background: #fafafa;
}

</style>
