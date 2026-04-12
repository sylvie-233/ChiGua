import dayjs from "dayjs"

/**
 * 格式化日期为YYYY-MM-DD格式
 */
export const formatDate = (date: Date): string => {
  return dayjs(date).format("YYYY-MM-DD")
}

/**
 * 格式化日期为YYYY-MM-DD格式
 */
export const formatDateWithChinese = (date: Date): string => {
  return dayjs(date).format("YYYY年MM月DD日")
}

/**
 * 格式化日期时间为YYYY-MM-DD HH:mm格式
 */
export const formatDateTime = (date: Date): string => {
  return dayjs(date).format("YYYY-MM-DD HH:mm")
}
