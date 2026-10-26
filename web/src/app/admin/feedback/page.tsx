'use client'

import { useState, useEffect } from 'react'
import { useRouter } from 'next/navigation'
import { get, put } from '@/lib/api'
import { toastSuccess, toastError } from '@/lib/toast'
import Link from 'next/link'

interface Feedback {
  id: number
  user_id: number
  user_nickname: string
  user_email: string
  type: string
  title: string
  content: string
  contact: string
  status: number
  reply: string
  created_at: string
}

const STATUS_MAP: Record<number, { label: string; color: string }> = {
  0: { label: '待处理', color: 'bg-yellow-50 text-yellow-600' },
  1: { label: '已处理', color: 'bg-green-50 text-green-600' },
  2: { label: '已忽略', color: 'bg-gray-100 text-gray-600' },
}

const TYPE_MAP: Record<string, { label: string; color: string }> = {
  bug: { label: 'Bug', color: 'bg-red-50 text-red-600' },
  suggestion: { label: '建议', color: 'bg-blue-50 text-blue-600' },
  other: { label: '其他', color: 'bg-gray-50 text-gray-600' },
}

export default function FeedbackAdminPage() {
  const router = useRouter()
  const [feedbacks, setFeedbacks] = useState<Feedback[]>([])
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [statusFilter, setStatusFilter] = useState<number | null>(null)
  const [typeFilter, setTypeFilter] = useState<string>('')
  const [loading, setLoading] = useState(false)
  const [selectedFeedback, setSelectedFeedback] = useState<Feedback | null>(null)
  const [replyText, setReplyText] = useState('')

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

  const fetchFeedbacks = async (p: number, status?: number | null, type?: string) => {
    setLoading(true)
    try {
      let url = `/admin/feedback?page=${p}&page_size=20`
      if (status !== null && status !== undefined) url += `&status=${status}`
      if (type) url += `&type=${type}`
      const res = await get<{ items: Feedback[]; total: number }>(url)
      if (res.code === 0) {
        setFeedbacks(res.data.items || [])
        setTotal(res.data.total)
      }
    } catch { toastError('获取反馈列表失败') }
    finally { setLoading(false) }
  }

  useEffect(() => { fetchFeedbacks(page, statusFilter, typeFilter) }, [page, statusFilter, typeFilter])

  const handleUpdateStatus = async (id: number, status: number) => {
    const action = status === 1 ? '标记已处理' : '忽略'
    if (!confirm(`确定${action}这条反馈吗？`)) return
    try {
      const res = await put(`/admin/feedback/${id}/status`, { status, reply: '' })
      if (res.code === 0) {
        toastSuccess(`已${action}`)
        setFeedbacks(feedbacks.map(f => f.id === id ? { ...f, status } : f))
        setSelectedFeedback(null)
      }
    } catch { toastError('操作失败') }
  }

  const handleReply = async (id: number) => {
    if (!replyText.trim()) {
      toastError('请输入回复内容')
      return
    }
    try {
      const res = await put(`/admin/feedback/${id}/status`, { status: 1, reply: replyText })
      if (res.code === 0) {
        toastSuccess('已回复')
        setFeedbacks(feedbacks.map(f => f.id === id ? { ...f, status: 1, reply: replyText } : f))
        setSelectedFeedback(null)
        setReplyText('')
      }
    } catch { toastError('回复失败') }
  }

  const totalPages = Math.ceil(total / 20)

  return (
    <div className="container mx-auto px-4 py-8 max-w-6xl">
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">反馈管理</h1>
          <p className="text-sm text-gray-500 mt-1">查看和处理用户反馈</p>
        </div>
        <Link href="/admin" className="btn-secondary text-sm">← 返回管理后台</Link>
      </div>

      {/* 筛选 */}
      <div className="card p-4 mb-6">
        <div className="flex items-center gap-4 flex-wrap">
          <span className="text-sm font-medium text-gray-700">状态筛选：</span>
          <div className="flex gap-2">
            <button
              onClick={() => { setStatusFilter(null); setPage(1) }}
              className={`px-3 py-1.5 rounded-lg text-sm font-medium transition-all ${
                statusFilter === null ? 'bg-blue-600 text-white' : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
              }`}
            >
              全部
            </button>
            {Object.entries(STATUS_MAP).map(([key, value]) => (
              <button
                key={key}
                onClick={() => { setStatusFilter(Number(key)); setPage(1) }}
                className={`px-3 py-1.5 rounded-lg text-sm font-medium transition-all ${
                  statusFilter === Number(key) ? 'bg-blue-600 text-white' : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
                }`}
              >
                {value.label}
              </button>
            ))}
          </div>

          <span className="text-sm font-medium text-gray-700 ml-4">类型筛选：</span>
          <div className="flex gap-2">
            <button
              onClick={() => { setTypeFilter(''); setPage(1) }}
              className={`px-3 py-1.5 rounded-lg text-sm font-medium transition-all ${
                typeFilter === '' ? 'bg-blue-600 text-white' : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
              }`}
            >
              全部
            </button>
            {Object.entries(TYPE_MAP).map(([key, value]) => (
              <button
                key={key}
                onClick={() => { setTypeFilter(key); setPage(1) }}
                className={`px-3 py-1.5 rounded-lg text-sm font-medium transition-all ${
                  typeFilter === key ? 'bg-blue-600 text-white' : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
                }`}
              >
                {value.label}
              </button>
            ))}
          </div>
          <span className="text-sm text-gray-500 ml-auto">共 {total} 条反馈</span>
        </div>
      </div>

      {/* 反馈列表 */}
      {loading ? (
        <div className="card p-12 text-center animate-pulse">
          <div className="h-8 bg-gray-200 rounded w-40 mx-auto mb-4" />
          <div className="h-4 bg-gray-200 rounded w-60 mx-auto" />
        </div>
      ) : feedbacks.length === 0 ? (
        <div className="card p-12 text-center">
          <div className="text-4xl mb-2">💬</div>
          <p className="text-gray-500">暂无反馈记录</p>
        </div>
      ) : (
        <div className="card overflow-hidden">
          <table className="w-full">
            <thead className="bg-gray-50 border-b border-gray-100">
              <tr>
                {['ID', '用户', '类型', '标题', '状态', '时间', '操作'].map(h => (
                  <th key={h} className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">{h}</th>
                ))}
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-100">
              {feedbacks.map((feedback) => {
                const status = STATUS_MAP[feedback.status]
                const type = TYPE_MAP[feedback.type] || TYPE_MAP.other
                return (
                  <tr key={feedback.id} className="hover:bg-gray-50 transition-colors cursor-pointer" onClick={() => setSelectedFeedback(feedback)}>
                    <td className="px-4 py-3 text-sm text-gray-500">{feedback.id}</td>
                    <td className="px-4 py-3">
                      <div className="text-sm">
                        <div className="text-gray-900">{feedback.user_nickname || `用户${feedback.user_id}`}</div>
                        <div className="text-gray-400 text-xs">{feedback.user_email}</div>
                      </div>
                    </td>
                    <td className="px-4 py-3">
                      <span className={`px-2 py-0.5 rounded-full text-xs ${type.color}`}>
                        {type.label}
                      </span>
                    </td>
                    <td className="px-4 py-3 text-sm text-gray-900 max-w-xs truncate">{feedback.title}</td>
                    <td className="px-4 py-3">
                      <span className={`px-2 py-0.5 rounded-full text-xs ${status.color}`}>
                        {status.label}
                      </span>
                    </td>
                    <td className="px-4 py-3 text-sm text-gray-400">{new Date(feedback.created_at).toLocaleDateString()}</td>
                    <td className="px-4 py-3">
                      {feedback.status === 0 && (
                        <div className="flex items-center gap-2" onClick={e => e.stopPropagation()}>
                          <button
                            onClick={() => handleUpdateStatus(feedback.id, 1)}
                            className="text-sm text-green-600 hover:text-green-800 font-medium"
                          >
                            已处理
                          </button>
                          <button
                            onClick={() => handleUpdateStatus(feedback.id, 2)}
                            className="text-sm text-gray-500 hover:text-gray-700 font-medium"
                          >
                            忽略
                          </button>
                        </div>
                      )}
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

      {/* 反馈详情弹窗 */}
      {selectedFeedback && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4" onClick={() => setSelectedFeedback(null)}>
          <div className="bg-white rounded-xl shadow-xl max-w-2xl w-full max-h-[90vh] overflow-y-auto" onClick={e => e.stopPropagation()}>
            <div className="p-6">
              <div className="flex items-center justify-between mb-4">
                <h2 className="text-xl font-bold text-gray-900">反馈详情</h2>
                <button onClick={() => setSelectedFeedback(null)} className="text-gray-400 hover:text-gray-600">✕</button>
              </div>

              <div className="space-y-4">
                <div>
                  <label className="text-sm text-gray-500">用户</label>
                  <p className="text-gray-900">{selectedFeedback.user_nickname || `用户${selectedFeedback.user_id}`} ({selectedFeedback.user_email})</p>
                </div>
                <div>
                  <label className="text-sm text-gray-500">类型</label>
                  <p><span className={`px-2 py-0.5 rounded-full text-xs ${TYPE_MAP[selectedFeedback.type]?.color || TYPE_MAP.other.color}`}>{TYPE_MAP[selectedFeedback.type]?.label || '其他'}</span></p>
                </div>
                <div>
                  <label className="text-sm text-gray-500">标题</label>
                  <p className="text-gray-900">{selectedFeedback.title}</p>
                </div>
                <div>
                  <label className="text-sm text-gray-500">内容</label>
                  <p className="text-gray-900 whitespace-pre-wrap">{selectedFeedback.content}</p>
                </div>
                {selectedFeedback.contact && (
                  <div>
                    <label className="text-sm text-gray-500">联系方式</label>
                    <p className="text-gray-900">{selectedFeedback.contact}</p>
                  </div>
                )}
                <div>
                  <label className="text-sm text-gray-500">状态</label>
                  <p><span className={`px-2 py-0.5 rounded-full text-xs ${STATUS_MAP[selectedFeedback.status]?.color}`}>{STATUS_MAP[selectedFeedback.status]?.label}</span></p>
                </div>
                {selectedFeedback.reply && (
                  <div>
                    <label className="text-sm text-gray-500">管理员回复</label>
                    <p className="text-gray-900 bg-gray-50 p-3 rounded-lg">{selectedFeedback.reply}</p>
                  </div>
                )}
                <div>
                  <label className="text-sm text-gray-500">提交时间</label>
                  <p className="text-gray-900">{new Date(selectedFeedback.created_at).toLocaleString()}</p>
                </div>
              </div>

              {/* 回复区域 */}
              {selectedFeedback.status === 0 && (
                <div className="mt-6 pt-4 border-t">
                  <label className="text-sm font-medium text-gray-700 mb-2 block">回复用户</label>
                  <textarea
                    value={replyText}
                    onChange={e => setReplyText(e.target.value)}
                    placeholder="输入回复内容..."
                    className="w-full p-3 border rounded-lg resize-none h-24 focus:outline-none focus:ring-2 focus:ring-blue-500"
                  />
                  <div className="flex gap-2 mt-3">
                    <button
                      onClick={() => handleReply(selectedFeedback.id)}
                      className="btn-primary text-sm px-4 py-2"
                    >
                      提交回复
                    </button>
                    <button
                      onClick={() => handleUpdateStatus(selectedFeedback.id, 2)}
                      className="btn-secondary text-sm px-4 py-2"
                    >
                      忽略
                    </button>
                  </div>
                </div>
              )}
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
