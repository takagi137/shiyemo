'use client'

import { useState, useEffect } from 'react'

interface FreedomDayCardProps {
  lastWorkDate?: string | null
  freedomDays?: number | null
  onSetDate?: (date: string) => void
  className?: string
}

export default function FreedomDayCard({ lastWorkDate, freedomDays, onSetDate, className = '' }: FreedomDayCardProps) {
  const [showDatePicker, setShowDatePicker] = useState(false)
  const [selectedDate, setSelectedDate] = useState(lastWorkDate || '')

  useEffect(() => {
    if (lastWorkDate) {
      setSelectedDate(lastWorkDate)
    }
  }, [lastWorkDate])

  const handleSave = () => {
    if (selectedDate && onSetDate) {
      onSetDate(selectedDate)
      setShowDatePicker(false)
    }
  }

  const handleClear = () => {
    if (onSetDate) {
      onSetDate('')
      setSelectedDate('')
      setShowDatePicker(false)
    }
  }

  return (
    <div className={`bg-gradient-to-br from-blue-50 to-indigo-50 rounded-xl p-6 border border-blue-100 ${className}`}>
      <div className="text-center">
        {freedomDays !== null && freedomDays !== undefined ? (
          <>
            <div className="text-5xl font-bold text-blue-600 mb-2">
              {freedomDays}
            </div>
            <div className="text-lg text-gray-700 mb-1">自由第 {freedomDays} 天</div>
            <p className="text-sm text-gray-500">
              无论如何，这段时间属于你。
            </p>
          </>
        ) : (
          <>
            <div className="text-4xl mb-3">🌅</div>
            <div className="text-lg text-gray-700 mb-2">开始你的自由之旅</div>
            <p className="text-sm text-gray-500 mb-4">
              设置你的&ldquo;最后工作日&rdquo;，记录你的自由时光。
            </p>
            <button
              onClick={() => setShowDatePicker(true)}
              className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors text-sm"
            >
              设置最后工作日
            </button>
          </>
        )}
      </div>

      {showDatePicker && (
        <div className="mt-4 pt-4 border-t border-blue-200">
          <div className="flex flex-col gap-3">
            <label className="text-sm text-gray-600">选择你的最后工作日：</label>
            <input
              type="date"
              value={selectedDate}
              onChange={(e) => setSelectedDate(e.target.value)}
              max={new Date().toISOString().split('T')[0]}
              className="px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
            <div className="flex gap-2">
              <button
                onClick={handleSave}
                disabled={!selectedDate}
                className="flex-1 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors text-sm disabled:opacity-50"
              >
                保存
              </button>
              {lastWorkDate && (
                <button
                  onClick={handleClear}
                  className="px-4 py-2 text-gray-600 border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors text-sm"
                >
                  清除
                </button>
              )}
              <button
                onClick={() => setShowDatePicker(false)}
                className="px-4 py-2 text-gray-600 border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors text-sm"
              >
                取消
              </button>
            </div>
          </div>
        </div>
      )}

      {freedomDays !== null && freedomDays !== undefined && !showDatePicker && (
        <div className="mt-4 pt-4 border-t border-blue-200 text-center">
          <button
            onClick={() => setShowDatePicker(true)}
            className="text-sm text-blue-600 hover:text-blue-800"
          >
            修改最后工作日
          </button>
        </div>
      )}
    </div>
  )
}
