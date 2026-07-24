<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { PlusOutlined } from '@ant-design/icons-vue'
import { Modal, message } from 'ant-design-vue'
import { getRoles, getRoleMenuIds, updateRoleMenus, createRole, updateRole, deleteRole } from '@/api/rbac'
import { getAllMenus } from '@/api/menu'
import { createPagination, zebraRow, emptyText } from '@/utils/table'
import { useUserStore } from '@/stores/user'
import type { Role } from '@/types'
import type { MenuItem } from '@/api/menu'

const userStore = useUserStore()
const loading = ref(false)
const roles = ref<Role[]>([])
const pagination = createPagination()
const searchText = ref('')

const columns = [
  { title: '角色名称', dataIndex: 'name', key: 'name', width: 150 },
  { title: '标识', dataIndex: 'code', key: 'code', width: 120 },
  { title: '描述', dataIndex: 'description', key: 'description', ellipsis: true },
  { title: '操作', key: 'actions', width: 200, align: 'center' as const }
]

const fetchRoles = async () => {
  loading.value = true
  try {
    const res = await getRoles({ page: pagination.current, pageSize: pagination.pageSize, keyword: searchText.value.trim() || undefined })
    if (res.code === 200) {
      roles.value = res.data.items
      pagination.total = res.data.total
    }
  } catch {
    message.error('获取角色列表失败')
  } finally {
    loading.value = false
  }
}

const handleSearch = () => { pagination.current = 1; fetchRoles() }

const handleTableChange = (p: { current: number; pageSize: number }) => {
  pagination.current = p.current; pagination.pageSize = p.pageSize; fetchRoles()
}

// ====== 新增/编辑弹窗 ======
const formVisible = ref(false)
const isEdit = ref(false)
const saving = ref(false)
const form = reactive({ id: 0, code: '', name: '', description: '' })

const openCreate = () => {
  isEdit.value = false
  form.id = 0; form.code = ''; form.name = ''; form.description = ''
  formVisible.value = true
}

const openEdit = (role: Role) => {
  isEdit.value = true
  form.id = role.id; form.code = role.code; form.name = role.name; form.description = role.description || ''
  formVisible.value = true
}

const handleSave = async () => {
  if (!form.name.trim()) { message.warning('请输入角色名称'); return }
  if (!isEdit.value && !form.code.trim()) { message.warning('请输入角色标识'); return }
  saving.value = true
  try {
    const res = isEdit.value
      ? await updateRole(form.id, { name: form.name, description: form.description })
      : await createRole({ code: form.code, name: form.name, description: form.description })
    if (res.code === 200) {
      message.success(isEdit.value ? '更新成功' : '创建成功')
      formVisible.value = false
      fetchRoles()
    } else {
      message.error(res.msg || '保存失败')
    }
  } catch { message.error('保存失败') }
  finally { saving.value = false }
}

const handleDelete = (role: Role) => {
  Modal.confirm({
    title: '确认删除', content: `确定删除角色"${role.name}"吗？`,
    okText: '删除', okButtonProps: { danger: true }, cancelText: '取消',
    async onOk() {
      const res = await deleteRole(role.id)
      if (res.code === 200) { message.success('已删除'); fetchRoles() }
      else { message.error(res.msg || '删除失败') }
    }
  })
}

// ====== 菜单权限抽屉 ======
const drawerVisible = ref(false)
const currentRole = ref<Role | null>(null)
const allMenus = ref<MenuItem[]>([])
const selectedMenuIds = ref<number[]>([])
const menuSaving = ref(false)

const buildMenuTree = (menus: MenuItem[], parentId = 0): MenuItem[] => {
  return menus.filter(m => m.parentId === parentId).sort((a, b) => a.sortOrder - b.sortOrder)
    .map(m => { const kids = buildMenuTree(menus, m.id); return kids.length ? { ...m, children: kids } : { ...m } })
}

const menuTree = ref<MenuItem[]>([])

const openPermissionEditor = async (role: Role) => {
  currentRole.value = role; drawerVisible.value = true
  try {
    const [menuRes, roleMenuRes] = await Promise.all([getAllMenus(), getRoleMenuIds(role.id)])
    if (menuRes.code === 200) { allMenus.value = menuRes.data; menuTree.value = buildMenuTree(menuRes.data) }
    if (roleMenuRes.code === 200) selectedMenuIds.value = roleMenuRes.data
  } catch { message.error('加载菜单失败') }
}

const toggleMenu = (menu: MenuItem, checked: boolean) => {
  const collectIds = (m: MenuItem): number[] => [m.id, ...(m.children?.flatMap(collectIds) || [])]
  const ids = collectIds(menu)
  if (checked) ids.forEach(id => { if (!selectedMenuIds.value.includes(id)) selectedMenuIds.value.push(id) })
  else selectedMenuIds.value = selectedMenuIds.value.filter(id => !ids.includes(id))
}

