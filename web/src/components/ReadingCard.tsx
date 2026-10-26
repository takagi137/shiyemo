'use client'

import { useState, useEffect } from 'react'
import Icon from '@/components/Icon'

interface Reading {
  id: number
  source: string
  dynasty: string
  author: string
  title: string
  content: string
  interpretation: string
  image_url: string
}

export default function ReadingCard() {
  const [reading, setReading] = useState<Reading | null>(null)
  const [loading, setLoading] = useState(true)
  const [seenIds, setSeenIds] = useState<number[]>([])

  useEffect(() => {
    fetchToday()
  }, [])

  const fetchToday = async () => {
    setLoading(true)
    try {
      const res = await fetch('/api/v1/readings/today')
      const json = await res.json()
      if (json.code === 0) {
        setReading(json.data)
        setSeenIds([json.data.id])
      }
    } catch (err) {
      console.error('Failed to fetch today reading:', err)
    } finally {
      setLoading(false)
    }
  }

  const fetchRandom = async () => {
    setLoading(true)
    try {
      const exclude = seenIds.join(',')
      const res = await fetch(`/api/v1/readings/random?exclude=${exclude}`)
      const json = await res.json()
      if (json.code === 0) {
        setReading(json.data)
        setSeenIds(prev => [...prev, json.data.id])
      }
    } catch (err) {
      console.error('Failed to fetch random reading:', err)
    } finally {
      setLoading(false)
    }
  }

  if (loading && !reading) {
    return (
      <div className="bg-surface/70 rounded-xl p-6 flex items-center justify-center h-48">
        <div className="text-ink-faint text-sm">加载中...</div>
      </div>
    )
  }

  if (!reading) {
    return (
      <div className="bg-surface/70 rounded-xl p-6 text-center">
        <p className="text-ink-soft text-sm italic">
          &ldquo;古人也曾在风雨中前行，你并不孤单。&rdquo;
        </p>
      </div>
    )
  }

  return (
    <div className="bg-surface/70 rounded-xl p-6 space-y-4">
      <h3 className="text-lg font-medium text-ink">今日阅读</h3>

      {/* 来源信息 */}
      <div className="flex items-center space-x-2 text-xs text-ink-soft">
        {reading.dynasty && (
          <span className="px-2 py-0.5 bg-insight-soft text-insight-ink rounded-full">
            {reading.dynasty}
          </span>
        )}
        <span>{reading.author}</span>
        {reading.source && <span>· {reading.source}</span>}
      </div>

      {/* 原文 */}
      <blockquote className="text-ink text-base leading-relaxed pl-4 border-l-2 border-accent-line whitespace-pre-wrap">
        {reading.content}
      </blockquote>

      {/* 背景解读 */}
      {reading.interpretation && (
        <div className="bg-accent-soft/50 rounded-lg p-3">
          <p className="text-xs text-ink-soft leading-relaxed flex gap-1.5">
            <Icon name="insight" size={14} className="mt-0.5" />
            <span>{reading.interpretation}</span>
          </p>
        </div>
      )}

      {/* 操作按钮 */}
      <div className="flex justify-end space-x-3">
        <button
          onClick={fetchRandom}
          className="px-4 py-1.5 text-xs text-accent bg-accent-soft rounded-full hover:bg-accent-line transition-colors"
        >
          下一篇 →
        </button>
      </div>
    </div>
  )
}
