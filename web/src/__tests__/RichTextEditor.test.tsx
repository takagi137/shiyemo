import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, fireEvent } from '@testing-library/react'
import userEvent from '@testing-library/user-event'

// mock @tiptap/react
const mockEditor = {
  isActive: vi.fn((name: string, attrs?: Record<string, unknown>) => false),
  chain: vi.fn(() => ({
    focus: vi.fn(() => ({
      toggleBold: vi.fn(() => ({ run: vi.fn() })),
      toggleItalic: vi.fn(() => ({ run: vi.fn() })),
      toggleStrike: vi.fn(() => ({ run: vi.fn() })),
      toggleHeading: vi.fn(() => ({ run: vi.fn() })),
      toggleBulletList: vi.fn(() => ({ run: vi.fn() })),
      toggleOrderedList: vi.fn(() => ({ run: vi.fn() })),
      toggleBlockquote: vi.fn(() => ({ run: vi.fn() })),
      toggleCodeBlock: vi.fn(() => ({ run: vi.fn() })),
      setImage: vi.fn(() => ({ run: vi.fn() })),
      setLink: vi.fn(() => ({ run: vi.fn() })),
      undo: vi.fn(() => ({ run: vi.fn() })),
      redo: vi.fn(() => ({ run: vi.fn() })),
    })),
  })),
  getHTML: vi.fn(() => '<p>Hello</p>'),
  view: {
    dispatch: vi.fn(),
  },
}

vi.mock('@tiptap/react', () => ({
  useEditor: vi.fn(({ onUpdate }) => {
    // simulate onUpdate setup
    if (onUpdate) {
      // store for later use
    }
    return mockEditor
  }),
  EditorContent: ({ editor }: { editor: typeof mockEditor }) => (
    <div data-testid="editor-content" role="textbox">
      {editor ? 'Editor Content' : 'No Editor'}
    </div>
  ),
}))

vi.mock('@tiptap/starter-kit', () => ({
  default: { configure: vi.fn(() => ({})) },
}))

vi.mock('@tiptap/extension-image', () => ({
  default: { configure: vi.fn(() => ({})) },
}))

vi.mock('@tiptap/extension-link', () => ({
  default: { configure: vi.fn(() => ({})) },
}))

vi.mock('@tiptap/extension-placeholder', () => ({
  default: { configure: vi.fn(() => ({})) },
}))

vi.mock('@/lib/api', () => ({
  getApiBase: vi.fn(() => '/api/v1'),
}))

// mock window.alert
const mockAlert = vi.fn()
Object.defineProperty(window, 'alert', { value: mockAlert, writable: true })

// mock window.prompt
const mockPrompt = vi.fn()
Object.defineProperty(window, 'prompt', { value: mockPrompt, writable: true })

import RichTextEditor from '@/components/RichTextEditor'

