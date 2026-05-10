<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { Layout, Menu, Button, Card, Table, Tag, Space, message } from 'ant-design-vue';
import { HomeOutlined, FileTextOutlined, UserOutlined, EditOutlined, DeleteOutlined, PlusOutlined } from '@ant-design/icons-vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '@/stores/auth';
import { articleApi } from '@/services';
import type { Article } from '@/types/api';

const { Sider, Content } = Layout;
const router = useRouter();
const authStore = useAuthStore();
const activeMenu = ref('articles');
const articles = ref<Article[]>([]);
const loading = ref(false);

const menuItems = [
  { key: 'home', label: '数据概览', icon: HomeOutlined },
  { key: 'articles', label: '文章管理', icon: FileTextOutlined },
  { key: 'profile', label: '个人设置', icon: UserOutlined }
];

const columns = [
  {
    title: '标题',
    dataIndex: 'title',
    key: 'title',
    ellipsis: true,
    style: { color: '#e5e7eb' }
  },
  {
    title: '分类',
    dataIndex: 'category',
    key: 'category',
    style: { color: '#e5e7eb' }
  },
  {
    title: '状态',
    dataIndex: 'status',
    key: 'status'
  },
  {
    title: '创建时间',
    dataIndex: 'created_at',
    key: 'created_at',
    style: { color: '#9ca3af' },
    width: 160
  },
  {
    title: '操作',
    key: 'action',
    width: 120
  }
];

const fetchArticles = async () => {
  loading.value = true;
  try {
    const response = await articleApi.getArticleList({ page: 1, pageSize: 10 });
    if (response.code === 200) {
      articles.value = response.data || [];
    }
  } catch (error) {
    message.error('获取文章列表失败');
  } finally {
    loading.value = false;
  }
};

const handleMenuClick = (key: string | number) => {
  activeMenu.value = String(key);
};

const handleEdit = (id: number) => {
  message.info(`编辑文章 ${id}`);
};

const handleDelete = (id: number) => {
  message.info(`删除文章 ${id}`);
};

const handleAddArticle = () => {
  message.info('添加文章');
};

const handleViewArticle = (id: number) => {
  router.push(`/article/${id}`);
};

const formatStatus = (status: number) => {
  return status === 1 ? '已发布' : '草稿';
};

const getStatusColor = (status: number) => {
  return status === 1 ? 'green' : 'orange';
};

onMounted(() => {
  fetchArticles();
});
</script>

<template>
  <div style="height: 100%; background-color: #262626; overflow: hidden;">
    <div class="h-full px-8 py-6">
      <div class="h-full max-w-6xl mx-auto">
        <Layout style="background: transparent; height: 100%;">
          <Sider width="180" style="background: #343232; border-radius: 12px; overflow: hidden;">
            <div class="text-white text-lg font-bold p-4 border-b border-gray-700">
              后台管理
            </div>
            <Menu
              mode="inline"
              :selected-keys="[activeMenu]"
              style="background: transparent; border-right: none;"
              @click="({ key }) => handleMenuClick(key)"
            >
              <Menu.Item
                v-for="item in menuItems"
                :key="item.key"
              >
                <component :is="item.icon" class="text-gray-400" />
                <span class="text-gray-300">{{ item.label }}</span>
              </Menu.Item>
            </Menu>
          </Sider>

          <Content style="background: transparent; padding-left: 20px; overflow-y: auto;">
            <div v-if="activeMenu === 'articles'">
              <Card
                style="background: #343232; border-color: #4a4848; border-radius: 12px;"
                title="文章管理"
              >
                <template #extra>
                  <Button
                    type="primary"
                    style="background: #6366f1; border-color: #6366f1; border-radius: 8px;"
                    @click="handleAddArticle"
                  >
                    <PlusOutlined /> 添加文章
                  </Button>
                </template>
                <Table
                  :data-source="articles"
                  :loading="loading"
                  :columns="columns"
                  :pagination="{ pageSize: 10, style: { color: '#9ca3af' } }"
                  row-key="id"
                  style="background: #343232;"
                >
                  <template #bodyCell="{ column, record }">
                    <template v-if="column.key === 'title'">
                      <a href="#" style="color: #818cf8;" @click.prevent="handleViewArticle(record.id)">
                        {{ record.title }}
                      </a>
                    </template>
                    <template v-else-if="column.key === 'category'">
                      {{ record.category?.name || '-' }}
                    </template>
                    <template v-else-if="column.key === 'status'">
                      <Tag :color="getStatusColor(record.status)">
                        {{ formatStatus(record.status) }}
                      </Tag>
                    </template>
                    <template v-else-if="column.key === 'action'">
                      <Space>
                        <Button
                          type="text"
                          style="color: #818cf8; padding: 4px 8px;"
                          @click="handleEdit(record.id)"
                        >
                          <EditOutlined /> 编辑
                        </Button>
                        <Button
                          type="text"
                          danger
                          style="padding: 4px 8px;"
                          @click="handleDelete(record.id)"
                        >
                          <DeleteOutlined /> 删除
                        </Button>
                      </Space>
                    </template>
                  </template>
                </Table>
              </Card>
            </div>

            <div v-else-if="activeMenu === 'home'">
              <Card style="background: #343232; border-color: #4a4848; border-radius: 12px;" title="数据概览">
                <div class="flex items-center gap-4 mb-8">
                  <div class="w-16 h-16 rounded-full bg-gradient-to-tr from-purple-600 to-red-500 flex items-center justify-center text-white text-2xl font-bold">
                    {{ authStore.user?.nickname?.charAt(0) || '?' }}
                  </div>
                  <div>
                    <h2 class="text-xl font-bold text-white">
                      {{ authStore.user?.nickname || '用户' }}
                    </h2>
                    <p class="text-gray-400">{{ authStore.user?.email }}</p>
                  </div>
                </div>
                <div class="grid grid-cols-3 gap-4">
                  <Card style="background: #262626; border-color: #3d3d3d; border-radius: 8px;" title="文章数量">
                    <p class="text-3xl font-bold text-white">{{ articles.length }}</p>
                  </Card>
                  <Card style="background: #262626; border-color: #3d3d3d; border-radius: 8px;" title="粉丝数">
                    <p class="text-3xl font-bold text-white">0</p>
                  </Card>
                  <Card style="background: #262626; border-color: #3d3d3d; border-radius: 8px;" title="获赞数">
                    <p class="text-3xl font-bold text-white">0</p>
                  </Card>
                </div>
              </Card>
            </div>

            <div v-else-if="activeMenu === 'profile'">
              <Card style="background: #343232; border-color: #4a4848; border-radius: 12px;" title="个人设置">
                <form class="space-y-6">
                  <div class="flex gap-4">
                    <label class="w-24 text-gray-400">昵称</label>
                    <input
                      type="text"
                      :value="authStore.user?.nickname"
                      class="flex-1 px-4 py-2 bg-gray-700 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-purple-500"
                    />
                  </div>
                  <div class="flex gap-4">
                    <label class="w-24 text-gray-400">用户名</label>
                    <input
                      type="text"
                      :value="authStore.user?.username"
                      disabled
                      class="flex-1 px-4 py-2 bg-gray-700 rounded-lg text-gray-500"
                    />
                  </div>
                  <div class="flex gap-4">
                    <label class="w-24 text-gray-400">邮箱</label>
                    <input
                      type="email"
                      :value="authStore.user?.email"
                      class="flex-1 px-4 py-2 bg-gray-700 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-purple-500"
                    />
                  </div>
                  <div class="flex gap-4">
                    <label class="w-24 text-gray-400">注册时间</label>
                    <span class="text-gray-400">
                      {{ authStore.user?.created_at || '-' }}
                    </span>
                  </div>
                  <div class="flex justify-end">
                    <Button type="primary" style="background: #6366f1; border-color: #6366f1; border-radius: 8px;">
                      保存修改
                    </Button>
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

