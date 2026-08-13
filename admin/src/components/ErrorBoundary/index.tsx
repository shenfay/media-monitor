import { Component, type ReactNode, type ErrorInfo } from 'react'
import { Button, Result } from 'antd'
import { withTranslation, type WithTranslation } from 'react-i18next'

interface Props extends WithTranslation {
  children: ReactNode
  fallback?: ReactNode
}

interface State {
  hasError: boolean
  error: Error | null
}

/**
 * React 错误边界
 * 捕获子组件树中的渲染错误，防止整个应用白屏。
 * 生产环境展示友好提示，开发环境保留控制台错误信息。
 */
export default withTranslation()(class ErrorBoundary extends Component<Props, State> {
  state: State = { hasError: false, error: null }

  static getDerivedStateFromError(error: Error): State {
    return { hasError: true, error }
  }

  componentDidCatch(error: Error, info: ErrorInfo) {
    // 生产环境可上报到监控平台
    console.error('[ErrorBoundary]', error, info)
  }

  handleReset = () => {
    this.setState({ hasError: false, error: null })
  }

  render() {
    if (this.state.hasError) {
      if (this.props.fallback) return this.props.fallback

      return (
        <Result
          status="error"
          title={this.props.t('pageError')}
          subTitle={import.meta.env.DEV ? this.state.error?.message : this.props.t('pageErrorSubtitle')}
          extra={
            <Button type="primary" onClick={this.handleReset}>
              {this.props.t('reload')}
            </Button>
          }
        />
      )
    }

    return this.props.children
  }
})