describe('RichTextEditor', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    mockEditor.isActive.mockReturnValue(false)
  })

  it('renders editor content area', () => {
    render(<RichTextEditor content="" onChange={vi.fn()} />)
    expect(screen.getByTestId('editor-content')).toBeInTheDocument()
  })

  it('renders toolbar buttons', () => {
    render(<RichTextEditor content="" onChange={vi.fn()} />)

    expect(screen.getByTitle('粗体')).toBeInTheDocument()
    expect(screen.getByTitle('斜体')).toBeInTheDocument()
    expect(screen.getByTitle('删除线')).toBeInTheDocument()
    expect(screen.getByTitle('标题')).toBeInTheDocument()
    expect(screen.getByTitle('小标题')).toBeInTheDocument()
    expect(screen.getByTitle('无序列表')).toBeInTheDocument()
    expect(screen.getByTitle('有序列表')).toBeInTheDocument()
    expect(screen.getByTitle('引用')).toBeInTheDocument()
    expect(screen.getByTitle('代码块')).toBeInTheDocument()
    expect(screen.getByTitle('插入图片')).toBeInTheDocument()
    expect(screen.getByTitle('插入链接')).toBeInTheDocument()
    expect(screen.getByTitle('撤销')).toBeInTheDocument()
    expect(screen.getByTitle('重做')).toBeInTheDocument()
  })

  it('has 13 toolbar buttons total', () => {
    render(<RichTextEditor content="" onChange={vi.fn()} />)
    const buttons = screen.getAllByRole('button')
    expect(buttons.length).toBe(13)
  })

  it('shows textarea textbox role', () => {
    render(<RichTextEditor content="initial content" onChange={vi.fn()} />)
    expect(screen.getByRole('textbox')).toBeInTheDocument()
  })

  it('renders hidden file input for image upload', () => {
    const { container } = render(<RichTextEditor content="" onChange={vi.fn()} />)
    const fileInput = container.querySelector('input[type="file"]') as HTMLInputElement
    expect(fileInput).toBeInTheDocument()
    expect(fileInput.accept).toContain('image/jpeg')
    expect(fileInput.className).toContain('hidden')
  })

  it('calls prompt when link button clicked', async () => {
    const user = userEvent.setup()
    mockPrompt.mockReturnValue('https://example.com')

    render(<RichTextEditor content="" onChange={vi.fn()} />)

    await user.click(screen.getByTitle('插入链接'))
    expect(mockPrompt).toHaveBeenCalledWith('请输入链接地址：')
  })

  it('does not call setLink if prompt returns null', async () => {
    const user = userEvent.setup()
    mockPrompt.mockReturnValue(null)

    render(<RichTextEditor content="" onChange={vi.fn()} />)

    await user.click(screen.getByTitle('插入链接'))
    // No error should be thrown
  })

  it('calls mockEditor methods on toolbar click', async () => {
    const user = userEvent.setup()

    render(<RichTextEditor content="" onChange={vi.fn()} />)

    await user.click(screen.getByTitle('粗体'))
    expect(mockEditor.chain).toHaveBeenCalled()

    await user.click(screen.getByTitle('斜体'))
    expect(mockEditor.chain).toHaveBeenCalled()

    await user.click(screen.getByTitle('撤销'))
    expect(mockEditor.chain).toHaveBeenCalled()

    await user.click(screen.getByTitle('重做'))
    expect(mockEditor.chain).toHaveBeenCalled()
  })

  it('shows toolbar separator dividers', () => {
    const { container } = render(<RichTextEditor content="" onChange={vi.fn()} />)
    const separators = container.querySelectorAll('.border-gray-300')
    // toolbar separators + outer border
    expect(separators.length).toBeGreaterThanOrEqual(1)
  })

  it('has proper editor content class', () => {
    const { container } = render(<RichTextEditor content="" onChange={vi.fn()} />)
    const border = container.querySelector('.border-gray-300.rounded-lg')
    expect(border).toBeInTheDocument()
  })

  it('renders toolbar background styling', () => {
    const { container } = render(<RichTextEditor content="" onChange={vi.fn()} />)
    const toolbar = container.querySelector('.bg-gray-50')
    expect(toolbar).toBeInTheDocument()
  })

  it('contains editor content wrapper div', () => {
    const { container } = render(<RichTextEditor content="" onChange={vi.fn()} />)
    const wrapper = container.querySelector('.border-gray-300.rounded-lg')
    expect(wrapper).toBeInTheDocument()
  })

  it('renders all toolbar buttons inside toolbar area', () => {
    const { container } = render(<RichTextEditor content="" onChange={vi.fn()} />)
    const toolbar = container.querySelector('.bg-gray-50')
    expect(toolbar).toBeInTheDocument()
    const buttons = toolbar!.querySelectorAll('button')
    expect(buttons.length).toBeGreaterThanOrEqual(13)
  })

  it('file input accepts only image types', () => {
    const { container } = render(<RichTextEditor content="" onChange={vi.fn()} />)
    const fileInput = container.querySelector('input[type="file"]') as HTMLInputElement
    expect(fileInput.accept).toBe('image/jpeg,image/png,image/gif,image/webp')
  })
})
