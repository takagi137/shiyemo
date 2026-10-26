'use client'

import { useState, useEffect } from 'react'

interface DailyImage {
  url: string
  description: string
  credit: string
}

interface DailyQuote {
  text: string
  author: string
}

interface DailyData {
  image: DailyImage
  quote: DailyQuote
}

export default function DailyScene() {
  const [data, setData] = useState<DailyData | null>(null)
  const [loading, setLoading] = useState(true)
  const [imageLoaded, setImageLoaded] = useState(false)

  useEffect(() => {
    fetchDaily()
  }, [])

  const fetchDaily = async () => {
    try {
      const res = await fetch('/api/v1/daily')
      const json = await res.json()
      if (json.code === 0) {
        setData(json.data)
      }
    } catch (err) {
      console.error('Failed to fetch daily content:', err)
    } finally {
      setLoading(false)
    }
  }

  const refreshQuote = async () => {
    try {
      const res = await fetch('/api/v1/daily/quote')
      const json = await res.json()
      if (json.code === 0 && data) {
        setData({ ...data, quote: json.data })
      }
    } catch (err) {
      console.error('Failed to refresh quote:', err)
    }
  }

  const refreshImage = async () => {
    try {
      const res = await fetch('/api/v1/daily/image')
      const json = await res.json()
      if (json.code === 0 && data) {
        // 仅当图片确实变化时才做淡入过渡：URL 相同时 <img> 不会重新触发 onLoad，
        // 若仍置为未加载，图片会保持透明（视觉上"消失"）。
        if (json.data.url !== data.image.url) {
          setImageLoaded(false)
        }
        setData({ ...data, image: json.data })
      }
    } catch (err) {
      console.error('Failed to refresh image:', err)
    }
  }

  if (loading) {
    return (
      <div className="relative w-full h-screen flex items-center justify-center bg-gradient-to-br from-blue-900 via-purple-900 to-indigo-900">
        <div className="text-white text-xl">加载中...</div>
      </div>
    )
  }

  return (
    <div className="relative w-full h-screen overflow-hidden">
      {/* 背景渐变（图片加载前/失败时显示） */}
      <div className="absolute inset-0 bg-gradient-to-br from-blue-900 via-purple-900 to-indigo-900" />

      {/* 风景图 */}
      {data?.image.url && (
        <div
          className={`absolute inset-0 transition-opacity duration-1000 ${
            imageLoaded ? 'opacity-100' : 'opacity-0'
          }`}
        >
          <img
            src={data.image.url}
            alt={data.image.description}
            className="w-full h-full object-cover"
            onLoad={() => setImageLoaded(true)}
          />
        </div>
      )}

      {/* 渐变遮罩（保证文字可读） */}
      <div className="absolute inset-0 bg-gradient-to-t from-black/60 via-black/20 to-transparent" />

      {/* 内容 */}
      <div className="relative z-10 flex flex-col items-center justify-center h-full px-6 text-center">
        {/* 温暖文案 */}
        <div className="max-w-2xl mx-auto mb-8">
          <p className="text-2xl md:text-3xl lg:text-4xl text-white font-light leading-relaxed drop-shadow-lg">
            {data?.quote.text || '今天的你，已经比昨天勇敢了一点。'}
          </p>
          {data?.quote.author && (
            <p className="mt-4 text-white/70 text-sm">
              —— {data.quote.author}
            </p>
          )}
        </div>

        {/* 刷新按钮组 */}
        <div className="flex space-x-4 mb-12">
          <button
            onClick={refreshImage}
            className="px-4 py-2 bg-white/20 hover:bg-white/30 text-white text-sm rounded-full backdrop-blur-sm transition-all duration-300"
          >
            🔄 换一张
          </button>
          <button
            onClick={refreshQuote}
            className="px-4 py-2 bg-white/20 hover:bg-white/30 text-white text-sm rounded-full backdrop-blur-sm transition-all duration-300"
          >
            💬 换一句
          </button>
        </div>

        {/* 进入社区按钮 */}
        <a
          href="/feed"
          className="px-8 py-3 bg-white/90 hover:bg-white text-gray-900 font-medium rounded-full transition-all duration-300 shadow-lg hover:shadow-xl"
        >
          进入社区 →
        </a>
      </div>

      {/* 图片版权信息 */}
      {data?.image.credit && (
        <div className="absolute bottom-4 right-4 text-white/50 text-xs">
          {data.image.credit}
        </div>
      )}
    </div>
  )
}
