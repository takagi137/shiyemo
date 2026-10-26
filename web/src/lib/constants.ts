// ============================================================
// constants.ts — 全局常量
// ============================================================

import type { IconName } from '@/components/Icon'

// 帖子类型配置
//
// 配色说明：5 种类型统一在同一明度与彩度上，仅靠色相区分（见 globals.css 的
// --color-*-soft/line/ink 令牌）。原方案用的是 rose/amber/blue/emerald/yellow
// 五个高饱和撞色，并排出现在信息流里会互相争夺注意力，且盖过帖子内容本身。
//
// icon 为线性图标名（见 components/Icon.tsx）。
// emoji 字段暂时保留：尚有页面直接引用，待全站替换完成后可移除。
export const TYPE_CONFIGS: Record<string, {
  label: string
  icon: IconName
  emoji: string
  color: string
}> = {
  emotion: {
    label: '情绪共鸣', icon: 'emotion', emoji: '🤝',
    color: 'bg-emotion-soft text-emotion-ink border-emotion-line',
  },
  insight: {
    label: '认知构建', icon: 'insight', emoji: '💡',
    color: 'bg-insight-soft text-insight-ink border-insight-line',
  },
  case: {
    label: '案例分享', icon: 'case', emoji: '📖',
    color: 'bg-case-soft text-case-ink border-case-line',
  },
  action: {
    label: '行动发起', icon: 'action', emoji: '🚀',
    color: 'bg-action-soft text-action-ink border-action-line',
  },
  achievement: {
    label: '成果展示', icon: 'achievement', emoji: '🏆',
    color: 'bg-achievement-soft text-achievement-ink border-achievement-line',
  },
}

export const TYPE_LABELS = TYPE_CONFIGS

export const TYPE_COLORS: Record<string, string> = Object.fromEntries(
  Object.entries(TYPE_CONFIGS).map(([k, v]) => [k, v.color])
)

export const TYPE_BORDER_COLORS: Record<string, string> = Object.fromEntries(
  Object.entries(TYPE_CONFIGS).map(([k, v]) => [k, v.color])
)

/** Feed / 搜索页的筛选按钮列表 */
export const TYPE_FILTER_OPTIONS: { key: string; label: string; icon: IconName; emoji: string }[] = [
  { key: '',          label: '全部',   icon: 'all',       emoji: '📋' },
  { key: 'following', label: '关注',   icon: 'following', emoji: '👥' },
  { key: 'cohort',    label: '同温层', icon: 'cohort',    emoji: '🔥' },
  ...Object.entries(TYPE_CONFIGS).map(([key, { label, icon, emoji }]) => ({ key, label, icon, emoji })),
]

/** 发帖页类型选择 */
export const POST_CREATE_TYPES = Object.entries(TYPE_CONFIGS).map(([key, v]) => ({
  key,
  ...v,
  desc: {
    emotion:     '想吐槽、想倾诉、想说"我今天很丧"，都可以',
    insight:     '踩过的坑、想通的事、对你有用的认知转变',
    case:        '面试经历、转行故事、失败教训，真实就好',
    action:      '想组队做点什么？找搭子、找方向，从这里开始',
    achievement: '学会做饭、跑完5公里、第一笔收入……任何成长都值得记录',
  }[key],
}))

// 心情类型配置
export const MOOD_OPTIONS = [
  { key: 'happy',   emoji: '😊', label: '开心', color: 'bg-insight-soft text-insight-ink border-insight-line' },
  { key: 'calm',    emoji: '😌', label: '平静', color: 'bg-case-soft text-case-ink border-case-line' },
  { key: 'sad',     emoji: '😢', label: '难过', color: 'bg-surface-sunk text-ink-soft border-line' },
  { key: 'angry',   emoji: '😠', label: '生气', color: 'bg-emotion-soft text-emotion-ink border-emotion-line' },
  { key: 'anxious', emoji: '😰', label: '焦虑', color: 'bg-achievement-soft text-achievement-ink border-achievement-line' },
]
