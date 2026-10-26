'use client'

import Link from 'next/link'
import Avatar from './Avatar'
import Icon from './Icon'
import { TYPE_LABELS, TYPE_BORDER_COLORS } from '@/lib/constants'
import { timeAgo } from '@/lib/utils'

interface Post {
  id: number
  title: string
  content: string
  type: string
  author_id: number
  author_name?: string
  author_avatar?: string
  like_count: number
  comment_count: number
  view_count: number
  is_top?: boolean
  is_featured?: boolean
  created_at: string
}

function stripHtml(html: string): string {
  return html.replace(/<[^>]+>/g, '').substring(0, 150)
}

export default function PostCard({ post }: { post: Post }) {
  const type = TYPE_LABELS[post.type] || TYPE_LABELS.emotion

  return (
    <div className="card p-5">
      {/* 头部：作者信息 + 标记 + 类型标签 */}
      <div className="flex items-center justify-between mb-3">
        <div className="flex items-center gap-3">
          <Avatar
            userId={post.author_id}
            nickname={post.author_name}
            avatar={post.author_avatar}
            size="sm"
          />
          <div>
            <Link
              href={`/users/${post.author_id}`}
              className="text-sm font-medium text-ink hover:text-accent transition-colors"
            >
              {post.author_name || '匿名用户'}
            </Link>
            <p className="text-xs text-ink-faint">{timeAgo(post.created_at)}</p>
          </div>
        </div>
        <div className="flex items-center gap-1.5">
          {post.is_top && (
            <span className="inline-flex items-center gap-1 text-xs text-emotion-ink">
              <Icon name="pin" size={12} />置顶
            </span>
          )}
          {post.is_featured && (
            <span className="inline-flex items-center gap-1 text-xs text-insight-ink">
              <Icon name="star" size={12} />精华
            </span>
          )}
          <span className={`inline-flex items-center gap-1.5 text-xs px-2.5 py-1 rounded-full border ${TYPE_BORDER_COLORS[post.type]}`}>
            <Icon name={type.icon} size={12} />
            {type.label}
          </span>
        </div>
      </div>

      {/* 内容 */}
      <Link href={`/posts/${post.id}`} className="block group">
        <h3 className="text-[17px] font-medium text-ink leading-snug group-hover:text-accent transition-colors mb-1.5 line-clamp-2">
          {post.title}
        </h3>
        <p className="text-sm text-ink-soft line-clamp-2 mb-4">
          {stripHtml(post.content)}
        </p>
      </Link>

      {/* 底部：互动数据 */}
      <div className="flex items-center gap-4 text-xs text-ink-faint">
        <span className="inline-flex items-center gap-1">
          <Icon name="heart" size={13} />
          {post.like_count}
        </span>
        <span className="inline-flex items-center gap-1">
          <Icon name="comment" size={13} />
          {post.comment_count}
        </span>
        <span className="inline-flex items-center gap-1">
          <Icon name="clock" size={13} />
          {post.view_count}
        </span>
      </div>
    </div>
  )
}
