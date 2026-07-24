<script setup lang="ts">
import { ref, reactive } from 'vue'
import { PlusOutlined, EditOutlined } from '@ant-design/icons-vue'
import { Modal, message } from 'ant-design-vue'
import { getCategoryList, createCategory, updateCategory, deleteCategory } from '@/api/category'
import { formatDate } from '@/utils/date'
import { createPagination, zebraRow, emptyText } from '@/utils/table'
import type { Category } from '@/types'

const searchText = ref('')
const loading = ref(false)
const categories = ref<Category[]>([])
const modalVisible = ref(false)
const modalLoading = ref(false)
const isEdit = ref(false)
const editingId = ref(0)

const pagination = createPagination()

const form = reactive({
  name: '',
  sortOrder: 0
})

const columns = [
  { title: '序号', dataIndex: 'index', key: 'index', width: 80, align: 'center' },
  { title: '分类名称', dataIndex: 'name', key: 'name', align: 'left' },
  { title: '排序', dataIndex: 'sortOrder', key: 'sortOrder', width: 80, align: 'center' },
  { title: '创建时间', dataIndex: 'createdAt', key: 'createdAt', width: 180, align: 'center' },
  { title: '操作', key: 'actions', width: 180, align: 'center', fixed: 'right' }
]

const fetchData = async () => {
  loading.value = true
  try {
    const params: { page: number; pageSize: number; keyword?: string } = { page: pagination.current, pageSize: pagination.pageSize }
    if (searchText.value.trim()) {
      params.keyword = searchText.value.trim()
    }
    const response = await getCategoryList(params)
    if (response.code === 200) {
      categories.value = response.data.items
      pagination.total = response.data.total
    }
  } catch (error) {
    console.error('获取分类列表失败:', error)
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
  isEdit.value = false
  editingId.value = 0
  form.name = ''
  form.sortOrder = 0
  modalVisible.value = true
}

const handleEdit = (record: Category) => {
  isEdit.value = true
  editingId.value = record.id
  form.name = record.name
  form.sortOrder = record.sortOrder || 0
  modalVisible.value = true
}

const handleDelete = async (id: number) => {
  Modal.confirm({
    title: '确认删除',
    content: '确定要删除这个分类吗？此操作不可撤销。',
    okText: '删除',
    okButtonProps: { danger: true },
    cancelText: '取消',
    async onOk() {
      try {
        const response = await deleteCategory(id)
        if (response.code === 200) {
          message.success('删除成功')
          fetchData()
        } else {
          message.error(response.msg || '删除失败')
        }
      } catch (error) {
        console.error('删除分类失败:', error)
        message.error('删除失败，请稍后重试')
      }
    }
  })
}

const handleModalSubmit = async () => {
  if (!form.name.trim()) {
    message.warning('请输入分类名称')
    return
  }
  modalLoading.value = true
  try {
    const response = isEdit.value
      ? await updateCategory(editingId.value, { name: form.name.trim(), sortOrder: form.sortOrder })
      : await createCategory({ name: form.name.trim(), sortOrder: form.sortOrder })
    if (response.code === 200) {
      message.success(isEdit.value ? '编辑成功' : '新增分类成功')
      modalVisible.value = false
      await fetchData()
    } else {
      message.error(response.msg || '保存失败')
    }
  } catch (error) {
    console.error(isEdit.value ? '更新分类失败:' : '创建分类失败:', error)
    message.error('保存失败，请稍后重试')
  } finally {
    modalLoading.value = false
  }
}

const handleModalCancel = () => {
  modalVisible.value = false
}

fetchData()
</script>

<template>
  <div>
    <a-card style="margin-bottom: 16px;">
      <div style="display: flex; justify-content: space-between; align-items: center; flex-wrap: wrap; gap: 12px;">
        <a-input-search v-model:value="searchText" placeholder="搜索分类名称..." allow-clear style="width: 240px;" @search="handleSearch" />
        <a-button type="primary" @click="handleAdd" v-permission="'category:create'">
          <PlusOutlined /> 新增分类
        </a-button>
      </div>
    </a-card>

    <a-card :body-style="{ padding: 0 }">
      <a-table
        bordered
        size="middle"
        :columns="columns"
        :data-source="categories"
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
        <template v-if="column.key === 'createdAt'">
          {{ formatDate(record.createdAt) }}
        </template>
        <template v-if="column.key === 'actions'">
          <a-space>
            <a-button type="primary" size="small" @click="handleEdit(record)" v-permission="'category:update'">
              <component :is="EditOutlined" />
              编辑
            </a-button>
            <a-button danger size="small" @click="handleDelete(record.id)" v-permission="'category:delete'">删除</a-button>
          </a-space>
        </template>
      </template>
    </a-table>
    </a-card>

    <a-modal
      v-model:visible="modalVisible"
      :title="isEdit ? '编辑分类' : '新增分类'"
      width="500px"
      @ok="handleModalSubmit"
      @cancel="handleModalCancel"
      :confirm-loading="modalLoading"
    >
      <a-form :model="form" :label-col="{ span: 4 }" :wrapper-col="{ span: 20 }">
        <a-form-item label="名称" name="name" :rules="[{ required: true, message: '请输入分类名称' }]">
          <a-input v-model:value="form.name" placeholder="请输入分类名称" />
        </a-form-item>
        <a-form-item label="排序">
          <a-input-number v-model:value="form.sortOrder" :min="0" style="width: 100%" />
        </a-form-item>
      </a-form>
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
