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
import type { Tag } from "@/types/tag"

const columns: TableColumnsType = [
  { title: "ID", dataIndex: "id", key: "id", width: 100 },
  { title: "名称", dataIndex: "name", key: "name" },
  { title: "创建时间", dataIndex: "created_at", key: "created_at", width: 180 },
  { title: "操作", key: "action", width: 150, fixed: "right" }
]

const tags = ref<Tag[]>([])
const loading = ref(false)
const modalVisible = ref(false)
const modalTitle = ref("新增标签")
const isEdit = ref(false)
const currentTag = ref<Tag | null>(null)

const formState = ref({
  name: ""
})

const loadTags = async () => {
  loading.value = true
  try {
    const res = await adminApi.tag.getAll()
    if (res.code === 200) {
      tags.value = res.data
    }
  } catch (error) {
    message.error("获取标签列表失败")
  } finally {
    loading.value = false
  }
}

const openModal = (record?: any) => {
  if (record) {
    isEdit.value = true
    currentTag.value = record as Tag
    modalTitle.value = "编辑标签"
    formState.value.name = record.name
  } else {
    isEdit.value = false
    currentTag.value = null
    modalTitle.value = "新增标签"
    formState.value.name = ""
  }
  modalVisible.value = true
}

const handleSubmit = async () => {
  if (!formState.value.name.trim()) {
    message.error("请填写标签名称")
    return
  }

  try {
    if (isEdit.value && currentTag.value) {
      await adminApi.tag.update(currentTag.value.id, {
        name: formState.value.name
      })
      message.success("更新成功")
    } else {
      await adminApi.tag.create({ name: formState.value.name })
      message.success("创建成功")
    }
    modalVisible.value = false
    loadTags()
  } catch (error: any) {
    if (error?.response?.data?.msg?.includes("已存在")) {
      message.error("标签名称已存在")
    } else {
      message.error(isEdit.value ? "更新失败" : "创建失败")
    }
  }
}

const handleDelete = async (id: number) => {
  try {
    await adminApi.tag.delete(id)
    message.success("删除成功")
    loadTags()
  } catch (error) {
    message.error("删除失败")
  }
}

onMounted(() => {
  loadTags()
})
</script>

<template>
  <div class="p-6">
    <div class="mb-4 flex justify-between items-center">
      <h2 class="text-xl font-semibold">标签管理</h2>
      <Button type="primary" @click="openModal()">
        <PlusOutlined /> 新增
      </Button>
    </div>

    <Table
      :columns="columns"
      :data-source="tags"
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
              title="确定删除该标签吗？"
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
        <Form.Item label="标签名称" required>
          <Input v-model:value="formState.name" placeholder="请输入标签名称" />
        </Form.Item>
      </Form>
    </Modal>
  </div>
</template>
