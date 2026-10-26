'use client'

import { useState, useEffect } from 'react'
import { useParams, useRouter } from 'next/navigation'
import { get, put } from '@/lib/api'
import { POST_CREATE_TYPES } from '@/lib/constants'
import Icon from '@/components/Icon'
import dynamic from 'next/dynamic'

const RichTextEditor = dynamic(() => import('@/components/RichTextEditor'), {
  ssr: false,
  loading: () => <div className="h-64 bg-gray-100 animate-pulse rounded-lg" />,
})

export default function EditPostPage() {
  const params = useParams()
  const router = useRouter()
  const postId = params.id

  const [type, setType] = useState('')
  const [title, setTitle] = useState('')
  const [content, setContent] = useState('')
  const [tags, setTags] = useState<string[]>([])
  const [tagInput, setTagInput] = useState('')
  const [loading, setLoading] = useState(false)
  const [fetchLoading, setFetchLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    const fetchPost = async () => {
      try {
        const res = await get<{
          id: number
          user_id: number
          type: string
          title: string
          content: string
          tags: string[]
        }>(`/posts/${postId}`)

        if (res.code === 0) {
          setType(res.data.type || '')
          setTitle(res.data.title || '')
          setContent(res.data.content || '')
          setTags(res.data.tags || [])
        } else {
          setError('帖子不存在')
        }
      } catch {
        setError('获取帖子失败')
      } finally {
        setFetchLoading(false)
      }
    }

    fetchPost()
  }, [postId])

  const handleAddTag = () => {
    const tag = tagInput.trim()
    if (tag && !tags.includes(tag) && tags.length < 5) {
      setTags([...tags, tag])
      setTagInput('')
    }
  }

  const handleRemoveTag = (tag: string) => {
    setTags(tags.filter(t => t !== tag))
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!type || !title.trim() || !content.trim()) {
      setError('请填写完整信息')
      return
    }

    setLoading(true)
    setError('')

    try {
      const res = await put(`/posts/${postId}`, {
        type,
        title: title.trim(),
        content: content.trim(),
        tags,
      })

      if (res.code === 0) {
        router.push(`/posts/${postId}`)
      } else {
        setError(res.message || '更新失败')
      }
    } catch {
      setError('网络错误，请稍后重试')
    } finally {
      setLoading(false)
    }
  }

  if (fetchLoading) {
    return (
      <div className="container mx-auto px-4 py-8 max-w-2xl">
        <div className="card p-8 animate-pulse">
          <div className="h-8 bg-gray-200 rounded w-1/3 mb-6" />
          <div className="h-12 bg-gray-200 rounded w-full mb-4" />
          <div className="h-48 bg-gray-200 rounded w-full" />
        </div>
      </div>
    )
  }

  if (error && !title) {
    return (
      <div className="container mx-auto px-4 py-20 text-center">
        <div className="text-5xl mb-4">😕</div>
        <h2 className="text-xl font-semibold text-gray-700 mb-2">{error}</h2>
        <button onClick={() => router.back()} className="btn-primary mt-4">返回</button>
      </div>
    )
  }

  return (
    <div className="container mx-auto px-4 py-8 max-w-2xl">
      {/* 面包屑 */}
      <div className="mb-6 text-sm text-gray-500">
        <button onClick={() => router.back()} className="hover:text-blue-600">返回</button>
        <span className="mx-2">›</span>
        <span className="text-gray-700">编辑帖子</span>
      </div>

      <form onSubmit={handleSubmit} className="card p-8">
        <h1 className="text-2xl font-bold text-gray-900 mb-6">编辑帖子</h1>

        {error && (
          <div className="mb-6 p-4 bg-red-50 border border-red-200 rounded-lg text-red-600 text-sm">
            {error}
          </div>
        )}

        {/* 类型选择 */}
        <div className="mb-6">
          <label className="block text-sm font-medium text-gray-700 mb-3">帖子类型</label>
          <div className="grid grid-cols-2 gap-3">
            {POST_CREATE_TYPES.map((t) => (
              <button
                key={t.key}
                type="button"
                onClick={() => setType(t.key)}
                className={`p-4 rounded-xl border-2 text-left transition-all ${
                  type === t.key
                    ? 'border-blue-500 bg-blue-50 shadow-md'
                    : 'border-gray-200 hover:border-gray-300 hover:bg-gray-50'
                }`}
              >
                <div className="flex items-center gap-2 mb-1">
                  <Icon name={t.icon} size={17} />
                  <span className="font-medium text-gray-900">{t.label}</span>
                </div>
                <p className="text-xs text-gray-500">{t.desc}</p>
              </button>
            ))}
          </div>
        </div>

        {/* 标题 */}
        <div className="mb-6">
          <label className="block text-sm font-medium text-gray-700 mb-2">标题</label>
          <input
            type="text"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            placeholder="一句话描述你想分享的内容"
            maxLength={200}
            className="input-field"
          />
        </div>

        {/* 内容 */}
        <div className="mb-6">
          <label className="block text-sm font-medium text-gray-700 mb-2">内容</label>
          <RichTextEditor
            content={content}
            onChange={setContent}
            placeholder="分享你的故事、想法或经验..."
          />
        </div>

        {/* 标签 */}
        <div className="mb-8">
          <label className="block text-sm font-medium text-gray-700 mb-2">标签（可选，最多5个）</label>
          <div className="flex flex-wrap gap-2 mb-3">
            {tags.map((tag) => (
              <span
                key={tag}
                className="inline-flex items-center gap-1 px-3 py-1 bg-blue-50 text-blue-700 border border-blue-200 rounded-full text-sm"
              >
                #{tag}
                <button
                  type="button"
                  onClick={() => handleRemoveTag(tag)}
                  className="ml-1 text-blue-400 hover:text-blue-600"
                >
                  ×
                </button>
              </span>
            ))}
          </div>
          {tags.length < 5 && (
            <div className="flex gap-2">
              <input
                type="text"
                value={tagInput}
                onChange={(e) => setTagInput(e.target.value)}
                placeholder="输入标签名称"
                maxLength={20}
                className="flex-1 input-field"
                onKeyDown={(e) => {
                  if (e.key === 'Enter') {
                    e.preventDefault()
                    handleAddTag()
                  }
                }}
              />
              <button
                type="button"
                onClick={handleAddTag}
                className="btn-secondary text-sm"
              >
                添加
              </button>
            </div>
          )}
        </div>

        {/* 提交按钮 */}
        <div className="flex gap-3">
          <button
            type="button"
            onClick={() => router.back()}
            className="btn-secondary flex-1"
          >
            取消
          </button>
          <button
            type="submit"
            disabled={loading}
            className="btn-primary flex-1"
          >
            {loading ? '保存中...' : '保存修改'}
          </button>
        </div>
      </form>
    </div>
  )
}
