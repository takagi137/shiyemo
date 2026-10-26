'use client'

import { useState, useEffect, useCallback } from 'react'
import { useParams } from 'next/navigation'
import Link from 'next/link'
import { get, post, del } from '@/lib/api'
import Avatar from '@/components/Avatar'
import PostCard from '@/components/PostCard'
import LoadingSpinner from '@/components/LoadingSpinner'

interface UserProfile {
  id: number
  nickname: string
  bio: string
  avatar: string
  company: string
  occupation: string
  follower_count: number
  following_count: number
  created_at: string
}

interface UserStats {
  post_count: number
  like_received: number
  follower_count: number
  following_count: number
}

interface Post {
  id: number
  title: string
  content: string
  type: string
  author_id: number
  author_name: string
  author_avatar: string
  like_count: number
  comment_count: number
  view_count: number
  is_top: boolean
  is_featured: boolean
  created_at: string
}

export default function UserProfilePage() {
  const params = useParams()
  const userId = params.id as string

  const [profile, setProfile] = useState<UserProfile | null>(null)
  const [stats, setStats] = useState<UserStats | null>(null)
  const [posts, setPosts] = useState<Post[]>([])
  const [page, setPage] = useState(1)
  const [hasMore, setHasMore] = useState(true)
  const [loading, setLoading] = useState(true)
  const [postsLoading, setPostsLoading] = useState(false)
  const [isFollowing, setIsFollowing] = useState(false)
  const [followLoading, setFollowLoading] = useState(false)
  const [currentUser, setCurrentUser] = useState<{ id: number } | null>(null)

  // 获取当前登录用户
  useEffect(() => {
    const userStr = localStorage.getItem('user')
    if (userStr) {
      try {
        setCurrentUser(JSON.parse(userStr))
      } catch { /* ignore */ }
    }
  }, [])

  // 获取用户信息
  useEffect(() => {
    if (!userId) return

    const fetchProfile = async () => {
      setLoading(true)
      try {
        const [profileRes, statsRes] = await Promise.all([
          get<UserProfile>(`/users/${userId}`),
          get<UserStats>(`/users/${userId}/stats`),
        ])

        if (profileRes.code === 0) setProfile(profileRes.data)
        if (statsRes.code === 0) setStats(statsRes.data)
      } catch (err) {
        console.error('获取用户信息失败:', err)
      } finally {
        setLoading(false)
      }
    }

    fetchProfile()
  }, [userId])

  // 获取用户帖子
  const fetchPosts = useCallback(async (pageNum: number) => {
    if (!userId) return
    setPostsLoading(true)
    try {
      const res = await get<{ items: Post[]; total: number; page: number; page_size: number }>(
        `/users/${userId}/posts?page=${pageNum}&page_size=20`
      )
      if (res.code === 0) {
        if (pageNum === 1) {
          setPosts(res.data.items || [])
        } else {
          setPosts(prev => [...prev, ...(res.data.items || [])])
        }
        setHasMore((res.data.items || []).length === 20)
      }
    } catch (err) {
      console.error('获取帖子失败:', err)
    } finally {
      setPostsLoading(false)
    }
  }, [userId])

  useEffect(() => {
    if (userId) {
      setPage(1)
      fetchPosts(1)
    }
  }, [userId, fetchPosts])

  // 加载更多
  const loadMore = () => {
    const nextPage = page + 1
    setPage(nextPage)
    fetchPosts(nextPage)
  }

  // 关注/取消关注
  const handleFollow = async () => {
    if (!currentUser || followLoading) return
    setFollowLoading(true)
    try {
      if (isFollowing) {
        await del(`/users/${userId}/follow`)
        setIsFollowing(false)
        setStats(prev => prev ? { ...prev, follower_count: prev.follower_count - 1 } : prev)
      } else {
        await post(`/users/${userId}/follow`)
        setIsFollowing(true)
        setStats(prev => prev ? { ...prev, follower_count: prev.follower_count + 1 } : prev)
      }
    } catch (err) {
      console.error('关注操作失败:', err)
    } finally {
      setFollowLoading(false)
    }
  }

  if (loading) {
    return <LoadingSpinner />
  }

  if (!profile) {
    return (
      <div className="max-w-4xl mx-auto px-4 py-12 text-center">
        <p className="text-gray-500">用户不存在</p>
        <Link href="/" className="text-blue-600 hover:underline mt-4 inline-block">返回首页</Link>
      </div>
    )
  }

  const isSelf = currentUser && currentUser.id === Number(userId)

  return (
    <div className="max-w-4xl mx-auto px-4 py-6">
      {/* 用户信息卡片 */}
      <div className="bg-white rounded-lg border border-gray-200 p-6 mb-6">
        <div className="flex items-start gap-6">
          {/* 头像 */}
          <Avatar
            userId={profile.id}
            nickname={profile.nickname}
            avatar={profile.avatar}
            size="lg"
            clickable={false}
          />

          {/* 信息 */}
          <div className="flex-1">
            <div className="flex items-center gap-4 mb-2">
              <h1 className="text-xl font-bold text-gray-900">{profile.nickname}</h1>
              {isSelf ? (
                <Link
                  href="/profile"
                  className="text-sm text-blue-600 hover:text-blue-700 border border-blue-200 px-3 py-1 rounded-full"
                >
                  编辑资料
                </Link>
              ) : currentUser ? (
                <button
                  onClick={handleFollow}
                  disabled={followLoading}
                  className={`text-sm px-4 py-1.5 rounded-full font-medium transition-colors ${
                    isFollowing
                      ? 'text-gray-700 border border-gray-300 hover:border-red-300 hover:text-red-600'
                      : 'bg-blue-600 text-white hover:bg-blue-700'
                  } disabled:opacity-50`}
                >
                  {followLoading ? '...' : isFollowing ? '已关注' : '关注'}
                </button>
              ) : null}
            </div>

            {profile.bio && (
              <p className="text-gray-600 mb-3">{profile.bio}</p>
            )}

            <div className="flex items-center gap-2 text-sm text-gray-500">
              {profile.occupation && <span>{profile.occupation}</span>}
              {profile.occupation && profile.company && <span>·</span>}
              {profile.company && <span>{profile.company}</span>}
            </div>
          </div>
        </div>

        {/* 统计数据 */}
        {stats && (
          <div className="grid grid-cols-4 gap-4 mt-6 pt-6 border-t border-gray-100">
            <div className="text-center">
              <div className="text-xl font-bold text-gray-900">{stats.post_count}</div>
              <div className="text-sm text-gray-500">帖子</div>
            </div>
            <div className="text-center">
              <div className="text-xl font-bold text-gray-900">{stats.like_received}</div>
              <div className="text-sm text-gray-500">获赞</div>
            </div>
            <div className="text-center">
              <div className="text-xl font-bold text-gray-900">{stats.following_count}</div>
              <div className="text-sm text-gray-500">关注</div>
            </div>
            <div className="text-center">
              <div className="text-xl font-bold text-gray-900">{stats.follower_count}</div>
              <div className="text-sm text-gray-500">粉丝</div>
            </div>
          </div>
        )}
      </div>

      {/* 帖子列表 */}
      <div>
        <h2 className="text-lg font-semibold text-gray-900 mb-4">TA 的帖子</h2>

        {posts.length === 0 && !postsLoading ? (
          <div className="bg-white rounded-lg border border-gray-200 p-12 text-center">
            <p className="text-gray-500">暂无帖子</p>
          </div>
        ) : (
          <div className="space-y-4">
            {posts.map(post => (
              <PostCard key={post.id} post={post} />
            ))}
          </div>
        )}

        {postsLoading && <LoadingSpinner />}

        {hasMore && !postsLoading && posts.length > 0 && (
          <div className="text-center mt-6">
            <button
              onClick={loadMore}
              className="text-blue-600 hover:text-blue-700 text-sm font-medium"
            >
              加载更多
            </button>
          </div>
        )}
      </div>
    </div>
  )
}
