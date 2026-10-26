import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import Button from '@/components/Button'

describe('Button', () => {
  it('renders children text', () => {
    render(<Button>保存</Button>)
    expect(screen.getByRole('button')).toHaveTextContent('保存')
  })

  it('applies primary variant styles by default', () => {
    render(<Button>Test</Button>)
    const btn = screen.getByRole('button')
    expect(btn.className).toContain('bg-blue-600')
  })

  it('applies secondary variant styles', () => {
    render(<Button variant="secondary">Test</Button>)
    expect(screen.getByRole('button').className).toContain('bg-gray-600')
  })

  it('applies danger variant styles', () => {
    render(<Button variant="danger">Delete</Button>)
    expect(screen.getByRole('button').className).toContain('bg-red-600')
  })

  it('applies size classes correctly', () => {
    const { rerender } = render(<Button size="sm">S</Button>)
    expect(screen.getByRole('button').className).toContain('px-3')
    rerender(<Button size="lg">L</Button>)
    expect(screen.getByRole('button').className).toContain('px-6')
  })

  it('calls onClick when clicked', async () => {
    const user = userEvent.setup()
    const onClick = vi.fn()
    render(<Button onClick={onClick}>Click</Button>)
    await user.click(screen.getByRole('button'))
    expect(onClick).toHaveBeenCalledOnce()
  })

  it('is disabled when disabled prop is true', () => {
    render(<Button disabled>Disabled</Button>)
    expect(screen.getByRole('button')).toBeDisabled()
  })

  it('is disabled when loading', () => {
    render(<Button loading>Loading</Button>)
    expect(screen.getByRole('button')).toBeDisabled()
  })

  it('shows spinner icon when loading', () => {
    render(<Button loading>Saving</Button>)
    const svg = screen.getByRole('button').querySelector('svg')
    expect(svg).toBeInTheDocument()
    expect(svg?.className.baseVal).toContain('animate-spin')
  })
})
