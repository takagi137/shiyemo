import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import LoginForm from '@/components/LoginForm'

// Mock next/navigation
vi.mock('next/navigation', () => ({
  useRouter: () => ({ push: vi.fn() }),
}))

// Mock @/lib/api
vi.mock('@/lib/api', () => ({
  post: vi.fn(),
}))

describe('LoginForm', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    localStorage.clear()
  })

  it('renders email and password inputs', () => {
    render(<LoginForm />)
    expect(screen.getByLabelText('邮箱')).toBeInTheDocument()
    expect(screen.getByLabelText('密码')).toBeInTheDocument()
  })

  it('renders submit button with "登录" text', () => {
    render(<LoginForm />)
    expect(screen.getByRole('button', { name: '登录' })).toBeInTheDocument()
  })

  it('renders register link', () => {
    render(<LoginForm />)
    expect(screen.getByText('立即注册')).toBeInTheDocument()
  })

  it('renders forgot password link', () => {
    render(<LoginForm />)
    expect(screen.getByText('忘记密码？')).toBeInTheDocument()
  })

  it('shows error message on failed login', async () => {
    const user = userEvent.setup()
    const { post } = await import('@/lib/api')
    vi.mocked(post).mockResolvedValue({ code: 40100, message: '邮箱或密码错误', data: null })

    render(<LoginForm />)

    await user.type(screen.getByLabelText('邮箱'), 'test@example.com')
    await user.type(screen.getByLabelText('密码'), 'wrongpass')
    await user.click(screen.getByRole('button', { name: '登录' }))

    await waitFor(() => {
      expect(screen.getByText('邮箱或密码错误')).toBeInTheDocument()
    })
  })

  it('stores token in localStorage on successful login', async () => {
    const user = userEvent.setup()
    const { post } = await import('@/lib/api')
    vi.mocked(post).mockResolvedValue({
      code: 0,
      message: 'ok',
      data: { token: 'fake-jwt-token', user: { id: 1, nickname: '测试' } },
    })

    render(<LoginForm />)

    await user.type(screen.getByLabelText('邮箱'), 'test@example.com')
    await user.type(screen.getByLabelText('密码'), 'correctpass')
    await user.click(screen.getByRole('button', { name: '登录' }))

    await waitFor(() => {
      expect(localStorage.getItem('token')).toBe('fake-jwt-token')
    })
  })
})
