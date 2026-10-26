'use client'

import { useState, useEffect } from 'react'
import { useRouter } from 'next/navigation'
import { get, post } from '@/lib/api'
import { TYPE_LABELS } from '@/lib/constants'
import { toastSuccess, toastError } from '@/lib/toast'
import Link from 'next/link'

interface PendingPost {
  id: number
  user_id: number
  user_name: string
  type: string
  title: string
  content: string
  created_at: string
}

export default function ReviewPage() {
  const router = useRouter()
  const [posts, setPosts] = useState<PendingPost[]>([])
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [loading, setLoading] = useState(false)
  const [expandedId, setExpandedId] = useState<number | null>(null)
  const [rejectReason, setRejectReason] = useState('')
  const [rejectingId, setRejectingId] = useState<number | null>(null)

  useEffect(() => {
    const token = localStorage.getItem('token')
    if (!token) { router.push('/login'); return }
    get<{ is_admin: boolean }>('/auth/profile').then(res => {
      if (res.code !== 0 || !res.data?.is_admin) {
        toastError('需要管理员权限')
        router.push('/')
      }
    }).catch(() => { router.push('/login') })
  }, [router])

  const fetchPendingPosts = async (p: number) => {
    setLoading(true)
    try {
      const res = await get<{ items: PendingPost[]; total: number }>(`/admin/posts/pending?page=${p}&page_size=10`)
      if (res.code === 0) {
        setPosts(res.data.items || [])
        setTotal(res.data.total)
      }
    } catch { toastError('获取待审核帖子失败') }
    finally { setLoading(false) }
  }

  useEffect(() => { fetchPendingPosts(page) }, [page])

  const handleApprove = async (id: number) => {
    if (!confirm('确定审核通过这条帖子吗？')) return
    try {
      const res = await post(`/admin/posts/${id}/approve`)
      if (res.code === 0) {
        toastSuccess('已通过审核')
        setPosts(posts.filter(p => p.id !== id))
        setTotal(t => t - 1)
      }
    } catch { toastError('操作失败') }
  }

  const handleReject = async (id: number) => {
    if (!rejectReason.trim()) {
      toastError('请填写拒绝原因')
      return
    }
    try {
      const res = await post(`/admin/posts/${id}/reject`, { reason: rejectReason })
      if (res.code === 0) {
        toastSuccess('已拒绝')
        setPosts(posts.filter(p => p.id !== id))
        setTotal(t => t - 1)
        setRejectingId(null)
        setRejectReason('')
      }
    } catch { toastError('操作失败') }
  }

  const totalPages = Math.ceil(total / 10)

  return (
    <div className="container mx-auto px-4 py-8 max-w-4xl">
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">帖子审核</h1>
          <p className="text-sm text-gray-500 mt-1">审核用户发布的帖子内容</p>
        </div>
        <Link href="/admin" className="btn-secondary text-sm">← 返回管理后台</Link>
      </div>

      <div className="mb-4">
        <p className="text-sm text-gray-500">
          待审核帖子：<span className="font-medium text-gray-700">{total}</span> 条
        </p>
      </div>

      {/* 帖子列表 */}
      {loading ? (
        <div className="card p-12 text-center animate-pulse">
          <div className="h-8 bg-gray-200 rounded w-40 mx-auto mb-4" />
          <div className="h-4 bg-gray-200 rounded w-60 mx-auto" />
        </div>
      ) : posts.length === 0 ? (
        <div className="card p-12 text-center">
          <div className="text-4xl mb-2">✅</div>
          <p className="text-gray-500">暂无待审核帖子</p>
        </div>
      ) : (
        <div className="space-y-4">
          {posts.map((post) => {
            const t = TYPE_LABELS[post.type]
            const isExpanded = expandedId === post.id
            const isRejecting = rejectingId === post.id

            return (
              <div key={post.id} className="card overflow-hidden">
                <div className="p-4">
                  {/* 头部信息 */}
                  <div className="flex items-start justify-between mb-3">
                    <div>
                      <div className="flex items-center gap-2 mb-1">
                        {t && <span className={`px-2 py-0.5 rounded-full text-xs ${t.color}`}>{t.label}</span>}
                        <span className="text-sm text-gray-500">#{post.id}</span>
                      </div>
                      <h3 className="text-lg font-semibold text-gray-900">{post.title}</h3>
                    </div>
                    <div className="text-right text-sm text-gray-500">
                      <div>{post.user_name || `用户${post.user_id}`}</div>
                      <div>{new Date(post.created_at).toLocaleString()}</div>
                    </div>
                  </div>

                  {/* 内容预览 */}
                  <div
                    className={`text-sm text-gray-600 mb-4 cursor-pointer ${isExpanded ? '' : 'line-clamp-3'}`}
                    onClick={() => setExpandedId(isExpanded ? null : post.id)}
                  >
                    {post.content}
                  </div>
                  {!isExpanded && post.content.length > 150 && (
                    <button
                      onClick={() => setExpandedId(post.id)}
                      className="text-sm text-blue-600 hover:text-blue-800 mb-4"
                    >
                      展开全文 ↓
                    </button>
                  )}

                  {/* 操作按钮 */}
                  <div className="flex items-center gap-3 pt-3 border-t border-gray-100">
                    <button
                      onClick={() => handleApprove(post.id)}
                      className="btn-primary text-sm px-4 py-2"
                    >
                      ✓ 通过
                    </button>
                    {!isRejecting ? (
                      <button
                        onClick={() => { setRejectingId(post.id); setRejectReason('') }}
                        className="btn-secondary text-sm px-4 py-2 text-red-600 hover:text-red-800"
                      >
                        ✗ 拒绝
                      </button>
                    ) : (
                      <div className="flex-1 flex items-center gap-2">
                        <input
                          type="text"
                          value={rejectReason}
                          onChange={e => setRejectReason(e.target.value)}
                          placeholder="请输入拒绝原因..."
                          className="flex-1 input-field text-sm"
                          onKeyDown={e => e.key === 'Enter' && handleReject(post.id)}
                        />
                        <button
                          onClick={() => handleReject(post.id)}
                          className="btn-primary text-sm px-4 py-2 bg-red-600 hover:bg-red-700"
                        >
                          确认拒绝
                        </button>
                        <button
                          onClick={() => { setRejectingId(null); setRejectReason('') }}
                          className="btn-secondary text-sm px-3 py-2"
                        >
                          取消
                        </button>
                      </div>
                    )}
                    <Link
                      href={`/posts/${post.id}`}
                      className="text-sm text-gray-500 hover:text-gray-700 ml-auto"
                      target="_blank"
                    >
                      查看原帖 →
                    </Link>
                  </div>
                </div>
              </div>
            )
          })}
        </div>
      )}

      {/* 分页 */}
      {totalPages > 1 && (
        <div className="mt-6 flex items-center justify-center gap-2">
          <button onClick={() => setPage(p => Math.max(1, p - 1))} disabled={page === 1} className="btn-secondary text-sm px-3 py-1.5 disabled:opacity-40">←</button>
          <span className="text-sm text-gray-500">{page} / {totalPages}</span>
          <button onClick={() => setPage(p => Math.min(totalPages, p + 1))} disabled={page === totalPages} className="btn-secondary text-sm px-3 py-1.5 disabled:opacity-40">→</button>
        </div>
      )}
    </div>
  )
}
