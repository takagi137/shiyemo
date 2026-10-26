'use client'

import Link from 'next/link'

interface AvatarProps {
  userId?: number
  nickname?: string
  avatar?: string
  size?: 'sm' | 'md' | 'lg'
  clickable?: boolean
  className?: string
}

const sizeMap = {
  sm: 'w-8 h-8 text-xs',
  md: 'w-10 h-10 text-sm',
  lg: 'w-16 h-16 text-lg',
}

export default function Avatar({
  userId,
  nickname,
  avatar,
  size = 'md',
  clickable = true,
  className = '',
}: AvatarProps) {
  const initial = nickname ? nickname.charAt(0).toUpperCase() : '?'
  const sizeClass = sizeMap[size]

  const content = avatar ? (
    <img
      src={avatar}
      alt={nickname || '用户'}
      className={`${sizeClass} rounded-full object-cover ${className}`}
    />
  ) : (
    <div
      className={`${sizeClass} rounded-full bg-blue-100 text-blue-600 flex items-center justify-center font-medium ${className}`}
    >
      {initial}
    </div>
  )

  if (clickable && userId) {
    return (
      <Link href={`/users/${userId}`} className="flex-shrink-0 hover:opacity-80 transition-opacity">
        {content}
      </Link>
    )
  }

  return content
}
