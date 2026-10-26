import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { render, screen, waitFor, act, fireEvent } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import AmbientPlayer from '@/components/AmbientPlayer'

// ============================================================
// Mocks
// ============================================================

// fetch mock —— 与 DailyTip.test.tsx 同款
const mockFetch = vi.fn()
global.fetch = mockFetch as unknown as typeof fetch

// Audio mock —— jsdom 没真实音频实现，需模拟
// 必须支持：play() 返回 Promise（成功或失败）、pause()、volume 赋值、loop 赋值
let lastAudio: MockAudio | null = null
class MockAudio {
  url: string
  volume = 1
  loop = false
  paused = true
  constructor(url: string) {
    this.url = url
    lastAudio = this
  }
  play(): Promise<void> {
    this.paused = false
    return mockPlayResult(this)
  }
  pause(): void {
    this.paused = true
  }
  // 兼容 audioRef.current 判空
}
type PlayBehaviour = 'resolve' | 'reject'
const mockPlayResult = vi.fn<(a: MockAudio) => Promise<void>>()
;(globalThis as any).Audio = MockAudio as unknown as typeof Audio

function setPlayResult(behaviour: PlayBehaviour): void {
  mockPlayResult.mockReset()
  mockPlayResult.mockImplementation(() => {
    if (behaviour === 'reject') return Promise.reject(new Error('load failed'))
    return Promise.resolve()
  })
}

const SAMPLE_SOUNDS = [
  { key: 'rain', name: '雨声', emoji: '🌧', url: 'https://example.com/rain.mp3' },
  { key: 'fire', name: '壁炉', emoji: '🔥', url: 'https://example.com/fire.mp3' },
  { key: 'wave', name: '海浪', emoji: '🌊', url: 'https://example.com/wave.mp3' },
]

function mockSoundsResponse(): ReturnType<typeof Promise.resolve> {
  return Promise.resolve({
    json: () => Promise.resolve({ code: 0, data: { items: SAMPLE_SOUNDS } }),
  })
}

// ============================================================
// Setup
// ============================================================

beforeEach(() => {
  mockFetch.mockReset()
  setPlayResult('resolve')
  localStorage.clear()
  lastAudio = null
})

afterEach(() => {
  vi.restoreAllMocks()
})

// ============================================================
// Tests
// ============================================================

