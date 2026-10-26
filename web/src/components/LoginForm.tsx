'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import Link from 'next/link'
import { post } from '@/lib/api'

export default function LoginForm() {
  const router = useRouter()
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)
    setError('')

    try {
      const res = await post<{ token: string; user: { id: number; nickname: string } }>('/auth/login/email', {
        email,
        password,
      })

      if (res.code === 0) {
        localStorage.setItem('token', res.data.token)
        localStorage.setItem('user', JSON.stringify(res.data.user))
        router.push('/feed')
      } else {
        setError(res.message || '登录失败')
      }
    } catch {
      setError('网络错误，请稍后重试')
    } finally {
      setLoading(false)
    }
  }

  return (
    <form className="space-y-5" onSubmit={handleSubmit}>
      {error && (
        <div className="bg-emotion-soft border border-emotion-line text-emotion-ink px-4 py-3 rounded-lg text-sm">
          {error}
        </div>
      )}

      <div>
        <label htmlFor="email" className="block text-sm font-medium text-ink-soft">
          邮箱
        </label>
        <input
          id="email"
          type="email"
          required
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          className="input-field"
          placeholder="your@email.com"
        />
      </div>

      <div>
        <label htmlFor="password" className="block text-sm font-medium text-ink-soft">
          密码
        </label>
        <input
          id="password"
          type="password"
          required
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          className="input-field"
          placeholder="请输入密码"
        />
      </div>

      <button type="submit" disabled={loading} className="btn-primary w-full">
        {loading ? (
          <span className="flex items-center justify-center gap-2">
            <svg className="animate-spin h-4 w-4" viewBox="0 0 24 24">
              <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" fill="none" />
              <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
            </svg>
            登录中...
          </span>
        ) : '登录'}
      </button>

      <div className="flex items-center justify-between text-sm">
        <Link href="/forgot-password" className="text-ink-faint hover:text-accent">
          忘记密码？
        </Link>
        <div>
          <span className="text-ink-faint">还没有账号？</span>
          <Link href="/register" className="text-accent hover:text-accent-hover font-medium ml-1">
            立即注册
          </Link>
        </div>
      </div>
    </form>
  )
}
