'use client'

import { useEditor, EditorContent } from '@tiptap/react'
import StarterKit from '@tiptap/starter-kit'
import Image from '@tiptap/extension-image'
import Link from '@tiptap/extension-link'
import Placeholder from '@tiptap/extension-placeholder'
import { useRef, useCallback } from 'react'
import { getApiBase } from '@/lib/api'

interface RichTextEditorProps {
  content: string
  onChange: (html: string) => void
  placeholder?: string
}

async function uploadImage(file: File): Promise<string> {
  const token = typeof window !== 'undefined' ? localStorage.getItem('token') : null
  const apiBase = getApiBase()

  const formData = new FormData()
  formData.append('file', file)

  const res = await fetch(`${apiBase}/upload/image`, {
    method: 'POST',
    headers: token ? { Authorization: `Bearer ${token}` } : {},
    body: formData,
  })

  const data = await res.json()
  if (data.code !== 0) {
    throw new Error(data.message || '上传失败')
  }

  const url = data.data.url as string
  if (url.startsWith('/')) {
    const base = apiBase.replace('/api/v1', '')
    return `${base}${url}`
  }
  return url
}

export default function RichTextEditor({ content, onChange, placeholder }: RichTextEditorProps) {
  const fileInputRef = useRef<HTMLInputElement>(null)

  const editor = useEditor({
    extensions: [
      StarterKit.configure({
        heading: { levels: [1, 2, 3] },
      }),
      Image.configure({
        HTMLAttributes: {
          class: 'max-w-full h-auto rounded-lg my-2',
        },
      }),
      Link.configure({
        openOnClick: false,
        HTMLAttributes: {
          class: 'text-blue-600 underline hover:text-blue-800',
        },
      }),
      Placeholder.configure({
        placeholder: placeholder || '分享你的想法...',
      }),
    ],
    content,
    onUpdate: ({ editor }) => {
      onChange(editor.getHTML())
    },
    editorProps: {
      attributes: {
        class: 'prose prose-sm sm:prose max-w-none focus:outline-none min-h-[200px] px-3 py-2',
      },
    },
  })

  const handleImageUpload = useCallback(
    async (file: File) => {
      if (!editor) return
      if (!file.type.startsWith('image/')) {
        alert('请选择图片文件')
        return
      }
      if (file.size > 10 * 1024 * 1024) {
        alert('图片大小不能超过 10MB')
        return
      }

      editor.chain().focus().setImage({ src: '', alt: '上传中...' }).run()

      try {
        const url = await uploadImage(file)
        const { state } = editor
        const tr = state.tr
        state.doc.descendants((node, nodePos) => {
          if (node.type.name === 'image' && node.attrs.src === '') {
            tr.setNodeMarkup(nodePos, undefined, { src: url, alt: file.name })
          }
        })
        editor.view.dispatch(tr)
      } catch (err) {
        const { state } = editor
        const tr = state.tr
        state.doc.descendants((node, nodePos) => {
          if (node.type.name === 'image' && node.attrs.src === '') {
            tr.delete(nodePos, nodePos + node.nodeSize)
          }
        })
        editor.view.dispatch(tr)
        alert(err instanceof Error ? err.message : '图片上传失败')
      }
    },
    [editor]
  )

  const handleFileChange = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>) => {
      const file = e.target.files?.[0]
      if (file) {
        handleImageUpload(file)
      }
      e.target.value = ''
    },
    [handleImageUpload]
  )

  const handlePaste = useCallback(
    (e: React.ClipboardEvent) => {
      const items = e.clipboardData?.items
      if (!items) return
      for (let i = 0; i < items.length; i++) {
        const item = items[i]
        if (item.type.startsWith('image/')) {
          e.preventDefault()
          const file = item.getAsFile()
          if (file) {
            handleImageUpload(file)
          }
          return
        }
      }
    },
    [handleImageUpload]
  )

  const handleDrop = useCallback(
    (e: React.DragEvent) => {
      const files = e.dataTransfer?.files
      if (!files?.length) return
      const file = files[0]
      if (file.type.startsWith('image/')) {
        e.preventDefault()
        handleImageUpload(file)
      }
    },
    [handleImageUpload]
  )

  if (!editor) return null

  return (
    <div className="border border-gray-300 rounded-lg overflow-hidden bg-white">
      {/* 工具栏 */}
      <div className="flex flex-wrap items-center gap-1 px-2 py-1.5 border-b border-gray-200 bg-gray-50">
        <ToolbarBtn
          active={editor.isActive('bold')}
          onClick={() => editor.chain().focus().toggleBold().run()}
          title="粗体"
        >
          <b>B</b>
        </ToolbarBtn>
        <ToolbarBtn
          active={editor.isActive('italic')}
          onClick={() => editor.chain().focus().toggleItalic().run()}
          title="斜体"
        >
          <i>I</i>
        </ToolbarBtn>
        <ToolbarBtn
          active={editor.isActive('strike')}
          onClick={() => editor.chain().focus().toggleStrike().run()}
          title="删除线"
        >
          <s>S</s>
        </ToolbarBtn>

        <div className="w-px h-5 bg-gray-300 mx-1" />

        <ToolbarBtn
          active={editor.isActive('heading', { level: 2 })}
          onClick={() => editor.chain().focus().toggleHeading({ level: 2 }).run()}
          title="标题"
        >
          H2
        </ToolbarBtn>
        <ToolbarBtn
          active={editor.isActive('heading', { level: 3 })}
          onClick={() => editor.chain().focus().toggleHeading({ level: 3 }).run()}
          title="小标题"
        >
          H3
        </ToolbarBtn>

        <div className="w-px h-5 bg-gray-300 mx-1" />

        <ToolbarBtn
          active={editor.isActive('bulletList')}
          onClick={() => editor.chain().focus().toggleBulletList().run()}
          title="无序列表"
        >
          •
        </ToolbarBtn>
        <ToolbarBtn
          active={editor.isActive('orderedList')}
          onClick={() => editor.chain().focus().toggleOrderedList().run()}
          title="有序列表"
        >
          1.
        </ToolbarBtn>
        <ToolbarBtn
          active={editor.isActive('blockquote')}
          onClick={() => editor.chain().focus().toggleBlockquote().run()}
          title="引用"
        >
          ❝
        </ToolbarBtn>
        <ToolbarBtn
          active={editor.isActive('codeBlock')}
          onClick={() => editor.chain().focus().toggleCodeBlock().run()}
          title="代码块"
        >
          {'</>'}
        </ToolbarBtn>

        <div className="w-px h-5 bg-gray-300 mx-1" />

        <ToolbarBtn onClick={() => fileInputRef.current?.click()} title="插入图片">
          🖼️
        </ToolbarBtn>

        <ToolbarBtn
          active={editor.isActive('link')}
          onClick={() => {
            const url = window.prompt('请输入链接地址：')
            if (url) {
              editor.chain().focus().setLink({ href: url }).run()
            }
          }}
          title="插入链接"
        >
          🔗
        </ToolbarBtn>

        <div className="w-px h-5 bg-gray-300 mx-1" />

        <ToolbarBtn onClick={() => editor.chain().focus().undo().run()} title="撤销">
          ↩
        </ToolbarBtn>
        <ToolbarBtn onClick={() => editor.chain().focus().redo().run()} title="重做">
          ↪
        </ToolbarBtn>
      </div>

      <input
        ref={fileInputRef}
        type="file"
        accept="image/jpeg,image/png,image/gif,image/webp"
        className="hidden"
        onChange={handleFileChange}
      />

      <div onPaste={handlePaste} onDrop={handleDrop} onDragOver={(e) => e.preventDefault()}>
        <EditorContent editor={editor} />
      </div>
    </div>
  )
}

function ToolbarBtn({
  active,
  onClick,
  title,
  children,
}: {
  active?: boolean
  onClick: () => void
  title: string
  children: React.ReactNode
}) {
  return (
    <button
      type="button"
      onClick={onClick}
      title={title}
      className={`w-8 h-8 flex items-center justify-center rounded text-sm transition-colors ${
        active
          ? 'bg-blue-100 text-blue-700'
          : 'text-gray-600 hover:bg-gray-200 hover:text-gray-900'
      }`}
    >
      {children}
    </button>
  )
}
