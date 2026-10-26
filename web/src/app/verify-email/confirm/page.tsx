'use client'

import { useEffect, useState, Suspense } from 'react'
import Link from 'next/link'
import { useSearchParams } from 'next/navigation'

function VerifyEmailConfirmContent() {
  const searchParams = useSearchParams()
  const token = searchParams.get('token')
  const [status, setStatus] = useState<'loading' | 'success' | 'error'>('loading')
  const [message, setMessage] = useState('')

  useEffect(() => {
    if (!token) {
      setStatus('error')
      setMessage('验证链接无效，缺少验证参数。')
      return
    }

    // 调用验证 API
    const verifyEmail = async () => {
      try {
        const response = await fetch(`/api/v1/auth/verify/${token}`)
        const data = await response.json()

        if (data.code === 0) {
          setStatus('success')
          setMessage('邮箱验证成功！现在你可以发布帖子、参与社区互动了。')
        } else {
          setStatus('error')
          setMessage(data.message || '验证失败，请重新尝试。')
        }
      } catch (error) {
        setStatus('error')
        setMessage('验证请求失败，请检查网络连接后重新尝试。')
      }
    }

    verifyEmail()
  }, [token])

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 via-indigo-50 to-purple-50 px-4">
      <div className="w-full max-w-md">
        {/* Logo */}
        <div className="text-center mb-8">
          <Link href="/" className="inline-block">
            <span className="text-3xl font-bold text-blue-600">失业么</span>
          </Link>
        </div>

        <div className="bg-white rounded-lg shadow-md p-8 text-center">
          {status === 'loading' && (
            <>
              <div className="mx-auto w-16 h-16 bg-blue-100 rounded-full flex items-center justify-center mb-6 animate-spin">
                <svg className="w-8 h-8 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                </svg>
              </div>
              <h1 className="text-2xl font-bold text-gray-900 mb-4">正在验证...</h1>
              <p className="text-gray-600">请稍候，正在验证你的邮箱。</p>
            </>
          )}

          {status === 'success' && (
            <>
              <div className="mx-auto w-16 h-16 bg-green-100 rounded-full flex items-center justify-center mb-6">
                <svg className="w-8 h-8 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                </svg>
              </div>
              <h1 className="text-2xl font-bold text-gray-900 mb-4">验证成功！</h1>
              <p className="text-gray-600 mb-6">{message}</p>
              <Link
                href="/feed"
                className="inline-block w-full bg-blue-600 text-white py-2 px-4 rounded-lg hover:bg-blue-700 transition-colors"
              >
                进入社区
              </Link>
            </>
          )}

          {status === 'error' && (
            <>
              <div className="mx-auto w-16 h-16 bg-red-100 rounded-full flex items-center justify-center mb-6">
                <svg className="w-8 h-8 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                </svg>
              </div>
              <h1 className="text-2xl font-bold text-gray-900 mb-4">验证失败</h1>
              <p className="text-gray-600 mb-6">{message}</p>
              <div className="space-y-4">
                <Link
                  href="/verify-email"
                  className="block w-full bg-blue-600 text-white py-2 px-4 rounded-lg hover:bg-blue-700 transition-colors"
                >
                  重新发送验证邮件
                </Link>
                <Link
                  href="/feed"
                  className="block w-full text-center text-blue-600 hover:text-blue-800 py-2"
                >
                  先去逛逛社区
                </Link>
              </div>
            </>
          )}
        </div>
      </div>
    </div>
  )
}

export default function VerifyEmailConfirmPage() {
  return (
    <Suspense fallback={
      <div className="min-h-screen flex items-center justify-center">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
      </div>
    }>
      <VerifyEmailConfirmContent />
    </Suspense>
  )
}
