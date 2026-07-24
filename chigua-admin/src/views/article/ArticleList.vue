<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { PlusOutlined, DownOutlined } from '@ant-design/icons-vue'
import { Modal, message } from 'ant-design-vue'
import { getArticleList, getArticleById, deleteArticle, approveArticle, rejectArticle, submitForReview, unpublishArticle } from '@/api/article'
import { formatDate } from '@/utils/date'
import { createPagination, zebraRow, emptyText } from '@/utils/table'
import type { ArticleResponse } from '@/types'
const router = useRouter()
const searchText = ref('')
const loading = ref(false)
const articles = ref<ArticleResponse[]>([])

const pagination = createPagination()

const statusMap: Record<number, { label: string; color: string }> = {
  0: { label: '草稿', color: 'default' },
  1: { label: '已发布', color: 'green' },
  2: { label: '已下架', color: 'red' },
  3: { label: '审核中', color: 'orange' },
}

const columns = [
  { title: '序号', dataIndex: 'index', key: 'index', width: 70, align: 'center' as const },
  { title: '标题', dataIndex: 'title', key: 'title', ellipsis: true, minWidth: 200 },
  { title: '分类', dataIndex: 'category', key: 'category', width: 120, align: 'center' as const },
  { title: '作者', dataIndex: 'author', key: 'author', width: 100, align: 'center' as const },
  { title: '状态', dataIndex: 'status', key: 'status', width: 100, align: 'center' as const },
  { title: '创建时间', dataIndex: 'createdAt', key: 'createdAt', width: 170, align: 'center' as const },
  { title: '操作', key: 'actions', width: 120, align: 'center' as const }
]

