import { reactive } from 'vue'

/** 中文分页配置 */
export const paginationLocale = {
  items_per_page: '条/页',
  jump_to: '跳至',
  jump_to_confirm: '确定',
  page: '页',
  prev_page: '上一页',
  next_page: '下一页',
  prev_5: '向前 5 页',
  next_5: '向后 5 页',
  prev_3: '向前 3 页',
  next_3: '向后 3 页'
}

/** 创建默认分页配置 */
export function createPagination(pageSize = 10) {
  return reactive({
    current: 1,
    pageSize,
    total: 0,
    showSizeChanger: true,
    showQuickJumper: true,
    showTotal: (total: number) => `共 ${total} 条记录`,
    locale: paginationLocale
  })
}

/** 斑马纹行样式 */
export function zebraRow(_record: unknown, index: number): string {
  return index % 2 === 1 ? 'row-striped' : ''
}

/** 表格空状态文案 */
export const emptyText = '暂无数据'
