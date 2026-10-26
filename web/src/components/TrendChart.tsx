'use client'

import { useState, useEffect } from 'react'

interface TrendPoint {
  date: string
  count: number
}

interface TrendSeries {
  label: string
  points: TrendPoint[]
}

interface TrendResponse {
  city_trends: TrendSeries[]
  industry_trends: TrendSeries[]
  updated_at: string
}

const COLORS = [
  '#3B82F6', '#8B5CF6', '#F59E0B', '#10B981', '#EF4444',
  '#06B6D4', '#F97316', '#84CC16', '#EC4899', '#6366F1',
]

// 简易 SVG 柱状图
function BarChart({ data, title }: { data: TrendSeries[]; title: string }) {
  if (!data || data.length === 0) {
    return (
      <div className="card p-6">
        <h3 className="text-sm font-medium text-gray-500 mb-4">{title}</h3>
        <div className="text-center text-gray-400 py-8">暂无数据</div>
      </div>
    )
  }

  // 取每个系列最新数据点
  const latestData = data.map(s => ({
    label: s.label,
    count: s.points.length > 0 ? s.points[s.points.length - 1].count : 0,
  }))

  const maxCount = Math.max(...latestData.map(d => d.count), 1)
  const barHeight = 32
  const chartHeight = latestData.length * (barHeight + 12) + 20

  return (
    <div className="card p-6">
      <h3 className="text-sm font-medium text-gray-500 mb-4">{title}</h3>
      <div style={{ height: chartHeight }}>
        {latestData.map((item, i) => {
          const width = (item.count / maxCount) * 100
          return (
            <div key={item.label} className="flex items-center gap-3 mb-3">
              <div className="w-20 text-xs text-gray-600 text-right truncate">{item.label}</div>
              <div className="flex-1 relative">
                <div
                  className="h-6 rounded-md transition-all duration-500"
                  style={{
                    width: `${width}%`,
                    backgroundColor: COLORS[i % COLORS.length],
                    minWidth: '24px',
                  }}
                />
                <span className="absolute right-0 top-0 h-6 flex items-center pr-2 text-xs font-medium text-gray-700">
                  {item.count}
                </span>
              </div>
            </div>
          )
        })}
      </div>
    </div>
  )
}

// 简易 SVG 折线图
function LineChart({ series, title }: { series: TrendSeries[]; title: string }) {
  if (!series || series.length === 0) {
    return (
      <div className="card p-6">
        <h3 className="text-sm font-medium text-gray-500 mb-4">{title}</h3>
        <div className="text-center text-gray-400 py-8">暂无数据</div>
      </div>
    )
  }

  const svgWidth = 500
  const svgHeight = 200
  const padding = { top: 20, right: 20, bottom: 30, left: 40 }
  const chartWidth = svgWidth - padding.left - padding.right
  const chartHeight = svgHeight - padding.top - padding.bottom

  // 收集所有日期
  const allDates = new Set<string>()
  series.forEach(s => s.points.forEach(p => allDates.add(p.date)))
  const dates = Array.from(allDates).sort()

  if (dates.length === 0) return null

  const maxCount = Math.max(...series.flatMap(s => s.points.map(p => p.count)), 1)

  const xScale = (i: number) => padding.left + (i / (dates.length - 1 || 1)) * chartWidth
  const yScale = (v: number) => padding.top + chartHeight - (v / maxCount) * chartHeight

  return (
    <div className="card p-6">
      <h3 className="text-sm font-medium text-gray-500 mb-4">{title}</h3>
      <svg viewBox={`0 0 ${svgWidth} ${svgHeight}`} className="w-full h-auto">
        {/* 网格线 */}
        {[0, 0.25, 0.5, 0.75, 1].map(ratio => {
          const y = padding.top + chartHeight * (1 - ratio)
          return (
            <g key={ratio}>
              <line x1={padding.left} y1={y} x2={svgWidth - padding.right} y2={y} stroke="#E5E7EB" strokeWidth="1" />
              <text x={padding.left - 8} y={y + 4} textAnchor="end" className="text-[10px] fill-gray-400">
                {Math.round(maxCount * ratio)}
              </text>
            </g>
          )
        })}

        {/* 数据线 */}
        {series.map((s, si) => {
          const points = dates.map((date, i) => {
            const point = s.points.find(p => p.date === date)
            return { x: xScale(i), y: yScale(point?.count || 0) }
          })
          const pathD = points.map((p, i) => `${i === 0 ? 'M' : 'L'} ${p.x} ${p.y}`).join(' ')
          return (
            <g key={s.label}>
              <path d={pathD} fill="none" stroke={COLORS[si % COLORS.length]} strokeWidth="2" />
              {points.map((p, i) => (
                <circle key={i} cx={p.x} cy={p.y} r="3" fill={COLORS[si % COLORS.length]} />
              ))}
            </g>
          )
        })}

        {/* X 轴日期 */}
        {dates.map((date, i) => {
          if (i % Math.ceil(dates.length / 5) !== 0 && i !== dates.length - 1) return null
          return (
            <text key={date} x={xScale(i)} y={svgHeight - 5} textAnchor="middle" className="text-[10px] fill-gray-400">
              {date.slice(5)}
            </text>
          )
        })}
      </svg>

      {/* 图例 */}
      <div className="flex flex-wrap gap-3 mt-4">
        {series.map((s, i) => (
          <div key={s.label} className="flex items-center gap-1.5 text-xs text-gray-600">
            <div className="w-3 h-3 rounded-full" style={{ backgroundColor: COLORS[i % COLORS.length] }} />
            {s.label}
          </div>
        ))}
      </div>
    </div>
  )
}

interface TrendChartProps {
  days?: number
}

export default function TrendChart({ days = 30 }: TrendChartProps) {
  const [data, setData] = useState<TrendResponse | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    const fetchTrends = async () => {
      setLoading(true)
      try {
        const res = await fetch(`/api/v1/stats/trends?days=${days}`)
        const json = await res.json()
        if (json.code === 0) {
          setData(json.data)
        }
      } catch (err) {
        console.error('获取趋势数据失败:', err)
      } finally {
        setLoading(false)
      }
    }
    fetchTrends()
  }, [days])

  if (loading) {
    return (
      <div className="space-y-4">
        {[...Array(2)].map((_, i) => (
          <div key={i} className="card p-6 animate-pulse">
            <div className="h-4 bg-gray-200 rounded w-32 mb-4" />
            <div className="h-40 bg-gray-100 rounded" />
          </div>
        ))}
      </div>
    )
  }

  if (!data) return null

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h2 className="text-lg font-bold text-gray-900">📊 趋势数据</h2>
        {data.updated_at && (
          <span className="text-xs text-gray-400">更新于 {data.updated_at}</span>
        )}
      </div>

      <BarChart data={data.city_trends} title="🏙️ 各城市用户数" />
      <BarChart data={data.industry_trends} title="💼 各行业用户数" />
      <LineChart series={data.city_trends} title="📈 城市趋势变化" />
      <LineChart series={data.industry_trends} title="📈 行业趋势变化" />
    </div>
  )
}
