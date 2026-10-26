import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import DailyTip from '@/components/DailyTip'

// mock fetch globally
const mockFetch = vi.fn()
global.fetch = mockFetch

describe('DailyTip', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('shows loading skeleton initially', () => {
    mockFetch.mockImplementation(() => new Promise(() => {})) // never resolves
    render(<DailyTip />)

    expect(document.querySelector('.animate-pulse')).toBeInTheDocument()
  })

  it('renders tip content after fetch', async () => {
    mockFetch.mockResolvedValue({
      json: () => Promise.resolve({
        code: 0,
        data: { id: 1, content: '保持积极心态很重要', category: 'mindset' },
      }),
    })

    render(<DailyTip />)

    await waitFor(() => {
      expect(screen.getByText('保持积极心态很重要')).toBeInTheDocument()
    })
  })

  it('shows correct category icon and name', async () => {
    mockFetch.mockResolvedValue({
      json: () => Promise.resolve({
        code: 0,
        data: { id: 1, content: '内容', category: 'finance' },
      }),
    })

    render(<DailyTip />)

    await waitFor(() => {
      expect(screen.getByText(/💰/)).toBeInTheDocument()
      expect(screen.getByText(/财务/)).toBeInTheDocument()
    })
  })

  it('shows default emoji for unknown category', async () => {
    mockFetch.mockResolvedValue({
      json: () => Promise.resolve({
        code: 0,
        data: { id: 1, content: '内容', category: 'unknown' },
      }),
    })

    render(<DailyTip />)

    await waitFor(() => {
      expect(screen.getByText(/💡/)).toBeInTheDocument()
      expect(screen.getByText(/通用/)).toBeInTheDocument()
    })
  })

  it('refresh button triggers re-fetch', async () => {
    const user = userEvent.setup()

    mockFetch
      .mockResolvedValueOnce({
        json: () => Promise.resolve({
          code: 0,
          data: { id: 1, content: '第一条贴士', category: 'life' },
        }),
      })
      .mockResolvedValueOnce({
        json: () => Promise.resolve({
          code: 0,
          data: { id: 2, content: '第二条贴士', category: 'policy' },
        }),
      })

    render(<DailyTip />)

    await waitFor(() => {
      expect(screen.getByText('第一条贴士')).toBeInTheDocument()
    })

    await user.click(screen.getByTitle('换一条'))

    await waitFor(() => {
      expect(screen.getByText('第二条贴士')).toBeInTheDocument()
    })

    expect(mockFetch).toHaveBeenCalledTimes(2)
  })

  it('returns null when response has no data', async () => {
    mockFetch.mockResolvedValue({
      json: () => Promise.resolve({ code: 0, data: null }),
    })

    const { container } = render(<DailyTip />)

    await waitFor(() => {
      expect(container.innerHTML).toBe('')
    })
  })

  it('handles fetch error gracefully', async () => {
    mockFetch.mockRejectedValue(new Error('Network error'))
    const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {})

    const { container } = render(<DailyTip />)

    await waitFor(() => {
      expect(container.innerHTML).toBe('')
    })

    consoleSpy.mockRestore()
  })

  it('fetches from correct endpoint', async () => {
    mockFetch.mockResolvedValue({
      json: () => Promise.resolve({ code: 0, data: null }),
    })

    render(<DailyTip />)

    await waitFor(() => {
      expect(mockFetch).toHaveBeenCalledWith('/api/v1/tips/random')
    })
  })
})
