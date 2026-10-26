'use client'

import { useEffect, useRef, useState, useCallback } from 'react'

interface WebSocketMessage {
  type: string
  content: Record<string, unknown>
}

interface UseWebSocketOptions {
  url: string
  token?: string
  onMessage?: (message: WebSocketMessage) => void
  onConnect?: () => void
  onDisconnect?: () => void
  autoReconnect?: boolean
  reconnectInterval?: number
}

export function useWebSocket({
  url,
  token,
  onMessage,
  onConnect,
  onDisconnect,
  autoReconnect = true,
  reconnectInterval = 3000,
}: UseWebSocketOptions) {
  const [isConnected, setIsConnected] = useState(false)
  const [lastMessage, setLastMessage] = useState<WebSocketMessage | null>(null)
  const wsRef = useRef<WebSocket | null>(null)
  const reconnectTimeoutRef = useRef<NodeJS.Timeout | null>(null)

  const connect = useCallback(() => {
    if (wsRef.current?.readyState === WebSocket.OPEN) return

    const wsUrl = token ? `${url}?token=${token}` : url
    const ws = new WebSocket(wsUrl)

    ws.onopen = () => {
      setIsConnected(true)
      onConnect?.()
    }

    ws.onmessage = (event) => {
      try {
        const message = JSON.parse(event.data) as WebSocketMessage
        setLastMessage(message)
        onMessage?.(message)
      } catch (err) {
        console.error('Failed to parse WebSocket message:', err)
      }
    }

    ws.onclose = () => {
      setIsConnected(false)
      onDisconnect?.()

      if (autoReconnect) {
        reconnectTimeoutRef.current = setTimeout(() => {
          connect()
        }, reconnectInterval)
      }
    }

    ws.onerror = (error) => {
      console.error('WebSocket error:', error)
    }

    wsRef.current = ws
  }, [url, token, onMessage, onConnect, onDisconnect, autoReconnect, reconnectInterval])

  const disconnect = useCallback(() => {
    if (reconnectTimeoutRef.current) {
      clearTimeout(reconnectTimeoutRef.current)
    }
    if (wsRef.current) {
      wsRef.current.close()
      wsRef.current = null
    }
  }, [])

  const send = useCallback((data: Record<string, unknown>) => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify(data))
    }
  }, [])

  useEffect(() => {
    connect()
    return () => disconnect()
  }, [connect, disconnect])

  return {
    isConnected,
    lastMessage,
    send,
    connect,
    disconnect,
  }
}

// 通知专用 hook
export function useNotificationWebSocket(onNotification: (notification: Record<string, unknown>) => void) {
  const token = typeof window !== 'undefined' ? localStorage.getItem('token') : undefined
  let wsUrl = 'ws://localhost/api/v1/ws'
  if (typeof window !== 'undefined') {
    const wsProtocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    wsUrl = `${wsProtocol}//${window.location.host}/api/v1/ws`
  }

  const { isConnected } = useWebSocket({
    url: wsUrl,
    token: token || undefined,
    onMessage: (message) => {
      if (message.type === 'notification') {
        onNotification(message.content)
      }
    },
  })

  return { isConnected }
}
