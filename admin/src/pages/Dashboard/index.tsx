import { useTranslation } from 'react-i18next'
import { Card, Row, Col, Statistic } from 'antd'
import {
  TeamOutlined,
  AimOutlined,
  FileTextOutlined,
  CheckCircleOutlined,
  StarOutlined,
  ShopOutlined,
} from '@ant-design/icons'

const statCardStyle = { borderRadius: 'var(--radius-md)', borderColor: 'var(--border-light)' } as React.CSSProperties

export default function Dashboard() {
  const { t } = useTranslation()

  return (
    <div style={{ padding: '20px 28px' }}>
      <Row gutter={[16, 16]}>
        <Col xs={24} sm={12} lg={8}>
          <Card style={statCardStyle}>
            <Statistic title={t('totalFamilies')} value={0} prefix={<TeamOutlined />} />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={8}>
          <Card style={statCardStyle}>
            <Statistic title={t('activeGoals')} value={0} prefix={<AimOutlined />} valueStyle={{ color: 'var(--blue)' }} />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={8}>
          <Card style={statCardStyle}>
            <Statistic title={t('todayCardSubmissions')} value={0} prefix={<FileTextOutlined />} />
          </Card>
        </Col>
      </Row>
      <Row gutter={[16, 16]} style={{ marginTop: 16 }}>
        <Col xs={24} sm={12} lg={8}>
          <Card style={statCardStyle}>
            <Statistic title={t('pendingAcceptance')} value={0} prefix={<CheckCircleOutlined />} valueStyle={{ color: 'var(--yellow)' }} />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={8}>
          <Card style={statCardStyle}>
            <Statistic title={t('todayPointsIssued')} value={0} prefix={<StarOutlined />} valueStyle={{ color: 'var(--green)' }} />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={8}>
          <Card style={statCardStyle}>
            <Statistic title={t('pendingExchanges')} value={0} prefix={<ShopOutlined />} valueStyle={{ color: 'var(--text-secondary)' }} />
          </Card>
        </Col>
      </Row>
    </div>
  )
}
