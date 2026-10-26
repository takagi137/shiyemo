import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import Input from '@/components/Input'

describe('Input', () => {
  it('renders an input element', () => {
    render(<Input placeholder="请输入" />)
    expect(screen.getByPlaceholderText('请输入')).toBeInTheDocument()
  })

  it('shows label when provided', () => {
    render(<Input label="用户名" />)
    expect(screen.getByText('用户名')).toBeInTheDocument()
    expect(screen.getByLabelText('用户名')).toBeInTheDocument()
  })

  it('shows error message when error prop is provided', () => {
    render(<Input label="邮箱" error="邮箱格式不正确" />)
    expect(screen.getByText('邮箱格式不正确')).toBeInTheDocument()
    expect(screen.getByLabelText('邮箱').className).toContain('border-red-500')
  })

  it('shows helper text when provided and no error', () => {
    render(<Input helperText="请输入6位以上密码" />)
    expect(screen.getByText('请输入6位以上密码')).toBeInTheDocument()
  })

  it('hides helper text when error is present', () => {
    render(<Input helperText="提示" error="错误" />)
    expect(screen.queryByText('提示')).not.toBeInTheDocument()
    expect(screen.getByText('错误')).toBeInTheDocument()
  })

  it('forwards input props correctly', () => {
    render(<Input type="email" maxLength={100} disabled />)
    const input = screen.getByRole('textbox')
    expect(input).toHaveAttribute('type', 'email')
    expect(input).toHaveAttribute('maxLength', '100')
    expect(input).toBeDisabled()
  })
})