describe('AmbientPlayer', () => {
  describe('拉取音效列表', () => {
    it('调用 /api/v1/sounds 端点', async () => {
      mockFetch.mockResolvedValueOnce(mockSoundsResponse())
      render(<AmbientPlayer />)
      await waitFor(() => {
        expect(mockFetch).toHaveBeenCalledWith('/api/v1/sounds')
      })
    })

    it('成功返回时为每个音效渲染按钮', async () => {
      mockFetch.mockResolvedValueOnce(mockSoundsResponse())
      render(<AmbientPlayer />)

      // 用 name 文本找按钮——按钮里包含 emoji + name
      await waitFor(() => {
        expect(screen.getByText('雨声')).toBeInTheDocument()
      })
      expect(screen.getByText('壁炉')).toBeInTheDocument()
      expect(screen.getByText('海浪')).toBeInTheDocument()
    })

    it('拉取失败时不抛错到 console（被 catch 接住）', async () => {
      const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {})
      mockFetch.mockRejectedValueOnce(new Error('network'))
      render(<AmbientPlayer />)
      // 等 fetch settle
      await waitFor(() => {
        expect(mockFetch).toHaveBeenCalled()
      })
      // 不能因为失败就崩；按钮区为空（不渲染任何音效按钮的 name）
      expect(screen.queryByText('雨声')).not.toBeInTheDocument()
      // catch 内调用了 console.error，避免误以为是未捕获的异常
      expect(consoleSpy).toHaveBeenCalled()
    })

    it('响应 code != 0 时不渲染按钮', async () => {
      mockFetch.mockResolvedValueOnce({
        json: () => Promise.resolve({ code: 1, message: 'fail', data: null }),
      })
      render(<AmbientPlayer />)
      await waitFor(() => expect(mockFetch).toHaveBeenCalled())
      expect(screen.queryByText('雨声')).not.toBeInTheDocument()
    })

    it('data.items 非数组时不渲染按钮', async () => {
      mockFetch.mockResolvedValueOnce({
        json: () => Promise.resolve({ code: 0, data: { items: 'oops' } }),
      })
      render(<AmbientPlayer />)
      await waitFor(() => expect(mockFetch).toHaveBeenCalled())
      expect(screen.queryByText('雨声')).not.toBeInTheDocument()
    })
  })

  describe('音量持久化', () => {
    it('localStorage 已有音量时恢复到 range value', async () => {
      localStorage.setItem('sanctuary-volume', '37')
      mockFetch.mockResolvedValueOnce(mockSoundsResponse())
      render(<AmbientPlayer />)

      const range = screen.getByRole('slider') as HTMLInputElement
      await waitFor(() => {
        expect(Number(range.value)).toBe(37)
      })
    })

    it('改变音量时写入 localStorage', async () => {
      mockFetch.mockResolvedValueOnce(mockSoundsResponse())
      render(<AmbientPlayer />)

      await waitFor(() => expect(screen.getByText('雨声')).toBeInTheDocument())
      const range = screen.getByRole('slider') as HTMLInputElement

      // 必须用 React 的 fireEvent.change —— 走 React onChange 合成系统
      await act(async () => {
        fireEvent.change(range, { target: { value: '80' } })
      })

      expect(localStorage.getItem('sanctuary-volume')).toBe('80')
    })
  })

  describe('播放控制', () => {
    it('点击音效 → 创建 Audio 对象并标记 playing', async () => {
      const user = userEvent.setup()
      mockFetch.mockResolvedValueOnce(mockSoundsResponse())
      render(<AmbientPlayer />)
      await waitFor(() => expect(screen.getByText('雨声')).toBeInTheDocument())

      await user.click(screen.getByText('雨声').closest('button')!)

      await waitFor(() => {
        expect(lastAudio).not.toBeNull()
        expect(lastAudio!.url).toBe('https://example.com/rain.mp3')
        expect(lastAudio!.loop).toBe(true)
      })
      expect(screen.getByText(/正在播放.*雨声/)).toBeInTheDocument()
      expect(localStorage.getItem('sanctuary-playing')).toBe('rain')
    })

    it('再次点击同一音效 → 暂停（不创建新 Audio）', async () => {
      const user = userEvent.setup()
      mockFetch.mockResolvedValueOnce(mockSoundsResponse())
      render(<AmbientPlayer />)
      await waitFor(() => expect(screen.getByText('雨声')).toBeInTheDocument())

      const rainBtn = screen.getByText('雨声').closest('button')!
      await user.click(rainBtn)
      await waitFor(() => expect(lastAudio).not.toBeNull())
      const first = lastAudio!

      await user.click(rainBtn)
      expect(first.paused).toBe(true)
      // "正在播放" 文字应消失
      expect(screen.queryByText(/正在播放/)).not.toBeInTheDocument()
    })

    it('切换到另一个音效 → 暂停当前并新建', async () => {
      const user = userEvent.setup()
      mockFetch.mockResolvedValueOnce(mockSoundsResponse())
      render(<AmbientPlayer />)
      await waitFor(() => expect(screen.getByText('雨声')).toBeInTheDocument())

      await user.click(screen.getByText('雨声').closest('button')!)
      await waitFor(() => expect(lastAudio!.url).toBe('https://example.com/rain.mp3'))
      const first = lastAudio!

      await user.click(screen.getByText('壁炉').closest('button')!)
      await waitFor(() => expect(lastAudio).not.toBe(first))
      expect(first.paused).toBe(true)
      expect(lastAudio!.url).toBe('https://example.com/fire.mp3')
      expect(screen.getByText(/正在播放.*壁炉/)).toBeInTheDocument()
    })

    it('Audio.play 抛错 → 显示"暂时放不出来"且不写 sanctuary-playing', async () => {
      const user = userEvent.setup()
      setPlayResult('reject')
      mockFetch.mockResolvedValueOnce(mockSoundsResponse())
      render(<AmbientPlayer />)
      await waitFor(() => expect(screen.getByText('雨声')).toBeInTheDocument())

      await user.click(screen.getByText('雨声').closest('button')!)

      await waitFor(() => {
        expect(screen.getByText(/暂时放不出来/)).toBeInTheDocument()
      })
      // 正在播放 不应出现
      expect(screen.queryByText(/正在播放/)).not.toBeInTheDocument()
      // 失败时移除 localStorage（防御性）
      expect(localStorage.getItem('sanctuary-playing')).toBeNull()
    })

    it('改变音量时同步到 audioRef.current.volume', async () => {
      const user = userEvent.setup()
      mockFetch.mockResolvedValueOnce(mockSoundsResponse())
      render(<AmbientPlayer />)
      await waitFor(() => expect(screen.getByText('雨声')).toBeInTheDocument())

      await user.click(screen.getByText('雨声').closest('button')!)
      await waitFor(() => expect(lastAudio).not.toBeNull())
      const audio = lastAudio!
      // 初始音量是 range.value (50) / 100 = 0.5
      expect(audio.volume).toBeCloseTo(0.5, 2)

      const range = screen.getByRole('slider') as HTMLInputElement
      await act(async () => {
        fireEvent.change(range, { target: { value: '20' } })
      })
      expect(audio.volume).toBeCloseTo(0.2, 2)
    })
  })

  describe('卸载清理', () => {
    it('组件卸载时调用 audioRef.current.pause()', async () => {
      mockFetch.mockResolvedValueOnce(mockSoundsResponse())
      const { unmount } = render(<AmbientPlayer />)
      await waitFor(() => expect(screen.getByText('雨声')).toBeInTheDocument())

      const user = userEvent.setup()
      await user.click(screen.getByText('雨声').closest('button')!)
      await waitFor(() => expect(lastAudio).not.toBeNull())
      const pauseSpy = vi.spyOn(lastAudio!, 'pause')

      unmount()
      expect(pauseSpy).toHaveBeenCalled()
    })
  })
})