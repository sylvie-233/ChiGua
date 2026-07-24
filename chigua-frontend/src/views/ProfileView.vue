<script setup lang="ts">
import { ref, onMounted, computed } from "vue"
import { Layout, Menu, Button, Card, message, Modal } from "ant-design-vue"
import { UserOutlined, EditOutlined, SaveOutlined, UploadOutlined } from "@ant-design/icons-vue"
import { Cropper } from "vue-advanced-cropper"
import "vue-advanced-cropper/dist/style.css"
import { useAuthStore } from "@/stores/auth"
import { userApi } from "@/services/user"
import { formatDateTimeFull } from "@/utils/dateFormat"

const { Sider, Content } = Layout
const authStore = useAuthStore()
const activeMenu = ref("profile")
const isEditing = ref(false)
const uploadingAvatar = ref(false)
const profileForm = ref({
  nickname: ""
})

// 裁剪相关
const cropModalVisible = ref(false)
const imageToCrop = ref<string | null>(null)
const cropperRef = ref<InstanceType<typeof Cropper> | null>(null)

const menuItems = computed(() => {
  const items = [{ key: "profile", label: "个人信息", icon: UserOutlined }]
  return items
})

const loadProfile = async () => {
  try {
    const response = await userApi.getProfile()
    if (response.code === 200 && response.data) {
      profileForm.value.nickname = response.data.nickname || ""
      if (authStore.user) {
        authStore.user.nickname = response.data.nickname || ""
        authStore.user.createdAt = response.data.createdAt || authStore.user.createdAt
        authStore.user.updateAt = response.data.updateAt || authStore.user.updateAt
        authStore.user.roles = response.data.roles || authStore.user.roles
        authStore.user.avatar = response.data.avatar || authStore.user.avatar
      }
    }
  } catch (error) {
    console.error("获取用户信息失败:", error)
  }
}

const toggleEdit = () => {
  isEditing.value = !isEditing.value
}

// 选择文件 → 打开裁剪弹窗
const handleFileSelect = (event: Event) => {
  const target = event.target as HTMLInputElement
  const file = target.files?.[0]
  if (!file) return

  const reader = new FileReader()
  reader.onload = (e) => {
    imageToCrop.value = e.target?.result as string
    cropModalVisible.value = true
  }
  reader.readAsDataURL(file)

  // 重置 input 以便再次选择同一文件
  target.value = ""
}

// 裁剪确认 → 上传
const handleCropConfirm = async () => {
  const cropper = cropperRef.value
  if (!cropper) {
    message.error("裁剪组件未就绪")
    return
  }

  uploadingAvatar.value = true
  try {
    const result = cropper.getResult()
    if (!result?.canvas) {
      message.error("获取裁剪结果失败")
      return
    }

    const canvas = result.canvas as HTMLCanvasElement
    // dataURL → blob（兼容性更好）
    const dataUrl = canvas.toDataURL("image/jpeg", 0.9)
    const res = await fetch(dataUrl)
    const blob = await res.blob()
    const file = new File([blob], "avatar.jpg", { type: "image/jpeg" })

    const uploadRes = await userApi.uploadFile(file)

    if (uploadRes.code === 200 && uploadRes.data) {
      const avatarUrl = uploadRes.data.fileUrl || uploadRes.data
      const avatarRes = await userApi.updateAvatar(avatarUrl)
      if (avatarRes.code === 200) {
        if (authStore.user) {
          authStore.user.avatar = avatarUrl
        }
        message.success("头像更新成功")
        cropModalVisible.value = false
        imageToCrop.value = null
      } else {
        message.error(avatarRes.msg || "更新头像失败")
      }
    } else {
      message.error(uploadRes.msg || "文件上传失败")
    }
  } catch (error) {
    console.error("头像上传:", error)
    message.error("头像上传失败")
  } finally {
    uploadingAvatar.value = false
  }
}

const handleCropCancel = () => {
  cropModalVisible.value = false
  imageToCrop.value = null
}

const handleSave = async () => {
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
                <!-- 头像区域 -->
                <div class="flex justify-center mb-8">
                  <div class="relative group">
                    <img
                      v-if="authStore.user?.avatar"
                      :src="authStore.user.avatar"
                      class="w-24 h-24 rounded-full object-cover border-2 border-gray-600"
                    />
                    <div
                      v-else
                      class="w-24 h-24 rounded-full flex items-center justify-center text-3xl font-bold text-white"
                      style="background-color: #1e8d77"
                    >
                      {{ authStore.user?.nickname?.charAt(0)?.toUpperCase() || "?" }}
                    </div>
                    <label
                      class="absolute inset-0 rounded-full bg-black/50 flex items-center justify-center text-white cursor-pointer opacity-0 group-hover:opacity-100 transition-opacity"
                    >
                      <UploadOutlined class="text-xl" />
                      <input
                        type="file"
                        accept="image/*"
                        class="hidden"
                        @change="handleFileSelect"
                      />
                    </label>
                  </div>
                </div>
                <form class="space-y-6">
                  <div class="flex gap-4 items-center">
                    <label class="w-24 text-gray-400">昵称</label>
                    <input
                      type="text"
                      v-model="profileForm.nickname"
                      :disabled="!isEditing"
                      maxlength="10"
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
                    <label class="w-24 text-gray-400 shrink-0">角色</label>
                    <div class="flex flex-wrap gap-2">
                      <span
                        v-for="role in authStore.user?.roles"
                        :key="role"
                        class="px-3 py-0.5 text-xs rounded-full border"
                        :class="role === 'admin'
                          ? 'text-red-400 border-red-500/40 bg-red-500/10'
                          : role === 'reviewer'
                            ? 'text-blue-400 border-blue-500/40 bg-blue-500/10'
                            : 'text-gray-300 border-gray-500/40 bg-gray-500/10'"
                      >
                        {{ role === 'admin' ? '管理员' : role === 'reviewer' ? '审核员' : '普通用户' }}
                      </span>
                      <span v-if="!authStore.user?.roles?.length" class="text-gray-500">-</span>
                    </div>
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

    <!-- 头像裁剪弹窗 -->
    <Modal
      v-model:open="cropModalVisible"
      title="裁剪头像"
      :footer="null"
      width="520px"
      :bodyStyle="{ padding: '16px' }"
      destroyOnClose
    >
      <div class="flex flex-col gap-4">
        <div class="bg-gray-800 rounded-lg overflow-hidden" style="height: 360px">
          <Cropper
            v-if="imageToCrop"
            ref="cropperRef"
            :src="imageToCrop"
            :stencil-props="{ aspectRatio: 1, movable: true, resizable: true }"
            :default-boundaries="'fit'"
          />
        </div>
        <div class="flex justify-end gap-3">
          <Button @click="handleCropCancel">取消</Button>
          <Button
            type="primary"
            :loading="uploadingAvatar"
            @click="handleCropConfirm"
          >
            确认
          </Button>
        </div>
      </div>
    </Modal>
  </div>
</template>
