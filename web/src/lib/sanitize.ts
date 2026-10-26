// ============================================================
// sanitize.ts — HTML 内容处理
// ============================================================

import DOMPurify from 'dompurify'

/**
 * 判断内容是否包含 HTML 标签
 */
export function isHTML(str: string): boolean {
  return /<[a-z][\s\S]*>/i.test(str)
}

/**
 * 去除 HTML 标签，返回纯文本（用于摘要预览）
 */
export function stripHTML(html: string): string {
  if (!isHTML(html)) return html
  if (typeof window === 'undefined') {
    return html.replace(/<[^>]*>/g, '')
  }
  const tmp = document.createElement('div')
  tmp.innerHTML = html
  return tmp.textContent || tmp.innerText || ''
}

/**
 * 清理 HTML 内容，防止 XSS
 */
export function sanitizeHTML(html: string): string {
  if (typeof window === 'undefined') {
    return html
  }
  return DOMPurify.sanitize(html, {
    ALLOWED_TAGS: [
      'p', 'br', 'strong', 'em', 's', 'h1', 'h2', 'h3', 'h4', 'h5', 'h6',
      'ul', 'ol', 'li', 'blockquote', 'pre', 'code', 'a', 'img',
    ],
    ALLOWED_ATTR: ['href', 'src', 'alt', 'title', 'class', 'target'],
  })
}
