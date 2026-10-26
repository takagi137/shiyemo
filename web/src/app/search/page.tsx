'use client'

import { useState, useEffect, useCallback, Suspense } from 'react'
import { useSearchParams, useRouter } from 'next/navigation'
import Link from 'next/link'
import { get } from '@/lib/api'
import { TYPE_LABELS } from '@/lib/constants'
import Avatar from '@/components/Avatar'
import Icon from '@/components/Icon'
import { stripHTML } from '@/lib/sanitize'

interface Post {
  id: number
  user_id: number
  type: string
  title: string
  content: string
  tags?: string[]
  like_count: number
  comment_count: number
  created_at: string
}

interface User {
  id: number
  nickname: string
  bio: string
  follower_count: number
}

const TABS = [
  { key: 'posts', label: '帖子', icon: '📝' },
  { key: 'users', label: '用户', icon: '👤' },
]

export default function SearchPage() {
  return (
    <Suspense fallback={
      <div className="container mx-auto px-4 py-12">
        <div className="animate-pulse h-12 max-w-2xl mx-auto rounded-xl bg-gray-200 mb-8" />
        <div className="space-y-4 max-w-2xl mx-auto">
          {[...Array(3)].map((_, i) => (
            <div key={i} className="card p-6">
              <div className="animate-pulse h-6 bg-gray-200 rounded w-3/4 mb-3" />
              <div className="animate-pulse h-4 bg-gray-200 rounded w-full" />
            </div>
          ))}
        </div>
      </div>
    }>
      <SearchContent />
    </Suspense>
  )
}

