<script setup lang="ts">
import { ref, onMounted, computed } from "vue"
import { Layout, Menu, Button, Card, message } from "ant-design-vue"
import { UserOutlined, EditOutlined, SaveOutlined } from "@ant-design/icons-vue"
import { useAuthStore } from "@/stores/auth"
import { userApi } from "@/services/user"
import { formatDateTimeFull } from "@/utils/dateFormat"

const { Sider, Content } = Layout
const authStore = useAuthStore()
const activeMenu = ref("profile")
const isEditing = ref(false)
const profileForm = ref({
  nickname: ""
})

const menuItems = computed(() => {
  const items = [{ key: "profile", label: "个人信息", icon: UserOutlined }]
  return items
})

const loadProfile = async () => {
  try {
    const response = await userApi.getProfile()
    if (response.code === 200 && response.data) {
      profileForm.value.nickname = response.data.nickname || ""
    }
  } catch (error) {
    console.error("获取用户信息失败:", error)
  }
}

const toggleEdit = () => {
  isEditing.value = !isEditing.value
}

const handleSave = async () => {
  // 昵称非空校验
  const trimmedNickname = profileForm.value.nickname.trim()
  if (!trimmedNickname) {
    message.error("昵称不能为空")
    return
  }

  try {
    const response = await userApi.updateProfile({ nickname: trimmedNickname })
    if (response.code === 200) {
      message.success("更新成功")
      if (authStore.user) {
        authStore.user.nickname = trimmedNickname
      }
      profileForm.value.nickname = trimmedNickname
      isEditing.value = false
    } else {
      message.error(response.msg || "更新失败")
    }
  } catch (error) {
    message.error("保存失败")
  }
}

const handleMenuClick = (info: { key: string | number }) => {
  const key = String(info.key)
  activeMenu.value = key
}

onMounted(() => {
  loadProfile()
})
</script>

<template>
  <div class="h-full bg-[#262626] overflow-hidden">
    <div class="h-full px-8 py-6">
      <div class="h-full max-w-6xl mx-auto">
        <Layout style="background: transparent; height: 100%">
          <Sider
            width="180"
            style="background: #343232; border-radius: 12px; overflow: hidden"
          >
            <div
              class="text-white text-lg font-bold p-4 border-b border-gray-700"
            >
              个人空间
            </div>
            <Menu
              mode="inline"
              :selected-keys="[activeMenu]"
              style="background: transparent; border-right: none"
              @click="handleMenuClick"
            >
              <Menu.Item v-for="item in menuItems" :key="item.key">
                <component
                  :is="item.icon"
                  class="text-gray-300 inline-block align-middle"
                />
                <span class="text-gray-300">{{ item.label }}</span>
              </Menu.Item>
            </Menu>
          </Sider>

          <Content
            style="
              background: transparent;
              padding-left: 20px;
              overflow-y: auto;
            "
          >
            <div v-if="activeMenu === 'profile'">
              <Card
                style="
                  background: #343232;
                  border-color: #4a4848;
                  border-radius: 12px;
                "
                title="个人信息"
              >
                <template #extra>
                  <Button
                    v-if="isEditing"
                    type="primary"
                    style="
                      background: #6366f1;
                      border-color: #6366f1;
                      border-radius: 8px;
                    "
                    @click="handleSave"
                  >
                    <SaveOutlined /> 保存
                  </Button>
                  <Button v-else type="default" @click="toggleEdit">
                    <EditOutlined /> 编辑
                  </Button>
                </template>
                <form class="space-y-6">
                  <div class="flex gap-4 items-center">
                    <label class="w-24 text-gray-400">昵称</label>
                    <input
                      type="text"
                      v-model="profileForm.nickname"
                      :disabled="!isEditing"
                      class="flex-1 px-4 py-2 bg-gray-700 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-purple-500 disabled:opacity-50"
                    />
                  </div>
                  <div class="flex gap-4 items-center">
                    <label class="w-24 text-gray-400">用户名</label>
                    <span class="text-gray-300">{{
                      authStore.user?.username || "-"
                    }}</span>
                  </div>
                  <div class="flex gap-4 items-center">
                    <label class="w-24 text-gray-400">角色</label>
                    <span class="text-gray-300">{{
                      authStore.user?.role === "admin" ? "管理员" : "普通用户"
                    }}</span>
                  </div>
                  <div class="flex gap-4 items-center">
                    <label class="w-24 text-gray-400">注册时间</label>
                    <span class="text-gray-400">
                      {{ formatDateTimeFull(authStore.user?.createdAt) }}
                    </span>
                  </div>
                </form>
              </Card>
            </div>
          </Content>
        </Layout>
      </div>
    </div>
  </div>
</template>
