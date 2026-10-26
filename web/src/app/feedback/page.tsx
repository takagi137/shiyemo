'use client'

import { useState, useEffect } from 'react'
import { useRouter } from 'next/navigation'
import { get, post } from '@/lib/api'
import { toastSuccess, toastError } from '@/lib/toast'
import Link from 'next/link'

interface FeedbackItem {
  id: number
  type: string
  title: string
  content: string
  contact: string
  status: number
  reply: string
  created_at: string
}

const TYPE_OPTIONS = [
  { value: 'bug', label: 'Bug 反馈', icon: '🐛' },
  { value: 'suggestion', label: '功能建议', icon: '💡' },
  { value: 'other', label: '其他反馈', icon: '💬' },
]

const STATUS_MAP: Record<number, { label: string; color: string }> = {
  0: { label: '待处理', color: 'text-yellow-600' },
  1: { label: '已处理', color: 'text-green-600' },
  2: { label: '已忽略', color: 'text-gray-500' },
}

export default function FeedbackPage() {
  const router = useRouter()
  const [feedbacks, setFeedbacks] = useState<FeedbackItem[]>([])
  const [loading, setLoading] = useState(false)
  const [showForm, setShowForm] = useState(true)
  const [formData, setFormData] = useState({
    type: 'suggestion',
    title: '',
    content: '',
    contact: '',
  })
  const [submitting, setSubmitting] = useState(false)

  useEffect(() => {
    const token = localStorage.getItem('token')
    if (!token) { router.push('/login'); return }
    fetchFeedbacks()
  }, [router])

  const fetchFeedbacks = async () => {
    setLoading(true)
    try {
      const res = await get<{ items: FeedbackItem[] }>('/feedback/mine')
      if (res.code === 0) {
        setFeedbacks(res.data.items || [])
      }
    } catch { toastError('获取反馈列表失败') }
    finally { setLoading(false) }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!formData.title.trim()) { toastError('请输入标题'); return }
    if (!formData.content.trim()) { toastError('请输入内容'); return }

    setSubmitting(true)
    try {
      const res = await post('/feedback', formData)
      if (res.code === 0) {
        toastSuccess('反馈提交成功')
        setShowForm(false)
        setFormData({ type: 'suggestion', title: '', content: '', contact: '' })
        fetchFeedbacks()
      }
    } catch { toastError('提交失败') }
    finally { setSubmitting(false) }
  }

  return (
    <div className="container mx-auto px-4 py-8 max-w-4xl">
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">意见反馈</h1>
          <p className="text-sm text-gray-500 mt-1">提交 Bug、建议或其他反馈</p>
        </div>
        <button
          onClick={() => setShowForm(!showForm)}
          className="inline-flex items-center justify-center rounded-lg bg-blue-600 px-5 py-2 text-sm font-medium text-white shadow-sm transition-all hover:bg-blue-700 active:scale-95"
        >
          {showForm ? '取消' : '+ 提交反馈'}
        </button>
      </div>

      {/* 反馈表单 */}
      {showForm && (
        <div className="card p-6 mb-6">
          <h2 className="text-lg font-semibold mb-4">提交新反馈</h2>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">反馈类型</label>
              <div className="flex gap-3">
                {TYPE_OPTIONS.map(opt => (
                  <button
                    key={opt.value}
                    type="button"
                    onClick={() => setFormData({ ...formData, type: opt.value })}
                    className={`flex items-center gap-2 px-4 py-2 rounded-lg border-2 transition-all ${
                      formData.type === opt.value
                        ? 'border-blue-500 bg-blue-50 text-blue-700'
                        : 'border-gray-200 hover:border-gray-300'
                    }`}
                  >
                    <span>{opt.icon}</span>
                    <span>{opt.label}</span>
                  </button>
                ))}
              </div>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">标题</label>
              <input
                type="text"
                value={formData.title}
                onChange={e => setFormData({ ...formData, title: e.target.value })}
                placeholder="简要描述你的反馈"
                className="input-field w-full"
                maxLength={200}
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">详细内容</label>
              <textarea
                value={formData.content}
                onChange={e => setFormData({ ...formData, content: e.target.value })}
                placeholder="请详细描述问题或建议..."
                className="w-full p-3 border rounded-lg resize-none h-32 focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">联系方式（可选）</label>
              <input
                type="text"
                value={formData.contact}
                onChange={e => setFormData({ ...formData, contact: e.target.value })}
                placeholder="邮箱或其他联系方式，方便我们回复你"
                className="input-field w-full"
              />
            </div>

            <div className="flex items-center gap-3 pt-2">
              <button
                type="submit"
                disabled={submitting}
                className="inline-flex items-center justify-center rounded-lg bg-blue-600 px-6 py-2.5 text-sm font-medium text-white shadow-sm transition-all hover:bg-blue-700 active:scale-95 disabled:cursor-not-allowed disabled:opacity-50"
              >
                {submitting ? (
                  <>
                    <svg className="-ml-0.5 mr-1.5 h-4 w-4 animate-spin text-white" viewBox="0 0 24 24" fill="none">
                      <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"/>
                      <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"/>
                    </svg>
                    提交中...
                  </>
                ) : '提交反馈'}
              </button>
              <button
                type="button"
                onClick={() => setShowForm(false)}
                className="inline-flex items-center justify-center rounded-lg px-6 py-2.5 text-sm font-medium text-gray-600 transition-all hover:bg-gray-100 active:scale-95"
              >
                取消
              </button>
            </div>
          </form>
        </div>
      )}

      {/* 我的反馈列表 */}
      <div className="mb-4">
        <h2 className="text-lg font-semibold text-gray-900">我的反馈</h2>
      </div>

      {loading ? (
        <div className="card p-12 text-center animate-pulse">
          <div className="h-8 bg-gray-200 rounded w-40 mx-auto mb-4" />
          <div className="h-4 bg-gray-200 rounded w-60 mx-auto" />
        </div>
      ) : feedbacks.length === 0 ? (
        <div className="card p-12 text-center">
          <div className="text-4xl mb-2">📝</div>
          <p className="text-gray-500 mb-4">暂无反馈记录</p>
          <button
            onClick={() => setShowForm(true)}
            className="inline-flex items-center justify-center rounded-lg bg-blue-600 px-5 py-2 text-sm font-medium text-white shadow-sm transition-all hover:bg-blue-700 active:scale-95"
          >
            提交第一条反馈
          </button>
        </div>
      ) : (
        <div className="space-y-4">
          {feedbacks.map(feedback => (
            <div key={feedback.id} className="card p-4">
              <div className="flex items-start justify-between mb-2">
                <div>
                  <span className="text-xs px-2 py-0.5 rounded-full bg-gray-100 text-gray-600 mr-2">
                    {TYPE_OPTIONS.find(t => t.value === feedback.type)?.label || '其他'}
                  </span>
                  <span className={`text-xs ${STATUS_MAP[feedback.status]?.color}`}>
                    {STATUS_MAP[feedback.status]?.label}
                  </span>
                </div>
                <span className="text-xs text-gray-400">
                  {new Date(feedback.created_at).toLocaleDateString()}
                </span>
              </div>
              <h3 className="font-medium text-gray-900 mb-2">{feedback.title}</h3>
              <p className="text-sm text-gray-600 mb-3">{feedback.content}</p>
              {feedback.reply && (
                <div className="bg-green-50 border-l-4 border-green-400 p-3 rounded">
                  <p className="text-sm font-medium text-green-800 mb-1">管理员回复</p>
                  <p className="text-sm text-green-700">{feedback.reply}</p>
                </div>
              )}
            </div>
          ))}
        </div>
      )}

      <div className="mt-6 text-center">
        <Link href="/" className="text-sm text-gray-500 hover:text-gray-700">← 返回首页</Link>
      </div>
    </div>
  )
}
