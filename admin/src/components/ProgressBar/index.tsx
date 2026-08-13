interface ProgressBarProps {
  percent: number
  status?: 'normal' | 'warning' | 'exception'
}

export default function ProgressBar({ percent, status = 'normal' }: ProgressBarProps) {
  const getColor = () => {
    if (status === 'exception') return 'var(--red)'
    if (status === 'warning') return 'var(--yellow)'
    return 'var(--green)'
  }

  return (
    <div
      style={{
        height: 6,
        background: 'var(--progress-bg)',
        borderRadius: 3,
        overflow: 'hidden',
        flex: 1,
      }}
    >
      <div
        style={{
          height: '100%',
          borderRadius: 3,
          background: getColor(),
          width: `${percent}%`,
          transition: 'width 0.3s',
        }}
      />
    </div>
  )
}
