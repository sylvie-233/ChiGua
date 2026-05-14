<script setup lang="ts">
import { ref, onMounted } from "vue"
import {
  Table,
  Button,
  Modal,
  Form,
  Input,
  message,
  Popconfirm,
  Space
} from "ant-design-vue"
import type { TableColumnsType } from "ant-design-vue"
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined
} from "@ant-design/icons-vue"
import { adminApi } from "@/services/admin"
import type { Category } from "@/types/category"

const columns: TableColumnsType = [
  { title: "ID", dataIndex: "id", key: "id", width: 100 },
  { title: "名称", dataIndex: "name", key: "name" },
  { title: "创建时间", dataIndex: "created_at", key: "created_at", width: 180 },
  { title: "操作", key: "action", width: 150, fixed: "right" }
]

const categories = ref<Category[]>([])
const loading = ref(false)
const modalVisible = ref(false)
const modalTitle = ref("新增分类")
const isEdit = ref(false)
const currentCategory = ref<Category | null>(null)

const formState = ref({
  name: ""
})

const loadCategories = async () => {
  loading.value = true
  try {
    const res = await adminApi.category.getAll()
    if (res.code === 200) {
      categories.value = res.data
    }
  } catch (error) {
    message.error("获取分类列表失败")
  } finally {
    loading.value = false
  }
}

const openModal = (record?: any) => {
  if (record) {
    isEdit.value = true
    currentCategory.value = record as Category
    modalTitle.value = "编辑分类"
    formState.value.name = record.name
  } else {
    isEdit.value = false
    currentCategory.value = null
    modalTitle.value = "新增分类"
    formState.value.name = ""
  }
  modalVisible.value = true
}

const handleSubmit = async () => {
  if (!formState.value.name.trim()) {
    message.error("请填写分类名称")
    return
  }

  try {
    if (isEdit.value && currentCategory.value) {
      await adminApi.category.update(currentCategory.value.id, {
        name: formState.value.name
      })
      message.success("更新成功")
    } else {
      await adminApi.category.create({ name: formState.value.name })
      message.success("创建成功")
    }
    modalVisible.value = false
    loadCategories()
  } catch (error: any) {
    if (error?.response?.data?.msg?.includes("已存在")) {
      message.error("分类名称已存在")
    } else {
      message.error(isEdit.value ? "更新失败" : "创建失败")
    }
  }
}

const handleDelete = async (id: number) => {
  try {
    await adminApi.category.delete(id)
    message.success("删除成功")
    loadCategories()
  } catch (error) {
    message.error("删除失败")
  }
}

onMounted(() => {
  loadCategories()
})
</script>

<template>
  <div class="p-6">
    <div class="mb-4 flex justify-between items-center">
      <h2 class="text-xl font-semibold">分类管理</h2>
      <Button type="primary" @click="openModal()">
        <PlusOutlined /> 新增
      </Button>
    </div>

    <Table
      :columns="columns"
      :data-source="categories"
      :loading="loading"
      :pagination="false"
      row-key="id"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'action'">
          <Space>
            <Button size="small" type="link" @click="openModal(record)">
              <EditOutlined /> 编辑
            </Button>
            <Popconfirm
              title="确定删除该分类吗？"
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

    <Modal v-model:open="modalVisible" :title="modalTitle" @ok="handleSubmit">
      <Form layout="vertical">
        <Form.Item label="分类名称" required>
          <Input v-model:value="formState.name" placeholder="请输入分类名称" />
        </Form.Item>
      </Form>
    </Modal>
  </div>
</template>