const getCheckedState = (menu: MenuItem) => {
  if (!menu.children?.length) return { checked: selectedMenuIds.value.includes(menu.id), indeterminate: false }
  const childIds = menu.children.map(c => c.id)
  const checkedCount = childIds.filter(id => selectedMenuIds.value.includes(id)).length
  return { checked: checkedCount === childIds.length, indeterminate: checkedCount > 0 && checkedCount < childIds.length }
}

const handleSavePermissions = async () => {
  if (!currentRole.value) return
  menuSaving.value = true
  try {
    const res = await updateRoleMenus(currentRole.value.id, selectedMenuIds.value)
    if (res.code === 200) { message.success('权限已更新'); drawerVisible.value = false }
    else message.error(res.msg || '保存失败')
  } catch { message.error('保存失败') }
  finally { menuSaving.value = false }
}

onMounted(() => fetchRoles())
</script>

<template>
  <div>
    <a-card style="margin-bottom: 16px;">
      <div style="display: flex; justify-content: space-between; align-items: center; flex-wrap: wrap; gap: 12px;">
        <a-input-search v-model:value="searchText" placeholder="搜索角色..." allow-clear style="width: 240px;" @search="handleSearch" />
        <a-button type="primary" @click="openCreate" v-if="userStore.hasPermission('role:manage')">
          <PlusOutlined /> 新增角色
        </a-button>
      </div>
    </a-card>

    <a-card :body-style="{ padding: 0 }">
      <a-table bordered size="middle" :columns="columns" :data-source="roles"
        :pagination="pagination" :loading="loading" :row-class-name="zebraRow"
        :locale="{ emptyText }" row-key="id" @change="handleTableChange">
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'actions'">
            <a-space>
              <a-button type="primary" size="small" @click="openPermissionEditor(record)" v-if="userStore.hasPermission('role:manage')">菜单权限</a-button>
              <a-button size="small" @click="openEdit(record)" v-if="userStore.hasPermission('role:manage')">编辑</a-button>
              <a-button size="small" danger @click="handleDelete(record)" v-if="userStore.hasPermission('role:manage')">删除</a-button>
            </a-space>
          </template>
        </template>
      </a-table>
    </a-card>

    <!-- 新增/编辑弹窗 -->
    <a-modal v-model:open="formVisible" :title="isEdit ? '编辑角色' : '新增角色'" @ok="handleSave" :confirm-loading="saving">
      <a-form :model="form" layout="vertical">
        <a-form-item label="标识" v-if="!isEdit">
          <a-input v-model:value="form.code" placeholder="如 editor" />
        </a-form-item>
        <a-form-item label="名称">
          <a-input v-model:value="form.name" placeholder="如 编辑员" />
        </a-form-item>
        <a-form-item label="描述">
          <a-textarea v-model:value="form.description" placeholder="角色描述" :rows="3" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 菜单权限抽屉 -->
    <a-drawer v-model:open="drawerVisible" :title="`菜单权限 - ${currentRole?.name || ''}`" :width="420" :footer-style="{ textAlign: 'right' }">
      <template #footer>
        <a-button style="margin-right: 8px;" @click="drawerVisible = false">取消</a-button>
        <a-button type="primary" :loading="menuSaving" @click="handleSavePermissions" v-if="userStore.hasPermission('role:manage')">保存</a-button>
      </template>
      <div style="font-size: 13px;">
        <template v-for="menu in menuTree" :key="menu.id">
          <div :style="{ marginBottom: '4px' }">
            <a-checkbox :checked="getCheckedState(menu).checked" :indeterminate="getCheckedState(menu).indeterminate" @change="(e: any) => toggleMenu(menu, e.target.checked)">
              <strong>{{ menu.title }}</strong>
            </a-checkbox>
          </div>
          <div v-if="menu.children?.length" style="margin-left: 24px; margin-bottom: 16px; padding-left: 12px; border-left: 1px solid #f0f0f0;">
            <template v-for="child in menu.children" :key="child.id">
              <div :style="{ marginBottom: '4px' }">
                <a-checkbox :checked="getCheckedState(child).checked" :indeterminate="getCheckedState(child).indeterminate" @change="(e: any) => toggleMenu(child, e.target.checked)">
                  <span :style="{ color: '#555' }">{{ child.title }}</span>
                </a-checkbox>
              </div>
              <div v-if="child.children?.length" style="margin-left: 24px; margin-bottom: 8px;">
                <a-checkbox v-for="sub in child.children" :key="sub.id" :checked="selectedMenuIds.includes(sub.id)" @change="(e: any) => toggleMenu(sub, e.target.checked)" style="margin-right: 16px; margin-bottom: 2px;">
                  <span :style="{ color: '#888', fontSize: '12px' }">{{ sub.title }}</span>
                </a-checkbox>
              </div>
            </template>
          </div>
        </template>
      </div>
    </a-drawer>
  </div>
</template>
