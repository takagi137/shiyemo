'use client'

import { useEffect, useState } from 'react'
import BreathingCircle from '@/components/BreathingCircle'
import AmbientPlayer from '@/components/AmbientPlayer'
import JournalEntry from '@/components/JournalEntry'
import ReadingCard from '@/components/ReadingCard'

export default function SanctuaryPage() {
  const [mounted, setMounted] = useState(false)

  useEffect(() => {
    // 静思空间对所有人开放，无需登录。
    // 需要登录的私密日记由 JournalEntry 组件内部按登录态自行处理。
    setMounted(true)
  }, [])

  if (!mounted) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gradient-to-b from-canvas to-accent-soft">
        <div className="text-ink-faint">加载中...</div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gradient-to-b from-canvas to-accent-soft">
      <div className="max-w-4xl mx-auto px-4 py-8">
        {/* 页面标题 */}
        <div className="text-center mb-10">
          <h1 className="text-2xl md:text-3xl font-light text-ink mb-2">
            静思空间
          </h1>
          <p className="text-sm text-ink-soft">
            放下手机，深呼吸，给自己一点时间
          </p>
        </div>

        {/* 四模块布局 */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-10">
          {/* 呼吸练习 */}
          <div className="bg-surface/50 backdrop-blur-sm rounded-2xl p-6 flex items-center justify-center min-h-[280px]">
            <BreathingCircle />
          </div>

          {/* 环境音 */}
          <div className="bg-surface/50 backdrop-blur-sm rounded-2xl p-6 flex items-center justify-center min-h-[280px]">
            <AmbientPlayer />
          </div>

          {/* 今日阅读 */}
          <div className="bg-surface/50 backdrop-blur-sm rounded-2xl p-6">
            <ReadingCard />
          </div>

          {/* 私密日记 */}
          <div className="bg-surface/50 backdrop-blur-sm rounded-2xl p-6">
            <JournalEntry />
          </div>
        </div>

        {/* 每日温暖文案 */}
        <div className="text-center mb-10">
          <p className="text-sm text-ink-faint italic">
            &ldquo;你不需要完美，你只需要真实。&rdquo;
          </p>
        </div>

        {/* 返回社区 */}
        <div className="text-center">
          <a
            href="/feed"
            className="inline-block px-6 py-2 text-sm text-accent bg-surface/70 hover:bg-surface rounded-full transition-all duration-300 shadow-sm hover:shadow-md"
          >
            返回社区 →
          </a>
        </div>
      </div>
    </div>
  )
}
