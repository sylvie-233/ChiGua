<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { PlusOutlined, EditOutlined } from '@ant-design/icons-vue'
import { Modal, message } from 'ant-design-vue'
import { getAllMenus, createMenu, updateMenu, deleteMenu } from '@/api/menu'
import type { MenuItem } from '@/api/menu'
import { useUserStore } from '@/stores/user'

const userStore = useUserStore()
const loading = ref(false)
const flatMenus = ref<MenuItem[]>([])
const expandedKeys = ref<(string | number)[]>([])

// 收集所有有子节点的 ID
const collectParentIds = (items: MenuItem[]): (string | number)[] => {
  return items.flatMap(m => m.children?.length ? [m.id, ...collectParentIds(m.children)] : [])
}
const treeData = ref<MenuItem[]>([])

// 构建树形数据
const buildTree = (items: MenuItem[], parentId = 0): MenuItem[] => {
  return items
    .filter(m => m.parentId === parentId)
    .sort((a, b) => a.sortOrder - b.sortOrder)
    .map(m => {
      const kids = buildTree(items, m.id)
      return kids.length ? { ...m, children: kids } : { ...m }
    })
}

const columns = [
  { title: '名称', dataIndex: 'title', key: 'title', width: 200 },
  { title: '路径', dataIndex: 'path', key: 'path', width: 200 },
  { title: '权限标识', dataIndex: 'permissionCode', key: 'permissionCode', width: 180 },
  { title: '图标', dataIndex: 'icon', key: 'icon', width: 150 },
  { title: '类型', dataIndex: 'menuType', key: 'menuType', width: 70 },
  { title: '排序', dataIndex: 'sortOrder', key: 'sortOrder', width: 70 },
  { title: '操作', key: 'actions', width: 150, align: 'center' as const }
]

const fetchData = async () => {
  loading.value = true
  try {
    const res = await getAllMenus()
    if (res.code === 200) {
      flatMenus.value = res.data
      treeData.value = buildTree(res.data)
      expandedKeys.value = collectParentIds(treeData.value)
    }
  } catch {
    message.error('获取菜单列表失败')
  } finally {
    loading.value = false
  }
}

// 表单
const modalVisible = ref(false)
const isEdit = ref(false)
const saving = ref(false)
const form = reactive({
  id: 0,
  parentId: 0,
  title: '',
  path: '',
  icon: '',
  permissionCode: '',
  menuType: 'M',
  sortOrder: 0,
  visible: true,
})

const openCreate = (parentId = 0) => {
  isEdit.value = false
  form.id = 0; form.parentId = parentId; form.title = ''; form.path = ''
  form.icon = ''; form.permissionCode = ''; form.menuType = 'M'
  form.sortOrder = 0; form.visible = true
  modalVisible.value = true
}

const openEdit = (menu: MenuItem) => {
  isEdit.value = true
  Object.assign(form, menu)
  modalVisible.value = true
}

const handleSave = async () => {
  if (!form.title.trim()) { message.warning('请输入名称'); return }
  saving.value = true
  try {
    const data = { ...form }
    const res = isEdit.value ? await updateMenu(form.id, data) : await createMenu(data)
    if (res.code === 200) {
      message.success(isEdit.value ? '更新成功' : '创建成功')
      modalVisible.value = false
      fetchData()
    } else {
      message.error(res.msg || '保存失败')
    }
  } catch { message.error('保存失败') }
  finally { saving.value = false }
}

const handleDelete = (id: number) => {
  Modal.confirm({
    title: '确认删除', content: '删除后子菜单也会失效，确定删除？',
    okText: '删除', okButtonProps: { danger: true }, cancelText: '取消',
    async onOk() {
      const res = await deleteMenu(id)
      if (res.code === 200) { message.success('已删除'); fetchData() }
      else { message.error(res.msg || '删除失败') }
    }
  })
}

const parentOptions = computed(() => {
  return [{ id: 0, title: '顶级' }, ...flatMenus.value.filter(m => m.menuType !== 'B')]
})

onMounted(fetchData)
</script>

<template>
  <div>
    <a-card style="margin-bottom: 16px;">
      <a-space>
        <a-button type="primary" @click="openCreate(0)" v-if="userStore.hasPermission('role:manage')">
          <PlusOutlined /> 新增菜单
        </a-button>
        <a-button @click="expandedKeys = collectParentIds(treeData)">展开全部</a-button>
        <a-button @click="expandedKeys = []">折叠全部</a-button>
      </a-space>
    </a-card>

    <a-card :body-style="{ padding: 0 }">
      <a-table bordered size="middle" :columns="columns"
        :data-source="treeData" :loading="loading"
        :pagination="false" row-key="id" v-model:expandedRowKeys="expandedKeys">
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'title'">{{ record.title }}</template>
          <template v-if="column.key === 'menuType'">
            <a-tag :color="record.menuType === 'M' ? 'blue' : 'orange'">{{ record.menuType === 'M' ? '菜单' : '按钮' }}</a-tag>
          </template>
          <template v-if="column.key === 'actions'">
            <a-space>
              <a-button size="small" @click="openCreate(record.id)" v-if="record.menuType === 'M' && userStore.hasPermission('role:manage')">添加子级</a-button>
              <a-button type="primary" size="small" @click="openEdit(record)" v-if="userStore.hasPermission('role:manage')"><EditOutlined /></a-button>
              <a-button size="small" danger @click="handleDelete(record.id)" v-if="userStore.hasPermission('role:manage')">删除</a-button>
            </a-space>
          </template>
        </template>
      </a-table>
    </a-card>

    <a-modal v-model:open="modalVisible" :title="isEdit ? '编辑菜单' : '新增菜单'" @ok="handleSave" :confirm-loading="saving" width="520px">
      <a-form :model="form" layout="vertical">
        <a-form-item label="上级菜单">
          <a-select v-model:value="form.parentId" style="width: 100%">
            <a-select-option v-for="p in parentOptions" :key="p.id" :value="p.id">{{ p.title }}</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="名称">
          <a-input v-model:value="form.title" />
        </a-form-item>
        <a-form-item label="路由路径">
          <a-input v-model:value="form.path" placeholder="如 /articles" />
        </a-form-item>
        <a-form-item label="图标">
          <a-input v-model:value="form.icon" placeholder="如 FileOutlined" />
        </a-form-item>
        <a-form-item label="权限标识">
          <a-input v-model:value="form.permissionCode" placeholder="如 article:list" />
        </a-form-item>
        <a-row :gutter="16">
          <a-col :span="8">
            <a-form-item label="类型">
              <a-select v-model:value="form.menuType">
                <a-select-option value="M">菜单</a-select-option>
                <a-select-option value="B">按钮</a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
          <a-col :span="8">
            <a-form-item label="排序">
              <a-input-number v-model:value="form.sortOrder" :min="0" style="width: 100%" />
            </a-form-item>
          </a-col>
          <a-col :span="8">
            <a-form-item label="可见">
              <a-switch v-model:checked="form.visible" />
            </a-form-item>
          </a-col>
        </a-row>
      </a-form>
    </a-modal>
  </div>
</template>
