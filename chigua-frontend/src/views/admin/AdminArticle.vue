<script setup lang="ts">
import { ref, onMounted, h } from "vue"
import {
  Table,
  Button,
  Modal,
  Form,
  Input,
  Select,
  message,
  Popconfirm,
  Tag as TagAnt,
  Space
} from "ant-design-vue"
import type { TableColumnsType, TablePaginationConfig } from "ant-design-vue"
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  CheckCircleOutlined,
  StopOutlined
} from "@ant-design/icons-vue"
import { adminApi } from "@/services/admin"
import type { Article } from "@/types/article"
import type { Category } from "@/types/category"
import type { Tag } from "@/types/tag"

const statusMap: Record<number, { text: string; color: string }> = {
  0: { text: "草稿", color: "default" },
  1: { text: "已发布", color: "success" },
  2: { text: "已下架", color: "error" },
  3: { text: "审核中", color: "processing" }
}

const getStatusTag = (status: number) => {
  const s = statusMap[status] || { text: "未知", color: "default" }
  return h(TagAnt, { color: s.color }, () => s.text)
}

const columns: TableColumnsType = [
  { title: "ID", dataIndex: "id", key: "id", width: 80 },
  { title: "标题", dataIndex: "title", key: "title", ellipsis: true },
  {
    title: "作者",
    dataIndex: ["author", "username"],
    key: "author",
    width: 120
  },
  {
    title: "分类",
    dataIndex: ["category", "name"],
    key: "category",
    width: 100
  },
  {
    title: "状态",
    dataIndex: "status",
    key: "status",
    width: 100,
    customRender: ({ text }: { text: number }) => getStatusTag(text)
  },
  { title: "创建时间", dataIndex: "created_at", key: "created_at", width: 180 },
  { title: "操作", key: "action", width: 200 }
]

const articles = ref<Article[]>([])
const categories = ref<Category[]>([])
const tags = ref<Tag[]>([])
const loading = ref(false)
const modalVisible = ref(false)
const modalTitle = ref("新增文章")
const isEdit = ref(false)
const currentArticle = ref<Article | null>(null)

const formState = ref({
  title: "",
  content: "",
  category_id: undefined as number | undefined,
  tag_ids: [] as number[],
  coverImage: ""
})

const pagination = ref<TablePaginationConfig>({
  current: 1,
  pageSize: 10,
  total: 0
})

const loadArticles = async () => {
  loading.value = true
  try {
    const res = await adminApi.article.list({
      page: pagination.value.current || 1,
      pageSize: pagination.value.pageSize || 10
    })
    if (res.code === 200) {
      articles.value = res.data.items
      pagination.value.total = res.data.total
    }
  } catch (error) {
    message.error("获取文章列表失败")
  } finally {
    loading.value = false
  }
}

const loadCategories = async () => {
  const res = await adminApi.category.getAll()
  if (res.code === 200) {
    categories.value = res.data
  }
}

const loadTags = async () => {
  const res = await adminApi.tag.getAll()
  if (res.code === 200) {
    tags.value = res.data
  }
}

const openModal = (record?: any) => {
  if (record) {
    isEdit.value = true
    currentArticle.value = record as Article
    modalTitle.value = "编辑文章"
    formState.value = {
      title: record.title,
      content: record.content,
      category_id: record.category_id,
      tag_ids: record.tags?.map((t: any) => t.id) || [],
      coverImage: record.coverImage || ""
    }
  } else {
    isEdit.value = false
    currentArticle.value = null
    modalTitle.value = "新增文章"
    formState.value = {
      title: "",
      content: "",
      category_id: undefined,
      tag_ids: [],
      coverImage: ""
    }
  }
  loadCategories()
  loadTags()
  modalVisible.value = true
}

