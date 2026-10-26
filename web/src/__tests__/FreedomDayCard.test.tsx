import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import FreedomDayCard from '@/components/FreedomDayCard'

describe('FreedomDayCard', () => {
  it('shows freedom days when provided', () => {
    render(<FreedomDayCard freedomDays={30} />)
    expect(screen.getByText('30')).toBeInTheDocument()
    expect(screen.getByText('自由第 30 天')).toBeInTheDocument()
  })

  it('shows encouraging text', () => {
    render(<FreedomDayCard freedomDays={1} />)
    expect(screen.getByText(/无论如何，这段时间属于你/)).toBeInTheDocument()
  })

  it('shows setup prompt when no freedom days', () => {
    render(<FreedomDayCard />)
    expect(screen.getByText('开始你的自由之旅')).toBeInTheDocument()
    expect(screen.getByText('设置最后工作日')).toBeInTheDocument()
  })

  it('shows emoji when no freedom days', () => {
    render(<FreedomDayCard />)
    expect(screen.getByText('🌅')).toBeInTheDocument()
  })

  it('opens date picker when setup button is clicked', async () => {
    const user = userEvent.setup()
    render(<FreedomDayCard />)

    await user.click(screen.getByText('设置最后工作日'))

    expect(screen.getByText('选择你的最后工作日：')).toBeInTheDocument()
    expect(screen.getByText('保存')).toBeInTheDocument()
    expect(screen.getByText('取消')).toBeInTheDocument()
  })

  it('opens date picker from "修改" button when freedom days exist', async () => {
    const user = userEvent.setup()
    render(<FreedomDayCard freedomDays={10} lastWorkDate="2024-01-01" />)

    await user.click(screen.getByText('修改最后工作日'))

    expect(screen.getByText('选择你的最后工作日：')).toBeInTheDocument()
    expect(screen.getByText('清除')).toBeInTheDocument()
  })

  it('shows clear button only when lastWorkDate exists', async () => {
    const user = userEvent.setup()
    render(<FreedomDayCard freedomDays={10} lastWorkDate="2024-01-01" />)

    await user.click(screen.getByText('修改最后工作日'))
    expect(screen.getByText('清除')).toBeInTheDocument()
  })

  it('does not show clear button without lastWorkDate', async () => {
    const user = userEvent.setup()
    render(<FreedomDayCard />)

    await user.click(screen.getByText('设置最后工作日'))
    expect(screen.queryByText('清除')).not.toBeInTheDocument()
  })

  it('calls onSetDate with selected date on save', async () => {
    const user = userEvent.setup()
    const onSetDate = vi.fn()
    render(<FreedomDayCard onSetDate={onSetDate} />)

    await user.click(screen.getByText('设置最后工作日'))

    const dateInput = screen.getByDisplayValue('')
    await user.type(dateInput, '2024-03-15')

    await user.click(screen.getByText('保存'))

    expect(onSetDate).toHaveBeenCalledWith('2024-03-15')
  })

  it('calls onSetDate with empty string on clear', async () => {
    const user = userEvent.setup()
    const onSetDate = vi.fn()
    render(<FreedomDayCard lastWorkDate="2024-01-01" freedomDays={30} onSetDate={onSetDate} />)

    await user.click(screen.getByText('修改最后工作日'))
    await user.click(screen.getByText('清除'))

    expect(onSetDate).toHaveBeenCalledWith('')
  })

  it('closes date picker on cancel', async () => {
    const user = userEvent.setup()
    render(<FreedomDayCard />)

    await user.click(screen.getByText('设置最后工作日'))
    expect(screen.getByText('选择你的最后工作日：')).toBeInTheDocument()

    await user.click(screen.getByText('取消'))
    expect(screen.queryByText('选择你的最后工作日：')).not.toBeInTheDocument()
  })

  it('shows pre-filled date when lastWorkDate is provided', async () => {
    const user = userEvent.setup()
    render(<FreedomDayCard lastWorkDate="2024-06-01" freedomDays={50} />)

    await user.click(screen.getByText('修改最后工作日'))

    const dateInput = screen.getByDisplayValue('2024-06-01')
    expect(dateInput).toBeInTheDocument()
  })

  it('applies custom className', () => {
    const { container } = render(<FreedomDayCard className="custom-class" />)
    expect(container.firstChild).toHaveAttribute('class')
    expect(container.firstChild?.textContent).toContain('开始你的自由之旅')
  })
})
