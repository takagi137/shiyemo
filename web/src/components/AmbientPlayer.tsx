'use client'

import { useState, useEffect, useRef } from 'react'

interface Sound {
  key: string
  name: string
  emoji: string
  url: string
}

export default function AmbientPlayer() {
  // 音效列表由后端维护（GET /api/v1/sounds），不再硬编码在前端。
  const [sounds, setSounds] = useState<Sound[]>([])
  const [playing, setPlaying] = useState<string | null>(null)
  const [volume, setVolume] = useState(50)
  // 播放失败的音效 key。必须让用户看见 —— 见 play() 里的说明。
  const [failed, setFailed] = useState<string | null>(null)
  const audioRef = useRef<HTMLAudioElement | null>(null)

  // 拉取音效列表
  useEffect(() => {
    fetch('/api/v1/sounds')
      .then((r) => r.json())
      .then((json) => {
        if (json.code === 0 && Array.isArray(json.data?.items)) {
          setSounds(json.data.items)
        }
      })
      .catch((err) => console.error('Failed to fetch sounds:', err))
  }, [])

  // 从 localStorage 恢复音量
  useEffect(() => {
    const saved = localStorage.getItem('sanctuary-volume')
    if (saved) setVolume(Number(saved))
  }, [])

  // 保存音量到 localStorage
  useEffect(() => {
    localStorage.setItem('sanctuary-volume', String(volume))
  }, [volume])

  // 离开页面时暂停
  useEffect(() => {
    return () => {
      if (audioRef.current) {
        audioRef.current.pause()
      }
    }
  }, [])

  const play = async (sound: Sound) => {
    if (playing === sound.key) {
      // 暂停
      audioRef.current?.pause()
      setPlaying(null)
      return
    }

    // 停止当前播放
    if (audioRef.current) {
      audioRef.current.pause()
    }
    setFailed(null)

    const audio = new Audio(sound.url)
    audio.volume = volume / 100
    audio.loop = true
    audioRef.current = audio

    // 只有真正播起来了才切到「播放中」。音频加载失败时显式提示，不假装在播。
    try {
      await audio.play()
      setPlaying(sound.key)
      localStorage.setItem('sanctuary-playing', sound.key)
    } catch {
      audioRef.current = null
      setPlaying(null)
      setFailed(sound.key)
      localStorage.removeItem('sanctuary-playing')
    }
  }

  const handleVolumeChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const v = Number(e.target.value)
    setVolume(v)
    if (audioRef.current) {
      audioRef.current.volume = v / 100
    }
  }

  return (
    <div className="flex flex-col items-center space-y-4">
      <h3 className="text-lg font-medium text-ink">环境音</h3>

      {/* 音效网格 */}
      <div className="grid grid-cols-3 gap-3 w-full max-w-xs">
        {sounds.map((sound) => (
          <button
            key={sound.key}
            onClick={() => play(sound)}
            className={`flex flex-col items-center p-3 rounded-xl transition-all duration-300 ${
              playing === sound.key
                ? 'bg-accent-soft text-accent shadow-md scale-105'
                : 'bg-surface/70 text-ink-soft hover:bg-surface/90'
            }`}
          >
            <span className="text-2xl mb-1">{sound.emoji}</span>
            <span className="text-xs">{sound.name}</span>
          </button>
        ))}
      </div>

      {/* 音量控制 */}
      <div className="flex items-center space-x-3 w-full max-w-xs">
        <span className="text-xs text-ink-faint">🔈</span>
        <input
          type="range"
          min="0"
          max="100"
          value={volume}
          onChange={handleVolumeChange}
          className="flex-1 h-1 bg-line rounded-lg appearance-none cursor-pointer accent-[var(--color-accent)]"
        />
        <span className="text-xs text-ink-faint">🔊</span>
      </div>

      {/* 播放状态 */}
      {playing && (
        <p className="text-xs text-accent">
          正在播放：{sounds.find((s) => s.key === playing)?.name}
        </p>
      )}

      {/* 播放失败 —— 宁可明说「放不出来」，也不要假装在播 */}
      {failed && (
        <p className="text-xs text-ink-soft text-center">
          「{sounds.find((s) => s.key === failed)?.name}」暂时放不出来
          <span className="block text-ink-faint mt-0.5">音频资源加载失败</span>
        </p>
      )}
    </div>
  )
}
