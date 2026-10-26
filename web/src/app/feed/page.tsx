'use client'

import { useState, useEffect, useCallback } from 'react'
import Link from 'next/link'
import { get } from '@/lib/api'
import { TYPE_FILTER_OPTIONS, TYPE_LABELS, TYPE_BORDER_COLORS, MOOD_OPTIONS } from '@/lib/constants'
import { timeAgo } from '@/lib/utils'
import { stripHTML } from '@/lib/sanitize'
import Icon from '@/components/Icon'

interface Post {
  id: number
  user_id: number
  type: string
  title: string
  content: string
  like_count: number
  comment_count: number
  is_top?: boolean
  is_featured?: boolean
  is_anonymous?: boolean
  created_at: string
}

interface CheckinItem {
  id: number
  user_id: number
  mood: string
  content: string
  images: string
  is_public: boolean
  created_at: string
}

/** 区段标题：一条细线 + 一个小标题，用来分隔打卡与帖子 */
function SectionLabel({ children }: { children: React.ReactNode }) {
  return (
    <div className="flex items-center gap-3 pt-2">
      <span className="text-xs font-medium tracking-wide text-ink-faint">{children}</span>
      <div className="flex-1 h-px bg-line-soft" />
    </div>
  )
}

function CheckinCard({ checkin }: { checkin: CheckinItem }) {
  const mood = MOOD_OPTIONS.find((m) => m.key === checkin.mood)
  const images: string[] = checkin.images ? JSON.parse(checkin.images) : []

  return (
    <div className="card p-5">
      <div className="flex items-center gap-3">
        {/* 心情用 emoji 是合理的——它表达的就是表情本身，而非功能图标 */}
        <span className="text-2xl leading-none" aria-hidden>{mood?.emoji || '😊'}</span>
        <div className="min-w-0">
          <div className="text-sm font-medium text-ink">今日打卡</div>
          <div className="text-xs text-ink-faint mt-0.5">心情 · {mood?.label || '未知'}</div>
        </div>
        <span className="ml-auto text-xs text-ink-faint shrink-0">{timeAgo(checkin.created_at)}</span>
      </div>

      {checkin.content && (
        <p className="mt-3 text-sm text-ink-soft">{checkin.content}</p>
      )}

      {images.length > 0 && (
        <div className="mt-3 flex gap-2">
          {images.map((url, i) => (
            <img
              key={i}
              src={url}
              alt={`打卡图片 ${i + 1}`}
              className="w-16 h-16 object-cover rounded-lg border border-line"
            />
          ))}
        </div>
      )}
    </div>
  )
}

function PostRow({ post, currentUserId }: { post: Post; currentUserId: number | null }) {
  const type = TYPE_LABELS[post.type]

  return (
    <Link href={`/posts/${post.id}`} className="card block p-5 group">
      {/* 标记行：只在有标记时才占位，避免空行破坏节奏 */}
      {(post.is_top || post.is_featured || (post.is_anonymous && currentUserId === post.user_id)) && (
        <div className="flex items-center gap-1.5 mb-2">
          {post.is_top && (
            <span className="inline-flex items-center gap-1 text-xs text-emotion-ink">
              <Icon name="pin" size={12} />置顶
            </span>
          )}
          {post.is_featured && (
            <span className="inline-flex items-center gap-1 text-xs text-insight-ink">
              <Icon name="star" size={12} />精华
            </span>
          )}
          {post.is_anonymous && currentUserId === post.user_id && (
            <span className="inline-flex items-center gap-1 text-xs text-ink-faint">
              <Icon name="mask" size={12} />匿名
            </span>
          )}
        </div>
      )}

      <h3 className="text-[17px] font-medium text-ink leading-snug group-hover:text-accent transition-colors line-clamp-2">
        {post.title}
      </h3>

      <p className="mt-2 text-sm text-ink-soft line-clamp-2">
        {stripHTML(post.content)}
      </p>

      {/* 元信息行：类型标签靠左，互动数据靠右，视觉重量比标题轻两级 */}
      <div className="mt-4 flex items-center gap-3 text-xs text-ink-faint">
        {type && (
          <span className={`inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full border ${TYPE_BORDER_COLORS[post.type]}`}>
            <Icon name={type.icon} size={12} />
            {type.label}
          </span>
        )}
        <span className="inline-flex items-center gap-1">
          <Icon name="heart" size={13} />{post.like_count}
        </span>
        <span className="inline-flex items-center gap-1">
          <Icon name="comment" size={13} />{post.comment_count}
        </span>
        <span className="ml-auto">{timeAgo(post.created_at)}</span>
      </div>
    </Link>
  )
}

