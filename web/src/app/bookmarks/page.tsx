'use client'

import { useState, useEffect, useCallback } from 'react'
import { useRouter } from 'next/navigation'
import Link from 'next/link'
import { get } from '@/lib/api'
import PostCard from '@/components/PostCard'
import LoadingSpinner from '@/components/LoadingSpinner'

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

export default function BookmarksPage() {
  const router = useRouter()
  const [posts, setPosts] = useState<Post[]>([])
  const [page, setPage] = useState(1)
  const [hasMore, setHasMore] = useState(true)
  const [loading, setLoading] = useState(true)
  const [userId, setUserId] = useState<number | null>(null)

  useEffect(() => {
    const userStr = localStorage.getItem('user')
    const token = localStorage.getItem('token')
    if (!userStr || !token) {
      router.push('/login')
      return
    }
    try {
      const user = JSON.parse(userStr)
      setUserId(user.id)
    } catch {
      router.push('/login')
    }
  }, [router])

  const fetchBookmarks = useCallback(async (pageNum: number) => {
    if (!userId) return
    setLoading(true)
    try {
      const res = await get<{ items: Post[]; total: number; page: number; page_size: number }>(
        `/users/${userId}/bookmarks?page=${pageNum}&size=20`
      )
      if (res.code === 0) {
        const items = res.data.items || []
        if (pageNum === 1) {
          setPosts(items)
        } else {
          setPosts(prev => [...prev, ...items])
        }
        setHasMore(items.length === 20)
      }
    } catch (err) {
      console.error('获取收藏失败:', err)
    } finally {
      setLoading(false)
    }
  }, [userId])

  useEffect(() => {
    if (userId) {
      setPage(1)
      fetchBookmarks(1)
    }
  }, [userId, fetchBookmarks])

  const loadMore = () => {
    const nextPage = page + 1
    setPage(nextPage)
    fetchBookmarks(nextPage)
  }

  return (
    <div className="max-w-4xl mx-auto px-4 py-6">
      <div className="flex items-center gap-3 mb-6">
        <Link href="/profile" className="text-gray-500 hover:text-gray-700">
          <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
          </svg>
        </Link>
        <h1 className="text-xl font-bold text-gray-900">我的收藏</h1>
      </div>

      {loading && posts.length === 0 ? (
        <LoadingSpinner />
      ) : posts.length === 0 ? (
        <div className="bg-white rounded-lg border border-gray-200 p-12 text-center">
          <p className="text-gray-500 mb-4">还没有收藏的帖子</p>
          <Link href="/feed" className="text-blue-600 hover:text-blue-700 text-sm font-medium">
            去看看动态
          </Link>
        </div>
      ) : (
        <>
          <div className="space-y-4">
            {posts.map(post => (
              <PostCard key={post.id} post={post} />
            ))}
          </div>

          {loading && <LoadingSpinner />}

          {hasMore && !loading && (
            <div className="text-center mt-6">
              <button
                onClick={loadMore}
                className="text-blue-600 hover:text-blue-700 text-sm font-medium"
              >
                加载更多
              </button>
            </div>
          )}
        </>
      )}
    </div>
  )
}
