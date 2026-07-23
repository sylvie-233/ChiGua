<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { PlusOutlined } from '@ant-design/icons-vue'
import { Modal } from 'ant-design-vue'
import { getArticleList, deleteArticle, updateArticleStatus } from '@/api/article'
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
  3: { label: '审核中', color: 'orange' }
}

const columns = [
  { title: '序号', dataIndex: 'index', key: 'index', width: 70, align: 'center' as const },
  { title: '标题', dataIndex: 'title', key: 'title', ellipsis: true, minWidth: 200 },
  { title: '分类', dataIndex: 'category', key: 'category', width: 120, align: 'center' as const },
  { title: '作者', dataIndex: 'author', key: 'author', width: 100, align: 'center' as const },
  { title: '状态', dataIndex: 'status', key: 'status', width: 100, align: 'center' as const },
  { title: '创建时间', dataIndex: 'createdAt', key: 'createdAt', width: 170, align: 'center' as const },
  { title: '操作', key: 'actions', width: 160, align: 'center' as const }
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

const handleEdit = (id: number) => {
  router.push(`/articles/${id}/edit`)
}

const handleDelete = async (id: number) => {
  Modal.confirm({
    title: '确认删除',
    content: '确定要删除这篇文章吗？此操作不可撤销。',
    okText: '删除',
    okType: 'danger',
    cancelText: '取消',
    async onOk() {
      try {
        const response = await deleteArticle(id)
        if (response.code === 200) {
          fetchData()
        }
      } catch (error) {
        console.error('删除文章失败:', error)
      }
    }
  })
}

const handleStatusChange = async (id: number, status: number) => {
  try {
    const response = await updateArticleStatus(id, status)
    if (response.code === 200) {
      fetchData()
    }
  } catch (error) {
    console.error('更新文章状态失败:', error)
  }
}

const onStatusMenuClick = (id: number, e: { key: string }) => {
  handleStatusChange(id, Number(e.key))
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
          <a-dropdown>
            <a-tag :color="statusMap[record.status].color" style="cursor: pointer;">
              {{ statusMap[record.status].label }}
            </a-tag>
            <template #overlay>
              <a-menu @click="onStatusMenuClick(record.id, $event)">
                <a-menu-item v-for="(item, key) in statusMap" :key="key" :disabled="Number(key) === record.status">
                  {{ item.label }}
                </a-menu-item>
              </a-menu>
            </template>
          </a-dropdown>
        </template>
        <template v-if="column.key === 'createdAt'">
          {{ formatDate(record.createdAt) }}
        </template>
        <template v-if="column.key === 'actions'">
          <a-space>
            <a-button size="small" @click="handleEdit(record.id)">编辑</a-button>
            <a-button size="small" type="danger" @click="handleDelete(record.id)">删除</a-button>
          </a-space>
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
