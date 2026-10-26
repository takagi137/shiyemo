'use client'

import { useState, useEffect } from 'react'
import { useRouter } from 'next/navigation'
import { get, post, put, del } from '@/lib/api'
import { toastSuccess, toastError } from '@/lib/toast'
import Link from 'next/link'

interface DailyImage {
  id: number
  url: string
  description: string
  credit: string
  sort_order: number
  is_active: boolean
  created_at: string
}

export default function AdminImagesPage() {
  const router = useRouter()
  const [images, setImages] = useState<DailyImage[]>([])
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [loading, setLoading] = useState(false)

  // 新增表单
  const [showForm, setShowForm] = useState(false)
  const [form, setForm] = useState({ url: '', description: '', credit: '', sort_order: 0 })
  const [submitting, setSubmitting] = useState(false)

  // 编辑表单
  const [editing, setEditing] = useState<DailyImage | null>(null)

  useEffect(() => {
    if (!localStorage.getItem('token')) router.push('/login')
  }, [router])

  const fetchImages = async (p: number) => {
    setLoading(true)
    try {
      const res = await get<{ items: DailyImage[]; total: number }>(`/admin/images?page=${p}&page_size=50`)
      if (res.code === 0) { setImages(res.data.items || []); setTotal(res.data.total) }
    } catch { toastError('获取图片列表失败') }
    finally { setLoading(false) }
  }

  useEffect(() => { fetchImages(page) }, [page])

  const handleCreate = async () => {
    if (!form.url.trim()) { toastError('请输入图片URL'); return }
    setSubmitting(true)
    try {
      const res = await post<DailyImage>('/admin/images', form)
      if (res.code === 0) {
        toastSuccess('添加成功')
        setForm({ url: '', description: '', credit: '', sort_order: 0 })
        setShowForm(false)
        fetchImages(page)
      } else { toastError(res.message || '添加失败') }
    } catch { toastError('添加失败') }
    finally { setSubmitting(false) }
  }

  const handleUpdate = async () => {
    if (!editing) return
    setSubmitting(true)
    try {
      const res = await put(`/admin/images/${editing.id}`, {
        url: editing.url,
        description: editing.description,
        credit: editing.credit,
        sort_order: editing.sort_order,
        is_active: editing.is_active,
      })
      if (res.code === 0) {
        toastSuccess('更新成功')
        setEditing(null)
        fetchImages(page)
      } else { toastError(res.message || '更新失败') }
    } catch { toastError('更新失败') }
    finally { setSubmitting(false) }
  }

  const handleDelete = async (id: number, desc: string) => {
    if (!confirm(`确定删除图片"${desc || id}"吗？`)) return
    try {
      const res = await del(`/admin/images/${id}`)
      if (res.code === 0) { toastSuccess('删除成功'); fetchImages(page) }
      else { toastError(res.message || '删除失败') }
    } catch { toastError('删除失败') }
  }

  const handleToggleActive = async (img: DailyImage) => {
    try {
      const res = await put(`/admin/images/${img.id}`, {
        url: img.url,
        description: img.description,
        credit: img.credit,
        sort_order: img.sort_order,
        is_active: !img.is_active,
      })
      if (res.code === 0) { toastSuccess(img.is_active ? '已隐藏' : '已启用'); fetchImages(page) }
    } catch { toastError('操作失败') }
  }

  const handleMoveUp = async (index: number) => {
    if (index === 0) return
    const newImages = [...images]
    ;[newImages[index - 1], newImages[index]] = [newImages[index], newImages[index - 1]]
    const ids = newImages.map(img => img.id)
    try {
      const res = await put('/admin/images/reorder', { ids })
      if (res.code === 0) { setImages(newImages) }
    } catch { toastError('排序失败') }
  }

  const handleMoveDown = async (index: number) => {
    if (index === images.length - 1) return
    const newImages = [...images]
    ;[newImages[index], newImages[index + 1]] = [newImages[index + 1], newImages[index]]
    const ids = newImages.map(img => img.id)
    try {
      const res = await put('/admin/images/reorder', { ids })
      if (res.code === 0) { setImages(newImages) }
    } catch { toastError('排序失败') }
  }

  return (
    <div className="container mx-auto px-4 py-8 max-w-5xl">
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">风景图管理</h1>
          <p className="text-sm text-gray-500 mt-1">管理首页每日风景图（共 {total} 张）</p>
        </div>
        <div className="flex gap-3">
          <button onClick={() => setShowForm(!showForm)} className="btn-primary text-sm">
            {showForm ? '取消' : '+ 添加图片'}
          </button>
          <Link href="/admin" className="btn-secondary text-sm">← 返回管理后台</Link>
        </div>
      </div>

      {/* 添加图片表单 */}
      {showForm && (
        <div className="card p-6 mb-6">
          <h2 className="text-base font-semibold text-gray-900 mb-4">添加新图片</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">图片URL *</label>
              <input
                type="url"
                value={form.url}
                onChange={e => setForm({ ...form, url: e.target.value })}
                placeholder="https://images.unsplash.com/..."
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-sm"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">描述</label>
              <input
                type="text"
                value={form.description}
                onChange={e => setForm({ ...form, description: e.target.value })}
                placeholder="如：山谷湖泊"
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-sm"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">版权信息</label>
              <input
                type="text"
                value={form.credit}
                onChange={e => setForm({ ...form, credit: e.target.value })}
                placeholder="Photo by xxx on Unsplash"
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-sm"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">排序（越小越前）</label>
              <input
                type="number"
                value={form.sort_order}
                onChange={e => setForm({ ...form, sort_order: parseInt(e.target.value) || 0 })}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-sm"
              />
            </div>
          </div>
          {/* URL 预览 */}
          {form.url && (
            <div className="mt-4">
              <img src={form.url} alt="预览" className="w-full h-48 object-cover rounded-lg" />
            </div>
          )}
          <div className="mt-4 flex justify-end">
            <button onClick={handleCreate} disabled={submitting} className="btn-primary text-sm">
              {submitting ? '添加中...' : '确认添加'}
            </button>
          </div>
        </div>
      )}

      {/* 图片列表 */}
      {loading ? (
        <div className="text-center py-12 text-gray-500">加载中...</div>
      ) : images.length === 0 ? (
        <div className="text-center py-12 text-gray-500">暂无图片，点击上方按钮添加</div>
      ) : (
        <div className="space-y-3">
          {images.map((img, index) => (
            <div key={img.id} className={`card p-4 flex items-center gap-4 ${!img.is_active ? 'opacity-50' : ''}`}>
              {/* 排序按钮 */}
              <div className="flex flex-col gap-1">
                <button
                  onClick={() => handleMoveUp(index)}
                  disabled={index === 0}
                  className="text-gray-400 hover:text-gray-600 disabled:opacity-30 text-xs"
                >▲</button>
                <button
                  onClick={() => handleMoveDown(index)}
                  disabled={index === images.length - 1}
                  className="text-gray-400 hover:text-gray-600 disabled:opacity-30 text-xs"
                >▼</button>
              </div>

              {/* 缩略图 */}
              <img src={img.url} alt={img.description} className="w-24 h-16 object-cover rounded-lg flex-shrink-0" />

              {/* 信息 */}
              <div className="flex-1 min-w-0">
                <div className="text-sm font-medium text-gray-900 truncate">{img.description || '无描述'}</div>
                <div className="text-xs text-gray-500 truncate mt-1">{img.credit}</div>
                <div className="text-xs text-gray-400 truncate mt-1">{img.url}</div>
              </div>

              {/* 状态 */}
              <span className={`text-xs px-2 py-1 rounded-full ${img.is_active ? 'bg-green-100 text-green-700' : 'bg-gray-100 text-gray-500'}`}>
                {img.is_active ? '启用' : '隐藏'}
              </span>

              {/* 操作 */}
              <div className="flex gap-2">
                <button
                  onClick={() => setEditing(img)}
                  className="text-blue-600 hover:text-blue-800 text-sm"
                >编辑</button>
                <button
                  onClick={() => handleToggleActive(img)}
                  className="text-yellow-600 hover:text-yellow-800 text-sm"
                >{img.is_active ? '隐藏' : '启用'}</button>
                <button
                  onClick={() => handleDelete(img.id, img.description)}
                  className="text-red-600 hover:text-red-800 text-sm"
                >删除</button>
              </div>
            </div>
          ))}
        </div>
      )}

      {/* 编辑弹窗 */}
      {editing && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
          <div className="bg-white rounded-xl p-6 w-full max-w-lg mx-4">
            <h2 className="text-lg font-bold mb-4">编辑图片</h2>
            <div className="space-y-3">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">图片URL</label>
                <input
                  type="url"
                  value={editing.url}
                  onChange={e => setEditing({ ...editing, url: e.target.value })}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">描述</label>
                <input
                  type="text"
                  value={editing.description}
                  onChange={e => setEditing({ ...editing, description: e.target.value })}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">版权信息</label>
                <input
                  type="text"
                  value={editing.credit}
                  onChange={e => setEditing({ ...editing, credit: e.target.value })}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm"
                />
              </div>
              <div className="grid grid-cols-2 gap-3">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">排序</label>
                  <input
                    type="number"
                    value={editing.sort_order}
                    onChange={e => setEditing({ ...editing, sort_order: parseInt(e.target.value) || 0 })}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm"
                  />
                </div>
                <div className="flex items-end">
                  <label className="flex items-center gap-2 cursor-pointer">
                    <input
                      type="checkbox"
                      checked={editing.is_active}
                      onChange={e => setEditing({ ...editing, is_active: e.target.checked })}
                      className="w-4 h-4 text-blue-600 rounded"
                    />
                    <span className="text-sm text-gray-700">启用</span>
                  </label>
                </div>
              </div>
              {/* 预览 */}
              {editing.url && (
                <img src={editing.url} alt="预览" className="w-full h-40 object-cover rounded-lg" />
              )}
            </div>
            <div className="flex justify-end gap-3 mt-6">
              <button onClick={() => setEditing(null)} className="btn-secondary text-sm">取消</button>
              <button onClick={handleUpdate} disabled={submitting} className="btn-primary text-sm">
                {submitting ? '保存中...' : '保存'}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
