<script setup lang="ts">
import { ref } from 'vue'
import { Modal, message } from 'ant-design-vue'
import { getCommentList, deleteComment } from '@/api/comment'
import { formatDate } from '@/utils/date'
import { createPagination, zebraRow, emptyText } from '@/utils/table'
import type { Comment } from '@/types'

const searchText = ref('')
const loading = ref(false)
const comments = ref<Comment[]>([])

const pagination = createPagination()

const columns = [
  { title: '序号', dataIndex: 'index', key: 'index', width: 80, align: 'center' as const },
  { title: '文章ID', dataIndex: 'articleId', key: 'articleId', width: 80, align: 'center' as const },
  { title: '评论内容', dataIndex: 'content', key: 'content', ellipsis: true },
  { title: '评论人', dataIndex: 'user', key: 'user', width: 120, align: 'center' as const },
  { title: '创建时间', dataIndex: 'createdAt', key: 'createdAt', width: 180, align: 'center' as const },
  { title: '操作', key: 'actions', width: 120, align: 'center' as const }
]

const fetchData = async () => {
  loading.value = true
  try {
    const params: { page: number; pageSize: number; keyword?: string } = { page: pagination.current, pageSize: pagination.pageSize }
    if (searchText.value.trim()) {
      params.keyword = searchText.value.trim()
    }
    const response = await getCommentList(params)
    if (response.code === 200) {
      comments.value = response.data.items
      pagination.total = response.data.total
    }
  } catch (error) {
    console.error('获取评论列表失败:', error)
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

const handleDelete = async (id: number) => {
  Modal.confirm({
    title: '确认删除',
    content: '确定要删除这条评论吗？此操作不可撤销。',
    okText: '删除',
    okButtonProps: { danger: true },
    cancelText: '取消',
    async onOk() {
      try {
        const response = await deleteComment(id)
        if (response.code === 200) {
          message.success('删除成功')
          fetchData()
        } else {
          message.error(response.msg || '删除失败')
        }
      } catch (error) {
        console.error('删除评论失败:', error)
        message.error('删除失败，请稍后重试')
      }
    }
  })
}

fetchData()
</script>

<template>
  <div>
    <a-card style="margin-bottom: 16px;">
      <div style="display: flex; justify-content: space-between; align-items: center; flex-wrap: wrap; gap: 12px;">
        <a-input-search v-model:value="searchText" placeholder="搜索评论内容..." allow-clear style="width: 240px;" @search="handleSearch" />
      </div>
    </a-card>

    <a-card :body-style="{ padding: 0 }">
      <a-table
        bordered
        size="middle"
        :columns="columns"
        :data-source="comments"
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
        <template v-if="column.key === 'user'">
          {{ record.user?.nickname || record.user?.username || '-' }}
        </template>
        <template v-if="column.key === 'createdAt'">
          {{ formatDate(record.createdAt) }}
        </template>
        <template v-if="column.key === 'actions'">
          <a-button size="small" danger @click="handleDelete(record.id)">删除</a-button>
        </template>
      </template>
    </a-table>
    </a-card>
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
