// API 基础配置
export function getApiBase(): string {
  if (process.env.NEXT_PUBLIC_API_URL) {
    return process.env.NEXT_PUBLIC_API_URL
  }
  // 客户端：使用相对路径，API 请求由 nginx /api/* 反代到后端
  if (typeof window !== 'undefined') {
    return '/api/v1'
  }
  // SSR：回退到 localhost（Docker 内部网络）
  return 'http://api:8080/api/v1'
}

// 统一响应结构
interface ApiResponse<T> {
  code: number
  message: string
  data: T
}

// 通用请求函数
async function request<T>(
  path: string,
  options?: RequestInit
): Promise<ApiResponse<T>> {
  const token = typeof window !== 'undefined' ? localStorage.getItem('token') : null
  
  const headers: HeadersInit = {
    'Content-Type': 'application/json',
  }
  
  if (token) {
    headers['Authorization'] = `Bearer ${token}`
  }

  try {
    const res = await fetch(`${getApiBase()}${path}`, {
      ...options,
      headers: {
        ...headers,
        ...options?.headers,
      },
    })

    const data: ApiResponse<T> = await res.json()

    if (data.code !== 0) {
      console.error('API Error:', data.message)
    }

    // 适配后端分页格式：后端返回 {list, page, size, has_more}，前端统一期望 {items, total}
    if (data.code === 0 && data.data && typeof data.data === 'object') {
      const d = data.data as Record<string, unknown>
      if ('list' in d && !('items' in d)) {
        d.items = d.list || []
        const list = (d.list as unknown[]) || []
        const pg = (d.page as number) || 1
        const sz = (d.size as number) || 10
        d.total = d.has_more ? pg * sz + 1 : (pg - 1) * sz + list.length
      }
    }

    return data
  } catch (err) {
    console.error('Network Error:', err)
    throw err
  }
}

// GET 请求
export function get<T>(path: string) {
  return request<T>(path)
}

// POST 请求
export function post<T>(path: string, data?: unknown) {
  return request<T>(path, {
    method: 'POST',
    body: data ? JSON.stringify(data) : undefined,
  })
}

// PUT 请求
export function put<T>(path: string, data?: unknown) {
  return request<T>(path, {
    method: 'PUT',
    body: data ? JSON.stringify(data) : undefined,
  })
}

// DELETE 请求
export function del<T>(path: string, data?: unknown) {
  return request<T>(path, {
    method: 'DELETE',
    body: data ? JSON.stringify(data) : undefined,
  })
}

// 文件下载（CSV 等二进制/文本响应）
export async function downloadFile(path: string, defaultFilename?: string) {
  const token = typeof window !== 'undefined' ? localStorage.getItem('token') : null
  const headers: Record<string, string> = {}
  if (token) headers['Authorization'] = `Bearer ${token}`

  const res = await fetch(`${getApiBase()}${path}`, { headers })
  if (!res.ok) throw new Error(`下载失败: ${res.statusText}`)

  // 从 Content-Disposition 提取文件名
  const disposition = res.headers.get('Content-Disposition') || ''
  const filenameMatch = disposition.match(/filename=([^;\s]+)/)
  const filename = filenameMatch?.[1] || defaultFilename || 'export.csv'

  const blob = await res.blob()
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = decodeURIComponent(filename)
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
  URL.revokeObjectURL(url)
}
