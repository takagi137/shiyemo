import Link from 'next/link'
import LoginForm from '@/components/LoginForm'

export default function LoginPage() {
  return (
    <div className="min-h-screen flex items-center justify-center bg-canvas px-4">
      <div className="w-full max-w-md">
        {/* Logo */}
        <div className="text-center mb-8">
          <Link href="/" className="inline-block">
            <span className="text-2xl font-semibold tracking-tight text-accent">失业么</span>
          </Link>
        </div>

        <div className="card p-8">
          <div className="mb-6">
            <h1 className="text-2xl font-semibold text-center text-ink">欢迎回来</h1>
            <p className="mt-1.5 text-center text-sm text-ink-faint">登录你的账号，继续重建之路</p>
          </div>

          <LoginForm />
        </div>
      </div>
    </div>
  )
}
