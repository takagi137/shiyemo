'use client'

import { useState, useEffect } from 'react'
import Link from 'next/link'
import { get } from '@/lib/api'
import Avatar from '@/components/Avatar'
import TrendChart from '@/components/TrendChart'

interface CohortStats {
  city_count: number
  industry_count: number
  time_range_count: number
  total_peers: number
}

interface Peer {
  id: number
  nickname: string
  avatar: string
  bio: string
  city: string
  industry: string
  left_job_at: string
}

export default function CohortPage() {
  const [stats, setStats] = useState<CohortStats | null>(null)
  const [peers, setPeers] = useState<Peer[]>([])
  const [loading, setLoading] = useState(true)
  const [activeTab, setActiveTab] = useState<'stats' | 'peers' | 'trends'>('stats')

  useEffect(() => {
    const fetchData = async () => {
      setLoading(true)
      try {
        const [statsRes, peersRes] = await Promise.all([
          get<CohortStats>('/cohort/stats'),
          get<{ items: Peer[] }>('/cohort/peers?page=1&page_size=20'),
        ])
        if (statsRes.code === 0) setStats(statsRes.data)
        if (peersRes.code === 0) setPeers(peersRes.data.items || [])
      } catch (err) {
        console.error('获取同温层数据失败:', err)
      } finally {
        setLoading(false)
      }
    }
    fetchData()
  }, [])

  if (loading) {
    return (
      <div className="container mx-auto px-4 py-8 max-w-4xl">
        <div className="animate-pulse space-y-6">
          <div className="h-8 bg-gray-200 rounded w-48"></div>
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
            {[...Array(4)].map((_, i) => (
              <div key={i} className="card p-6">
                <div className="h-8 bg-gray-200 rounded w-16 mb-2"></div>
                <div className="h-4 bg-gray-200 rounded w-24"></div>
              </div>
            ))}
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="container mx-auto px-4 py-8 max-w-4xl">
      {/* 标题 */}
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900 mb-2">🔥 同温层</h1>
        <p className="text-gray-500">找到和你相似的人，一起前行</p>
      </div>

      {/* Tab 切换 */}
      <div className="flex gap-2 mb-6">
        <button
          onClick={() => setActiveTab('stats')}
          className={`px-4 py-2 rounded-lg text-sm font-medium transition-all ${
            activeTab === 'stats'
              ? 'bg-blue-600 text-white shadow-md'
              : 'bg-white text-gray-600 hover:bg-gray-50 border border-gray-200'
          }`}
        >
          📊 数据统计
        </button>
        <button
          onClick={() => setActiveTab('peers')}
          className={`px-4 py-2 rounded-lg text-sm font-medium transition-all ${
            activeTab === 'peers'
              ? 'bg-blue-600 text-white shadow-md'
              : 'bg-white text-gray-600 hover:bg-gray-50 border border-gray-200'
          }`}
        >
          👥 同温层伙伴
        </button>
        <button
          onClick={() => setActiveTab('trends')}
          className={`px-4 py-2 rounded-lg text-sm font-medium transition-all ${
            activeTab === 'trends'
              ? 'bg-blue-600 text-white shadow-md'
              : 'bg-white text-gray-600 hover:bg-gray-50 border border-gray-200'
          }`}
        >
          📊 趋势数据
        </button>
      </div>

      {/* 统计卡片 */}
      {activeTab === 'stats' && stats && (
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-8">
          <div className="card p-6 text-center">
            <div className="text-3xl font-bold text-blue-600">{stats.city_count}</div>
            <div className="text-sm text-gray-500 mt-1">🏙️ 同城伙伴</div>
          </div>
          <div className="card p-6 text-center">
            <div className="text-3xl font-bold text-purple-600">{stats.industry_count}</div>
            <div className="text-sm text-gray-500 mt-1">💼 同行伙伴</div>
          </div>
          <div className="card p-6 text-center">
            <div className="text-3xl font-bold text-amber-600">{stats.time_range_count}</div>
            <div className="text-sm text-gray-500 mt-1">⏰ 同时期</div>
          </div>
          <div className="card p-6 text-center">
            <div className="text-3xl font-bold text-emerald-600">{stats.total_peers}</div>
            <div className="text-sm text-gray-500 mt-1">👥 总伙伴</div>
          </div>
        </div>
      )}

      {/* 同温层伙伴列表 */}
      {activeTab === 'peers' && (
        <div className="space-y-4">
          {peers.length === 0 ? (
            <div className="card p-12 text-center">
              <div className="text-4xl mb-3">🔍</div>
              <p className="text-gray-500 text-lg">暂无匹配的同温层伙伴</p>
              <p className="text-gray-400 text-sm mt-1">完善个人资料（城市、行业、离职时间）可以帮助匹配</p>
              <Link href="/profile" className="inline-block mt-4 btn-primary">
                去完善资料
              </Link>
            </div>
          ) : (
            peers.map((peer) => (
              <Link
                key={peer.id}
                href={`/profile?id=${peer.id}`}
                className="card p-4 flex items-center gap-4 hover:shadow-md transition-shadow"
              >
                <Avatar nickname={peer.nickname} avatar={peer.avatar} size="lg" />
                <div className="flex-1 min-w-0">
                  <div className="font-medium text-gray-900">{peer.nickname}</div>
                  {peer.bio && <div className="text-sm text-gray-500 mt-1 line-clamp-1">{peer.bio}</div>}
                  <div className="flex items-center gap-3 mt-2 text-xs text-gray-400">
                    {peer.city && <span>🏙️ {peer.city}</span>}
                    {peer.industry && <span>💼 {peer.industry}</span>}
                    {peer.left_job_at && (
                      <span>📅 离职 {new Date(peer.left_job_at).toLocaleDateString()}</span>
                    )}
                  </div>
                </div>
                <div className="text-blue-600 text-sm">查看 →</div>
              </Link>
            ))
          )}
        </div>
      )}

      {/* 趋势数据 */}
      {activeTab === 'trends' && (
        <TrendChart days={30} />
      )}

      {/* 快捷入口 */}
      <div className="mt-8 grid grid-cols-1 md:grid-cols-2 gap-4">
        <Link href="/feed?type=cohort" className="card p-6 hover:shadow-md transition-shadow">
          <div className="text-2xl mb-2">📝</div>
          <div className="font-medium text-gray-900">同温层帖子</div>
          <div className="text-sm text-gray-500 mt-1">看看同温层的伙伴们在说什么</div>
        </Link>
        <Link href="/profile" className="card p-6 hover:shadow-md transition-shadow">
          <div className="text-2xl mb-2">✏️</div>
          <div className="font-medium text-gray-900">完善资料</div>
          <div className="text-sm text-gray-500 mt-1">填写城市、行业、离职时间，匹配更精准</div>
        </Link>
      </div>
    </div>
  )
}
