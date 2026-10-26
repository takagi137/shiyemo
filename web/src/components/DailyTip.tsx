'use client'

import { useState, useEffect } from 'react'

interface Tip {
  id: number
  content: string
  category: string
}

interface DailyTipProps {
  className?: string
}

export default function DailyTip({ className = '' }: DailyTipProps) {
  const [tip, setTip] = useState<Tip | null>(null)
  const [loading, setLoading] = useState(true)
  const [refreshing, setRefreshing] = useState(false)

  const fetchTip = async () => {
    try {
      const response = await fetch('/api/v1/tips/random')
      const data = await response.json()
      if (data.code === 0 && data.data) {
        setTip(data.data)
      }
    } catch (error) {
      console.error('获取贴士失败:', error)
    } finally {
      setLoading(false)
      setRefreshing(false)
    }
  }

  useEffect(() => {
    fetchTip()
  }, [])

  const handleRefresh = () => {
    setRefreshing(true)
    fetchTip()
  }

  const categoryIcons: Record<string, string> = {
    mindset: '🧠',
    finance: '💰',
    life: '🌱',
    policy: '📋',
    general: '💡',
  }

  const categoryNames: Record<string, string> = {
    mindset: '心态',
    finance: '财务',
    life: '生活',
    policy: '政策',
    general: '通用',
  }

  if (loading) {
    return (
      <div className={`bg-gradient-to-br from-amber-50 to-orange-50 rounded-xl p-6 border border-amber-100 ${className}`}>
        <div className="animate-pulse">
          <div className="h-4 bg-amber-200 rounded w-1/3 mb-4"></div>
          <div className="h-6 bg-amber-200 rounded w-full mb-2"></div>
          <div className="h-6 bg-amber-200 rounded w-2/3"></div>
        </div>
      </div>
    )
  }

  if (!tip) {
    return null
  }

  return (
    <div className={`bg-gradient-to-br from-amber-50 to-orange-50 rounded-xl p-6 border border-amber-100 ${className}`}>
      <div className="flex items-center justify-between mb-4">
        <div className="flex items-center gap-2">
          <span className="text-xl">{categoryIcons[tip.category] || '💡'}</span>
          <span className="text-sm font-medium text-amber-700">
            今日贴士 · {categoryNames[tip.category] || '通用'}
          </span>
        </div>
        <button
          onClick={handleRefresh}
          disabled={refreshing}
          className="text-amber-600 hover:text-amber-800 transition-colors"
          title="换一条"
        >
          <svg
            className={`w-5 h-5 ${refreshing ? 'animate-spin' : ''}`}
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"
            />
          </svg>
        </button>
      </div>
      <p className="text-gray-800 text-lg leading-relaxed">{tip.content}</p>
    </div>
  )
}