<style scoped>
:deep(.ant-menu-item-selected) {
  background-color: #262626 !important;
  border-right: 2px solid #6b7280;
}

:deep(.ant-menu-item:hover) {
  background-color: #262626 !important;
}

:deep(.ant-card-head-title) {
  color: #d1d5db;
}

:deep(.ant-pagination-item) {
  background-color: #262626;
  border-color: #3d3d3d;
}

:deep(.ant-pagination-item a) {
  color: #9ca3af;
}

:deep(.ant-pagination-item-active) {
  background-color: #6366f1 !important;
  border-color: #6366f1 !important;
}

:deep(.ant-pagination-item-active a) {
  color: white !important;
}

/* 表格header样式 */
:deep(.ant-table-thead > tr > th) {
  background-color: #262626 !important;
  color: #d1d5db !important;
  border-bottom: 1px solid #4a4848 !important;
}

/* 表格body背景 */
:deep(.ant-table-tbody) {
  background-color: #343232 !important;
}

/* 表格容器背景 */
:deep(.ant-table-container) {
  background-color: #343232 !important;
}

/* 表格行样式 - 覆盖所有状态 */
:deep(.ant-table-tbody > tr) {
  background-color: #343232 !important;
  transition: none !important;
}

:deep(.ant-table-tbody > tr):active {
  background-color: #343232 !important;
}

:deep(.ant-table-tbody > tr):focus {
  background-color: #343232 !important;
}

/* 表格悬浮行样式 */
:deep(.ant-table-tbody > tr:hover) {
  background-color: #2a2a2a !important;
  transition: none !important;
}

:deep(.ant-table-tbody > tr:hover) {
  background-color: #2a2a2a !important;
}

:deep(.ant-table-tbody > tr:hover > td) {
  background-color: #2a2a2a !important;
}

/* 表格单元格样式 */
:deep(.ant-table-tbody > tr > td) {
  color: #e5e7eb;
  border-bottom: 1px solid #4a4848;
  background-color: transparent !important;
}

:deep(.ant-table-tbody > tr:hover > td) {
  background-color: transparent !important;
}

/* 禁用表格行过渡动画 */
:deep(.ant-table-row) {
  transition: none !important;
}

:deep(.ant-table-row-hover) {
  background-color: #2a2a2a !important;
}

/* 表格容器禁止横向滚动 */
:deep(.ant-table) {
  overflow-x: hidden;
}

/* 表格内容容器禁止横向滚动 */
:deep(.ant-table-container) {
  overflow-x: hidden;
}

/* 操作按钮区域 */
:deep(.ant-table-tbody > tr > td:last-child) {
  padding: 8px 4px !important;
}

/* Space组件间距 */
:deep(.ant-space) {
  gap: 4px !important;
}

/* 操作按钮样式 */
:deep(.ant-table-tbody .ant-btn-text) {
  padding: 2px 6px !important;
  font-size: 12px;
}
</style>
