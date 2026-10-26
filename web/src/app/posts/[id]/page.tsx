'use client'

import { useState, useEffect } from 'react'
import { useParams } from 'next/navigation'
import Link from 'next/link'
import { get, post as apiPost } from '@/lib/api'
import { TYPE_LABELS } from '@/lib/constants'
import { timeAgo } from '@/lib/utils'
import Avatar from '@/components/Avatar'
import Icon from '@/components/Icon'
import { isHTML, sanitizeHTML } from '@/lib/sanitize'

interface Author {
  id: number
  nickname: string
  avatar: string
}

interface Comment {
  id: number
  content: string
  created_at: string
  author: Author
  replies?: Comment[]
}

interface PostDetail {
  id: number
  user_id: number
  type: string
  title: string
  content: string
  like_count: number
  comment_count: number
  is_anonymous: boolean
  created_at: string
}

export default function PostDetailPage() {
  const params = useParams()
  const postId = params.id

  const [post, setPost] = useState<PostDetail | null>(null)
  const [comments, setComments] = useState<Comment[]>([])
  const [newComment, setNewComment] = useState('')
  const [loading, setLoading] = useState(true)
  const [liked, setLiked] = useState(false)
  const [bookmarked, setBookmarked] = useState(false)
  const [currentUserId, setCurrentUserId] = useState<number | null>(null)
  const [shareMsg, setShareMsg] = useState('')
  const [showReport, setShowReport] = useState(false)
  const [reportReason, setReportReason] = useState('')
  const [reportLoading, setReportLoading] = useState(false)

  useEffect(() => {
    const stored = localStorage.getItem('user')
    if (stored) {
      try { setCurrentUserId(JSON.parse(stored).id) } catch { /* ignore */ }
    }

    const fetchPost = async () => {
      try {
        const res = await get<PostDetail>(`/posts/${postId}`)
        if (res.code === 0) setPost(res.data)
      } catch (err) {
        console.error('获取帖子详情失败:', err)
      }
    }

    // 初始化当前用户对本帖的点赞/收藏态（仅登录用户；游客默认未点赞未收藏）
    const fetchState = async () => {
      if (typeof window === 'undefined' || !localStorage.getItem('token')) return
      try {
        const res = await get<{ liked: boolean; bookmarked: boolean }>(`/posts/${postId}/interaction`)
        if (res.code === 0) {
          setLiked(res.data.liked)
          setBookmarked(res.data.bookmarked)
        }
      } catch { /* 忽略：状态初始化失败不影响浏览 */ }
    }

    const fetchComments = async () => {
      try {
        const res = await get<{ items: Comment[] }>(`/posts/${postId}/comments`)
        if (res.code === 0) setComments(res.data.items || [])
      } catch (err) {
        console.error('获取评论失败:', err)
      }
    }

    Promise.all([fetchPost(), fetchComments(), fetchState()]).finally(() => setLoading(false))
  }, [postId])

  const handleLike = async () => {
    try {
      // 采用服务器返回的权威状态与计数，不再前端"乐观 ±1"（会在状态不同步时漂移）
      const res = await apiPost<{ liked: boolean; like_count: number }>(`/posts/${postId}/like`)
      if (res.code === 0) {
        setLiked(res.data.liked)
        if (post) setPost({ ...post, like_count: res.data.like_count })
      }
    } catch (err) {
      console.error('点赞失败:', err)
    }
  }

  const handleBookmark = async () => {
    try {
      const res = await apiPost<{ bookmarked: boolean }>(`/posts/${postId}/bookmark`)
      if (res.code === 0) setBookmarked(res.data.bookmarked)
    } catch (err) {
      console.error('收藏失败:', err)
    }
  }

  const handleShare = async () => {
    const url = window.location.href
    try {
      await navigator.clipboard.writeText(url)
      setShareMsg('链接已复制')
      setTimeout(() => setShareMsg(''), 2000)
    } catch {
      setShareMsg('复制失败')
      setTimeout(() => setShareMsg(''), 2000)
    }
  }

  const handleReport = async () => {
    if (!reportReason.trim() || reportLoading) return
    setReportLoading(true)
    try {
      const res = await apiPost('/reports', {
        target_type: 'post',
        target_id: Number(postId),
        reason: reportReason.trim()
      })
      if (res.code === 0) {
        setShowReport(false)
        setReportReason('')
        alert('举报成功')
      } else {
        alert(res.message || '举报失败')
      }
    } catch {
      alert('网络错误')
    } finally {
      setReportLoading(false)
    }
  }

  const handleSubmitComment = async () => {
    if (!newComment.trim()) return
    try {
      const res = await apiPost(`/posts/${postId}/comments`, { content: newComment })
      if (res.code === 0) {
        setNewComment('')
        const commentsRes = await get<{ items: Comment[] }>(`/posts/${postId}/comments`)
        if (commentsRes.code === 0) setComments(commentsRes.data.items || [])
        if (post) setPost({ ...post, comment_count: post.comment_count + 1 })
      }
    } catch (err) {
      console.error('发表评论失败:', err)
    }
  }

  if (loading) {
    return (
      <div className="container mx-auto px-4 py-8 max-w-3xl">
        <div className="card p-8">
          <div className="animate-pulse">
            <div className="h-8 bg-surface-sunk rounded w-3/4 mb-4" />
            <div className="h-4 bg-surface-sunk rounded w-1/4 mb-8" />
            <div className="space-y-3">
              <div className="h-4 bg-surface-sunk rounded w-full" />
              <div className="h-4 bg-surface-sunk rounded w-full" />
              <div className="h-4 bg-surface-sunk rounded w-2/3" />
            </div>
          </div>
        </div>
      </div>
    )
  }

  if (!post) {
    return (
      <div className="container mx-auto px-4 py-20 text-center">
        <div className="text-5xl mb-4">😕</div>
        <h2 className="text-xl font-semibold text-ink-soft mb-2">帖子不存在</h2>
        <Link href="/feed" className="text-accent hover:underline">返回社区</Link>
      </div>
    )
  }

  const typeInfo = TYPE_LABELS[post.type]

  return (
    <div className="container mx-auto px-4 py-8 max-w-3xl">
      {/* 面包屑 */}
      <div className="mb-4 text-sm text-ink-faint">
        <Link href="/feed" className="hover:text-accent">社区</Link>
        <span className="mx-2">›</span>
        <span className="text-ink-soft">{post.title}</span>
      </div>

      {/* 帖子详情 */}
      <article className="card p-8">
        <div className="flex items-center gap-2 mb-4">
          {typeInfo && (
            <span className={`inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full border text-xs font-medium ${typeInfo.color}`}>
              <Icon name={typeInfo.icon} size={12} /> {typeInfo.label}
            </span>
          )}
          <span className="text-sm text-ink-faint">{timeAgo(post.created_at)}</span>
        </div>

        <h1 className="text-2xl md:text-3xl font-bold text-ink mb-6">
          {post.is_anonymous && (
            <span className="inline-flex items-center gap-1 px-2 py-1 bg-achievement-soft text-achievement-ink border border-achievement-line rounded-full text-xs font-medium mr-2">
              <Icon name="mask" size={12} className="mr-1" />匿名
            </span>
          )}
          {post.title}
        </h1>

        <div className="text-ink-soft leading-relaxed">
          {isHTML(post.content) ? (
            <div
              className="prose prose-sm max-w-none"
              dangerouslySetInnerHTML={{ __html: sanitizeHTML(post.content) }}
            />
          ) : (
            <div className="whitespace-pre-wrap">{post.content}</div>
          )}
        </div>

        {/* 互动栏 */}
        <div className="mt-8 pt-6 border-t border-line-soft flex items-center gap-1 flex-wrap">
          <button
            onClick={handleLike}
            className={`inline-flex items-center gap-1.5 px-4 py-2 rounded-lg text-sm font-medium transition-all ${
              liked
                ? 'bg-case-soft text-accent border border-case-line'
                : 'text-ink-faint hover:bg-surface-sunk border border-transparent'
            }`}
          >
            <Icon name="heart" size={15} />{post.like_count}
          </button>
          <button
            onClick={handleBookmark}
            className={`inline-flex items-center gap-1.5 px-4 py-2 rounded-lg text-sm font-medium transition-all ${
              bookmarked
                ? 'bg-insight-soft text-insight-ink border border-insight-line'
                : 'text-ink-faint hover:bg-surface-sunk border border-transparent'
            }`}
          >
            <Icon name="bookmark" size={15} />收藏
          </button>
          <span className="inline-flex items-center gap-1.5 px-4 py-2 text-sm text-ink-faint">
            <Icon name="comment" size={15} />{post.comment_count} 评论
          </span>
          <button
            onClick={handleShare}
            className="inline-flex items-center gap-1.5 px-4 py-2 rounded-lg text-sm font-medium text-ink-faint hover:bg-surface-sunk hover:text-accent transition-all"
          >
            <Icon name="share" size={15} />{shareMsg || '分享'}
          </button>
          {currentUserId && currentUserId !== post.user_id && (
            <button
              onClick={() => setShowReport(true)}
              className="inline-flex items-center gap-1.5 px-4 py-2 rounded-lg text-sm font-medium text-ink-faint hover:bg-surface-sunk hover:text-emotion-ink transition-all"
            >
              ⚑ 举报
            </button>
          )}
          {currentUserId === post.user_id && (
            <Link
              href={`/posts/${post.id}/edit`}
              className="inline-flex items-center gap-1.5 px-4 py-2 rounded-lg text-sm font-medium text-ink-faint hover:bg-surface-sunk hover:text-accent transition-all ml-auto"
            >
              ✏️ 编辑
            </Link>
          )}
        </div>
      </article>

      {/* 评论区 */}
      <div className="mt-8">
        <h2 className="text-lg font-bold text-ink mb-4">评论 ({comments.length})</h2>

        {/* 发表评论 */}
        <div className="card p-4 mb-6">
          <textarea
            value={newComment}
            onChange={(e) => setNewComment(e.target.value)}
            placeholder="写下你的想法..."
            className="w-full p-3 border border-line rounded-lg resize-none focus:outline-none focus:ring-2 focus:ring-accent/20 focus:border-accent transition-all"
            rows={3}
          />
          <div className="mt-3 flex justify-end">
            <button
              onClick={handleSubmitComment}
              className="btn-primary text-sm px-5 py-2"
            >
              发表评论
            </button>
          </div>
        </div>

        {/* 评论列表 */}
        {comments.length === 0 ? (
          <div className="card p-8 text-center text-ink-faint">
            <div className="flex justify-center mb-2 text-ink-faint"><Icon name="comment" size={28} /></div>
            <p>还没有评论，来抢沙发吧</p>
          </div>
        ) : (
          <div className="space-y-3">
            {comments.map((comment) => (
              <div key={comment.id} className="card p-5">
                <div className="flex items-center gap-3 mb-3">
                  <Avatar nickname={comment.author.nickname} avatar={comment.author.avatar} size="sm" />
                  <div>
                    <span className="text-sm font-medium text-ink">{comment.author.nickname}</span>
                    <span className="text-xs text-ink-faint ml-2">{timeAgo(comment.created_at)}</span>
                  </div>
                </div>
                <p className="text-ink-soft text-sm leading-relaxed ml-11">{comment.content}</p>

                {/* 回复 */}
                {comment.replies && comment.replies.length > 0 && (
                  <div className="mt-4 ml-11 pl-4 border-l-2 border-line-soft space-y-3">
                    {comment.replies.map((reply) => (
                      <div key={reply.id}>
                        <div className="flex items-center gap-2 mb-1">
                          <span className="text-sm font-medium text-ink">{reply.author.nickname}</span>
                          <span className="text-xs text-ink-faint">{timeAgo(reply.created_at)}</span>
                        </div>
                        <p className="text-sm text-ink-soft">{reply.content}</p>
                      </div>
                    ))}
                  </div>
                )}
              </div>
            ))}
          </div>
        )}
      </div>

      {/* 举报弹窗 */}
      {showReport && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
          <div className="bg-white rounded-xl p-6 max-w-md w-full mx-4">
            <h3 className="text-lg font-bold text-ink mb-4">举报帖子</h3>
            <textarea
              value={reportReason}
              onChange={(e) => setReportReason(e.target.value)}
              placeholder="请描述举报原因..."
              rows={4}
              className="w-full p-3 border border-line rounded-lg resize-none focus:outline-none focus:ring-2 focus:ring-accent/20 focus:border-accent"
              maxLength={500}
            />
            <p className="text-xs text-ink-faint mt-1">{reportReason.length}/500</p>
            <div className="flex gap-3 mt-4">
              <button
                onClick={() => { setShowReport(false); setReportReason('') }}
                className="btn-secondary flex-1"
              >
                取消
              </button>
              <button
                onClick={handleReport}
                disabled={!reportReason.trim() || reportLoading}
                className="btn-primary flex-1 bg-emotion-ink hover:opacity-90 disabled:opacity-50"
              >
                {reportLoading ? '提交中...' : '提交举报'}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
