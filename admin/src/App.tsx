import { useTranslation } from 'react-i18next'
import { ConfigProvider } from 'antd'
import zhCN from 'antd/locale/zh_CN'
import enUS from 'antd/locale/en_US'
import { RouterProvider } from 'react-router-dom'
import router from '@/router'
import type { Locale } from 'antd/es/locale'

const antdLocales: Record<string, Locale> = {
  'zh-CN': zhCN,
  'en-US': enUS,
}

export default function App() {
  const { i18n } = useTranslation()
  const antdLocale = antdLocales[i18n.language] || zhCN

  return (
    <ConfigProvider
      locale={antdLocale}
      theme={{
        token: {
          colorPrimary: '#2b2b2b',
          colorError: '#e74c3c',
          borderRadius: 8,
          borderRadiusSM: 6,
          fontFamily: '-apple-system, "PingFang SC", "Helvetica Neue", "Microsoft YaHei", sans-serif',
          fontSize: 13,
          fontSizeSM: 12,
          lineHeight: 1.6,
          colorText: '#2b2b2b',
          colorTextSecondary: '#6b6258',
          colorTextTertiary: '#b0a89a',
          colorTextQuaternary: '#c4bdb0',
          colorBorder: '#e8e2d8',
          colorBorderSecondary: '#efeae2',
          colorFill: '#f5f2ed',
          colorFillSecondary: '#faf8f5',
          colorBgContainer: '#ffffff',
          colorBgLayout: '#ffffff',
          controlHeight: 34,
          wireframe: false,
        },
        components: {
          Button: {
            borderRadius: 8,
            fontWeight: 500,
            controlHeight: 34,
            colorPrimary: '#2b2b2b',
            colorPrimaryHover: '#4d4d4d',
            colorPrimaryActive: '#4d4d4d',
            colorError: '#e74c3c',
            colorErrorHover: '#c0392b',
            colorErrorActive: '#c0392b',
            defaultBorderColor: '#e8e2d8',
            defaultColor: '#6b6258',
            defaultHoverBg: '#f5f2ed',
            defaultHoverBorderColor: '#d4cdc0',
            defaultHoverColor: '#6b6258',
            primaryShadow: 'none',
          },
          Input: {
            borderRadius: 8,
            controlHeight: 34,
            colorBorder: '#e8e2d8',
            activeBorderColor: '#2b2b2b',
            hoverBorderColor: '#2b2b2b',
            colorBgContainer: '#ffffff',
            activeShadow: 'none',
          },
          Select: {
            borderRadius: 8,
            controlHeight: 34,
            colorBorder: '#e8e2d8',
            colorBgContainer: '#ffffff',
            optionSelectedBg: '#E4E0D8',
          },
          Table: {
            borderRadius: 8,
            fontSize: 13,
            headerBg: '#faf8f5',
            headerColor: '#8a8276',
            headerSplitColor: '#efeae2',
            borderColor: '#f5f2ed',
            rowHoverBg: '#faf8f5',
            cellPaddingBlock: 13,
            cellPaddingInline: 14,
            headerBorderRadius: 8,
            fontWeightStrong: 600,
          },
          Menu: {
            itemBorderRadius: 8,
            itemMarginBlock: 2,
            itemMarginInline: 8,
            itemColor: '#6b6258',
            itemSelectedBg: '#E4E0D8',
            itemSelectedColor: '#2b2b2b',
            itemHoverBg: '#E4E0D8',
            itemHoverColor: '#2b2b2b',
            subMenuItemBg: 'transparent',
            iconSize: 16,
          },
          Card: {
            borderRadius: 12,
            colorBorderSecondary: '#efeae2',
            headerFontSize: 14,
          },
          Modal: {
            borderRadius: 12,
          },
          Tag: {
            borderRadiusSM: 6,
            defaultBg: 'transparent',
            colorBorder: 'transparent',
          },
          Pagination: {
            borderRadius: 6,
            itemActiveBg: '#2b2b2b',
            colorPrimary: '#2b2b2b',
            colorPrimaryHover: '#4d4d4d',
          },
          Breadcrumb: {
            fontSize: 13,
            separatorColor: '#b0a89a',
            itemColor: '#b0a89a',
            linkColor: '#b0a89a',
            lastItemColor: '#6b6258',
          },
          Form: {
            labelFontSize: 13,
          },
          Dropdown: {
            borderRadius: 8,
          },
          Layout: {
            bodyBg: '#ffffff',
            headerBg: '#ffffff',
            siderBg: '#F5F3EF',
          },
        },
      }}
    >
      <RouterProvider router={router} />
    </ConfigProvider>
  )
}
