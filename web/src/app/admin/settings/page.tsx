'use client'

import { useState, useEffect } from 'react'
import { useRouter } from 'next/navigation'
import Link from 'next/link'
import { get, put, post } from '@/lib/api'
import { toastSuccess, toastError } from '@/lib/toast'

interface SMTPSettings {
  enabled: boolean
  host: string
  port: number
  username: string
  password: string
  password_set: boolean
  from: string
  email_verify_required: boolean
}

export default function AdminSettingsPage() {
  const router = useRouter()
  const [form, setForm] = useState<SMTPSettings | null>(null)
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  const [testing, setTesting] = useState(false)
  const [testTo, setTestTo] = useState('')

  useEffect(() => {
    if (!localStorage.getItem('token')) router.push('/login')
  }, [router])

  useEffect(() => {
    const fetchSettings = async () => {
      try {
        const res = await get<SMTPSettings>('/admin/settings/smtp')
        if (res.code === 0) setForm(res.data)
        else toastError(res.message || '获取配置失败')
      } catch {
        toastError('获取配置失败')
      } finally {
        setLoading(false)
      }
    }
    fetchSettings()
  }, [])

  const update = <K extends keyof SMTPSettings>(key: K, value: SMTPSettings[K]) => {
    setForm((prev) => (prev ? { ...prev, [key]: value } : prev))
  }

  const handleSave = async () => {
    if (!form) return
    setSaving(true)
    try {
      const res = await put('/admin/settings/smtp', {
        enabled: form.enabled,
        host: form.host,
        port: Number(form.port) || 0,
        username: form.username,
        password: form.password, // 留空表示不修改
        from: form.from,
        email_verify_required: form.email_verify_required,
      })
      if (res.code === 0) {
        toastSuccess('配置已保存，立即生效')
        // 重新拉取：密码不回显，且后端可能因 SMTP 关闭而重置强制验证开关
        const fresh = await get<SMTPSettings>('/admin/settings/smtp')
        if (fresh.code === 0) setForm(fresh.data)
      } else {
        toastError(res.message || '保存失败')
      }
    } catch {
      toastError('保存失败')
    } finally {
      setSaving(false)
    }
  }

  const handleTest = async () => {
    const to = testTo.trim()
    if (!to) {
      toastError('请填写接收测试邮件的邮箱')
      return
    }
    setTesting(true)
    try {
      const res = await post('/admin/settings/smtp/test', { to })
      if (res.code === 0) toastSuccess(`测试邮件已发送至 ${to}`)
      else toastError(res.message || '发送失败')
    } catch {
      toastError('发送失败，请检查配置')
    } finally {
      setTesting(false)
    }
  }

  if (loading) {
    return (
      <div className="container mx-auto px-4 py-8 max-w-3xl">
        <div className="card p-12 text-center animate-pulse text-ink-faint">加载中…</div>
      </div>
    )
  }

  if (!form) {
    return (
      <div className="container mx-auto px-4 py-8 max-w-3xl">
        <div className="card p-12 text-center text-ink-soft">配置加载失败</div>
      </div>
    )
  }

  return (
    <div className="container mx-auto px-4 py-8 max-w-3xl">
      <div className="mb-6 flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-semibold text-ink">系统设置</h1>
          <p className="text-sm text-ink-soft mt-1">邮件服务配置，保存后立即生效，无需重启服务</p>
        </div>
        <Link href="/admin" className="text-sm text-ink-soft hover:text-accent transition-colors">
          ← 返回管理后台
        </Link>
      </div>

      {/* SMTP 开关 */}
      <div className="card p-6 mb-4">
        <label className="flex items-start gap-3 cursor-pointer">
          <input
            type="checkbox"
            checked={form.enabled}
            onChange={(e) => update('enabled', e.target.checked)}
            className="mt-1 w-4 h-4 accent-[var(--color-accent)]"
          />
          <span>
            <span className="font-medium text-ink">启用 SMTP 邮件发送</span>
            <span className="block text-sm text-ink-soft mt-0.5">
              关闭时，验证邮件与密码重置邮件都不会真实投递，仅记录日志。
            </span>
          </span>
        </label>
      </div>

      {/* SMTP 参数 */}
      <div className="card p-6 mb-4 space-y-4">
        <h2 className="font-medium text-ink">SMTP 服务器</h2>

        <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
          <div className="sm:col-span-2">
            <label className="block text-sm font-medium text-ink-soft">服务器地址</label>
            <input
              type="text"
              value={form.host}
              onChange={(e) => update('host', e.target.value)}
              placeholder="smtp.qq.com"
              className="input-field"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-ink-soft">端口</label>
            <input
              type="number"
              value={form.port || ''}
              onChange={(e) => update('port', Number(e.target.value))}
              placeholder="465"
              className="input-field"
            />
          </div>
        </div>

        <div>
          <label className="block text-sm font-medium text-ink-soft">用户名</label>
          <input
            type="text"
            value={form.username}
            onChange={(e) => update('username', e.target.value)}
            placeholder="you@example.com"
            className="input-field"
          />
        </div>

        <div>
          <label className="block text-sm font-medium text-ink-soft">
            密码 / 授权码
            {form.password_set && (
              <span className="ml-2 text-xs font-normal text-action-ink">已配置（留空则不修改）</span>
            )}
          </label>
          <input
            type="password"
            value={form.password}
            onChange={(e) => update('password', e.target.value)}
            placeholder={form.password_set ? '••••••••（留空保持不变）' : '请输入 SMTP 授权码'}
            className="input-field"
            autoComplete="new-password"
          />
          <p className="text-xs text-ink-faint mt-1">
            出于安全考虑，已保存的密码不会回显。
          </p>
        </div>

        <div>
          <label className="block text-sm font-medium text-ink-soft">发件人地址</label>
          <input
            type="email"
            value={form.from}
            onChange={(e) => update('from', e.target.value)}
            placeholder="noreply@yourdomain.com"
            className="input-field"
          />
        </div>
      </div>

      {/* 邮箱验证策略 */}
      <div className="card p-6 mb-4">
        <label className={`flex items-start gap-3 ${form.enabled ? 'cursor-pointer' : 'cursor-not-allowed opacity-60'}`}>
          <input
            type="checkbox"
            checked={form.email_verify_required}
            disabled={!form.enabled}
            onChange={(e) => update('email_verify_required', e.target.checked)}
            className="mt-1 w-4 h-4 accent-[var(--color-accent)]"
          />
          <span>
            <span className="font-medium text-ink">强制邮箱验证后才能发帖</span>
            <span className="block text-sm text-ink-soft mt-0.5">
              开启后，用户必须先点击验证邮件中的链接，才能发帖和评论。
            </span>
          </span>
        </label>

        {!form.enabled && (
          <div className="mt-3 text-sm bg-insight-soft text-insight-ink border border-insight-line rounded-lg px-3 py-2">
            ⚠️ SMTP 未启用时该选项被强制关闭 —— 否则用户收不到验证邮件，将永远无法发帖。
          </div>
        )}
      </div>

      {/* 保存 */}
      <div className="flex justify-end mb-8">
        <button onClick={handleSave} disabled={saving} className="btn-primary">
          {saving ? '保存中…' : '保存配置'}
        </button>
      </div>

      {/* 测试发信 */}
      <div className="card p-6">
        <h2 className="font-medium text-ink mb-1">发送测试邮件</h2>
        <p className="text-sm text-ink-soft mb-4">
          用当前已保存的配置真实投递一封测试邮件。即使上方「启用 SMTP」处于关闭状态也会发送，
          方便你在正式开启前先验证配置是否正确。
        </p>
        <div className="flex flex-col sm:flex-row gap-3">
          <input
            type="email"
            value={testTo}
            onChange={(e) => setTestTo(e.target.value)}
            placeholder="接收测试邮件的邮箱"
            className="input-field flex-1 mt-0"
          />
          <button onClick={handleTest} disabled={testing} className="btn-secondary whitespace-nowrap">
            {testing ? '发送中…' : '发送测试邮件'}
          </button>
        </div>
      </div>
    </div>
  )
}
