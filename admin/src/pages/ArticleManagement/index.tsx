import { useState, useEffect, useCallback } from 'react'
import { useTranslation } from 'react-i18next'
import { Table, Tag, Button, Select, Input, message, Drawer, Descriptions, Typography } from 'antd'
import { EyeOutlined } from '@ant-design/icons'
import DataPanel, { FilterSearch } from '@/components/DataPanel'
import { getArticles, getArticle, getSources, type Article, type Source } from '@/services/crawl'

const { Paragraph } = Typography

function formatTime(dateStr: string | null): string {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  if (isNaN(d.getTime())) return '-'
  return d.toLocaleString('zh-CN', { year: 'numeric', month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit', hour12: false })
}

export default function ArticleManagement() {
  const { t } = useTranslation()
  const [articles, setArticles] = useState<Article[]>([])
  const [sources, setSources] = useState<Source[]>([])
  const [loading, setLoading] = useState(false)
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(20)

  // Filters
  const [sourceFilter, setSourceFilter] = useState('')
  const [platformFilter, setPlatformFilter] = useState('')
  const [keyword, setKeyword] = useState('')

  // Detail
  const [detailOpen, setDetailOpen] = useState(false)
  const [detailArticle, setDetailArticle] = useState<Article | null>(null)
  const [detailLoading, setDetailLoading] = useState(false)

  const sourceMap = Object.fromEntries(sources.map(s => [s.id, s.name]))

  const fetchData = useCallback(async () => {
    setLoading(true)
    try {
      const params: Record<string, unknown> = {
        limit: pageSize,
        offset: (page - 1) * pageSize,
      }
      if (sourceFilter) params.source_id = sourceFilter
      if (platformFilter) params.platform = platformFilter
      if (keyword) params.keyword = keyword

      const res = await getArticles(params as Parameters<typeof getArticles>[0])
      setArticles(res.articles || [])
      setTotal(res.total || 0)
    } catch { /* ignore */ }
    finally { setLoading(false) }
  }, [page, pageSize, sourceFilter, platformFilter, keyword])

  useEffect(() => {
    getSources().then(setSources).catch(() => {})
  }, [])

  useEffect(() => { fetchData() }, [fetchData])

  const handleViewDetail = async (id: string) => {
    setDetailOpen(true)
    setDetailLoading(true)
    try {
      const article = await getArticle(id)
      setDetailArticle(article)
    } catch {
      message.error(t('crawlFetchDetailFailed'))
    } finally {
      setDetailLoading(false)
    }
  }

  const handleSearch = () => {
    setPage(1)
    fetchData()
  }

  const columns = [
    {
      title: t('crawlTitle'), dataIndex: 'title', key: 'title', ellipsis: true,
      render: (v: string, record: Article) => (
        <a onClick={() => handleViewDetail(record.id)} style={{ fontWeight: 500 }}>
          {v || t('crawlNoTitle')}
        </a>
      ),
    },
    {
      title: t('crawlDataSource'), dataIndex: 'source_id', key: 'data_source', width: 150, ellipsis: true,
      render: (_: string, record: Article) => sourceMap[record.source_id] || record.source_name || '-',
    },
    {
      title: t('crawlSource'), dataIndex: 'source_id', key: 'source_id', width: 120, ellipsis: true,
      render: (v: string, record: Article) => record.source_name || sourceMap[v] || v?.slice(0, 8) || '-',
    },
    { title: t('crawlPlatform'), dataIndex: 'platform', key: 'platform', width: 80, render: (v: string) => v ? <Tag>{v}</Tag> : '-' },
    { title: t('crawlAuthor'), dataIndex: 'author', key: 'author', width: 100, ellipsis: true, render: (v: string) => v || '-' },
    { title: t('crawlLanguage'), dataIndex: 'language', key: 'language', width: 60, render: (v: string) => v || '-' },
    { title: t('crawlPublishTime'), dataIndex: 'published_at', key: 'published_at', width: 150, render: (v: string | null) => formatTime(v) },
    { title: t('crawlFetchedAt'), dataIndex: 'fetched_at', key: 'fetched_at', width: 150, render: (v: string) => formatTime(v) },
    {
      title: t('actions'), key: 'actions', width: 110,
      render: (_: unknown, record: Article) => (
        <Button type="link" size="small" onClick={() => handleViewDetail(record.id)}>
          <EyeOutlined style={{ marginRight: 4 }} />{t('crawlView')}
        </Button>
      ),
    },
  ]

  return (
    <div>
      <DataPanel
        title={t('crawlArticleMgmt')}
        filters={
          <>
            <FilterSearch value={keyword} onChange={setKeyword} placeholder={t('crawlSearchTitle')} onSearch={handleSearch} />
            <Select
              value={sourceFilter} onChange={v => { setSourceFilter(v); setPage(1) }} style={{ width: 150 }}
              placeholder={t('crawlAllSources')} allowClear
              options={[{ value: '', label: t('crawlAllSources') }, ...sources.map(s => ({ value: s.id, label: s.name }))]}
            />
            <Select
              value={platformFilter} onChange={v => { setPlatformFilter(v); setPage(1) }} style={{ width: 120 }}
              placeholder={t('crawlAllPlatforms')} allowClear
              options={[
                { value: '', label: t('crawlAllPlatforms') },
                { value: 'news', label: t('crawlNews') },
                { value: 'social', label: t('crawlSocial') },
                { value: 'social_overseas', label: t('crawlSocialOverseas') },
              ]}
            />
            </>
        }
      >
        <Table
          dataSource={articles}
          columns={columns}
          rowKey="id"
          loading={loading}
          size="small"
          pagination={{
            current: page,
            pageSize,
            total,
            showSizeChanger: true,
            showTotal: total => t('total', { total }),
            onChange: (p, ps) => { setPage(p); if (ps !== pageSize) { setPageSize(ps); setPage(1) } },
          }}
        />
      </DataPanel>

      <Drawer title={t('crawlArticleDetail')} open={detailOpen} onClose={() => { setDetailOpen(false); setDetailArticle(null) }} width={600}>
        {detailLoading ? (
          <div style={{ textAlign: 'center', padding: 40 }}>{t('loading')}</div>
        ) : detailArticle ? (
          <div style={{ display: 'flex', flexDirection: 'column', gap: 16 }}>
            <Descriptions column={2} bordered size="small">
              <Descriptions.Item label={t('crawlTitle')} span={2}>{detailArticle.title}</Descriptions.Item>
              <Descriptions.Item label={t('crawlSource')}>{detailArticle.source_name || sourceMap[detailArticle.source_id] || '-'}</Descriptions.Item>
              <Descriptions.Item label={t('crawlPlatform')}>{detailArticle.platform || '-'}</Descriptions.Item>
              <Descriptions.Item label={t('crawlAuthor')}>{detailArticle.author || '-'}</Descriptions.Item>
              <Descriptions.Item label={t('crawlLanguage')}>{detailArticle.language || '-'}</Descriptions.Item>
              <Descriptions.Item label={t('crawlPublishTime')}>{formatTime(detailArticle.published_at)}</Descriptions.Item>
              <Descriptions.Item label={t('crawlFetchedAt')}>{formatTime(detailArticle.fetched_at)}</Descriptions.Item>
              <Descriptions.Item label={t('crawlUrl')} span={2}>
                {detailArticle.url ? <a href={detailArticle.url} target="_blank" rel="noopener noreferrer">{detailArticle.url}</a> : '-'}
              </Descriptions.Item>
            </Descriptions>

            {detailArticle.summary && (
              <div>
                <h4 style={{ marginBottom: 8 }}>{t('crawlSummary')}</h4>
                <Paragraph style={{ background: 'var(--bg-gray)', padding: 12, borderRadius: 6, whiteSpace: 'pre-wrap' }}>
                  {detailArticle.summary}
                </Paragraph>
              </div>
            )}

            {detailArticle.body && (
              <div>
                <h4 style={{ marginBottom: 8 }}>{t('crawlBody')}</h4>
                {detailArticle.body_format === 'html' ? (
                  <div
                    style={{ border: '1px solid var(--border-light)', borderRadius: 6, padding: 12, maxHeight: 400, overflow: 'auto' }}
                    dangerouslySetInnerHTML={{ __html: detailArticle.body }}
                  />
                ) : (
                  <Paragraph style={{ whiteSpace: 'pre-wrap', maxHeight: 400, overflow: 'auto', border: '1px solid var(--border-light)', borderRadius: 6, padding: 12 }}>
                    {detailArticle.body}
                  </Paragraph>
                )}
              </div>
            )}

            {detailArticle.media && (
              <div>
                <h4 style={{ marginBottom: 8 }}>{t('crawlMedia')}</h4>
                <pre style={{ background: 'var(--bg-gray)', padding: 12, borderRadius: 6, fontSize: 12, overflow: 'auto', maxHeight: 200 }}>
                  {JSON.stringify(detailArticle.media, null, 2)}
                </pre>
              </div>
            )}

            {detailArticle.interactions && (
              <div>
                <h4 style={{ marginBottom: 8 }}>{t('crawlInteractions')}</h4>
                <pre style={{ background: 'var(--bg-gray)', padding: 12, borderRadius: 6, fontSize: 12, overflow: 'auto', maxHeight: 200 }}>
                  {JSON.stringify(detailArticle.interactions, null, 2)}
                </pre>
              </div>
            )}
          </div>
        ) : null}
      </Drawer>
    </div>
  )
}
