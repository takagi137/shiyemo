'use client'

import { useState, useEffect } from 'react'
import Link from 'next/link'

interface UserInfo {
  id: number
  nickname: string
  avatar_url?: string
}

export default function Navbar() {
  const [isMenuOpen, setIsMenuOpen] = useState(false)
  const [user, setUser] = useState<UserInfo | null>(null)

  useEffect(() => {
    const stored = localStorage.getItem('user')
    if (stored) {
      try { setUser(JSON.parse(stored)) } catch { /* ignore */ }
    }
  }, [])

  const handleLogout = () => {
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    setUser(null)
  }

  return (
    <nav className="bg-surface/85 backdrop-blur-sm border-b border-line sticky top-0 z-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between h-16">
          {/* Logo */}
          <div className="flex items-center">
            <Link href="/" className="flex-shrink-0">
              <span className="text-lg font-semibold tracking-tight text-accent">失业么</span>
            </Link>
          </div>

          {/* Desktop Navigation */}
          <div className="hidden md:flex md:items-center md:space-x-8">
            <Link href="/feed" className="text-ink-soft hover:text-ink px-3 py-2 text-sm font-medium transition-colors">
              动态
            </Link>
            <Link href="/posts/create" className="text-ink-soft hover:text-ink px-3 py-2 text-sm font-medium transition-colors">
              发帖
            </Link>
            <Link href="/checkin" className="text-ink-soft hover:text-ink px-3 py-2 text-sm font-medium transition-colors">
              打卡
            </Link>
            <Link href="/cohort" className="text-ink-soft hover:text-ink px-3 py-2 text-sm font-medium transition-colors">
              同温层
            </Link>
            <Link href="/search" className="text-ink-soft hover:text-ink px-3 py-2 text-sm font-medium transition-colors">
              搜索
            </Link>
            <Link href="/feedback" className="text-ink-soft hover:text-ink px-3 py-2 text-sm font-medium transition-colors">
              反馈
            </Link>
            <Link href="/sanctuary" className="text-ink-soft hover:text-ink px-3 py-2 text-sm font-medium transition-colors">
              静思
            </Link>
          </div>

          {/* Desktop Auth */}
          <div className="hidden md:flex md:items-center md:space-x-4">
            {user ? (
              <>
                <Link href="/profile" className="flex items-center gap-2 text-sm text-ink-soft hover:text-ink transition-colors">
                  <div className="w-8 h-8 rounded-full bg-accent-soft flex items-center justify-center text-accent font-medium text-sm">
                    {user.nickname.charAt(0)}
                  </div>
                  <span>{user.nickname}</span>
                </Link>
                <Link href="/feed" className="text-ink-soft hover:text-ink px-3 py-2 text-sm font-medium transition-colors">
                  动态
                </Link>
                <Link href="/bookmarks" className="text-ink-soft hover:text-ink px-3 py-2 text-sm font-medium transition-colors">
                  收藏
                </Link>
                <button
                  onClick={handleLogout}
                  className="text-ink-faint hover:text-ink px-3 py-2 text-sm font-medium transition-colors"
                >
                  退出
                </button>
              </>
            ) : (
              <>
                <Link href="/login" className="text-ink-soft hover:text-ink px-4 py-2 text-sm font-medium transition-colors">
                  登录
                </Link>
                <Link href="/register" className="btn-primary text-sm px-4 py-2">
                  注册
                </Link>
              </>
            )}
          </div>

          {/* Mobile Menu Button */}
          <div className="flex items-center md:hidden">
            <button
              onClick={() => setIsMenuOpen(!isMenuOpen)}
              className="inline-flex items-center justify-center p-2 rounded-md text-ink-soft hover:text-ink hover:bg-surface-sunk focus:outline-none focus:ring-2 focus:ring-inset focus:ring-accent"
              aria-expanded={isMenuOpen}
            >
              <span className="sr-only">打开主菜单</span>
              {isMenuOpen ? (
                <svg className="block h-6 w-6" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                </svg>
              ) : (
                <svg className="block h-6 w-6" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" />
                </svg>
              )}
            </button>
          </div>
        </div>
      </div>

      {/* Mobile Menu */}
      {isMenuOpen && (
        <div className="md:hidden">
          <div className="px-2 pt-2 pb-3 space-y-1 sm:px-3">
            <Link href="/feed" className="block px-3 py-2 rounded-md text-base font-medium text-ink-soft hover:text-ink hover:bg-surface-sunk transition-colors">
              动态
            </Link>
            <Link href="/posts/create" className="block px-3 py-2 rounded-md text-base font-medium text-ink-soft hover:text-ink hover:bg-surface-sunk transition-colors">
              发帖
            </Link>
            <Link href="/checkin" className="block px-3 py-2 rounded-md text-base font-medium text-ink-soft hover:text-ink hover:bg-surface-sunk transition-colors">
              打卡
            </Link>
            <Link href="/cohort" className="block px-3 py-2 rounded-md text-base font-medium text-ink-soft hover:text-ink hover:bg-surface-sunk transition-colors">
              同温层
            </Link>
            <Link href="/search" className="block px-3 py-2 rounded-md text-base font-medium text-ink-soft hover:text-ink hover:bg-surface-sunk transition-colors">
              搜索
            </Link>
            <Link href="/sanctuary" className="block px-3 py-2 rounded-md text-base font-medium text-ink-soft hover:text-ink hover:bg-surface-sunk transition-colors">
              静思空间
            </Link>
          </div>
          <div className="pt-4 pb-3 border-t border-line">
            {user ? (
              <div className="px-5 space-y-3">
                <div className="flex items-center gap-3">
                  <Link href="/profile" className="flex items-center gap-2 text-sm font-medium text-ink-soft">
                    <div className="w-8 h-8 rounded-full bg-accent-soft flex items-center justify-center text-accent font-medium text-sm">
                      {user.nickname.charAt(0)}
                    </div>
                    <span>{user.nickname}</span>
                  </Link>
                </div>
                <div className="flex gap-2">
                  <Link href="/feed" className="flex-1 text-center text-ink-soft hover:text-ink px-4 py-2 text-sm font-medium transition-colors">
                    动态
                  </Link>
                  <Link href="/bookmarks" className="flex-1 text-center text-ink-soft hover:text-ink px-4 py-2 text-sm font-medium transition-colors">
                    收藏
                  </Link>
                </div>
                <button
                  onClick={handleLogout}
                  className="w-full text-center text-ink-soft hover:text-ink px-4 py-2 text-sm font-medium transition-colors border border-line rounded-lg"
                >
                  退出登录
                </button>
              </div>
            ) : (
              <div className="flex items-center px-5 space-x-3">
                <Link href="/login" className="flex-1 text-center text-ink-soft hover:text-ink px-4 py-2 text-sm font-medium transition-colors border border-line rounded-lg">
                  登录
                </Link>
                <Link href="/register" className="flex-1 text-center btn-primary text-sm px-4 py-2">
                  注册
                </Link>
              </div>
            )}
          </div>
        </div>
      )}
    </nav>
  )
}
