'use client'

import { useState, useEffect } from 'react'
import { useRouter } from 'next/navigation'
import { get, put } from '@/lib/api'
import { toastSuccess, toastError } from '@/lib/toast'
import Link from 'next/link'

interface Report {
  id: number
  reporter_id: number
  reporter_name: string
  target_type: string
  target_id: number
  target_title: string
  reason: string
  status: number
  created_at: string
  handled_at?: string
}

const STATUS_MAP: Record<number, { label: string; color: string }> = {
  1: { label: '待处理', color: 'bg-yellow-50 text-yellow-600' },
  2: { label: '已处理', color: 'bg-green-50 text-green-600' },
  3: { label: '已忽略', color: 'bg-gray-100 text-gray-600' },
}

const TARGET_TYPE_MAP: Record<string, string> = {
  post: '帖子',
  comment: '评论',
  user: '用户',
}

export default function ReportsPage() {
  const router = useRouter()
  const [reports, setReports] = useState<Report[]>([])
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [statusFilter, setStatusFilter] = useState<number | null>(null)
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

  const fetchReports = async (p: number, status?: number | null) => {
    setLoading(true)
    try {
      let url = `/admin/reports?page=${p}&size=20`
      if (status) url += `&status=${status}`
      const res = await get<{ items: Report[]; total: number }>(url)
      if (res.code === 0) {
        setReports(res.data.items || [])
        setTotal(res.data.total)
      }
    } catch { toastError('获取举报列表失败') }
    finally { setLoading(false) }
  }

  useEffect(() => { fetchReports(page, statusFilter) }, [page, statusFilter])

  const handleReport = async (id: number, status: number) => {
    const action = status === 2 ? '标记已处理' : '忽略'
    if (!confirm(`确定${action}这条举报吗？`)) return
    try {
      const res = await put(`/admin/reports/${id}`, { status })
      if (res.code === 0) {
        toastSuccess(`已${action}`)
        setReports(reports.map(r => r.id === id ? { ...r, status } : r))
      }
    } catch { toastError('操作失败') }
  }

  const totalPages = Math.ceil(total / 20)

  return (
    <div className="container mx-auto px-4 py-8 max-w-6xl">
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">举报管理</h1>
          <p className="text-sm text-gray-500 mt-1">处理用户举报的内容</p>
        </div>
        <Link href="/admin" className="btn-secondary text-sm">← 返回管理后台</Link>
      </div>

      {/* 筛选 */}
      <div className="card p-4 mb-6">
        <div className="flex items-center gap-4">
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
          <span className="text-sm text-gray-500 ml-auto">共 {total} 条举报</span>
        </div>
      </div>

      {/* 举报列表 */}
      {loading ? (
        <div className="card p-12 text-center animate-pulse">
          <div className="h-8 bg-gray-200 rounded w-40 mx-auto mb-4" />
          <div className="h-4 bg-gray-200 rounded w-60 mx-auto" />
        </div>
      ) : reports.length === 0 ? (
        <div className="card p-12 text-center">
          <div className="text-4xl mb-2">🛡️</div>
          <p className="text-gray-500">暂无举报记录</p>
        </div>
      ) : (
        <div className="card overflow-hidden">
          <table className="w-full">
            <thead className="bg-gray-50 border-b border-gray-100">
              <tr>
                {['ID', '举报人', '举报对象', '原因', '状态', '时间', '操作'].map(h => (
                  <th key={h} className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">{h}</th>
                ))}
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-100">
              {reports.map((report) => {
                const status = STATUS_MAP[report.status]
                const targetType = TARGET_TYPE_MAP[report.target_type] || report.target_type
                return (
                  <tr key={report.id} className="hover:bg-gray-50 transition-colors">
                    <td className="px-4 py-3 text-sm text-gray-500">{report.id}</td>
                    <td className="px-4 py-3 text-sm text-gray-900">{report.reporter_name || `用户${report.reporter_id}`}</td>
                    <td className="px-4 py-3">
                      <div className="text-sm">
                        <span className="text-gray-500">{targetType}</span>
                        <span className="mx-1">·</span>
                        <Link
                          href={report.target_type === 'post' ? `/posts/${report.target_id}` : '#'}
                          className="text-blue-600 hover:text-blue-800 hover:underline"
                        >
                          {report.target_title || `#${report.target_id}`}
                        </Link>
                      </div>
                    </td>
                    <td className="px-4 py-3 text-sm text-gray-600 max-w-xs truncate">{report.reason}</td>
                    <td className="px-4 py-3">
                      <span className={`px-2 py-0.5 rounded-full text-xs ${status.color}`}>
                        {status.label}
                      </span>
                    </td>
                    <td className="px-4 py-3 text-sm text-gray-400">{new Date(report.created_at).toLocaleDateString()}</td>
                    <td className="px-4 py-3">
                      {report.status === 1 && (
                        <div className="flex items-center gap-2">
                          <button
                            onClick={() => handleReport(report.id, 2)}
                            className="text-sm text-green-600 hover:text-green-800 font-medium"
                          >
                            已处理
                          </button>
                          <button
                            onClick={() => handleReport(report.id, 3)}
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
    </div>
  )
}
