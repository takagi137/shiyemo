import { describe, it, expect } from 'vitest'
import { render } from '@testing-library/react'
import Icon, { type IconName } from '@/components/Icon'

// 所有声明的 IconName 必须能渲染（防止 PATHS 漏配 / IconName 漂移）
const ALL_NAMES: IconName[] = [
  'emotion', 'insight', 'case', 'action', 'achievement',
  'all', 'following', 'cohort',
  'heart', 'comment', 'clock', 'pin', 'star', 'mask', 'bookmark', 'share',
]

describe('Icon', () => {
  it.each(ALL_NAMES)('renders svg for %s', (name) => {
    const { container } = render(<Icon name={name} />)
    const svg = container.querySelector('svg')
    expect(svg).not.toBeNull()
    expect(svg!.tagName.toLowerCase()).toBe('svg')
  })

  it('uses default size of 16', () => {
    const { container } = render(<Icon name="heart" />)
    const svg = container.querySelector('svg')!
    expect(svg.getAttribute('width')).toBe('16')
    expect(svg.getAttribute('height')).toBe('16')
  })

  it('respects custom size', () => {
    const { container } = render(<Icon name="heart" size={24} />)
    const svg = container.querySelector('svg')!
    expect(svg.getAttribute('width')).toBe('24')
    expect(svg.getAttribute('height')).toBe('24')
  })

  it('passes className through to svg', () => {
    const { container } = render(<Icon name="heart" className="text-accent w-5" />)
    const svg = container.querySelector('svg')!
    // jsdom 下 SVGElement.className 是 SVGAnimatedString，需用 getAttribute
    const cls = svg.getAttribute('class') ?? ''
    expect(cls).toContain('text-accent')
    expect(cls).toContain('w-5')
    // 基础类也应保留
    expect(cls).toContain('inline-block')
    expect(cls).toContain('shrink-0')
  })

  it('uses currentColor for stroke (inherits parent color)', () => {
    const { container } = render(<Icon name="star" />)
    const svg = container.querySelector('svg')!
    expect(svg.getAttribute('stroke')).toBe('currentColor')
  })

  it('renders the correct viewBox for 24x24 grid', () => {
    const { container } = render(<Icon name="pin" />)
    const svg = container.querySelector('svg')!
    expect(svg.getAttribute('viewBox')).toBe('0 0 24 24')
  })

  describe('accessibility', () => {
    it('is aria-hidden by default (decorative)', () => {
      const { container } = render(<Icon name="comment" />)
      const svg = container.querySelector('svg')!
      expect(svg.getAttribute('aria-hidden')).toBe('true')
      expect(svg.getAttribute('role')).toBeNull()
    })

    it('exposes role=img and aria-label when label provided', () => {
      const { container } = render(<Icon name="comment" label="3 条评论" />)
      const svg = container.querySelector('svg')!
      expect(svg.getAttribute('role')).toBe('img')
      expect(svg.getAttribute('aria-label')).toBe('3 条评论')
      expect(svg.getAttribute('aria-hidden')).toBeNull()
    })
  })

  describe('paths', () => {
    it('emits one path per entry in PATHS', () => {
      // insight 有 3 条 path，其余 1~5 条不等。用一个明确的图标验证至少 1 条
      const { container } = render(<Icon name="share" />)
      const paths = container.querySelectorAll('path')
      expect(paths.length).toBeGreaterThanOrEqual(1)
      // 每个 path 都有 d 属性
      paths.forEach((p) => {
        expect(p.getAttribute('d')).toBeTruthy()
      })
    })

    it('renders the right number of paths for multi-path icons', () => {
      // share 5 条（与 PATHS 表对齐，防回归）
      const { container } = render(<Icon name="share" />)
      expect(container.querySelectorAll('path').length).toBe(5)
    })
  })
})