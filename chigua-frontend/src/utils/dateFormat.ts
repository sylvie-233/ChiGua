import dayjs from "dayjs"
import "dayjs/locale/zh-cn"

// 设置中文语言
dayjs.locale("zh-cn")

/**
 * 格式化日期为YYYY-MM-DD格式
 */
export const formatDate = (
  date: Date | number | string | undefined
): string => {
  if (!date) {
    return "-"
  }
  return dayjs(date).format("YYYY-MM-DD")
}

/**
 * 格式化日期为YYYY-MM-DD格式（中文）
 */
export const formatDateWithChinese = (
  date: Date | number | string | undefined
): string => {
  if (!date) {
    return "-"
  }
  return dayjs(date).format("YYYY年MM月DD日")
}

/**
 * 格式化日期时间为YYYY-MM-DD HH:mm格式
 */
export const formatDateTime = (
  date: Date | number | string | undefined
): string => {
  if (!date) {
    return "-"
  }
  return dayjs(date).format("YYYY-MM-DD HH:mm")
}

/**
 * 格式化日期时间为YYYY-MM-DD HH:mm:ss格式
 */
export const formatDateTimeFull = (
  date: Date | number | string | undefined
): string => {
  if (!date) {
    return "-"
  }
  return dayjs(date).format("YYYY-MM-DD HH:mm:ss")
}

/**
 * 格式化时间为HH:mm:ss格式
 */
export const formatTime = (
  date: Date | number | string | undefined
): string => {
  if (!date) {
    return "-"
  }
  return dayjs(date).format("HH:mm:ss")
}

/**
 * 计算两个日期之间的天数
 */
export const diffDays = (
  startDate: Date | number | string,
  endDate: Date | number | string
): number => {
  return dayjs(endDate).diff(dayjs(startDate), "day")
}

/**
 * 判断是否是今天
 */
export const isToday = (date: Date | number | string | undefined): boolean => {
  if (!date) {
    return false
  }
  return dayjs(date).isSame(dayjs(), "day")
}
