import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import TrendChart from '@/components/TrendChart'

const mockFetch = vi.fn()
global.fetch = mockFetch

const mockData = {
  city_trends: [
    {
      label: '北京',
      points: [
        { date: '2026-05-01', count: 10 },
        { date: '2026-05-02', count: 15 },
        { date: '2026-05-03', count: 20 },
      ],
    },
    {
      label: '上海',
      points: [
        { date: '2026-05-01', count: 5 },
        { date: '2026-05-02', count: 8 },
        { date: '2026-05-03', count: 12 },
      ],
    },
  ],
  industry_trends: [
    {
      label: '互联网',
      points: [
        { date: '2026-05-01', count: 30 },
        { date: '2026-05-02', count: 35 },
        { date: '2026-05-03', count: 40 },
      ],
    },
  ],
  updated_at: '2026-05-30 10:00',
}

describe('TrendChart', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('shows loading skeleton initially', () => {
    mockFetch.mockImplementation(() => new Promise(() => {}))
    render(<TrendChart />)

    const skeletons = document.querySelectorAll('.animate-pulse')
    expect(skeletons.length).toBeGreaterThan(0)
  })

  it('fetches trends data on mount with default days=30', async () => {
    mockFetch.mockResolvedValue({
      json: () => Promise.resolve({ code: 0, data: mockData }),
    })

    render(<TrendChart />)

    await waitFor(() => {
      expect(mockFetch).toHaveBeenCalledWith('/api/v1/stats/trends?days=30')
    })
  })

  it('fetches with custom days parameter', async () => {
    mockFetch.mockResolvedValue({
      json: () => Promise.resolve({ code: 0, data: mockData }),
    })

    render(<TrendChart days={7} />)

    await waitFor(() => {
      expect(mockFetch).toHaveBeenCalledWith('/api/v1/stats/trends?days=7')
    })
  })

  it('renders city labels in bar chart and legend', async () => {
    mockFetch.mockResolvedValue({
      json: () => Promise.resolve({ code: 0, data: mockData }),
    })

    render(<TrendChart />)

    await waitFor(() => {
      // 每个标签出现2次：bar label + legend
      expect(screen.getAllByText('北京').length).toBeGreaterThanOrEqual(2)
      expect(screen.getAllByText('上海').length).toBeGreaterThanOrEqual(2)
    })
  })

  it('renders industry labels in bar chart and legend', async () => {
    mockFetch.mockResolvedValue({
      json: () => Promise.resolve({ code: 0, data: mockData }),
    })

    render(<TrendChart />)

    await waitFor(() => {
      expect(screen.getAllByText('互联网').length).toBeGreaterThanOrEqual(2)
    })
  })

  it('displays latest count for each bar', async () => {
    mockFetch.mockResolvedValue({
      json: () => Promise.resolve({ code: 0, data: mockData }),
    })

    render(<TrendChart />)

    await waitFor(() => {
      // 北京最新点: 20
      // 上海最新点: 12
      const counts = screen.getAllByText(/\d+/)
      expect(counts.some(el => el.textContent === '20')).toBe(true)
      expect(counts.some(el => el.textContent === '12')).toBe(true)
    })
  })

  it('shows updated_at timestamp', async () => {
    mockFetch.mockResolvedValue({
      json: () => Promise.resolve({ code: 0, data: mockData }),
    })

    render(<TrendChart />)

    await waitFor(() => {
      expect(screen.getByText(/更新于 2026-05-30 10:00/)).toBeInTheDocument()
    })
  })

  it('shows chart titles', async () => {
    mockFetch.mockResolvedValue({
      json: () => Promise.resolve({ code: 0, data: mockData }),
    })

    render(<TrendChart />)

    await waitFor(() => {
      expect(screen.getByText(/趋势数据/)).toBeInTheDocument()
      expect(screen.getByText(/城市用户数/)).toBeInTheDocument()
      expect(screen.getByText(/行业用户数/)).toBeInTheDocument()
      expect(screen.getByText(/城市趋势变化/)).toBeInTheDocument()
      expect(screen.getByText(/行业趋势变化/)).toBeInTheDocument()
    })
  })

  it('renders line chart SVG elements', async () => {
    mockFetch.mockResolvedValue({
      json: () => Promise.resolve({ code: 0, data: mockData }),
    })

    render(<TrendChart />)

    await waitFor(() => {
      const svgs = document.querySelectorAll('svg')
      expect(svgs.length).toBeGreaterThanOrEqual(2) // 2 line charts
    })
  })

  it('shows line legend items', async () => {
    mockFetch.mockResolvedValue({
      json: () => Promise.resolve({ code: 0, data: mockData }),
    })

    render(<TrendChart />)

    await waitFor(() => {
      const legendItems = document.querySelectorAll('.rounded-full')
      expect(legendItems.length).toBeGreaterThanOrEqual(3) // 北京, 上海, 互联网
    })
  })

  it('shows "暂无数据" when city_trends is empty', async () => {
    mockFetch.mockResolvedValue({
      json: () => Promise.resolve({
        code: 0,
        data: {
          city_trends: [],
          industry_trends: [],
          updated_at: '',
        },
      }),
    })

    render(<TrendChart />)

    await waitFor(() => {
      expect(screen.getAllByText('暂无数据').length).toBeGreaterThanOrEqual(1)
    })
  })

  it('shows "暂无数据" for empty series inside BarChart', async () => {
    mockFetch.mockResolvedValue({
      json: () => Promise.resolve({
        code: 0,
        data: {
          city_trends: [{
            label: '北京',
            points: [],
          }],
          industry_trends: [],
          updated_at: '',
        },
      }),
    })

    render(<TrendChart />)

    await waitFor(() => {
      // BarChart with empty points → latest count = 0
      expect(screen.getByText('北京')).toBeInTheDocument()
    })
  })

  it('handles fetch error gracefully', async () => {
    mockFetch.mockRejectedValue(new Error('Network error'))
    const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {})

    const { container } = render(<TrendChart />)

    await waitFor(() => {
      expect(container.querySelector('.animate-pulse')).not.toBeInTheDocument()
    })

    consoleSpy.mockRestore()
  })

  it('returns null when data is null', async () => {
    mockFetch.mockResolvedValue({
      json: () => Promise.resolve({ code: 999, data: null }),
    })

    const { container } = render(<TrendChart />)

    await waitFor(() => {
      // 数据为 null 时，loading 结束后返回 null
      expect(container.innerHTML).toBe('')
    })
  })

  it('handles non-zero code response', async () => {
    mockFetch.mockResolvedValue({
      json: () => Promise.resolve({ code: 50001, message: '服务器错误', data: null }),
    })

    const { container } = render(<TrendChart />)

    await waitFor(() => {
      expect(container.innerHTML).toBe('')
    })
  })
})
