<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import { getAllReviewRecords } from '@/api/article'
import { formatDate } from '@/utils/date'
import { createPagination, zebraRow, emptyText } from '@/utils/table'
import type { ArticleReviewRecordWithTitle } from '@/types'

const router = useRouter()
const loading = ref(false)
const records = ref<ArticleReviewRecordWithTitle[]>([])
const pagination = createPagination()

const columns = [
  { title: '序号', dataIndex: 'index', key: 'index', width: 70, align: 'center' as const },
  { title: '文章标题', dataIndex: 'articleTitle', key: 'articleTitle', ellipsis: true, minWidth: 200 },
  { title: '审核人', dataIndex: 'reviewer', key: 'reviewer', width: 120, align: 'center' as const },
  { title: '操作', dataIndex: 'action', key: 'action', width: 100, align: 'center' as const },
  { title: '审核意见', dataIndex: 'comment', key: 'comment', width: 200, ellipsis: true },
  { title: '审核时间', dataIndex: 'createdAt', key: 'createdAt', width: 170, align: 'center' as const },
]

const fetchData = async () => {
  loading.value = true
  try {
    const res = await getAllReviewRecords({
      page: pagination.current,
      pageSize: pagination.pageSize
    })
    if (res.code === 200) {
      records.value = res.data.items
      pagination.total = res.data.total
    }
  } catch {
    message.error('获取审核记录失败')
  } finally {
    loading.value = false
  }
}

const handleTableChange = (paginationInfo: { current: number; pageSize: number }) => {
  pagination.current = paginationInfo.current
  pagination.pageSize = paginationInfo.pageSize
  fetchData()
}

const handleArticleClick = (articleId: number) => {
  router.push(`/articles/${articleId}/edit`)
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
        :data-source="records"
        :pagination="pagination"
        :loading="loading"
        :row-class-name="zebraRow"
        :locale="{ emptyText: emptyText || '暂无审核记录' }"
        row-key="id"
        @change="handleTableChange"
      >
        <template #bodyCell="{ column, record, index }">
          <template v-if="column.key === 'index'">
            {{ (pagination.current - 1) * pagination.pageSize + index + 1 }}
          </template>
          <template v-if="column.key === 'articleTitle'">
            <a @click="handleArticleClick(record.articleId)">{{ record.articleTitle }}</a>
          </template>
          <template v-if="column.key === 'reviewer'">
            {{ record.reviewer?.nickname || record.reviewer?.username || '-' }}
          </template>
          <template v-if="column.key === 'action'">
            <a-tag :color="record.action === 'approve' ? 'green' : record.action === 'reject' ? 'red' : 'orange'">
              {{ record.action === 'approve' ? '通过' : record.action === 'reject' ? '驳回' : '下架' }}
            </a-tag>
          </template>
          <template v-if="column.key === 'comment'">
            {{ record.comment || '-' }}
          </template>
          <template v-if="column.key === 'createdAt'">
            {{ formatDate(record.createdAt) }}
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
