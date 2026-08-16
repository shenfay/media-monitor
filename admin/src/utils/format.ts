/**
 * 通用格式化工具函数
 *
 * 集中管理日期格式化、耗时计算、相对时间等逻辑，
 * 避免各页面重复定义 formatTime 等函数。
 */

/**
 * 格式化日期时间为本地化字符串
 * @param dateStr 日期字符串（ISO 8601 或可解析格式）
 * @param lang    语言标识（如 'zh-CN'、'en-US'），不传则使用浏览器默认
 * @param opts    可选覆盖 toLocaleString 的选项
 */
export function formatTime(
  dateStr: string | null,
  lang?: string,
  opts?: Intl.DateTimeFormatOptions,
): string {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  if (isNaN(d.getTime())) return '-'
  return d.toLocaleString(
    lang,
    opts ?? {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      hour12: false,
    },
  )
}

/**
 * 格式化日期时间（含秒），用于操作日志等需要精确时间的场景
 */
export function formatTimeWithSeconds(
  dateStr: string | null,
  lang?: string,
): string {
  return formatTime(dateStr, lang, {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
    hour12: false,
  })
}

/**
 * 格式化日期时间（仅月日时分），用于 Dashboard 等空间紧凑的场景
 */
export function formatTimeShort(
  dateStr: string | null,
  lang?: string,
): string {
  return formatTime(dateStr, lang, {
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    hour12: false,
  })
}

/**
 * 计算并格式化时间跨度（耗时）
 * @param start 开始时间
 * @param end   结束时间（null 则以当前时间计算）
 * @param t     i18n 翻译函数
 */
export function formatDuration(
  start: string | null,
  end: string | null,
  t: (key: string, opts?: Record<string, unknown>) => string,
): string {
  if (!start) return '-'
  const s = new Date(start).getTime()
  const e = end ? new Date(end).getTime() : Date.now()
  const diff = Math.floor((e - s) / 1000)
  if (diff < 60) return t('crawlDurSec', { n: diff })
  if (diff < 3600)
    return t('crawlDurMinSec', {
      m: Math.floor(diff / 60),
      s: diff % 60,
    })
  return t('crawlDurHourMin', {
    h: Math.floor(diff / 3600),
    m: Math.floor((diff % 3600) / 60),
  })
}

/**
 * 计算并格式化「多久之前」的相对时间
 * @param dateStr 日期字符串
 * @param t       i18n 翻译函数
 */
export function timeAgo(
  dateStr: string,
  t: (key: string, opts?: Record<string, unknown>) => string,
): string {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  if (isNaN(d.getTime())) return '-'
  const diff = Math.floor((Date.now() - d.getTime()) / 1000)
  if (diff < 60) return t('workerTimeAgoSec', { n: diff })
  if (diff < 3600) return t('workerTimeAgoMin', { n: Math.floor(diff / 60) })
  if (diff < 86400)
    return t('workerTimeAgoHour', { n: Math.floor(diff / 3600) })
  return t('workerTimeAgoDay', { n: Math.floor(diff / 86400) })
}
