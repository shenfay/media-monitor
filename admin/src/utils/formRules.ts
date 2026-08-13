import type { Rule } from 'antd/es/form'
import i18n from '@/locales'

/**
 * 统一表单校验规则
 * 各页面直接引用，避免重复定义
 * 使用 getter 延迟求值，确保语言切换后规则文本跟随更新
 */

export const requiredRule: Rule = { required: true, message: i18n.t('required') }

export function getEmailRules(): Rule[] {
  return [
    { required: true, message: i18n.t('pleaseEnter', { field: i18n.t('email') }) },
    { type: 'email', message: i18n.t('emailInvalid') },
  ]
}

export function getPasswordRules(): Rule[] {
  return [
    { required: true, message: i18n.t('pleaseEnter', { field: i18n.t('password') }) },
    { min: 8, message: i18n.t('passwordMinLength') },
  ]
}

export function getNameRules(): Rule[] {
  return [
    { required: true, message: i18n.t('pleaseEnter', { field: i18n.t('name') }) },
    { max: 50, message: i18n.t('nameMaxLength') },
  ]
}

/**
 * @deprecated 使用 getEmailRules() 代替
 * 保留兼容旧引用，但新代码应使用函数形式
 */
export const emailRules: Rule[] = getEmailRules()

/**
 * @deprecated 使用 getPasswordRules() 代替
 */
export const passwordRules: Rule[] = getPasswordRules()

/**
 * @deprecated 使用 getNameRules() 代替
 */
export const nameRules: Rule[] = getNameRules()
