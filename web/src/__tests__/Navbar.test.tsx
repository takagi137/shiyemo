import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import Navbar from '@/components/Navbar'

// mock next/link
vi.mock('next/link', () => {
  return {
    default: ({ children, href, ...props }: any) => (
      <a href={href} {...props}>{children}</a>
    ),
  }
})

describe('Navbar', () => {
  it('renders brand name', () => {
    render(<Navbar />)
    expect(screen.getByText('失业么')).toBeInTheDocument()
  })

  it('renders brand link to home', () => {
    render(<Navbar />)
    const link = screen.getByText('失业么').closest('a')
    expect(link?.getAttribute('href')).toBe('/')
  })

  it('renders desktop navigation links', () => {
    render(<Navbar />)
    expect(screen.getByText('动态')).toBeInTheDocument()
    expect(screen.getByText('发帖')).toBeInTheDocument()
    expect(screen.getByText('打卡')).toBeInTheDocument()
    expect(screen.getByText('同温层')).toBeInTheDocument()
    expect(screen.getByText('搜索')).toBeInTheDocument()
    expect(screen.getByText('反馈')).toBeInTheDocument()
  })

  it('renders login and register buttons', () => {
    render(<Navbar />)
    expect(screen.getByText('登录')).toBeInTheDocument()
    expect(screen.getByText('注册')).toBeInTheDocument()
  })

  it('mobile menu button toggles mobile menu', async () => {
    const user = userEvent.setup()
    render(<Navbar />)

    // Mobile toggle button is the only button in Navbar
    const toggleButton = screen.getByRole('button')
    expect(toggleButton).toBeInTheDocument()

    // Click to open mobile menu
    await user.click(toggleButton)

    // After click, mobile menu should render extra nav links (desktop + mobile)
    const dynamicLinks = screen.getAllByText('动态')
    expect(dynamicLinks.length).toBeGreaterThanOrEqual(2)
  })

  it('renders all navigation link hrefs correctly', () => {
    render(<Navbar />)
    const links = screen.getAllByRole('link')
    const hrefs = links.map(l => (l as HTMLAnchorElement).getAttribute('href'))

    expect(hrefs).toContain('/')
    expect(hrefs).toContain('/feed')
    expect(hrefs).toContain('/posts/create')
    expect(hrefs).toContain('/checkin')
    expect(hrefs).toContain('/cohort')
    expect(hrefs).toContain('/search')
    expect(hrefs).toContain('/feedback')
    expect(hrefs).toContain('/login')
    expect(hrefs).toContain('/register')
  })
})