function SearchContent() {
  const searchParams = useSearchParams()
  const router = useRouter()

  const [query, setQuery] = useState(searchParams.get('q') || '')
  const [tab, setTab] = useState<'posts' | 'users'>('posts')
  const [posts, setPosts] = useState<Post[]>([])
  const [users, setUsers] = useState<User[]>([])
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [loading, setLoading] = useState(false)
  const [searched, setSearched] = useState(false)
  const [hotTags, setHotTags] = useState<string[]>([])

  useEffect(() => {
    get<{ id: number; name: string; count: number }[]>('/tags').then(res => {
      if (res.code === 0) setHotTags((res.data || []).slice(0, 10).map(t => t.name))
    }).catch(() => {})
  }, [])

  const doSearch = useCallback(async (q: string, t: string, p: number) => {
    if (!q.trim()) return
    setLoading(true)
    setSearched(true)
    try {
      const endpoint = t === 'posts' ? '/search/posts' : '/search/users'
      const res = await get<{ items: (Post | User)[]; total: number }>(`${endpoint}?q=${encodeURIComponent(q)}&page=${p}&size=10`)
      if (res.code === 0) {
        if (t === 'posts') setPosts(res.data.items as Post[])
        else setUsers(res.data.items as User[])
        setTotal(res.data.total)
      }
    } catch (err) {
      console.error('搜索失败:', err)
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    const q = searchParams.get('q')
    if (q) {
      setQuery(q)
      doSearch(q, tab, page)
    }
  }, [searchParams, doSearch, tab, page])

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault()
    if (!query.trim()) return
    setPage(1)
    router.push(`/search?q=${encodeURIComponent(query.trim())}`)
    doSearch(query.trim(), tab, 1)
  }

  const handleTabChange = (t: 'posts' | 'users') => {
    setTab(t)
    setPage(1)
    if (query.trim()) doSearch(query.trim(), t, 1)
  }

  const totalPages = Math.ceil(total / 10)

  return (
    <div className="container mx-auto px-4 py-8 max-w-3xl">
      {/* 搜索框 */}
      <form onSubmit={handleSearch} className="mb-8">
        <div className="flex rounded-xl overflow-hidden shadow-sm border border-gray-200 focus-within:border-blue-500 focus-within:ring-2 focus-within:ring-blue-500/20 transition-all bg-white">
          <input
            type="text"
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            placeholder="搜索帖子或用户..."
            className="flex-1 px-5 py-3.5 border-0 focus:outline-none text-base"
          />
          <button type="submit" className="px-6 bg-blue-600 text-white hover:bg-blue-700 transition-colors font-medium">
            搜索
          </button>
        </div>
      </form>

      {/* Tab 切换 */}
      {searched && (
        <div className="mb-6 flex gap-1 border-b border-gray-200">
          {TABS.map((t) => (
            <button
              key={t.key}
              onClick={() => handleTabChange(t.key as 'posts' | 'users')}
              className={`flex items-center gap-1.5 px-5 py-3 text-sm font-medium border-b-2 transition-colors ${
                tab === t.key
                  ? 'border-blue-600 text-blue-600'
                  : 'border-transparent text-gray-500 hover:text-gray-700'
              }`}
            >
              <span>{t.icon}</span>
              {t.label}
            </button>
          ))}
        </div>
      )}

      {/* 搜索结果 */}
      <div>
        {loading ? (
          <div className="space-y-4">
            {[...Array(3)].map((_, i) => (
              <div key={i} className="card p-6 animate-pulse">
                <div className="h-6 bg-gray-200 rounded w-3/4 mb-3" />
                <div className="h-4 bg-gray-200 rounded w-full mb-2" />
                <div className="h-4 bg-gray-200 rounded w-1/2" />
              </div>
            ))}
          </div>
        ) : !searched ? (
          <div className="text-center py-20">
            <div className="text-5xl mb-4">🔍</div>
            <p className="text-lg text-gray-500">输入关键词搜索帖子或用户</p>
            <p className="text-sm text-gray-400 mt-1">试试搜索「转型」「AI」「副业」</p>
            {hotTags.length > 0 && (
              <div className="mt-6">
                <p className="text-sm text-gray-400 mb-3">热门标签</p>
                <div className="flex flex-wrap justify-center gap-2">
                  {hotTags.map(tag => (
                    <button
                      key={tag}
                      onClick={() => { setQuery(tag); router.push(`/search?q=${encodeURIComponent(tag)}`); doSearch(tag, tab, 1) }}
                      className="px-3 py-1.5 bg-white border border-gray-200 rounded-full text-sm text-gray-600 hover:border-blue-300 hover:text-blue-600 hover:bg-blue-50 transition-all"
                    >
                      #{tag}
                    </button>
                  ))}
                </div>
              </div>
            )}
          </div>
        ) : total === 0 ? (
          <div className="text-center py-16">
            <div className="text-5xl mb-4">😮</div>
            <p className="text-lg text-gray-600 mb-1">没有找到相关{tab === 'posts' ? '帖子' : '用户'}</p>
            <p className="text-sm text-gray-400">换个关键词试试？</p>
          </div>
        ) : (
          <>
            <p className="text-sm text-gray-500 mb-4">共找到 <span className="font-medium text-gray-700">{total}</span> 个结果</p>

            {tab === 'posts' ? (
              <div className="space-y-3">
                {posts.map((post) => {
                  const typeInfo = TYPE_LABELS[post.type]
                  return (
                    <Link key={post.id} href={`/posts/${post.id}`} className="block card p-5 group">
                      <h3 className="text-base font-semibold text-gray-900 group-hover:text-blue-600 transition-colors mb-2 line-clamp-1">
                        {post.title}
                      </h3>
                      <p className="text-sm text-gray-500 mb-3 line-clamp-2">{stripHTML(post.content)}</p>
                      <div className="flex items-center gap-3 text-xs text-gray-400">
                        {typeInfo && (
                          <span className={`inline-flex items-center gap-1 px-2 py-0.5 rounded-full border ${typeInfo.color}`}>
                            <Icon name={typeInfo.icon} size={11} /> {typeInfo.label}
                          </span>
                        )}
                        {post.tags && post.tags.length > 0 && post.tags.slice(0, 3).map(tag => (
                          <span key={tag} className="px-2 py-0.5 bg-gray-100 text-gray-500 rounded-full">#{tag}</span>
                        ))}
                        <span>👍 {post.like_count}</span>
                        <span>💬 {post.comment_count}</span>
                      </div>
                    </Link>
                  )
                })}
              </div>
            ) : (
              <div className="space-y-3">
                {users.map((user) => (
                  <Link key={user.id} href={`/profile?id=${user.id}`} className="flex items-center gap-4 card p-5 group">
                    <Avatar nickname={user.nickname} size="lg" />
                    <div className="min-w-0">
                      <h3 className="font-semibold text-gray-900 group-hover:text-blue-600 transition-colors">{user.nickname}</h3>
                      {user.bio && <p className="text-sm text-gray-500 mt-0.5 truncate">{user.bio}</p>}
                      <p className="text-xs text-gray-400 mt-1">{user.follower_count} 粉丝</p>
                    </div>
                  </Link>
                ))}
              </div>
            )}

            {totalPages > 1 && (
              <div className="mt-8 flex items-center justify-center gap-2">
                <button
                  onClick={() => { const p = Math.max(1, page - 1); setPage(p); doSearch(query.trim(), tab, p) }}
                  disabled={page === 1}
                  className="btn-secondary text-sm px-4 py-2 disabled:opacity-40"
                >
                  ← 上一页
                </button>
                <span className="text-sm text-gray-500 px-3">{page} / {totalPages}</span>
                <button
                  onClick={() => { const p = Math.min(totalPages, page + 1); setPage(p); doSearch(query.trim(), tab, p) }}
                  disabled={page === totalPages}
                  className="btn-secondary text-sm px-4 py-2 disabled:opacity-40"
                >
                  下一页 →
                </button>
              </div>
            )}
          </>
        )}
      </div>
    </div>
  )
}
