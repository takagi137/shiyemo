'use client'

import { useState, useEffect } from 'react'
import { useRouter } from 'next/navigation'
import Link from 'next/link'
import { get, post } from '@/lib/api'
import { MOOD_OPTIONS } from '@/lib/constants'

export default function CheckinPage() {
  const router = useRouter()
  const [mood, setMood] = useState('')
  const [content, setContent] = useState('')
  const [images, setImages] = useState<string[]>([])
  const [isPublic, setIsPublic] = useState(true)
  const [loading, setLoading] = useState(false)
  const [uploading, setUploading] = useState(false)
  const [error, setError] = useState('')
  const [todayCheckin, setTodayCheckin] = useState<any>(null)
  const [checkingToday, setCheckingToday] = useState(true)

  // 检查今天是否已打卡
  useEffect(() => {
    const checkToday = async () => {
      try {
        const res = await get<any>('/checkin/today')
        if (res.code === 0 && res.data) {
          setTodayCheckin(res.data)
        }
      } catch (err) {
        console.error('检查今日打卡状态失败:', err)
      } finally {
        setCheckingToday(false)
      }
    }
    checkToday()
  }, [])

  // 处理图片上传
  const handleImageUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const files = e.target.files
    if (!files || files.length === 0) return

    // 检查图片数量限制
    if (images.length + files.length > 3) {
      setError('最多只能上传3张图片')
      return
    }

    setUploading(true)
    setError('')

    try {
      const uploadedUrls: string[] = []
      
      for (let i = 0; i < files.length; i++) {
        const file = files[i]
        const formData = new FormData()
        formData.append('file', file)

        const res = await fetch('/api/v1/upload/image', {
          method: 'POST',
          headers: {
            'Authorization': `Bearer ${localStorage.getItem('token')}`,
          },
          body: formData,
        })

        const data = await res.json()
        if (data.code === 0) {
          uploadedUrls.push(data.data.url)
        } else {
          setError(data.message || '图片上传失败')
          break
        }
      }

      if (uploadedUrls.length > 0) {
        setImages([...images, ...uploadedUrls])
      }
    } catch (err) {
      setError('图片上传失败，请稍后重试')
    } finally {
      setUploading(false)
      // 清空 input
      e.target.value = ''
    }
  }

  // 删除图片
  const handleRemoveImage = (index: number) => {
    setImages(images.filter((_, i) => i !== index))
  }

  // 提交打卡
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')

    if (!mood) {
      setError('请选择今天的心情')
      return
    }

    setLoading(true)
    try {
      const res = await post<any>('/checkin', {
        mood,
        content: content.trim(),
        images,
        is_public: isPublic,
      })

      if (res.code === 0) {
        setTodayCheckin(res.data)
        // 成功后刷新页面状态
      } else {
        setError(res.message || '打卡失败')
      }
    } catch (err) {
      setError('网络错误，请稍后重试')
    } finally {
      setLoading(false)
    }
  }

  // 加载中状态
  if (checkingToday) {
    return (
      <div className="container mx-auto px-4 py-8 max-w-2xl">
        <div className="card p-8 text-center">
          <div className="animate-pulse">
            <div className="h-8 bg-gray-200 rounded w-48 mx-auto mb-4"></div>
            <div className="h-4 bg-gray-200 rounded w-64 mx-auto"></div>
          </div>
        </div>
      </div>
    )
  }

  // 今天已打卡状态
  if (todayCheckin) {
    const moodOption = MOOD_OPTIONS.find(m => m.key === todayCheckin.mood)
    return (
      <div className="container mx-auto px-4 py-8 max-w-2xl">
        <div className="mb-6 text-sm text-gray-500">
          <Link href="/" className="hover:text-blue-600">首页</Link>
          <span className="mx-2">›</span>
          <span className="text-gray-700">今日打卡</span>
        </div>

        <div className="card p-8 text-center">
          <div className="text-6xl mb-4">{moodOption?.emoji || '😊'}</div>
          <h1 className="text-2xl font-bold text-gray-900 mb-2">今日已打卡</h1>
          <p className="text-gray-500 mb-6">
            你今天的心情是：<span className="font-medium">{moodOption?.label || '未知'}</span>
          </p>
          
          {todayCheckin.content && (
            <div className="bg-gray-50 rounded-lg p-4 mb-6 text-left">
              <p className="text-gray-700">{todayCheckin.content}</p>
            </div>
          )}

          {todayCheckin.images && JSON.parse(todayCheckin.images || '[]').length > 0 && (
            <div className="flex gap-2 justify-center mb-6">
              {JSON.parse(todayCheckin.images).map((url: string, index: number) => (
                <img
                  key={index}
                  src={url}
                  alt={`打卡图片 ${index + 1}`}
                  className="w-20 h-20 object-cover rounded-lg"
                />
              ))}
            </div>
          )}

          <div className="flex gap-4 justify-center">
            <Link href="/checkin/history" className="btn-secondary">
              查看历史
            </Link>
            <Link href="/feed" className="btn-primary">
              返回社区
            </Link>
          </div>
        </div>
      </div>
    )
  }

  // 打卡表单
  return (
    <div className="container mx-auto px-4 py-8 max-w-2xl">
      {/* 面包屑 */}
      <div className="mb-6 text-sm text-gray-500">
        <Link href="/" className="hover:text-blue-600">首页</Link>
        <span className="mx-2">›</span>
        <span className="text-gray-700">今日打卡</span>
      </div>

      <div className="card p-8">
        <h1 className="text-2xl font-bold text-gray-900 mb-6">记录今天的心情</h1>

        <form onSubmit={handleSubmit} className="space-y-6">
          {error && (
            <div className="bg-red-50 border border-red-200 text-red-600 px-4 py-3 rounded-lg text-sm">
              {error}
            </div>
          )}

          {/* 心情选择 */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-3">
              今天的心情 <span className="text-red-500">*</span>
            </label>
            <div className="grid grid-cols-5 gap-3">
              {MOOD_OPTIONS.map((option) => (
                <button
                  key={option.key}
                  type="button"
                  onClick={() => setMood(option.key)}
                  className={`p-4 rounded-xl border-2 text-center transition-all duration-200 ${
                    mood === option.key
                      ? `${option.color} border-current shadow-sm`
                      : 'border-gray-200 bg-white hover:border-gray-300 hover:shadow-sm'
                  }`}
                >
                  <div className="text-3xl mb-2">{option.emoji}</div>
                  <div className="text-sm font-medium text-gray-900">{option.label}</div>
                </button>
              ))}
            </div>
          </div>

          {/* 文字内容 */}
          <div>
            <label htmlFor="content" className="block text-sm font-medium text-gray-700 mb-1">
              想说点什么 <span className="text-gray-400 font-normal">（可选）</span>
            </label>
            <textarea
              id="content"
              value={content}
              onChange={(e) => setContent(e.target.value)}
              maxLength={500}
              rows={4}
              className="input-field resize-none"
              placeholder="记录今天的感受..."
            />
            <div className="text-right text-xs text-gray-400 mt-1">{content.length}/500</div>
          </div>

          {/* 图片上传 */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              添加图片 <span className="text-gray-400 font-normal">（可选，最多3张）</span>
            </label>
            
            <div className="flex flex-wrap gap-3 mb-3">
              {images.map((url, index) => (
                <div key={index} className="relative group">
                  <img
                    src={url}
                    alt={`上传图片 ${index + 1}`}
                    className="w-24 h-24 object-cover rounded-lg"
                  />
                  <button
                    type="button"
                    onClick={() => handleRemoveImage(index)}
                    className="absolute -top-2 -right-2 w-6 h-6 bg-red-500 text-white rounded-full flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity"
                  >
                    ×
                  </button>
                </div>
              ))}
              
              {images.length < 3 && (
                <label className="w-24 h-24 border-2 border-dashed border-gray-300 rounded-lg flex flex-col items-center justify-center cursor-pointer hover:border-blue-400 hover:bg-blue-50 transition-colors">
                  <input
                    type="file"
                    accept="image/*"
                    multiple
                    onChange={handleImageUpload}
                    className="hidden"
                    disabled={uploading}
                  />
                  {uploading ? (
                    <div className="animate-spin rounded-full h-6 w-6 border-b-2 border-blue-500"></div>
                  ) : (
                    <>
                      <svg className="w-8 h-8 text-gray-400 mb-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
                      </svg>
                      <span className="text-xs text-gray-500">添加</span>
                    </>
                  )}
                </label>
              )}
            </div>
          </div>

          {/* 是否公开 */}
          <div className="flex items-center gap-3">
            <button
              type="button"
              onClick={() => setIsPublic(!isPublic)}
              className={`relative inline-flex h-6 w-11 items-center rounded-full transition-colors ${
                isPublic ? 'bg-blue-600' : 'bg-gray-200'
              }`}
            >
              <span
                className={`inline-block h-4 w-4 transform rounded-full bg-white transition-transform ${
                  isPublic ? 'translate-x-6' : 'translate-x-1'
                }`}
              />
            </button>
            <span className="text-sm text-gray-700">
              {isPublic ? '公开打卡（在社区时间线展示）' : '私密打卡（仅自己可见）'}
            </span>
          </div>

          {/* 提交 */}
          <div className="flex items-center justify-between pt-4 border-t border-gray-100">
            <Link href="/" className="text-sm text-gray-500 hover:text-gray-700">
              ← 返回首页
            </Link>
            <button
              type="submit"
              disabled={loading || !mood}
              className="btn-primary px-8 py-2.5"
            >
              {loading ? '打卡中...' : '完成打卡'}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}
