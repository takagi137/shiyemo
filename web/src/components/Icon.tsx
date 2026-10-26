/**
 * 统一线性图标
 *
 * 替代原先散落在各页面的 emoji（🤝💡📖🚀🏆👍💬…）。
 * emoji 的问题：各平台字形差异大、彩色饱和度高会盖过内容、风格偏幼态，
 * 与「安静、克制」的产品调性不符。
 *
 * 路径取自 lucide 图标集（ISC License），24×24 网格、线性描边。
 * 颜色继承 currentColor，尺寸由 size 控制。
 */

export type IconName =
  | 'emotion'
  | 'insight'
  | 'case'
  | 'action'
  | 'achievement'
  | 'all'
  | 'following'
  | 'cohort'
  | 'heart'
  | 'comment'
  | 'clock'
  | 'pin'
  | 'star'
  | 'mask'
  | 'bookmark'
  | 'share'

const PATHS: Record<IconName, string[]> = {
  // 情绪共鸣 — 心
  emotion: [
    'M19 14c1.49-1.46 3-3.21 3-5.5A5.5 5.5 0 0 0 16.5 3c-1.76 0-3 .5-4.5 2-1.5-1.5-2.74-2-4.5-2A5.5 5.5 0 0 0 2 8.5c0 2.3 1.5 4.05 3 5.5l7 7Z',
  ],
  // 认知构建 — 灯泡
  insight: [
    'M15 14c.2-1 .7-1.7 1.5-2.5 1-.9 1.5-2.2 1.5-3.5A6 6 0 0 0 6 8c0 1 .2 2.2 1.5 3.5.7.7 1.3 1.5 1.5 2.5',
    'M9 18h6',
    'M10 22h4',
  ],
  // 案例分享 — 摊开的书
  case: [
    'M12 7v14',
    'M3 18a1 1 0 0 1-1-1V4a1 1 0 0 1 1-1h5a4 4 0 0 1 4 4 4 4 0 0 1 4-4h5a1 1 0 0 1 1 1v13a1 1 0 0 1-1 1h-6a3 3 0 0 0-3 3 3 3 0 0 0-3-3z',
  ],
  // 行动发起 — 纸飞机
  action: [
    'M14.5 21.7a.5.5 0 0 0 .9 0l6.5-19a.5.5 0 0 0-.64-.64l-19 6.5a.5.5 0 0 0 0 .94l7.9 3.18a2 2 0 0 1 1.11 1.11z',
    'M21.9 2.1 10.9 13',
  ],
  // 成果展示 — 幼苗（比奖杯更贴合「学会做饭也是成果」的产品主张）
  achievement: [
    'M7 20h10',
    'M10 20c5.5-2.5.8-6.4 3-10',
    'M9.5 9.4c1.1.8 1.8 2.2 2.3 3.7-2 .4-3.5.4-4.8-.3-1.2-.6-2.3-1.9-3-4.2 2.8-.5 4.4 0 5.5.8z',
    'M14.1 6a7 7 0 0 0-1.1 4c1.9-.1 3.3-.6 4.3-1.4 1-1 1.6-2.3 1.7-4.6-2.7.1-4 1-4.9 2z',
  ],
  // 全部 — 层叠
  all: [
    'M4 6h16',
    'M4 12h16',
    'M4 18h10',
  ],
  // 关注 — 双人
  following: [
    'M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2',
    'M22 21v-2a4 4 0 0 0-3-3.87',
    'M16 3.1a4 4 0 0 1 0 7.75',
  ],
  // 同温层 — 同心圆
  cohort: [
    'M12 12m-9 0a9 9 0 1 0 18 0a9 9 0 1 0-18 0',
    'M12 12m-4 0a4 4 0 1 0 8 0a4 4 0 1 0-8 0',
  ],
  heart: [
    'M19 14c1.49-1.46 3-3.21 3-5.5A5.5 5.5 0 0 0 16.5 3c-1.76 0-3 .5-4.5 2-1.5-1.5-2.74-2-4.5-2A5.5 5.5 0 0 0 2 8.5c0 2.3 1.5 4.05 3 5.5l7 7Z',
  ],
  comment: [
    'M7.9 20A9 9 0 1 0 4 16.1L2 22z',
  ],
  clock: [
    'M12 12m-10 0a10 10 0 1 0 20 0a10 10 0 1 0-20 0',
    'M12 6v6l4 2',
  ],
  pin: [
    'M12 17v5',
    'M9 10.76a2 2 0 0 1-1.11 1.79l-1.78.9A2 2 0 0 0 5 15.24V16a1 1 0 0 0 1 1h12a1 1 0 0 0 1-1v-.76a2 2 0 0 0-1.11-1.79l-1.78-.9A2 2 0 0 1 15 10.76V7a1 1 0 0 1 1-1 2 2 0 0 0 0-4H8a2 2 0 0 0 0 4 1 1 0 0 1 1 1z',
  ],
  star: [
    'M11.5 3.1a.5.5 0 0 1 1 0l2.3 4.7 5.2.8a.5.5 0 0 1 .3.8l-3.8 3.7.9 5.2a.5.5 0 0 1-.7.5L12 16.3l-4.7 2.5a.5.5 0 0 1-.7-.5l.9-5.2-3.8-3.7a.5.5 0 0 1 .3-.8l5.2-.8z',
  ],
  mask: [
    'M12 3c-4 0-8 1-8 1v5c0 4.5 3 8.5 8 12 5-3.5 8-7.5 8-12V4s-4-1-8-1z',
    'M9 10h.01',
    'M15 10h.01',
  ],
  // 收藏 — 书签
  bookmark: [
    'M19 21l-7-5-7 5V5a2 2 0 0 1 2-2h10a2 2 0 0 1 2 2z',
  ],
  // 分享 — 三点连线
  share: [
    'M18 8m-3 0a3 3 0 1 0 6 0a3 3 0 1 0-6 0',
    'M6 12m-3 0a3 3 0 1 0 6 0a3 3 0 1 0-6 0',
    'M18 19m-3 0a3 3 0 1 0 6 0a3 3 0 1 0-6 0',
    'M8.59 13.51l6.83 3.98',
    'M15.41 6.51l-6.82 3.98',
  ],
}

interface IconProps {
  name: IconName
  size?: number
  className?: string
  /** 用于纯装饰图标时隐藏于读屏软件；有语义时传 label */
  label?: string
}

export default function Icon({ name, size = 16, className = '', label }: IconProps) {
  return (
    <svg
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth={1.75}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={`inline-block shrink-0 ${className}`}
      role={label ? 'img' : undefined}
      aria-label={label}
      aria-hidden={label ? undefined : true}
    >
      {PATHS[name].map((d) => (
        <path key={d} d={d} />
      ))}
    </svg>
  )
}
