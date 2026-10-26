'use client'

import { useState, useEffect } from 'react'
import { useRouter } from 'next/navigation'
import { get, del, put, downloadFile } from '@/lib/api'
import { TYPE_LABELS } from '@/lib/constants'
import { toastSuccess, toastError } from '@/lib/toast'
import Link from 'next/link'

interface Post {
  id: number
  user_id: number
  type: string
  title: string
  content: string
  like_count: number
  comment_count: number
  created_at: string
  is_top?: boolean
  is_featured?: boolean
}

interface User {
  id: number
  nickname: string
  email: string
  phone: string
  is_admin: boolean
  status: number
  created_at: string
}

interface Stats {
  total_users: number
  total_posts: number
  total_comments: number
  active_users_today: number
}

export default function AdminPage() {
  const router = useRouter()
  const [tab, setTab] = useState<'stats' | 'posts' | 'users'>('stats')
  const [stats, setStats] = useState<Stats | null>(null)
  const [posts, setPosts] = useState<Post[]>([])
  const [users, setUsers] = useState<User[]>([])
  const [postsTotal, setPostsTotal] = useState(0)
  const [usersTotal, setUsersTotal] = useState(0)
  const [postsPage, setPostsPage] = useState(1)
  const [usersPage, setUsersPage] = useState(1)
  const [loading, setLoading] = useState(false)
  const [authChecked, setAuthChecked] = useState(false)
  const [userSearch, setUserSearch] = useState('')
  const [userStatusFilter, setUserStatusFilter] = useState<number | null>(null)
  const [previewPost, setPreviewPost] = useState<Post | null>(null)

  useEffect(() => {
    const token = localStorage.getItem('token')
    if (!token) { router.push('/login'); return }
    get<{ is_admin: boolean }>('/auth/profile').then(res => {
      if (res.code !== 0 || !res.data?.is_admin) {
        toastError('需要管理员权限')
        router.push('/')
      } else {
        setAuthChecked(true)
      }
    }).catch(() => { router.push('/login') })
  }, [router])

  const fetchStats = async () => {
    setLoading(true)
    try {
      const res = await get<Stats>('/admin/stats')
      if (res.code === 0) setStats(res.data)
    } catch { toastError('获取统计数据失败') }
    finally { setLoading(false) }
  }

  const fetchPosts = async (page: number) => {
    setLoading(true)
    try {
      const res = await get<{ items: Post[]; total: number }>(`/admin/posts?page=${page}&page_size=10`)
      if (res.code === 0) { setPosts(res.data.items || []); setPostsTotal(res.data.total) }
    } catch { toastError('获取帖子列表失败') }
    finally { setLoading(false) }
  }

  const fetchUsers = async (page: number, search?: string, status?: number | null) => {
    setLoading(true)
    try {
      let url = `/admin/users?page=${page}&page_size=10`
      if (search) url += `&search=${encodeURIComponent(search)}`
      if (status !== null && status !== undefined) url += `&status=${status}`
      const res = await get<{ items: User[]; total: number }>(url)
      if (res.code === 0) { setUsers(res.data.items || []); setUsersTotal(res.data.total) }
    } catch { toastError('获取用户列表失败') }
    finally { setLoading(false) }
  }

  useEffect(() => {
    if (tab === 'stats') fetchStats()
    else if (tab === 'posts') fetchPosts(postsPage)
    else if (tab === 'users') fetchUsers(usersPage, userSearch, userStatusFilter)
  }, [tab, postsPage, usersPage, userSearch, userStatusFilter])

  const handleDeletePost = async (id: number) => {
    if (!confirm('确定删除这条帖子吗？')) return
    try {
      const res = await del(`/admin/posts/${id}`)
      if (res.code === 0) {
        toastSuccess('帖子已删除')
        setPosts(posts.filter(p => p.id !== id))
        setPostsTotal(t => t - 1)
      }
    } catch { toastError('删除失败') }
  }

  const handleToggleTop = async (post: Post) => {
    const newIsTop = !post.is_top
    try {
      const res = await put(`/admin/posts/${post.id}/top`, { is_top: newIsTop })
      if (res.code === 0) {
        toastSuccess(newIsTop ? '已置顶' : '已取消置顶')
        setPosts(posts.map(p => p.id === post.id ? { ...p, is_top: newIsTop } : p))
      }
    } catch { toastError('操作失败') }
  }

  const handleToggleFeatured = async (post: Post) => {
    const newIsFeatured = !post.is_featured
    try {
      const res = await put(`/admin/posts/${post.id}/featured`, { is_featured: newIsFeatured })
      if (res.code === 0) {
        toastSuccess(newIsFeatured ? '已加精' : '已取消加精')
        setPosts(posts.map(p => p.id === post.id ? { ...p, is_featured: newIsFeatured } : p))
      }
    } catch { toastError('操作失败') }
  }

  const handleToggleUserStatus = async (user: User) => {
    const newStatus = user.status === 1 ? 0 : 1
    const action = newStatus === 0 ? '禁用' : '启用'
    if (!confirm(`确定${action}用户 ${user.nickname} 吗？`)) return
    try {
      const res = await put(`/admin/users/${user.id}/status`, { status: newStatus })
      if (res.code === 0) {
        toastSuccess(`用户已${action}`)
        setUsers(users.map(u => u.id === user.id ? { ...u, status: newStatus } : u))
      }
    } catch { toastError('操作失败') }
  }

  const tabs = [
    { key: 'stats' as const, label: '统计概览', icon: '📊' },
    { key: 'posts' as const, label: '帖子管理', icon: '📝' },
    { key: 'users' as const, label: '用户管理', icon: '👥' },
  ]

  if (!authChecked) return (
    <div className="flex items-center justify-center py-20">
      <p className="text-gray-500">验证管理员权限...</p>
    </div>
  )

  return (
    <div className="container mx-auto px-4 py-8 max-w-6xl">
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold text-gray-900">管理后台</h1>
        <div className="flex items-center gap-2">
          <Link href="/admin/moderation" className="btn-secondary text-sm flex items-center gap-1.5">
            🛡️ 敏感词
          </Link>
          <Link href="/admin/reports" className="btn-secondary text-sm flex items-center gap-1.5">
            🚨 举报管理
          </Link>
          <Link href="/admin/review" className="btn-secondary text-sm flex items-center gap-1.5">
            📋 帖子审核
          </Link>
          <Link href="/admin/comments" className="btn-secondary text-sm flex items-center gap-1.5">
            💬 评论管理
          </Link>
          <Link href="/admin/feedback" className="btn-secondary text-sm flex items-center gap-1.5">
            📝 反馈管理
          </Link>
          <Link href="/admin/images" className="btn-secondary text-sm flex items-center gap-1.5">
            🖼️ 风景图管理
          </Link>
          <Link href="/admin/settings" className="btn-secondary text-sm flex items-center gap-1.5">
            ⚙️ 系统设置
          </Link>
        </div>
      </div>

      {/* Tab 切换 */}
      <div className="flex gap-1 mb-6">
        {tabs.map((t) => (
          <button
            key={t.key}
            onClick={() => setTab(t.key)}
            className={`flex items-center gap-1.5 px-5 py-2.5 rounded-lg text-sm font-medium transition-all ${
              tab === t.key
                ? 'bg-blue-600 text-white shadow-md shadow-blue-200'
                : 'bg-white text-gray-600 hover:bg-gray-50 border border-gray-200'
            }`}
          >
            <span>{t.icon}</span>
            {t.label}
          </button>
        ))}
      </div>

      {loading ? (
        <div className="card p-12 text-center animate-pulse">
          <div className="h-8 bg-gray-200 rounded w-40 mx-auto mb-4" />
          <div className="h-4 bg-gray-200 rounded w-60 mx-auto" />
        </div>
      ) : (
        <>
          {/* 统计概览 */}
          {tab === 'stats' && stats && (
            <>
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
                {[
                  { label: '总用户数', value: stats.total_users, icon: '👥', color: 'from-blue-500 to-indigo-500' },
                  { label: '总帖子数', value: stats.total_posts, icon: '📝', color: 'from-emerald-500 to-teal-500' },
                  { label: '总评论数', value: stats.total_comments, icon: '💬', color: 'from-purple-500 to-violet-500' },
                  { label: '今日活跃', value: stats.active_users_today, icon: '🔥', color: 'from-orange-500 to-red-500' },
                ].map(s => (
                  <div key={s.label} className="card p-6">
                    <div className="flex items-center gap-3 mb-3">
                      <div className={`w-10 h-10 rounded-xl bg-gradient-to-br ${s.color} flex items-center justify-center text-lg`}>
                        {s.icon}
                      </div>
                      <span className="text-sm text-gray-500">{s.label}</span>
                    </div>
                    <div className="text-3xl font-bold text-gray-900">{(s.value ?? 0).toLocaleString()}</div>
                  </div>
                ))}
              </div>

              {/* 数据导出 */}
              <div className="card p-6 mt-6">
                <div className="flex items-center gap-3 mb-4">
                  <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-slate-500 to-gray-600 flex items-center justify-center text-lg">
                    📥
                  </div>
                  <div>
                    <h3 className="text-sm font-semibold text-gray-900">数据导出</h3>
                    <p className="text-xs text-gray-500">导出 CSV 格式数据文件</p>
                  </div>
                </div>
                <div className="flex items-center gap-3 flex-wrap">
                  <button
                    onClick={() => downloadFile('/admin/export/users', 'users.csv').catch(() => toastError('导出用户数据失败'))}
                    className="px-4 py-2 bg-white border border-gray-200 rounded-lg text-sm font-medium text-gray-700 hover:bg-gray-50 hover:border-gray-300 transition-all flex items-center gap-2"
                  >
                    <span>👥</span> 导出用户数据
                  </button>
                  <button
                    onClick={() => downloadFile('/admin/export/posts', 'posts.csv').catch(() => toastError('导出帖子数据失败'))}
                    className="px-4 py-2 bg-white border border-gray-200 rounded-lg text-sm font-medium text-gray-700 hover:bg-gray-50 hover:border-gray-300 transition-all flex items-center gap-2"
                  >
                    <span>📝</span> 导出帖子数据
                  </button>
                </div>
              </div>
            </>
          )}

          {/* 帖子管理 */}
          {tab === 'posts' && (
            <>
              <p className="text-sm text-gray-500 mb-4">共 <span className="font-medium text-gray-700">{postsTotal}</span> 条帖子</p>
              <div className="card overflow-hidden">
                <table className="w-full">
                  <thead className="bg-gray-50 border-b border-gray-100">
                    <tr>
                      {['ID', '标题', '类型', '互动', '时间', '操作'].map(h => (
                        <th key={h} className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">{h}</th>
                      ))}
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-gray-100">
                    {posts.map((post) => {
                      const t = TYPE_LABELS[post.type]
                      return (
                        <tr key={post.id} className="hover:bg-gray-50 transition-colors">
                          <td className="px-4 py-3 text-sm text-gray-500">{post.id}</td>
                          <td className="px-4 py-3 max-w-xs">
                            <button
                              onClick={() => setPreviewPost(post)}
                              className="text-sm font-medium text-blue-600 hover:text-blue-800 hover:underline truncate block text-left max-w-full"
                            >
                              {post.is_top && '📌 '}{post.is_featured && '⭐ '}{post.title}
                            </button>
                          </td>
                          <td className="px-4 py-3">
                            {t && <span className={`px-2 py-0.5 rounded-full text-xs ${t.color}`}>{t.label}</span>}
                          </td>
                          <td className="px-4 py-3 text-sm text-gray-500">👍{post.like_count} 💬{post.comment_count}</td>
                          <td className="px-4 py-3 text-sm text-gray-400">{new Date(post.created_at).toLocaleDateString()}</td>
                          <td className="px-4 py-3">
                            <div className="flex items-center gap-2">
                              <button
                                onClick={() => handleToggleTop(post)}
                                className={`text-sm font-medium ${post.is_top ? 'text-orange-600 hover:text-orange-800' : 'text-gray-500 hover:text-gray-700'}`}
                                title={post.is_top ? '取消置顶' : '置顶'}
                              >
                                {post.is_top ? '📌 取消置顶' : '📌 置顶'}
                              </button>
                              <button
                                onClick={() => handleToggleFeatured(post)}
                                className={`text-sm font-medium ${post.is_featured ? 'text-yellow-600 hover:text-yellow-800' : 'text-gray-500 hover:text-gray-700'}`}
                                title={post.is_featured ? '取消加精' : '加精'}
                              >
                                {post.is_featured ? '⭐ 取消加精' : '⭐ 加精'}
                              </button>
                              <button onClick={() => handleDeletePost(post.id)} className="text-sm text-red-600 hover:text-red-800 font-medium">删除</button>
                            </div>
                          </td>
                        </tr>
                      )
                    })}
                  </tbody>
                </table>
              </div>
              {Math.ceil(postsTotal / 10) > 1 && (
                <div className="mt-4 flex items-center justify-center gap-2">
                  <button onClick={() => setPostsPage(p => Math.max(1, p - 1))} disabled={postsPage === 1} className="btn-secondary text-sm px-3 py-1.5 disabled:opacity-40">←</button>
                  <span className="text-sm text-gray-500">{postsPage} / {Math.ceil(postsTotal / 10)}</span>
                  <button onClick={() => setPostsPage(p => Math.min(Math.ceil(postsTotal / 10), p + 1))} disabled={postsPage >= Math.ceil(postsTotal / 10)} className="btn-secondary text-sm px-3 py-1.5 disabled:opacity-40">→</button>
                </div>
              )}
            </>
          )}

          {/* 用户管理 */}
          {tab === 'users' && (
            <>
              {/* 用户搜索和筛选 */}
              <div className="card p-4 mb-4">
                <div className="flex items-center gap-4 flex-wrap">
                  <div className="flex-1 min-w-[200px]">
                    <input
                      type="text"
                      value={userSearch}
                      onChange={e => { setUserSearch(e.target.value); setUsersPage(1) }}
                      placeholder="搜索昵称/邮箱/手机号"
                      className="input-field text-sm w-full"
                    />
                  </div>
                  <div className="flex items-center gap-2">
                    <span className="text-sm text-gray-600">状态：</span>
                    <select
                      value={userStatusFilter ?? ''}
                      onChange={e => { setUserStatusFilter(e.target.value ? Number(e.target.value) : null); setUsersPage(1) }}
                      className="input-field text-sm"
                    >
                      <option value="">全部</option>
                      <option value="1">正常</option>
                      <option value="0">禁用</option>
                    </select>
                  </div>
                  <span className="text-sm text-gray-500">共 {usersTotal} 位用户</span>
                </div>
              </div>
              <div className="card overflow-hidden">
                <table className="w-full">
                  <thead className="bg-gray-50 border-b border-gray-100">
                    <tr>
                      {['ID', '昵称', '邮箱', '角色', '状态', '注册时间', '操作'].map(h => (
                        <th key={h} className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">{h}</th>
                      ))}
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-gray-100">
                    {users.map((user) => (
                      <tr key={user.id} className="hover:bg-gray-50 transition-colors">
                        <td className="px-4 py-3 text-sm text-gray-500">{user.id}</td>
                        <td className="px-4 py-3 text-sm font-medium text-gray-900">{user.nickname}</td>
                        <td className="px-4 py-3 text-sm text-gray-500">{user.email || '-'}</td>
                        <td className="px-4 py-3">
                          {user.is_admin
                            ? <span className="px-2 py-0.5 rounded-full text-xs bg-purple-50 text-purple-600">管理员</span>
                            : <span className="px-2 py-0.5 rounded-full text-xs bg-gray-100 text-gray-600">普通用户</span>
                          }
                        </td>
                        <td className="px-4 py-3">
                          {user.status === 1
                            ? <span className="px-2 py-0.5 rounded-full text-xs bg-green-50 text-green-600">正常</span>
                            : <span className="px-2 py-0.5 rounded-full text-xs bg-red-50 text-red-600">禁用</span>
                          }
                        </td>
                        <td className="px-4 py-3 text-sm text-gray-400">{new Date(user.created_at).toLocaleDateString()}</td>
                        <td className="px-4 py-3">
                          {!user.is_admin && (
                            <button
                              onClick={() => handleToggleUserStatus(user)}
                              className={`text-sm font-medium ${user.status === 1 ? 'text-red-600 hover:text-red-800' : 'text-green-600 hover:text-green-800'}`}
                            >
                              {user.status === 1 ? '禁用' : '启用'}
                            </button>
                          )}
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
              {Math.ceil(usersTotal / 10) > 1 && (
                <div className="mt-4 flex items-center justify-center gap-2">
                  <button onClick={() => setUsersPage(p => Math.max(1, p - 1))} disabled={usersPage === 1} className="btn-secondary text-sm px-3 py-1.5 disabled:opacity-40">←</button>
                  <span className="text-sm text-gray-500">{usersPage} / {Math.ceil(usersTotal / 10)}</span>
                  <button onClick={() => setUsersPage(p => Math.min(Math.ceil(usersTotal / 10), p + 1))} disabled={usersPage >= Math.ceil(usersTotal / 10)} className="btn-secondary text-sm px-3 py-1.5 disabled:opacity-40">→</button>
                </div>
              )}
            </>
          )}
        </>
      )}

      {/* 帖子预览弹窗 */}
      {previewPost && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4" onClick={() => setPreviewPost(null)}>
          <div className="bg-white rounded-xl max-w-2xl w-full max-h-[80vh] overflow-hidden shadow-2xl" onClick={e => e.stopPropagation()}>
            <div className="p-6">
              <div className="flex items-center justify-between mb-4">
                <h2 className="text-xl font-bold text-gray-900">
                  {previewPost.is_top && '📌 '}{previewPost.is_featured && '⭐ '}{previewPost.title}
                </h2>
                <button
                  onClick={() => setPreviewPost(null)}
                  className="text-gray-400 hover:text-gray-600 text-2xl"
                >
                  ×
                </button>
              </div>
              <div className="flex items-center gap-3 mb-4 text-sm text-gray-500">
                <span>ID: {previewPost.id}</span>
                <span>·</span>
                <span>{TYPE_LABELS[previewPost.type]?.label || previewPost.type}</span>
                <span>·</span>
                <span>👍 {previewPost.like_count}</span>
                <span>💬 {previewPost.comment_count}</span>
                <span>·</span>
                <span>{new Date(previewPost.created_at).toLocaleString()}</span>
              </div>
              <div className="prose max-w-none overflow-y-auto max-h-[50vh]">
                <p className="text-gray-700 whitespace-pre-wrap">{previewPost.content}</p>
              </div>
              <div className="mt-6 flex justify-end gap-3">
                <Link
                  href={`/posts/${previewPost.id}`}
                  className="btn-secondary text-sm"
                  target="_blank"
                >
                  查看原帖 →
                </Link>
                <button
                  onClick={() => setPreviewPost(null)}
                  className="btn-primary text-sm"
                >
                  关闭
                </button>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
