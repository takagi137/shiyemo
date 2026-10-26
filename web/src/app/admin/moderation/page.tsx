'use client'

import { useState, useEffect } from 'react'
import { useRouter } from 'next/navigation'
import { get, post, del } from '@/lib/api'
import { toastSuccess, toastError } from '@/lib/toast'
import Link from 'next/link'

interface ModerationWord {
  id: number
  word: string
  created_at: string
}

export default function ModerationPage() {
  const router = useRouter()
  const [words, setWords] = useState<ModerationWord[]>([])
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [loading, setLoading] = useState(false)
  const [newWord, setNewWord] = useState('')
  const [adding, setAdding] = useState(false)

  useEffect(() => {
    if (!localStorage.getItem('token')) router.push('/login')
  }, [router])

  const fetchWords = async (p: number) => {
    setLoading(true)
    try {
      const res = await get<{ items: ModerationWord[]; total: number }>(`/admin/moderation/words?page=${p}&page_size=20`)
      if (res.code === 0) { setWords(res.data.items || []); setTotal(res.data.total) }
    } catch { toastError('获取敏感词列表失败') }
    finally { setLoading(false) }
  }

  useEffect(() => { fetchWords(page) }, [page])

  const handleAddWord = async () => {
    const word = newWord.trim()
    if (!word) { toastError('请输入敏感词'); return }
    setAdding(true)
    try {
      const res = await post<ModerationWord>('/admin/moderation/words', { word })
      if (res.code === 0) { toastSuccess('添加成功'); setNewWord(''); fetchWords(page) }
      else { toastError(res.message || '添加失败') }
    } catch { toastError('添加失败') }
    finally { setAdding(false) }
  }

  const handleDeleteWord = async (id: number, word: string) => {
    if (!confirm(`确定删除敏感词"${word}"吗？`)) return
    try {
      const res = await del(`/admin/moderation/words/${id}`)
      if (res.code === 0) { toastSuccess('删除成功'); fetchWords(page) }
      else { toastError(res.message || '删除失败') }
    } catch { toastError('删除失败') }
  }

  const totalPages = Math.ceil(total / 20)

  return (
    <div className="container mx-auto px-4 py-8 max-w-4xl">
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">敏感词管理</h1>
          <p className="text-sm text-gray-500 mt-1">管理帖子内容审核的敏感词库</p>
        </div>
        <Link href="/admin" className="btn-secondary text-sm">← 返回管理后台</Link>
      </div>

      {/* 添加敏感词 */}
      <div className="card p-6 mb-6">
        <h2 className="text-base font-semibold text-gray-900 mb-4">添加敏感词</h2>
        <div className="flex gap-3">
          <input
            type="text"
            value={newWord}
            onChange={e => setNewWord(e.target.value)}
            onKeyDown={e => e.key === 'Enter' && handleAddWord()}
            placeholder="输入敏感词，按回车添加"
            className="flex-1 input-field"
            maxLength={50}
          />
          <button onClick={handleAddWord} disabled={adding || !newWord.trim()} className="btn-primary disabled:opacity-50">
            {adding ? '添加中...' : '添加'}
          </button>
        </div>
        <p className="text-xs text-gray-400 mt-2">发布帖子时会自动检测并拦截包含敏感词的内容</p>
      </div>

      {/* 敏感词列表 */}
      <div className="card overflow-hidden">
        <div className="px-6 py-4 border-b border-gray-100 flex items-center justify-between">
          <h2 className="text-base font-semibold text-gray-900">敏感词库</h2>
          <span className="text-sm text-gray-500">共 {total} 个</span>
        </div>

        {loading ? (
          <div className="p-8 text-center text-gray-500">加载中...</div>
        ) : words.length === 0 ? (
          <div className="p-12 text-center">
            <div className="text-4xl mb-2">🛡️</div>
            <p className="text-gray-500">暂无敏感词</p>
          </div>
        ) : (
          <table className="w-full">
            <thead className="bg-gray-50 border-b border-gray-100">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">ID</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">敏感词</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">添加时间</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">操作</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-100">
              {words.map(w => (
                <tr key={w.id} className="hover:bg-gray-50 transition-colors">
                  <td className="px-6 py-3 text-sm text-gray-400">{w.id}</td>
                  <td className="px-6 py-3 text-sm font-medium text-gray-900">{w.word}</td>
                  <td className="px-6 py-3 text-sm text-gray-400">{new Date(w.created_at).toLocaleString()}</td>
                  <td className="px-6 py-3">
                    <button onClick={() => handleDeleteWord(w.id, w.word)} className="text-sm text-red-600 hover:text-red-800 font-medium">删除</button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}

        {totalPages > 1 && (
          <div className="px-6 py-4 border-t border-gray-100 flex items-center justify-center gap-2">
            <button onClick={() => setPage(p => Math.max(1, p - 1))} disabled={page === 1} className="btn-secondary text-sm px-3 py-1.5 disabled:opacity-40">←</button>
            <span className="text-sm text-gray-500">{page} / {totalPages}</span>
            <button onClick={() => setPage(p => Math.min(totalPages, p + 1))} disabled={page === totalPages} className="btn-secondary text-sm px-3 py-1.5 disabled:opacity-40">→</button>
          </div>
        )}
      </div>

      {/* 说明 */}
      <div className="mt-6 card p-6 bg-blue-50 border-blue-200">
        <h3 className="font-semibold text-blue-900 mb-2 flex items-center gap-1.5">💡 使用说明</h3>
        <ul className="text-sm text-blue-800 space-y-1.5 ml-5 list-disc">
          <li>敏感词库用于自动审核用户发布的帖子内容</li>
          <li>包含敏感词的帖子会被自动拦截并提示修改</li>
          <li>内置敏感词库覆盖政治、色情、暴力、赌博、毒品等类别</li>
        </ul>
      </div>
    </div>
  )
}