export default function FeedPage() {
  const [posts, setPosts] = useState<Post[]>([])
  const [checkins, setCheckins] = useState<CheckinItem[]>([])
  const [type, setType] = useState('')
  const [sort, setSort] = useState<'latest' | 'hot'>('latest')
  const [page, setPage] = useState(1)
  const [total, setTotal] = useState(0)
  const [loading, setLoading] = useState(false)
  const [currentUserId, setCurrentUserId] = useState<number | null>(null)

  useEffect(() => {
    const stored = localStorage.getItem('user')
    if (stored) {
      try { setCurrentUserId(JSON.parse(stored).id) } catch { /* ignore */ }
    }
  }, [])

  const fetchPosts = useCallback(async () => {
    setLoading(true)
    try {
      const params = new URLSearchParams({ page: String(page), page_size: '10', sort })
      let endpoint = '/posts'

      if (type === 'following') {
        endpoint = '/posts/following'
      } else if (type === 'cohort') {
        endpoint = '/cohort/posts'
      } else if (type) {
        params.set('type', type)
      }

      const [postsRes, checkinsRes] = await Promise.all([
        get<{ items: Post[]; total: number }>(`${endpoint}?${params}`),
        type === '' ? get<{ items: CheckinItem[] }>('/checkin/timeline?page=1&page_size=5') : Promise.resolve({ code: -1, data: { items: [] } }),
      ])

      if (postsRes.code === 0) {
        setPosts(postsRes.data.items || [])
        setTotal(postsRes.data.total)
      }
      if (checkinsRes.code === 0) {
        setCheckins(checkinsRes.data.items || [])
      }
    } catch (err) {
      console.error('获取帖子失败:', err)
    } finally {
      setLoading(false)
    }
  }, [type, sort, page])

  useEffect(() => { fetchPosts() }, [fetchPosts])

  const totalPages = Math.ceil(total / 10)
  const showCheckins = page === 1 && type === '' && checkins.length > 0

  return (
    <div className="mx-auto max-w-2xl px-4 py-10">
      {/* 筛选：横向滚动而非换行，避免在窄屏挤成两三行把内容推下去 */}
      <div className="flex items-center gap-3 mb-6">
        <div className="scroll-x-hidden flex-1 min-w-0 overflow-x-auto -mx-1 px-1">
          <div className="flex gap-1.5 w-max">
            {TYPE_FILTER_OPTIONS.map((t) => {
              const active = type === t.key
              return (
                <button
                  key={t.key}
                  onClick={() => { setType(t.key); setPage(1) }}
                  className={`inline-flex items-center gap-1.5 px-3 py-1.5 rounded-full text-[13px] whitespace-nowrap border transition-colors ${
                    active
                      ? 'bg-accent text-surface border-accent'
                      : 'bg-surface text-ink-soft border-line hover:border-ink-faint'
                  }`}
                >
                  <Icon name={t.icon} size={13} />
                  {t.label}
                </button>
              )
            })}
          </div>
        </div>

        <div className="flex items-center gap-0.5 shrink-0 rounded-full border border-line bg-surface p-0.5">
          {(['latest', 'hot'] as const).map((s) => (
            <button
              key={s}
              onClick={() => { setSort(s); setPage(1) }}
              className={`px-3 py-1 rounded-full text-xs transition-colors ${
                sort === s ? 'bg-surface-sunk text-ink font-medium' : 'text-ink-faint hover:text-ink-soft'
              }`}
            >
              {s === 'latest' ? '最新' : '最热'}
            </button>
          ))}
        </div>
      </div>

      {/* 列表 */}
      {loading ? (
        <div className="space-y-3">
          {[...Array(3)].map((_, i) => (
            <div key={i} className="card p-5">
              <div className="animate-pulse space-y-3">
                <div className="h-4 rounded bg-surface-sunk w-3/5" />
                <div className="h-3 rounded bg-surface-sunk w-full" />
                <div className="h-3 rounded bg-surface-sunk w-4/5" />
                <div className="h-5 rounded-full bg-surface-sunk w-24 mt-2" />
              </div>
            </div>
          ))}
        </div>
      ) : posts.length === 0 && checkins.length === 0 ? (
        <div className="card px-6 py-16 text-center">
          <p className="text-ink-soft">这里还很安静</p>
          <p className="mt-1.5 text-sm text-ink-faint">成为第一个开口的人吧，不必说得漂亮。</p>
          <Link href="/posts/create" className="btn-primary mt-6 text-sm">
            写点什么
          </Link>
        </div>
      ) : (
        <div className="space-y-3">
          {showCheckins && (
            <>
              <SectionLabel>今日打卡</SectionLabel>
              {checkins.slice(0, 3).map((c) => (
                <CheckinCard key={`checkin-${c.id}`} checkin={c} />
              ))}
              <SectionLabel>社区帖子</SectionLabel>
            </>
          )}

          {posts.map((post) => (
            <PostRow key={post.id} post={post} currentUserId={currentUserId} />
          ))}
        </div>
      )}

      {/* 分页 */}
      {totalPages > 1 && (
        <div className="mt-8 flex items-center justify-center gap-3 text-sm">
          <button
            onClick={() => setPage((p) => Math.max(1, p - 1))}
            disabled={page === 1}
            className="px-3 py-1.5 rounded-lg border border-line bg-surface text-ink-soft disabled:opacity-40 hover:bg-surface-sunk transition-colors"
          >
            上一页
          </button>
          <span className="text-ink-faint text-xs">{page} / {totalPages}</span>
          <button
            onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
            disabled={page === totalPages}
            className="px-3 py-1.5 rounded-lg border border-line bg-surface text-ink-soft disabled:opacity-40 hover:bg-surface-sunk transition-colors"
          >
            下一页
          </button>
        </div>
      )}
    </div>
  )
}
