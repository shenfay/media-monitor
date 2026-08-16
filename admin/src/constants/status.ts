/**
 * 全局共享状态常量
 *
 * 集中管理各页面的状态颜色、状态翻译 Key、平台标签等，
 * 避免跨文件重复定义。
 */

/* ─── 任务状态（TaskRun） ─── */
export const TASK_STATUS_COLORS: Record<string, string> = {
  queued: 'default',
  running: 'processing',
  success: 'success',
  partial: 'warning',
  failed: 'error',
}

export const TASK_STATUS_KEYS: Record<string, string> = {
  queued: 'crawlQueued',
  running: 'crawlRunning',
  success: 'success',
  partial: 'crawlPartial',
  failed: 'failed',
}

/* ─── 文章状态（Article） ─── */
export const ARTICLE_STATUS_COLORS: Record<string, string> = {
  pending: 'orange',
  completed: 'green',
  failed: 'red',
}

/**
 * 返回文章状态的翻译标签映射
 * 需要在组件内调用以获取最新的 t() 绑定
 */
export function getArticleStatusLabel(
  t: (key: string) => string,
): Record<string, string> {
  return {
    pending: t('crawlStatusPending'),
    completed: t('crawlStatusCompleted'),
    failed: t('crawlStatusFailed'),
  }
}

/* ─── Worker 状态 ─── */
export const WORKER_STATUS_COLORS: Record<string, string> = {
  online: 'success',
  offline: 'default',
  paused: 'warning',
}

export const WORKER_STATUS_KEYS: Record<string, string> = {
  online: 'workerOnline',
  offline: 'workerOffline',
  paused: 'workerPaused',
}

/* ─── 平台类型标签 ─── */
/**
 * 返回平台类型的翻译标签映射
 * 需要在组件内调用以获取最新的 t() 绑定
 */
export function getPlatformLabels(
  t: (key: string) => string,
): Record<string, string> {
  return {
    news: t('crawlNews'),
    social: t('crawlSocial'),
    social_overseas: t('crawlSocialOverseas'),
  }
}
