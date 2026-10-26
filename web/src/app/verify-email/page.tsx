'use client'

import { useState } from 'react'
import Link from 'next/link'
import { getApiBase } from '@/lib/api'

export default function VerifyEmailPage() {
  const [sending, setSending] = useState(false)
  const [sent, setSent] = useState(false)
  const [error, setError] = useState('')

  const handleResend = async () => {
    setSending(true)
    setError('')

    const token = localStorage.getItem('token')
    if (!token) {
      setError('请先登录')
      setSending(false)
      return
    }

    try {
      const res = await fetch(`${getApiBase()}/auth/verify/send`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
      })
      const data = await res.json()
      if (data.code === 0) {
        setSent(true)
      } else {
        setError(data.message || '发送失败，请稍后重试')
      }
    } catch {
      setError('网络错误，请检查连接')
    } finally {
      setSending(false)
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 via-indigo-50 to-purple-50 px-4">
      <div className="w-full max-w-md">
        <div className="text-center mb-8">
          <Link href="/" className="inline-block">
            <span className="text-3xl font-bold text-blue-600">失业么</span>
          </Link>
        </div>

        <div className="bg-white rounded-lg shadow-md p-8 text-center">
          <div className="mx-auto w-16 h-16 bg-blue-100 rounded-full flex items-center justify-center mb-6">
            <svg className="w-8 h-8 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
            </svg>
          </div>

          <h1 className="text-2xl font-bold text-gray-900 mb-4">验证你的邮箱</h1>
          <p className="text-gray-600 mb-6">
            我们已向你的注册邮箱发送了一封验证邮件。
            <br />
            请查看收件箱（或垃圾邮件文件夹），点击链接完成验证。
          </p>

          {sent && (
            <div className="mb-4 p-3 bg-green-50 border border-green-200 rounded-lg text-green-700 text-sm">
              验证邮件已重新发送，请查收。
            </div>
          )}

          {error && (
            <div className="mb-4 p-3 bg-red-50 border border-red-200 rounded-lg text-red-700 text-sm">
              {error}
            </div>
          )}

          <div className="space-y-4">
            <button
              onClick={handleResend}
              disabled={sending || sent}
              className="w-full bg-blue-600 text-white py-2 px-4 rounded-lg hover:bg-blue-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {sending ? '发送中...' : sent ? '已发送' : '重新发送验证邮件'}
            </button>

            <Link
              href="/feed"
              className="block w-full text-center text-blue-600 hover:text-blue-800 py-2"
            >
              先去逛逛社区
            </Link>
          </div>

          <div className="mt-6 pt-6 border-t border-gray-200">
            <p className="text-sm text-gray-500">
              未验证邮箱也可以浏览社区，但发帖需要先完成验证。
            </p>
          </div>
        </div>
      </div>
    </div>
  )
}
