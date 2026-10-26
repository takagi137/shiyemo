'use client'

import { useState, useEffect } from 'react'
import { get } from '@/lib/api'

interface MoodCurveItem {
  date: string
  mood: string
  emoji: string
}

interface MoodCurveProps {
  userId: number
}

export default function MoodCurve({ userId }: MoodCurveProps) {
  const [data, setData] = useState<MoodCurveItem[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    const fetchMoodCurve = async () => {
      try {
        const res = await get<MoodCurveItem[]>('/checkin/mood-curve')
        if (res.code === 0) {
          setData(res.data || [])
        } else {
          setError(res.message || '获取心情数据失败')
        }
      } catch (err) {
        setError('网络错误，请稍后重试')
      } finally {
        setLoading(false)
      }
    }

    fetchMoodCurve()
  }, [userId])

  if (loading) {
    return (
      <div className="card p-6">
        <h3 className="text-lg font-semibold text-gray-900 mb-4">📊 近 30 天心情曲线</h3>
        <div className="animate-pulse">
          <div className="h-48 bg-gray-200 rounded-lg"></div>
        </div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="card p-6">
        <h3 className="text-lg font-semibold text-gray-900 mb-4">📊 近 30 天心情曲线</h3>
        <div className="text-center py-8 text-gray-500">
          <p>{error}</p>
        </div>
      </div>
    )
  }

  if (data.length === 0) {
    return (
      <div className="card p-6">
        <h3 className="text-lg font-semibold text-gray-900 mb-4">📊 近 30 天心情曲线</h3>
        <div className="text-center py-8">
          <div className="text-4xl mb-3">📈</div>
          <p className="text-gray-500">暂无打卡数据</p>
          <p className="text-sm text-gray-400 mt-1">开始记录你的心情吧</p>
        </div>
      </div>
    )
  }

  // 计算心情统计
  const moodStats = data.reduce((acc, item) => {
    if (item.mood !== 'none') {
      acc[item.mood] = (acc[item.mood] || 0) + 1
    }
    return acc
  }, {} as Record<string, number>)

  const totalCheckins = Object.values(moodStats).reduce((a, b) => a + b, 0)

  // 获取心情颜色
  const getMoodColor = (mood: string) => {
    const colors: Record<string, string> = {
      happy: 'bg-yellow-400',
      calm: 'bg-blue-400',
      sad: 'bg-gray-400',
      angry: 'bg-red-400',
      anxious: 'bg-purple-400',
      none: 'bg-gray-100',
    }
    return colors[mood] || 'bg-gray-100'
  }

  // 获取心情标签
  const getMoodLabel = (mood: string) => {
    const labels: Record<string, string> = {
      happy: '开心',
      calm: '平静',
      sad: '难过',
      angry: '生气',
      anxious: '焦虑',
      none: '未打卡',
    }
    return labels[mood] || '未知'
  }

  // 格式化日期显示
  const formatDate = (dateStr: string) => {
    const date = new Date(dateStr)
    return `${date.getMonth() + 1}/${date.getDate()}`
  }

  // 获取今天和昨天的日期
  const today = new Date().toISOString().split('T')[0]
  const yesterday = new Date(Date.now() - 86400000).toISOString().split('T')[0]

  return (
    <div className="card p-6">
      <h3 className="text-lg font-semibold text-gray-900 mb-4">📊 近 30 天心情曲线</h3>
      
      {/* 心情日历网格 */}
      <div className="mb-6">
        <div className="flex flex-wrap gap-1.5">
          {data.map((item, index) => {
            const isToday = item.date === today
            const isYesterday = item.date === yesterday
            const dateLabel = isToday ? '今天' : isYesterday ? '昨天' : formatDate(item.date)
            
            return (
              <div
                key={index}
                className="relative group"
                title={`${dateLabel}: ${getMoodLabel(item.mood)} ${item.emoji}`}
              >
                <div
                  className={`w-8 h-8 rounded-md flex items-center justify-center text-sm transition-all ${
                    item.mood === 'none'
                      ? 'bg-gray-50 border border-gray-200'
                      : `${getMoodColor(item.mood)} text-white shadow-sm`
                  } ${isToday ? 'ring-2 ring-blue-500 ring-offset-1' : ''}`}
                >
                  {item.mood !== 'none' ? item.emoji : ''}
                </div>
                
                {/* Tooltip */}
                <div className="absolute bottom-full left-1/2 transform -translate-x-1/2 mb-2 px-2 py-1 bg-gray-900 text-white text-xs rounded opacity-0 group-hover:opacity-100 transition-opacity whitespace-nowrap pointer-events-none z-10">
                  {dateLabel}: {getMoodLabel(item.mood)}
                  {item.mood !== 'none' && ` ${item.emoji}`}
                </div>
              </div>
            )
          })}
        </div>
        
        {/* 图例 */}
        <div className="flex flex-wrap gap-3 mt-4 text-xs text-gray-500">
          {Object.entries(moodStats).map(([mood, count]) => (
            <div key={mood} className="flex items-center gap-1">
              <div className={`w-3 h-3 rounded-sm ${getMoodColor(mood)}`}></div>
              <span>{getMoodLabel(mood)}: {count}天</span>
            </div>
          ))}
        </div>
      </div>

      {/* 统计卡片 */}
      <div className="grid grid-cols-2 gap-4">
        <div className="bg-gradient-to-br from-blue-50 to-blue-100 rounded-lg p-4">
          <div className="text-2xl font-bold text-blue-700">{totalCheckins}</div>
          <div className="text-sm text-blue-600">累计打卡</div>
        </div>
        <div className="bg-gradient-to-br from-green-50 to-green-100 rounded-lg p-4">
          <div className="text-2xl font-bold text-green-700">
            {data.filter(item => item.mood !== 'none').length}
          </div>
          <div className="text-sm text-green-600">打卡天数</div>
        </div>
      </div>

      {/* 心情分布 */}
      {totalCheckins > 0 && (
        <div className="mt-4 pt-4 border-t border-gray-100">
          <h4 className="text-sm font-medium text-gray-700 mb-3">心情分布</h4>
          <div className="space-y-2">
            {Object.entries(moodStats)
              .sort(([, a], [, b]) => b - a)
              .map(([mood, count]) => {
                const percentage = Math.round((count / totalCheckins) * 100)
                return (
                  <div key={mood} className="flex items-center gap-2">
                    <div className="w-16 text-sm text-gray-600">{getMoodLabel(mood)}</div>
                    <div className="flex-1 bg-gray-100 rounded-full h-2">
                      <div
                        className={`h-2 rounded-full ${getMoodColor(mood)}`}
                        style={{ width: `${percentage}%` }}
                      ></div>
                    </div>
                    <div className="w-12 text-right text-sm text-gray-500">{percentage}%</div>
                  </div>
                )
              })}
          </div>
        </div>
      )}
    </div>
  )
}
