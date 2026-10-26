// ============================================================
// toast.ts — 轻量级全局 Toast
// ============================================================

type ToastFn = (message: string, type?: string) => void

let registeredToast: ToastFn | null = null

export function registerToast(fn: ToastFn) {
  registeredToast = fn
}

export function toastSuccess(message: string) {
  registeredToast?.(message, 'success')
}

export function toastError(message: string) {
  registeredToast?.(message, 'error')
}

export function toastWarning(message: string) {
  registeredToast?.(message, 'warning')
}
