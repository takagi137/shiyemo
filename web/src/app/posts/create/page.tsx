'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import Link from 'next/link'
import { post } from '@/lib/api'
import { POST_CREATE_TYPES } from '@/lib/constants'
import Icon from '@/components/Icon'
import dynamic from 'next/dynamic'

const RichTextEditor = dynamic(() => import('@/components/RichTextEditor'), {
  ssr: false,
  loading: () => <div className="h-64 bg-surface-sunk animate-pulse rounded-lg" />,
})

export default function CreatePostPage() {
  const router = useRouter()
  const [type, setType] = useState('')
  const [title, setTitle] = useState('')
  const [content, setContent] = useState('')
  const [tagInput, setTagInput] = useState('')
  const [tags, setTags] = useState<string[]>([])
  const [isAnonymous, setIsAnonymous] = useState(false)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  const handleAddTag = () => {
    const tag = tagInput.trim()
    if (tag && !tags.includes(tag) && tags.length < 5) {
      setTags([...tags, tag])
      setTagInput('')
    }
  }

  const handleRemoveTag = (tag: string) => setTags(tags.filter(t => t !== tag))

  const handleTagKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter') {
      e.preventDefault()
      handleAddTag()
    }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    if (!type) { setError('请选择帖子类型'); return }
    if (!title.trim()) { setError('请输入标题'); return }
    if (!content.trim()) { setError('请输入内容'); return }

    setLoading(true)
    try {
      const res = await post<{ id: number }>('/posts', {
        type, title: title.trim(), content: content.trim(), tags, is_anonymous: isAnonymous,
      })
      if (res.code === 0) {
        router.push(`/posts/${res.data.id}`)
      } else if (res.code === 40301) {
        // 邮箱未验证，跳转到验证页面
        router.push('/verify-email')
      } else {
        setError(res.message || '发布失败')
      }
    } catch {
      setError('网络错误，请稍后重试')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="container mx-auto px-4 py-8 max-w-2xl">
      {/* 面包屑 */}
      <div className="mb-6 text-sm text-ink-faint">
        <Link href="/feed" className="hover:text-accent">社区</Link>
        <span className="mx-2">›</span>
        <span className="text-ink-soft">发布帖子</span>
      </div>

      <div className="card p-8">
        <h1 className="text-2xl font-bold text-ink mb-6">发布帖子</h1>

        <form onSubmit={handleSubmit} className="space-y-6">
          {error && (
            <div className="bg-emotion-soft border border-emotion-line text-emotion-ink px-4 py-3 rounded-lg text-sm">
              {error}
            </div>
          )}

          {/* 帖子类型 */}
          <div>
            <label className="block text-sm font-medium text-ink-soft mb-3">
              选择类型 <span className="text-emotion-ink">*</span>
            </label>
            <div className="grid grid-cols-2 gap-3">
              {POST_CREATE_TYPES.map((t) => (
                <button
                  key={t.key}
                  type="button"
                  onClick={() => setType(t.key)}
                  className={`p-4 rounded-xl border-2 text-left transition-all duration-200 ${
                    type === t.key
                      ? `${t.color} border-current shadow-sm`
                      : 'border-line bg-white hover:border-line hover:shadow-sm'
                  }`}
                >
                  <div className="flex items-center gap-2 mb-1">
                    <Icon name={t.icon} size={18} />
                    <span className="font-medium text-ink">{t.label}</span>
                  </div>
                  <p className="text-sm text-ink-faint ml-8">{t.desc}</p>
                </button>
              ))}
            </div>
          </div>

          {/* 标题 */}
          <div>
            <label htmlFor="title" className="block text-sm font-medium text-ink-soft mb-1">
              标题 <span className="text-emotion-ink">*</span>
            </label>
            <input
              id="title"
              type="text"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              maxLength={200}
              className="input-field"
              placeholder="请输入标题"
            />
            <div className="text-right text-xs text-ink-faint mt-1">{title.length}/200</div>
          </div>

          {/* 内容 */}
          <div>
            <label className="block text-sm font-medium text-ink-soft mb-1">
              内容 <span className="text-emotion-ink">*</span>
            </label>
            <RichTextEditor
              content={content}
              onChange={setContent}
              placeholder="分享你的想法..."
            />
          </div>

          {/* 标签 */}
          <div>
            <label className="block text-sm font-medium text-ink-soft mb-1">
              标签 <span className="text-ink-faint font-normal">（可选，最多5个）</span>
            </label>
            <div className="flex items-center gap-2 mb-2">
              <input
                type="text"
                value={tagInput}
                onChange={(e) => setTagInput(e.target.value)}
                onKeyDown={handleTagKeyDown}
                className="flex-1 input-field"
                placeholder="输入标签后按回车"
                maxLength={20}
              />
              <button
                type="button"
                onClick={handleAddTag}
                disabled={!tagInput.trim() || tags.length >= 5}
                className="btn-secondary text-sm px-4 py-2.5 disabled:opacity-50"
              >
                添加
              </button>
            </div>
            {tags.length > 0 && (
              <div className="flex flex-wrap gap-2">
                {tags.map((tag) => (
                  <span
                    key={tag}
                    className="inline-flex items-center gap-1 px-3 py-1 bg-case-soft text-case-ink border border-case-line rounded-full text-sm"
                  >
                    #{tag}
                    <button type="button" onClick={() => handleRemoveTag(tag)} className="ml-1 text-case-ink/60 hover:text-case-ink">
                      ×
                    </button>
                  </span>
                ))}
              </div>
            )}
          </div>

          {/* 匿名发布 */}
          <div className="flex items-center gap-3">
            <button
              type="button"
              onClick={() => setIsAnonymous(!isAnonymous)}
              className={`relative inline-flex h-6 w-11 items-center rounded-full transition-colors ${
                isAnonymous ? 'bg-accent' : 'bg-surface-sunk'
              }`}
            >
              <span
                className={`inline-block h-4 w-4 transform rounded-full bg-white transition-transform ${
                  isAnonymous ? 'translate-x-6' : 'translate-x-1'
                }`}
              />
            </button>
            <span className="text-sm text-ink-soft">
              {isAnonymous ? '匿名发布（其他用户看不到你的真实信息）' : '公开发布'}
            </span>
          </div>

          {/* 提交 */}
          <div className="flex items-center justify-between pt-4 border-t border-line-soft">
            <Link href="/feed" className="text-sm text-ink-faint hover:text-ink-soft">
              ← 返回社区
            </Link>
            <button
              type="submit"
              disabled={loading || !type || !title.trim() || !content.trim()}
              className="btn-primary px-8 py-2.5"
            >
              {loading ? '发布中...' : '发布'}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}
