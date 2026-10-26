import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import Avatar from '@/components/Avatar'

describe('Avatar', () => {
  it('renders initial letter when no avatar image provided', () => {
    const { container } = render(<Avatar nickname="小明" />)
    expect(container.textContent).toBe('小')
  })

  it('renders question mark when no nickname', () => {
    const { container } = render(<Avatar />)
    expect(container.textContent).toBe('?')
  })

  it('renders img tag when avatar URL is provided', () => {
    render(<Avatar nickname="用户" avatar="https://example.com/avatar.jpg" />)
    const img = screen.getByRole('img')
    expect(img).toHaveAttribute('src', 'https://example.com/avatar.jpg')
    expect(img).toHaveAttribute('alt', '用户')
  })

  it('applies size classes correctly', () => {
    const { container } = render(<Avatar nickname="用户" size="sm" />)
    const el = container.firstChild as HTMLElement
    expect(el.className).toContain('w-8')
  })

  it('renders large avatar correctly', () => {
    render(<Avatar nickname="用户" avatar="https://example.com/a.jpg" size="lg" />)
    const img = screen.getByRole('img')
    expect(img.className).toContain('w-16')
  })
})
