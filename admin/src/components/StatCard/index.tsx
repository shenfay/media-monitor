import { ArrowUpOutlined, ArrowDownOutlined } from '@ant-design/icons'

interface StatCardProps {
  label: string
  value: string | number
  change?: number
  changeLabel?: string
  color?: string
}

export default function StatCard({ label, value, change, changeLabel, color = 'var(--brand-dark)' }: StatCardProps) {
  const isPositive = change !== undefined && change >= 0
  const isNegative = change !== undefined && change < 0

  return (
    <div
      style={{
        background: 'var(--main-bg)',
        border: '1px solid var(--border-color)',
        borderRadius: 'var(--radius-md)',
        padding: 20,
        transition: 'box-shadow 0.2s',
        cursor: 'default',
      }}
      onMouseEnter={e => {
        e.currentTarget.style.boxShadow = 'var(--shadow-md)'
      }}
      onMouseLeave={e => {
        e.currentTarget.style.boxShadow = 'none'
      }}
    >
      <div
        style={{
          fontSize: 13,
          color: 'var(--text-muted)',
          marginBottom: 8,
          display: 'flex',
          alignItems: 'center',
          gap: 6,
        }}
      >
        <span
          style={{
            width: 8,
            height: 8,
            borderRadius: '50%',
            background: color,
            display: 'inline-block',
          }}
        />
        {label}
      </div>
      <div
        style={{
          fontSize: 28,
          fontWeight: 600,
          color: 'var(--text-primary)',
          lineHeight: 1.2,
        }}
      >
        {value}
      </div>
      {(change !== undefined || changeLabel) && (
        <div
          style={{
            fontSize: 12,
            marginTop: 6,
            display: 'flex',
            alignItems: 'center',
            gap: 4,
            color: isPositive ? 'var(--green)' : isNegative ? 'var(--red)' : 'var(--text-muted)',
          }}
        >
          {isPositive && <ArrowUpOutlined style={{ fontSize: 10 }} />}
          {isNegative && <ArrowDownOutlined style={{ fontSize: 10 }} />}
          {change !== undefined && Math.abs(change)}
          {changeLabel && ` ${changeLabel}`}
        </div>
      )}
    </div>
  )
}