const fetchData = async () => {
  loading.value = true
  try {
    const params: { page: number; pageSize: number; keyword?: string } = {
      page: pagination.current,
      pageSize: pagination.pageSize
    }
    if (searchText.value.trim()) {
      params.keyword = searchText.value.trim()
    }
    const response = await getArticleList(params)
    if (response.code === 200) {
      articles.value = response.data.items
      pagination.total = response.data.total
    }
  } catch (error) {
    console.error('获取文章列表失败:', error)
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.current = 1
  fetchData()
}

const handleTableChange = (paginationInfo: { current: number; pageSize: number }) => {
  pagination.current = paginationInfo.current
  pagination.pageSize = paginationInfo.pageSize
  fetchData()
}

const handleAdd = () => {
  router.push('/articles/new')
}

const handleDelete = async (id: number) => {
  Modal.confirm({
    title: '确认删除',
    content: '确定要删除这篇文章吗？此操作不可撤销。',
    okText: '删除',
    okButtonProps: { danger: true },
    cancelText: '取消',
    async onOk() {
      try {
        const response = await deleteArticle(id)
        if (response.code === 200) {
          message.success('删除成功')
          fetchData()
        } else {
          message.error(response.msg || '删除失败')
        }
      } catch (error) {
        console.error('删除文章失败:', error)
        message.error('删除失败，请稍后重试')
      }
    }
  })
}

// 审核相关
const rejectModalVisible = ref(false)
const rejectComment = ref('')
const currentRejectId = ref(0)

// 下架
const unpublishModalVisible = ref(false)
const unpublishComment = ref('')
const currentUnpublishId = ref(0)

const handleUnpublishClick = (id: number) => {
  currentUnpublishId.value = id
  unpublishComment.value = ''
  unpublishModalVisible.value = true
}

const handleUnpublishConfirm = async () => {
  if (!unpublishComment.value.trim()) {
    message.warning('请输入下架原因')
    return
  }
  try {
    const res = await unpublishArticle(currentUnpublishId.value, unpublishComment.value)
    if (res.code === 200) {
      message.success('已下架')
      unpublishModalVisible.value = false
      fetchData()
    } else {
      message.error(res.msg || '操作失败')
    }
  } catch {
    message.error('操作失败')
  }
}

const handleApprove = async (id: number) => {
  try {
    const res = await approveArticle(id)
    if (res.code === 200) {
      message.success('审核通过')
      fetchData()
    } else {
      message.error(res.msg || '操作失败')
    }
  } catch {
    message.error('操作失败')
  }
}

const handleReject = (id: number) => {
  currentRejectId.value = id
  rejectComment.value = ''
  rejectModalVisible.value = true
}

const handleRejectConfirm = async () => {
  if (!rejectComment.value.trim()) {
    message.warning('请输入驳回理由')
    return
  }
  try {
    const res = await rejectArticle(currentRejectId.value, rejectComment.value)
    if (res.code === 200) {
      message.success('已驳回')
      rejectModalVisible.value = false
      fetchData()
    } else {
      message.error(res.msg || '操作失败')
    }
  } catch {
    message.error('操作失败')
  }
}

const submittingId = ref(0)

// ====== 提交审核：前端校验 ======
const validateModalVisible = ref(false)
const validateErrors = ref<string[]>([])
const validateArticleId = ref(0)

const handleAction = (record: ArticleResponse, key: string) => {
  switch (key) {
    case 'edit':
      router.push(`/articles/${record.id}/edit`)
      break
    case 'delete':
      handleDelete(record.id)
      break
    case 'submit':
      handleSubmitValidate(record.id)
      break
    case 'approve':
      handleApprove(record.id)
      break
    case 'reject':
      handleReject(record.id)
      break
    case 'unpublish':
      handleUnpublishClick(record.id)
      break
  }
}

const handleSubmitValidate = async (id: number) => {
  submittingId.value = id
  try {
    // 先获取文章详情进行前端校验
    const articleRes = await getArticleById(id)
    if (articleRes.code !== 200) {
      message.error('获取文章信息失败')
      return
    }

    const article = articleRes.data
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

    if (errors.length > 0) {
      validateErrors.value = errors
      validateArticleId.value = id
      validateModalVisible.value = true
      return
    }

    // 校验通过，提交审核
    const res = await submitForReview(id)
    if (res.code === 200) {
      message.success('已提交审核')
      fetchData()
    } else {
      message.error(res.msg || '提交失败')
    }
  } catch {
    message.error('提交失败')
  } finally {
    submittingId.value = 0
  }
}

fetchData()
</script>

<template>
  <div>
    <a-card style="margin-bottom: 16px;">
      <div style="display: flex; justify-content: space-between; align-items: center; flex-wrap: wrap; gap: 12px;">
        <a-input-search v-model:value="searchText" placeholder="搜索文章标题..." allow-clear style="width: 260px;" @search="handleSearch" />
        <a-button type="primary" @click="handleAdd">
          <PlusOutlined /> 新增文章
        </a-button>
      </div>
    </a-card>

    <a-card :body-style="{ padding: 0 }">
      <a-table
        bordered
        size="middle"
        :columns="columns"
        :data-source="articles"
        :pagination="pagination"
        :loading="loading"
        :row-class-name="zebraRow"
        :locale="{ emptyText }"
        row-key="id"
        @change="handleTableChange"
      >
      <template #bodyCell="{ column, record, index }">
        <template v-if="column.key === 'index'">
          {{ (pagination.current - 1) * pagination.pageSize + index + 1 }}
        </template>
        <template v-if="column.key === 'category'">
          {{ record.category?.name || '-' }}
        </template>
        <template v-if="column.key === 'author'">
          {{ record.author?.nickname || '-' }}
        </template>
        <template v-if="column.key === 'status'">
          <a-tag :color="statusMap[record.status].color">
            {{ statusMap[record.status].label }}
          </a-tag>
        </template>
        <template v-if="column.key === 'createdAt'">
          {{ formatDate(record.createdAt) }}
        </template>
        <template v-if="column.key === 'actions'">
          <a-dropdown>
            <a-button size="small">操作 <DownOutlined /></a-button>
            <template #overlay>
              <a-menu @click="(e: any) => handleAction(record, e.key)">
                <a-menu-item key="edit">编辑</a-menu-item>
                <a-menu-item v-if="record.status === 0 || record.status === 2" key="submit">提交审核</a-menu-item>
                <a-menu-item v-if="record.status === 1" key="unpublish">下架</a-menu-item>
                <a-menu-item v-if="record.status === 3" key="approve">
                  <span style="color: #52c41a;">审核通过</span>
                </a-menu-item>
                <a-menu-item v-if="record.status === 3" key="reject">
                  <span style="color: #ff4d4f;">驳回</span>
                </a-menu-item>
                <a-menu-divider />
                <a-menu-item key="delete" danger>删除</a-menu-item>
              </a-menu>
            </template>
          </a-dropdown>
        </template>
      </template>
    </a-table>
    </a-card>

    <!-- 驳回弹窗 -->
    <a-modal v-model:open="rejectModalVisible" title="驳回意见" @ok="handleRejectConfirm">
      <a-textarea v-model:value="rejectComment" placeholder="请输入驳回原因..." :rows="4" />
    </a-modal>

    <!-- 下架弹窗 -->
    <a-modal v-model:open="unpublishModalVisible" title="下架原因" @ok="handleUnpublishConfirm">
      <a-textarea v-model:value="unpublishComment" placeholder="请输入下架原因..." :rows="4" />
    </a-modal>

    <!-- 校验失败弹窗 -->
    <a-modal v-model:open="validateModalVisible" title="无法提交审核" :footer="null">
      <p style="margin-bottom: 8px; color: #666;">以下必填项未完善，请编辑文章后重试：</p>
      <ul style="padding-left: 20px; color: #ff4d4f;">
        <li v-for="err in validateErrors" :key="err">{{ err }}</li>
      </ul>
      <div style="margin-top: 16px; text-align: right;">
        <a-button type="primary" @click="validateModalVisible = false">知道了</a-button>
      </div>
    </a-modal>
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
