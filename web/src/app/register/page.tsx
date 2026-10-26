import Link from 'next/link'
import RegisterForm from '@/components/RegisterForm'

export default function RegisterPage() {
  return (
    <div className="min-h-screen flex items-center justify-center bg-canvas px-4 py-10">
      <div className="w-full max-w-md">
        {/* Logo */}
        <div className="text-center mb-8">
          <Link href="/" className="inline-block">
            <span className="text-2xl font-semibold tracking-tight text-accent">失业么</span>
          </Link>
        </div>

        <div className="card p-8">
          <div className="mb-6">
            <h1 className="text-2xl font-semibold text-center text-ink">加入失业么</h1>
            <p className="mt-1.5 text-center text-sm text-ink-faint">开始你的能力重建之路</p>
          </div>

          <RegisterForm />
        </div>

        <p className="mt-6 text-center text-xs text-ink-faint">
          注册即表示你同意我们的{' '}
          <Link href="/terms" className="text-accent hover:underline">
            服务条款
          </Link>{' '}
          和{' '}
          <Link href="/privacy" className="text-accent hover:underline">
            隐私政策
          </Link>
        </p>
      </div>
    </div>
  )
}
