'use client'

import { useState, useEffect } from 'react'
import { useRouter } from 'next/navigation'
import { get, put } from '@/lib/api'
import { toastSuccess, toastError } from '@/lib/toast'
import Link from 'next/link'

interface Comment {
  id: number
  post_id: number
  post_title: string
  user_id: number
  user_name: string
  content: string
  status: number
  created_at: string
}

const STATUS_MAP: Record<number, { label: string; color: string }> = {
  1: { label: '正常', color: 'bg-green-50 text-green-600' },
  0: { label: '已屏蔽', color: 'bg-red-50 text-red-600' },
}

export default function CommentsPage() {
  const router = useRouter()
  const [comments, setComments] = useState<Comment[]>([])
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [postIdFilter, setPostIdFilter] = useState('')
  const [loading, setLoading] = useState(false)

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

  const fetchComments = async (p: number, postId?: string) => {
    setLoading(true)
    try {
      let url = `/admin/comments?page=${p}&page_size=20`
      if (postId) url += `&post_id=${postId}`
      const res = await get<{ items: Comment[]; total: number }>(url)
      if (res.code === 0) {
        setComments(res.data.items || [])
        setTotal(res.data.total)
      }
    } catch { toastError('获取评论列表失败') }
    finally { setLoading(false) }
  }

  useEffect(() => { fetchComments(page, postIdFilter) }, [page, postIdFilter])

  const handleToggleStatus = async (comment: Comment) => {
    const isHidden = comment.status === 0
    const action = isHidden ? '恢复' : '屏蔽'
    if (!confirm(`确定${action}这条评论吗？`)) return
    try {
      const endpoint = isHidden ? `/admin/comments/${comment.id}/restore` : `/admin/comments/${comment.id}/hide`
      const res = await put(endpoint)
      if (res.code === 0) {
        toastSuccess(`已${action}`)
        setComments(comments.map(c => c.id === comment.id ? { ...c, status: isHidden ? 1 : 0 } : c))
      }
    } catch { toastError('操作失败') }
  }

  const totalPages = Math.ceil(total / 20)

  return (
    <div className="container mx-auto px-4 py-8 max-w-6xl">
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">评论管理</h1>
          <p className="text-sm text-gray-500 mt-1">管理用户评论内容</p>
        </div>
        <Link href="/admin" className="btn-secondary text-sm">← 返回管理后台</Link>
      </div>

      {/* 筛选 */}
      <div className="card p-4 mb-6">
        <div className="flex items-center gap-4">
          <span className="text-sm font-medium text-gray-700">按帖子筛选：</span>
          <input
            type="text"
            value={postIdFilter}
            onChange={e => { setPostIdFilter(e.target.value); setPage(1) }}
            placeholder="输入帖子ID"
            className="input-field text-sm w-40"
          />
          {postIdFilter && (
            <button
              onClick={() => { setPostIdFilter(''); setPage(1) }}
              className="text-sm text-gray-500 hover:text-gray-700"
            >
              清除筛选
            </button>
          )}
          <span className="text-sm text-gray-500 ml-auto">共 {total} 条评论</span>
        </div>
      </div>

      {/* 评论列表 */}
      {loading ? (
        <div className="card p-12 text-center animate-pulse">
          <div className="h-8 bg-gray-200 rounded w-40 mx-auto mb-4" />
          <div className="h-4 bg-gray-200 rounded w-60 mx-auto" />
        </div>
      ) : comments.length === 0 ? (
        <div className="card p-12 text-center">
          <div className="text-4xl mb-2">💬</div>
          <p className="text-gray-500">暂无评论记录</p>
        </div>
      ) : (
        <div className="card overflow-hidden">
          <table className="w-full">
            <thead className="bg-gray-50 border-b border-gray-100">
              <tr>
                {['ID', '评论内容', '评论人', '所属帖子', '状态', '时间', '操作'].map(h => (
                  <th key={h} className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">{h}</th>
                ))}
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-100">
              {comments.map((comment) => {
                const status = STATUS_MAP[comment.status] || STATUS_MAP[1]
                return (
                  <tr key={comment.id} className="hover:bg-gray-50 transition-colors">
                    <td className="px-4 py-3 text-sm text-gray-500">{comment.id}</td>
                    <td className="px-4 py-3 max-w-xs">
                      <p className="text-sm text-gray-600 truncate">{comment.content}</p>
                    </td>
                    <td className="px-4 py-3 text-sm text-gray-900">{comment.user_name || `用户${comment.user_id}`}</td>
                    <td className="px-4 py-3">
                      <Link
                        href={`/posts/${comment.post_id}`}
                        className="text-sm text-blue-600 hover:text-blue-800 hover:underline"
                        target="_blank"
                      >
                        {comment.post_title || `#${comment.post_id}`}
                      </Link>
                    </td>
                    <td className="px-4 py-3">
                      <span className={`px-2 py-0.5 rounded-full text-xs ${status.color}`}>
                        {status.label}
                      </span>
                    </td>
                    <td className="px-4 py-3 text-sm text-gray-400">{new Date(comment.created_at).toLocaleDateString()}</td>
                    <td className="px-4 py-3">
                      <button
                        onClick={() => handleToggleStatus(comment)}
                        className={`text-sm font-medium ${
                          comment.status === 1
                            ? 'text-red-600 hover:text-red-800'
                            : 'text-green-600 hover:text-green-800'
                        }`}
                      >
                        {comment.status === 1 ? '屏蔽' : '恢复'}
                      </button>
                    </td>
                  </tr>
                )
              })}
            </tbody>
          </table>
        </div>
      )}

      {/* 分页 */}
      {totalPages > 1 && (
        <div className="mt-4 flex items-center justify-center gap-2">
          <button onClick={() => setPage(p => Math.max(1, p - 1))} disabled={page === 1} className="btn-secondary text-sm px-3 py-1.5 disabled:opacity-40">←</button>
          <span className="text-sm text-gray-500">{page} / {totalPages}</span>
          <button onClick={() => setPage(p => Math.min(totalPages, p + 1))} disabled={page === totalPages} className="btn-secondary text-sm px-3 py-1.5 disabled:opacity-40">→</button>
        </div>
      )}
    </div>
  )
}
