'use client'

import { useState, useEffect, useRef, useCallback } from 'react'

type Phase = 'inhale' | 'hold' | 'exhale'

// 采用「共振式呼吸 / coherent breathing」标准：约 6 次/分钟，吸气 4s、呼气 6s，
// 呼气长于吸气以激活副交感神经、帮助放松；不设满肺屏息 —— 原 4-4-6 里的 4s
// 满肺 hold（且呼气后无空肺停顿）是「窒息感」的主因，去掉后呼吸像平缓的波浪。
const PHASES: { phase: Phase; label: string; duration: number }[] = [
  { phase: 'inhale', label: '吸气...', duration: 4000 },
  { phase: 'exhale', label: '呼气...', duration: 6000 },
]

export default function BreathingCircle() {
  const [isActive, setIsActive] = useState(false)
  const [currentPhase, setCurrentPhase] = useState(0)
  const [scale, setScale] = useState(1)
  const [cycles, setCycles] = useState(0)
  const timerRef = useRef<NodeJS.Timeout | null>(null)

  const runPhase = useCallback((phaseIndex: number) => {
    const { phase, duration } = PHASES[phaseIndex]

    // 设置目标缩放
    if (phase === 'inhale') {
      setScale(1.6)
    } else if (phase === 'exhale') {
      setScale(1)
    }
    // hold 保持当前缩放

    timerRef.current = setTimeout(() => {
      const nextIndex = (phaseIndex + 1) % PHASES.length
      setCurrentPhase(nextIndex)
      if (nextIndex === 0) {
        setCycles(c => c + 1)
      }
      runPhase(nextIndex)
    }, duration)
  }, [])

  useEffect(() => {
    if (isActive) {
      setCurrentPhase(0)
      runPhase(0)
    }

    return () => {
      if (timerRef.current) {
        clearTimeout(timerRef.current)
        timerRef.current = null
      }
    }
  }, [isActive, runPhase])

  const toggle = () => {
    setIsActive(!isActive)
    if (isActive) {
      setScale(1)
      setCurrentPhase(0)
    }
  }

  return (
    <div className="flex flex-col items-center space-y-6">
      <h3 className="text-lg font-medium text-ink">呼吸练习</h3>

      {/* 呼吸圆圈 */}
      <div className="relative w-48 h-48 flex items-center justify-center">
        {/* 外圈光晕 */}
        <div
          className="absolute w-48 h-48 rounded-full bg-accent-line/50 transition-all ease-in-out"
          style={{
            transform: `scale(${scale})`,
            transitionDuration: `${PHASES[currentPhase].duration}ms`,
          }}
        />
        {/* 内圈 */}
        <div
          className="relative w-32 h-32 rounded-full bg-gradient-to-br from-accent to-accent-hover shadow-lg flex items-center justify-center transition-all ease-in-out"
          style={{
            transform: `scale(${scale})`,
            transitionDuration: `${PHASES[currentPhase].duration}ms`,
          }}
        >
          <span className="text-white text-sm font-medium drop-shadow">
            {isActive ? PHASES[currentPhase].label : '开始'}
          </span>
        </div>
      </div>

      {/* 控制按钮 */}
      <button
        onClick={toggle}
        className={`px-6 py-2 rounded-full text-sm font-medium transition-all duration-300 ${
          isActive
            ? 'bg-emotion-soft text-emotion-ink hover:bg-emotion-line'
            : 'bg-accent-soft text-accent hover:bg-accent-line'
        }`}
      >
        {isActive ? '停止' : '开始呼吸'}
      </button>

      {/* 循环计数 */}
      {cycles > 0 && (
        <p className="text-xs text-ink-faint">已完成 {cycles} 个循环</p>
      )}
    </div>
  )
}
