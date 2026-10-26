import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import PostCard from '@/components/PostCard'

// mock next/link to render as <a>
vi.mock('next/link', () => {
  return {
    default: ({ children, href, ...props }: any) => (
      <a href={href} {...props}>{children}</a>
    ),
  }
})

// mock Avatar
vi.mock('@/components/Avatar', () => {
  return {
    default: ({ nickname, size }: any) => (
      <div data-testid="avatar" data-size={size}>{nickname}</div>
    ),
  }
})

const MOCK_POST = {
  id: 1,
  title: '如何度过低谷期',
  content: '<p>这是一段很长的内容</p>',
  type: 'emotion',
  author_id: 10,
  author_name: '小明',
  author_avatar: '',
  like_count: 12,
  comment_count: 5,
  view_count: 100,
  created_at: new Date(Date.now() - 30 * 60 * 1000).toISOString(),
}

describe('PostCard', () => {
  it('renders title', () => {
    render(<PostCard post={MOCK_POST} />)
    expect(screen.getByText('如何度过低谷期')).toBeInTheDocument()
  })

  it('renders author name', () => {
    render(<PostCard post={MOCK_POST} />)
    expect(screen.getAllByText('小明').length).toBeGreaterThanOrEqual(1)
  })

  it('renders like/comment/view counts', () => {
    render(<PostCard post={MOCK_POST} />)
    expect(screen.getByText('12')).toBeInTheDocument()
    expect(screen.getByText('5')).toBeInTheDocument()
    expect(screen.getByText('100')).toBeInTheDocument()
  })

  it('shows emotion type badge', () => {
    render(<PostCard post={MOCK_POST} />)
    // 视觉重做后 emoji 换成 Icon 组件，标签文字与图标分离，只校验文字
    expect(screen.getByText('情绪共鸣')).toBeInTheDocument()
  })

  it('shows insight type badge', () => {
    render(<PostCard post={{ ...MOCK_POST, type: 'insight' }} />)
    expect(screen.getByText('认知构建')).toBeInTheDocument()
  })

  it('shows 置顶 badge when is_top is true', () => {
    render(<PostCard post={{ ...MOCK_POST, is_top: true }} />)
    expect(screen.getByText('置顶')).toBeInTheDocument()
  })

  it('shows 精华 badge when is_featured is true', () => {
    render(<PostCard post={{ ...MOCK_POST, is_featured: true }} />)
    expect(screen.getByText('精华')).toBeInTheDocument()
  })

  it('does not show badges when flags are false', () => {
    render(<PostCard post={{ ...MOCK_POST, is_top: false, is_featured: false }} />)
    expect(screen.queryByText('置顶')).not.toBeInTheDocument()
    expect(screen.queryByText('精华')).not.toBeInTheDocument()
  })

  it('renders link to post detail', () => {
    render(<PostCard post={MOCK_POST} />)
    const link = screen.getByText('如何度过低谷期').closest('a')
    expect(link?.getAttribute('href')).toBe('/posts/1')
  })

  it('renders link to author profile', () => {
    render(<PostCard post={MOCK_POST} />)
    // Author name appears in both the link and the mock Avatar, use the <a> link specifically
    const links = screen.getAllByText('小明')
    const authorLink = links.find(el => el.tagName === 'A')
    expect(authorLink?.getAttribute('href')).toBe('/users/10')
  })

  it('passes size to Avatar', () => {
    render(<PostCard post={MOCK_POST} />)
    const avatar = screen.getByTestId('avatar')
    expect(avatar.getAttribute('data-size')).toBe('sm')
  })

  it('falls back to anonymous user when author_name is missing', () => {
    render(<PostCard post={{ ...MOCK_POST, author_name: undefined }} />)
    expect(screen.getByText('匿名用户')).toBeInTheDocument()
  })

  it('falls back to emotion type for unknown type', () => {
    render(<PostCard post={{ ...MOCK_POST, type: 'unknown' }} />)
    expect(screen.getByText('情绪共鸣')).toBeInTheDocument()
  })
})
