'use client'

import { useState, useEffect, useCallback } from 'react'
import Link from 'next/link'
import { get, put } from '@/lib/api'
import { timeAgo } from '@/lib/utils'
import { useNotificationWebSocket } from '@/lib/useWebSocket'

interface Notification {
  id: number
  type: string
  content: string
  target_id: number | null
  target_type: string | null
  is_read: boolean
  created_at: string
}

const TYPE_ICONS: Record<string, { icon: string; color: string }> = {
  like: { icon: '👍', color: 'bg-blue-50 text-blue-500' },
  comment: { icon: '💬', color: 'bg-green-50 text-green-500' },
  reply: { icon: '↩️', color: 'bg-purple-50 text-purple-500' },
  follow: { icon: '👤', color: 'bg-amber-50 text-amber-500' },
}

export default function NotificationsPage() {
  const [allItems, setAllItems] = useState<Notification[]>([])
  const [page, setPage] = useState(1)
  const [unreadCount, setUnreadCount] = useState(0)
  const [loading, setLoading] = useState(true)
  const [typeFilter, setTypeFilter] = useState<string>('all')
  const [selectMode, setSelectMode] = useState(false)
  const [selectedIds, setSelectedIds] = useState<Set<number>>(new Set())

  const filteredItems = typeFilter === 'all' ? allItems : allItems.filter(n => n.type === typeFilter)
  const pageSize = 20
  const totalFiltered = filteredItems.length
  const totalPages = Math.ceil(totalFiltered / pageSize)
  const pagedItems = filteredItems.slice((page - 1) * pageSize, page * pageSize)
  const unreadFiltered = filteredItems.filter(n => !n.is_read).length

  const { isConnected } = useNotificationWebSocket((notification) => {
    setAllItems(prev => [notification as unknown as Notification, ...prev])
    setUnreadCount((notification.unread_count as number) || 0)
  })

  // suppress unused warning
  void isConnected

  const loadNotifications = useCallback(async () => {
    setLoading(true)
    try {
      const res = await get<{ items: Notification[]; total: number }>(`/notifications?page=1&page_size=200`)
      if (res.code === 0) {
        setAllItems(res.data.items || [])
      }
    } finally {
      setLoading(false)
    }
  }, [])

  const loadUnreadCount = useCallback(async () => {
    const res = await get<{ count: number }>('/notifications/unread-count')
    if (res.code === 0) setUnreadCount(res.data.count)
  }, [])

  useEffect(() => { loadNotifications() }, [loadNotifications])
  useEffect(() => { setPage(1) }, [typeFilter])
  useEffect(() => { loadUnreadCount() }, [loadUnreadCount])

  const handleMarkAsRead = async (id: number) => {
    const res = await put(`/notifications/${id}/read`)
    if (res.code === 0) {
      setAllItems(allItems.map(n => n.id === id ? { ...n, is_read: true } : n))
      setUnreadCount(c => Math.max(0, c - 1))
    }
  }

  const handleMarkAllAsRead = async () => {
    const res = await put('/notifications/read-all')
    if (res.code === 0) {
      setAllItems(allItems.map(n => ({ ...n, is_read: true })))
      setUnreadCount(0)
    }
  }

  const handleBatchMarkAsRead = async () => {
    const ids = Array.from(selectedIds)
    await Promise.all(ids.map(id => put(`/notifications/${id}/read`)))
    setAllItems(allItems.map(n => selectedIds.has(n.id) ? { ...n, is_read: true } : n))
    setUnreadCount(c => Math.max(0, c - ids.length))
    setSelectedIds(new Set())
    setSelectMode(false)
  }

  return (
    <div className="container mx-auto px-4 py-8 max-w-2xl">
      {/* 头部操作栏 */}
      <div className="flex items-center justify-between mb-4">
        <h1 className="text-2xl font-bold text-gray-900 flex items-center gap-2">
          通知
          {unreadCount > 0 && (
            <span className="text-sm bg-red-500 text-white rounded-full px-2.5 py-0.5 font-medium">{unreadFiltered}</span>
          )}
        </h1>
        <div className="flex items-center gap-3">
          {selectMode ? (
            <>
              <button onClick={() => { setSelectMode(false); setSelectedIds(new Set()) }} className="text-sm text-gray-500 hover:text-gray-700">取消</button>
              {selectedIds.size > 0 && (
                <button onClick={handleBatchMarkAsRead} className="text-sm text-blue-600 hover:text-blue-500 font-medium">批量已读 ({selectedIds.size})</button>
              )}
            </>
          ) : (
            <>
              {unreadFiltered > 0 && (
                <button onClick={handleMarkAllAsRead} className="text-sm text-blue-600 hover:text-blue-500 font-medium">全部已读</button>
              )}
              <button onClick={() => setSelectMode(true)} className="text-sm text-gray-500 hover:text-gray-700">选择</button>
            </>
          )}
        </div>
      </div>

      {/* 类型筛选 */}
      <div className="flex gap-1 mb-4 overflow-x-auto pb-1">
        {[{ key: 'all', label: '全部', icon: '📢' },
          { key: 'like', label: '点赞', icon: '👍' },
          { key: 'comment', label: '评论', icon: '💬' },
          { key: 'reply', label: '回复', icon: '↩️' },
          { key: 'follow', label: '关注', icon: '👤' }].map(t => (
          <button
            key={t.key}
            onClick={() => setTypeFilter(t.key)}
            className={`flex items-center gap-1 px-3 py-2 rounded-lg text-sm font-medium transition-all whitespace-nowrap ${
              typeFilter === t.key
                ? 'bg-blue-600 text-white shadow-md'
                : 'bg-white text-gray-600 hover:bg-gray-50 border border-gray-200'
            }`}
          >
            {t.icon} {t.label}
          </button>
        ))}
      </div>

      {loading ? (
        <div className="space-y-3">
          {[...Array(5)].map((_, i) => (
            <div key={i} className="card p-4 flex items-start gap-3">
              <div className="w-10 h-10 rounded-full bg-gray-200 flex-shrink-0 animate-pulse" />
              <div className="flex-1">
                <div className="h-4 bg-gray-200 rounded w-3/4 mb-2 animate-pulse" />
                <div className="h-3 bg-gray-200 rounded w-1/4 animate-pulse" />
              </div>
            </div>
          ))}
        </div>
      ) : pagedItems.length === 0 ? (
        <div className="card p-16 text-center">
          <div className="text-5xl mb-3">🔔</div>
          <p className="text-gray-500 text-lg">暂无通知</p>
          <p className="text-sm text-gray-400 mt-1">有新消息时会在这里通知你</p>
        </div>
      ) : (
        <div className="space-y-2">
          {pagedItems.map((n) => {
            const link = n.target_type === 'post' && n.target_id ? `/posts/${n.target_id}` : null
            const typeInfo = TYPE_ICONS[n.type] || { icon: '📢', color: 'bg-gray-50 text-gray-500' }
            return (
              <div
                key={n.id}
                className={`card p-4 flex items-start gap-3 ${!n.is_read ? 'border-l-4 border-blue-500 bg-blue-50/30' : ''} ${selectedIds.has(n.id) ? 'ring-2 ring-blue-400' : ''}`}
              >
                {selectMode && (
                  <input
                    type="checkbox"
                    checked={selectedIds.has(n.id)}
                    onChange={() => {
                      const next = new Set(selectedIds)
                      if (next.has(n.id)) next.delete(n.id)
                      else next.add(n.id)
                      setSelectedIds(next)
                    }}
                    className="mt-2 w-4 h-4 text-blue-600 rounded flex-shrink-0"
                  />
                )}
                <div className={`w-10 h-10 rounded-full flex items-center justify-center text-lg flex-shrink-0 ${typeInfo.color}`}>
                  {typeInfo.icon}
                </div>
                <div className="flex-1 min-w-0">
                  {link ? (
                    <Link href={link} className="text-sm text-gray-800 hover:text-blue-600 leading-relaxed">
                      {n.content}
                    </Link>
                  ) : (
                    <p className="text-sm text-gray-800 leading-relaxed">{n.content}</p>
                  )}
                  <p className="text-xs text-gray-400 mt-1">{timeAgo(n.created_at)}</p>
                </div>
                {!n.is_read && (
                  <button
                    onClick={() => handleMarkAsRead(n.id)}
                    className="text-xs text-blue-600 hover:text-blue-500 font-medium whitespace-nowrap mt-1"
                  >
                    已读
                  </button>
                )}
              </div>
            )
          })}
        </div>
      )}

      {totalPages > 1 && (
        <div className="mt-6 flex items-center justify-center gap-2">
          <button
            onClick={() => setPage(p => Math.max(1, p - 1))}
            disabled={page === 1}
            className="btn-secondary text-sm px-4 py-2 disabled:opacity-40"
          >
            ← 上一页
          </button>
          <span className="text-sm text-gray-500 px-3">{page} / {totalPages}</span>
          <button
            onClick={() => setPage(p => Math.min(totalPages, p + 1))}
            disabled={page === totalPages}
            className="btn-secondary text-sm px-4 py-2 disabled:opacity-40"
          >
            下一页 →
          </button>
        </div>
      )}
    </div>
  )
}
