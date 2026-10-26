'use client'

import { useState, useEffect } from 'react'

interface Journal {
  id: number
  content: string
  created_at: string
}

export default function JournalEntry() {
  const [journals, setJournals] = useState<Journal[]>([])
  const [content, setContent] = useState('')
  const [loading, setLoading] = useState(false)
  const [saving, setSaving] = useState(false)
  const [showHistory, setShowHistory] = useState(false)
  const [loggedIn, setLoggedIn] = useState(false)

  useEffect(() => {
    const token = typeof window !== 'undefined' ? localStorage.getItem('token') : null
    setLoggedIn(!!token)
    if (token) fetchJournals()
  }, [])

  const fetchJournals = async () => {
    setLoading(true)
    try {
      const token = localStorage.getItem('token')
      if (!token) return

      const res = await fetch('/api/v1/journal?page_size=20', {
        headers: { Authorization: `Bearer ${token}` },
      })
      const json = await res.json()
      if (json.code === 0 && json.data?.items) {
        setJournals(json.data.items)
      }
    } catch (err) {
      console.error('Failed to fetch journals:', err)
    } finally {
      setLoading(false)
    }
  }

  const saveJournal = async () => {
    if (!content.trim()) return

    setSaving(true)
    try {
      const token = localStorage.getItem('token')
      if (!token) return

      const res = await fetch('/api/v1/journal', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({ content: content.trim() }),
      })
      const json = await res.json()
      if (json.code === 0) {
        setContent('')
        fetchJournals()
        setShowHistory(true)
      }
    } catch (err) {
      console.error('Failed to save journal:', err)
    } finally {
      setSaving(false)
    }
  }

  const deleteJournal = async (id: number) => {
    try {
      const token = localStorage.getItem('token')
      if (!token) return

      const res = await fetch(`/api/v1/journal/${id}`, {
        method: 'DELETE',
        headers: { Authorization: `Bearer ${token}` },
      })
      const json = await res.json()
      if (json.code === 0) {
        setJournals(journals.filter(j => j.id !== id))
      }
    } catch (err) {
      console.error('Failed to delete journal:', err)
    }
  }

  const formatDate = (dateStr: string) => {
    const d = new Date(dateStr)
    const month = d.getMonth() + 1
    const day = d.getDate()
    const hours = String(d.getHours()).padStart(2, '0')
    const minutes = String(d.getMinutes()).padStart(2, '0')
    return `${month}月${day}日 ${hours}:${minutes}`
  }

  return (
    <div className="flex flex-col space-y-4">
      <h3 className="text-lg font-medium text-ink">私密日记</h3>

      {!loggedIn ? (
        <div className="bg-surface/70 rounded-xl p-4 text-center space-y-3">
          <p className="text-sm text-ink-soft">私密日记只有你自己能看到，登录后即可书写与回看。</p>
          <a
            href="/login?redirect=/sanctuary"
            className="inline-block px-4 py-1.5 rounded-full text-sm font-medium bg-accent text-white hover:bg-accent-hover transition-all duration-300"
          >
            登录后书写
          </a>
        </div>
      ) : (
      <>
      {/* 输入区域 */}
      <div className="bg-surface/70 rounded-xl p-4 space-y-3">
        <textarea
          value={content}
          onChange={(e) => setContent(e.target.value)}
          placeholder="今天想说些什么..."
          className="w-full h-24 p-3 text-sm text-ink bg-transparent border-none resize-none focus:outline-none placeholder-ink-faint"
          maxLength={2000}
        />
        <div className="flex items-center justify-between">
          <span className="text-xs text-ink-faint">{content.length}/2000</span>
          <button
            onClick={saveJournal}
            disabled={saving || !content.trim()}
            className={`px-4 py-1.5 rounded-full text-sm font-medium transition-all duration-300 ${
              saving || !content.trim()
                ? 'bg-surface-sunk text-ink-faint cursor-not-allowed'
                : 'bg-accent text-white hover:bg-accent-hover'
            }`}
          >
            {saving ? '保存中...' : '保存'}
          </button>
        </div>
      </div>

      {/* 历史记录 */}
      <button
        onClick={() => setShowHistory(!showHistory)}
        className="text-xs text-ink-soft hover:text-ink self-start"
      >
        {showHistory ? '收起历史' : `查看历史 (${journals.length})`}
      </button>

      {showHistory && (
        <div className="space-y-2 max-h-64 overflow-y-auto">
          {loading ? (
            <p className="text-xs text-ink-faint text-center py-4">加载中...</p>
          ) : journals.length === 0 ? (
            <p className="text-xs text-ink-faint text-center py-4">还没有日记</p>
          ) : (
            journals.map((journal) => (
              <div
                key={journal.id}
                className="bg-surface/50 rounded-lg p-3 group"
              >
                <p className="text-sm text-ink-soft whitespace-pre-wrap">{journal.content}</p>
                <div className="flex items-center justify-between mt-2">
                  <span className="text-xs text-ink-faint">
                    {formatDate(journal.created_at)}
                  </span>
                  <button
                    onClick={() => deleteJournal(journal.id)}
                    className="text-xs text-emotion-ink opacity-0 group-hover:opacity-100 transition-opacity hover:underline"
                  >
                    删除
                  </button>
                </div>
              </div>
            ))
          )}
        </div>
      )}
      </>
      )}
    </div>
  )
}