const handleSubmit = async () => {
  if (!formState.value.title) {
    message.error("请填写标题")
    return
  }
  if (!formState.value.category_id) {
    message.error("请选择分类")
    return
  }

  const submitData = {
    title: formState.value.title,
    content: formState.value.content,
    category_id: formState.value.category_id,
    tag_ids: formState.value.tag_ids,
    coverImage: formState.value.coverImage
  }

  try {
    if (isEdit.value && currentArticle.value) {
      await adminApi.article.update(currentArticle.value.id, submitData)
      message.success("更新成功")
    } else {
      await adminApi.article.create(submitData)
      message.success("创建成功")
    }
    modalVisible.value = false
    loadArticles()
  } catch (error) {
    message.error(isEdit.value ? "更新失败" : "创建失败")
  }
}

const handleDelete = async (id: number) => {
  try {
    await adminApi.article.delete(id)
    message.success("删除成功")
    loadArticles()
  } catch (error) {
    message.error("删除失败")
  }
}

const handlePublish = async (id: number) => {
  try {
    await adminApi.article.publish(id)
    message.success("发布成功")
    loadArticles()
  } catch (error) {
    message.error("发布失败")
  }
}

const handleStatusChange = async (id: number, status: number) => {
  try {
    await adminApi.article.updateStatus(id, status)
    message.success("状态更新成功")
    loadArticles()
  } catch (error) {
    message.error("状态更新失败")
  }
}

const handleTableChange = (pag: TablePaginationConfig) => {
  if (pag.current) pagination.value.current = pag.current
  if (pag.pageSize) pagination.value.pageSize = pag.pageSize
  loadArticles()
}

onMounted(() => {
  loadArticles()
})
</script>

<template>
  <div class="p-6">
    <div class="mb-4 flex justify-between items-center">
      <h2 class="text-xl font-semibold">文章管理</h2>
      <Button type="primary" @click="openModal()">
        <PlusOutlined /> 新增
      </Button>
    </div>

    <Table
      :columns="columns"
      :data-source="articles"
      :loading="loading"
      :pagination="pagination"
      :scroll="{ x: 1000 }"
      row-key="id"
      @change="handleTableChange"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'action'">
          <Space>
            <Button size="small" type="link" @click="openModal(record)">
              <EditOutlined /> 编辑
            </Button>
            <Button
              v-if="record.status === 0 || record.status === 2"
              size="small"
              type="link"
              @click="handlePublish(record.id)"
            >
              <CheckCircleOutlined /> 发布
            </Button>
            <Button
              v-if="record.status === 1"
              size="small"
              type="link"
              danger
              @click="handleStatusChange(record.id, 2)"
            >
              <StopOutlined /> 下架
            </Button>
            <Popconfirm
              title="确定删除这篇文章吗？"
              @confirm="handleDelete(record.id)"
            >
              <Button size="small" type="link" danger>
                <DeleteOutlined /> 删除
              </Button>
            </Popconfirm>
          </Space>
        </template>
      </template>
    </Table>

    <Modal
      v-model:open="modalVisible"
      :title="modalTitle"
      width="800px"
      @ok="handleSubmit"
    >
      <Form layout="vertical">
        <Form.Item label="标题" required>
          <Input v-model:value="formState.title" placeholder="请输入文章标题" />
        </Form.Item>
        <Form.Item label="封面图片">
          <Input
            v-model:value="formState.coverImage"
            placeholder="请输入封面图片URL"
          />
        </Form.Item>
        <Form.Item label="分类" required>
          <Select
            v-model:value="formState.category_id"
            placeholder="请选择分类"
          >
            <Select.Option
              v-for="cat in categories"
              :key="cat.id"
              :value="cat.id"
            >
              {{ cat.name }}
            </Select.Option>
          </Select>
        </Form.Item>
        <Form.Item label="标签">
          <Select
            v-model:value="formState.tag_ids"
            mode="multiple"
            placeholder="请选择标签"
          >
            <Select.Option v-for="tag in tags" :key="tag.id" :value="tag.id">
              {{ tag.name }}
            </Select.Option>
          </Select>
        </Form.Item>
        <Form.Item label="内容">
          <Input.TextArea
            v-model:value="formState.content"
            placeholder="请输入文章内容"
            :rows="10"
          />
        </Form.Item>
      </Form>
    </Modal>
  </div>
</template>
