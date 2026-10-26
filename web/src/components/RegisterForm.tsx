'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import Link from 'next/link'
import { post } from '@/lib/api'

interface RegisterFormData {
  email: string
  password: string
  confirmPassword: string
  nickname: string
  city: string
  industry: string
  leftJobAt: string
}

export default function RegisterForm() {
  const router = useRouter()
  const [formData, setFormData] = useState<RegisterFormData>({
    email: '',
    password: '',
    confirmPassword: '',
    nickname: '',
    city: '',
    industry: '',
    leftJobAt: '',
  })
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    const { name, value } = e.target
    setFormData(prev => ({ ...prev, [name]: value }))
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)
    setError('')

    // 验证密码匹配
    if (formData.password !== formData.confirmPassword) {
      setError('两次输入的密码不一致')
      setLoading(false)
      return
    }

    // 验证密码长度
    if (formData.password.length < 6) {
      setError('密码长度至少6位')
      setLoading(false)
      return
    }

    // 验证昵称长度
    if (formData.nickname.length < 2 || formData.nickname.length > 50) {
      setError('昵称长度需要在2-50个字符之间')
      setLoading(false)
      return
    }

    try {
      const res = await post<{ token: string; user: { id: number; nickname: string } }>('/auth/register/email', {
        email: formData.email,
        password: formData.password,
        nickname: formData.nickname,
        city: formData.city || undefined,
        industry: formData.industry || undefined,
        left_job_at: formData.leftJobAt || undefined,
      })

      if (res.code === 0) {
        localStorage.setItem('token', res.data.token)
        localStorage.setItem('user', JSON.stringify(res.data.user))
        router.push('/feed')
      } else {
        setError(res.message || '注册失败')
      }
    } catch {
      setError('网络错误，请稍后重试')
    } finally {
      setLoading(false)
    }
  }

  return (
    <form className="space-y-4" onSubmit={handleSubmit}>
      {error && (
        <div className="bg-emotion-soft border border-emotion-line text-emotion-ink px-4 py-3 rounded-lg text-sm">
          {error}
        </div>
      )}

      {/* 邮箱 */}
      <div>
        <label htmlFor="email" className="block text-sm font-medium text-ink-soft">
          邮箱 <span className="text-emotion-ink">*</span>
        </label>
        <input
          id="email"
          name="email"
          type="email"
          required
          value={formData.email}
          onChange={handleChange}
          className="input-field"
          placeholder="your@email.com"
        />
      </div>

      {/* 昵称 */}
      <div>
        <label htmlFor="nickname" className="block text-sm font-medium text-ink-soft">
          昵称 <span className="text-emotion-ink">*</span>
        </label>
        <input
          id="nickname"
          name="nickname"
          type="text"
          required
          minLength={2}
          maxLength={50}
          value={formData.nickname}
          onChange={handleChange}
          className="input-field"
          placeholder="给自己取个名字"
        />
      </div>

      {/* 密码 */}
      <div>
        <label htmlFor="password" className="block text-sm font-medium text-ink-soft">
          密码 <span className="text-emotion-ink">*</span>
        </label>
        <input
          id="password"
          name="password"
          type="password"
          required
          minLength={6}
          value={formData.password}
          onChange={handleChange}
          className="input-field"
          placeholder="至少6位"
        />
      </div>

      {/* 确认密码 */}
      <div>
        <label htmlFor="confirmPassword" className="block text-sm font-medium text-ink-soft">
          确认密码 <span className="text-emotion-ink">*</span>
        </label>
        <input
          id="confirmPassword"
          name="confirmPassword"
          type="password"
          required
          value={formData.confirmPassword}
          onChange={handleChange}
          className="input-field"
          placeholder="再次输入密码"
        />
      </div>

      {/* 可选信息折叠区 */}
      <details className="group">
        <summary className="flex items-center justify-between cursor-pointer text-sm text-ink-faint hover:text-ink-soft">
          <span>完善个人信息（可选）</span>
          <svg className="h-5 w-5 transform group-open:rotate-180 transition-transform" viewBox="0 0 20 20" fill="currentColor">
            <path fillRule="evenodd" d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z" clipRule="evenodd" />
          </svg>
        </summary>
        <div className="mt-3 space-y-4">
          {/* 城市 */}
          <div>
            <label htmlFor="city" className="block text-sm font-medium text-ink-soft">
              所在城市
            </label>
            <input
              id="city"
              name="city"
              type="text"
              value={formData.city}
              onChange={handleChange}
              className="input-field"
              placeholder="如：北京"
            />
          </div>

          {/* 行业 */}
          <div>
            <label htmlFor="industry" className="block text-sm font-medium text-ink-soft">
              所在行业
            </label>
            <select
              id="industry"
              name="industry"
              value={formData.industry}
              onChange={handleChange}
              className="input-field"
            >
              <option value="">请选择行业</option>
              <option value="互联网">互联网</option>
              <option value="金融">金融</option>
              <option value="教育">教育</option>
              <option value="医疗">医疗</option>
              <option value="制造业">制造业</option>
              <option value="服务业">服务业</option>
              <option value="其他">其他</option>
            </select>
          </div>

          {/* 离职时间 */}
          <div>
            <label htmlFor="leftJobAt" className="block text-sm font-medium text-ink-soft">
              离职时间
            </label>
            <input
              id="leftJobAt"
              name="leftJobAt"
              type="date"
              value={formData.leftJobAt}
              onChange={handleChange}
              className="input-field"
            />
          </div>
        </div>
      </details>

      {/* 提交按钮 */}
      <button
        type="submit"
        disabled={loading}
        className="btn-primary w-full"
      >
        {loading ? (
          <span className="flex items-center justify-center gap-2">
            <svg className="animate-spin h-4 w-4" viewBox="0 0 24 24">
              <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" fill="none" />
              <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
            </svg>
            注册中...
          </span>
        ) : '注册'}
      </button>

      {/* 登录链接 */}
      <div className="text-center text-sm">
        <span className="text-ink-faint">已有账号？</span>
        <Link href="/login" className="text-accent hover:text-accent-hover font-medium ml-1">
          立即登录
        </Link>
      </div>
    </form>
  )
}
