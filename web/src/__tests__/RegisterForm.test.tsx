import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import RegisterForm from '@/components/RegisterForm'

// mock next/link
vi.mock('next/link', () => {
  return {
    default: ({ children, href, ...props }: any) => (
      <a href={href} {...props}>{children}</a>
    ),
  }
})

// mock api
const mockPost = vi.fn()
vi.mock('@/lib/api', () => ({
  post: (...args: any[]) => mockPost(...args),
}))

// mock next/navigation
const mockPush = vi.fn()
vi.mock('next/navigation', () => ({
  useRouter: () => ({ push: mockPush }),
}))

describe('RegisterForm', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    localStorage.clear()
  })

  it('renders all required form fields', () => {
    render(<RegisterForm />)
    expect(screen.getByLabelText(/邮箱/)).toBeInTheDocument()
    expect(screen.getByLabelText(/昵称/)).toBeInTheDocument()
    expect(screen.getByPlaceholderText('至少6位')).toBeInTheDocument()
    expect(screen.getByPlaceholderText('再次输入密码')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: '注册' })).toBeInTheDocument()
  })

  it('renders login link', () => {
    render(<RegisterForm />)
    expect(screen.getByText('立即登录')).toBeInTheDocument()
  })

  it('shows error when passwords do not match', async () => {
    const user = userEvent.setup()
    render(<RegisterForm />)

    await user.type(screen.getByLabelText(/邮箱/), 'test@test.com')
    await user.type(screen.getByLabelText(/昵称/), '测试用户')
    await user.type(screen.getByPlaceholderText('至少6位'), '123456')
    await user.type(screen.getByPlaceholderText('再次输入密码'), '654321')
    await user.click(screen.getByRole('button', { name: '注册' }))

    expect(screen.getByText('两次输入的密码不一致')).toBeInTheDocument()
    expect(mockPost).not.toHaveBeenCalled()
  })

  it('shows error when password is too short', async () => {
    const user = userEvent.setup()
    render(<RegisterForm />)

    await user.type(screen.getByLabelText(/邮箱/), 'test@test.com')
    await user.type(screen.getByLabelText(/昵称/), '测试用户')
    await user.type(screen.getByPlaceholderText('至少6位'), '123')
    await user.type(screen.getByPlaceholderText('再次输入密码'), '123')
    await user.click(screen.getByRole('button', { name: '注册' }))

    expect(screen.getByText('密码长度至少6位')).toBeInTheDocument()
    expect(mockPost).not.toHaveBeenCalled()
  })

  it('shows error when nickname is too short', async () => {
    const user = userEvent.setup()
    render(<RegisterForm />)

    await user.type(screen.getByLabelText(/邮箱/), 'test@test.com')
    await user.type(screen.getByLabelText(/昵称/), 'A')
    await user.type(screen.getByPlaceholderText('至少6位'), '123456')
    await user.type(screen.getByPlaceholderText('再次输入密码'), '123456')
    await user.click(screen.getByRole('button', { name: '注册' }))

    expect(screen.getByText('昵称长度需要在2-50个字符之间')).toBeInTheDocument()
    expect(mockPost).not.toHaveBeenCalled()
  })

  it('calls API on valid submission', async () => {
    const user = userEvent.setup()
    mockPost.mockResolvedValue({ code: 0, data: { token: 'abc', user: { id: 1, nickname: '测试' } } })

    render(<RegisterForm />)

    await user.type(screen.getByLabelText(/邮箱/), 'test@test.com')
    await user.type(screen.getByLabelText(/昵称/), '测试用户')
    await user.type(screen.getByPlaceholderText('至少6位'), '123456')
    await user.type(screen.getByPlaceholderText('再次输入密码'), '123456')
    await user.click(screen.getByRole('button', { name: '注册' }))

    await waitFor(() => {
      expect(mockPost).toHaveBeenCalledWith('/auth/register/email', expect.objectContaining({
        email: 'test@test.com',
        nickname: '测试用户',
      }))
    })
  })

  it('stores token and navigates on success', async () => {
    const user = userEvent.setup()
    mockPost.mockResolvedValue({ code: 0, data: { token: 'jwt-token', user: { id: 1, nickname: '测试' } } })

    render(<RegisterForm />)

    await user.type(screen.getByLabelText(/邮箱/), 'test@test.com')
    await user.type(screen.getByLabelText(/昵称/), '测试用户')
    await user.type(screen.getByPlaceholderText('至少6位'), '123456')
    await user.type(screen.getByPlaceholderText('再次输入密码'), '123456')
    await user.click(screen.getByRole('button', { name: '注册' }))

    await waitFor(() => {
      expect(localStorage.getItem('token')).toBe('jwt-token')
      expect(mockPush).toHaveBeenCalledWith('/feed')
    })
  })

  it('shows server error message', async () => {
    const user = userEvent.setup()
    mockPost.mockResolvedValue({ code: 40001, message: '邮箱已注册' })

    render(<RegisterForm />)

    await user.type(screen.getByLabelText(/邮箱/), 'test@test.com')
    await user.type(screen.getByLabelText(/昵称/), '测试用户')
    await user.type(screen.getByPlaceholderText('至少6位'), '123456')
    await user.type(screen.getByPlaceholderText('再次输入密码'), '123456')
    await user.click(screen.getByRole('button', { name: '注册' }))

    await waitFor(() => {
      expect(screen.getByText('邮箱已注册')).toBeInTheDocument()
    })
  })

  it('shows network error on failure', async () => {
    const user = userEvent.setup()
    mockPost.mockRejectedValue(new Error('Network error'))

    render(<RegisterForm />)

    await user.type(screen.getByLabelText(/邮箱/), 'test@test.com')
    await user.type(screen.getByLabelText(/昵称/), '测试用户')
    await user.type(screen.getByPlaceholderText('至少6位'), '123456')
    await user.type(screen.getByPlaceholderText('再次输入密码'), '123456')
    await user.click(screen.getByRole('button', { name: '注册' }))

    await waitFor(() => {
      expect(screen.getByText('网络错误，请稍后重试')).toBeInTheDocument()
    })
  })

  it('disables submit button while loading', async () => {
    const user = userEvent.setup()
    let resolvePost: any
    mockPost.mockImplementation(() => new Promise(r => { resolvePost = r }))

    render(<RegisterForm />)

    await user.type(screen.getByLabelText(/邮箱/), 'test@test.com')
    await user.type(screen.getByLabelText(/昵称/), '测试用户')
    await user.type(screen.getByPlaceholderText('至少6位'), '123456')
    await user.type(screen.getByPlaceholderText('再次输入密码'), '123456')
    await user.click(screen.getByRole('button', { name: '注册' }))

    await waitFor(() => {
      expect(screen.getByRole('button', { name: /注册中/ })).toBeDisabled()
    })

    // Resolve to clean up
    resolvePost({ code: 0, data: { token: '', user: {} } })
  })
})
