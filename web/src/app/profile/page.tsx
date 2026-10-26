'use client'

import { useState, useEffect } from 'react'
import { useRouter } from 'next/navigation'
import Link from 'next/link'
import { get } from '@/lib/api'
import Avatar from '@/components/Avatar'
import Icon from '@/components/Icon'
import LoadingSpinner from '@/components/LoadingSpinner'

interface UserProfile {
  id: number
  nickname: string
  email: string
  bio: string
  avatar: string
  company: string
  occupation: string
  follower_count: number
  following_count: number
}

interface UserStats {
  post_count: number
  like_received: number
  follower_count: number
  following_count: number
}

export default function ProfilePage() {
  const router = useRouter()
  const [profile, setProfile] = useState<UserProfile | null>(null)
  const [stats, setStats] = useState<UserStats | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    const userStr = localStorage.getItem('user')
    const token = localStorage.getItem('token')
    if (!userStr || !token) {
      router.push('/login')
      return
    }

    try {
      const user = JSON.parse(userStr)
      const fetchProfile = async () => {
        const [profileRes, statsRes] = await Promise.all([
          get<UserProfile>(`/users/${user.id}`),
          get<UserStats>(`/users/${user.id}/stats`),
        ])
        if (profileRes.code === 0) setProfile(profileRes.data)
        if (statsRes.code === 0) setStats(statsRes.data)
        setLoading(false)
      }
      fetchProfile()
    } catch {
      router.push('/login')
    }
  }, [router])

  const handleLogout = () => {
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    router.push('/')
  }

  if (loading) return <LoadingSpinner />
  if (!profile) return null

  return (
    <div className="max-w-4xl mx-auto px-4 py-6">
      {/* 用户信息卡片 */}
      <div className="bg-white rounded-lg border border-line p-6 mb-6">
        <div className="flex items-start gap-6">
          <Avatar
            userId={profile.id}
            nickname={profile.nickname}
            avatar={profile.avatar}
            size="lg"
            clickable={false}
          />
          <div className="flex-1">
            <h1 className="text-xl font-bold text-ink mb-1">{profile.nickname}</h1>
            <p className="text-sm text-ink-faint mb-2">{profile.email}</p>
            {profile.bio && <p className="text-ink-soft mb-2">{profile.bio}</p>}
            <div className="flex items-center gap-2 text-sm text-ink-faint">
              {profile.occupation && <span>{profile.occupation}</span>}
              {profile.occupation && profile.company && <span>·</span>}
              {profile.company && <span>{profile.company}</span>}
            </div>
          </div>
        </div>

        {stats && (
          <div className="grid grid-cols-4 gap-4 mt-6 pt-6 border-t border-line-soft">
            <div className="text-center">
              <div className="text-xl font-bold text-ink">{stats.post_count}</div>
              <div className="text-sm text-ink-faint">帖子</div>
            </div>
            <div className="text-center">
              <div className="text-xl font-bold text-ink">{stats.like_received}</div>
              <div className="text-sm text-ink-faint">获赞</div>
            </div>
            <div className="text-center">
              <div className="text-xl font-bold text-ink">{stats.following_count}</div>
              <div className="text-sm text-ink-faint">关注</div>
            </div>
            <div className="text-center">
              <div className="text-xl font-bold text-ink">{stats.follower_count}</div>
              <div className="text-sm text-ink-faint">粉丝</div>
            </div>
          </div>
        )}
      </div>

      {/* 功能入口 */}
      <div className="bg-white rounded-lg border border-line divide-y divide-line-soft">
        <Link
          href={`/users/${profile.id}`}
          className="flex items-center justify-between px-6 py-4 hover:bg-surface-sunk transition-colors"
        >
          <div className="flex items-center gap-3">
            <span className="text-lg">👤</span>
            <span className="text-ink">我的主页</span>
          </div>
          <svg className="w-5 h-5 text-ink-faint" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
          </svg>
        </Link>

        <Link
          href="/bookmarks"
          className="flex items-center justify-between px-6 py-4 hover:bg-surface-sunk transition-colors"
        >
          <div className="flex items-center gap-3">
            <Icon name="bookmark" size={18} className="text-ink-faint" />
            <span className="text-ink">我的收藏</span>
          </div>
          <svg className="w-5 h-5 text-ink-faint" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
          </svg>
        </Link>

        <Link
          href="/notifications"
          className="flex items-center justify-between px-6 py-4 hover:bg-surface-sunk transition-colors"
        >
          <div className="flex items-center gap-3">
            <span className="text-lg">🔔</span>
            <span className="text-ink">消息通知</span>
          </div>
          <svg className="w-5 h-5 text-ink-faint" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
          </svg>
        </Link>

        <Link
          href="/checkin"
          className="flex items-center justify-between px-6 py-4 hover:bg-surface-sunk transition-colors"
        >
          <div className="flex items-center gap-3">
            <span className="text-lg">✅</span>
            <span className="text-ink">每日打卡</span>
          </div>
          <svg className="w-5 h-5 text-ink-faint" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
          </svg>
        </Link>
      </div>

      {/* 退出登录 */}
      <div className="mt-6">
        <button
          onClick={handleLogout}
          className="w-full bg-white border border-emotion-line text-emotion-ink py-3 rounded-lg hover:bg-emotion-soft transition-colors font-medium"
        >
          退出登录
        </button>
      </div>
    </div>
  )
}
