import { reactive } from 'vue'

export type ToastType = 'success' | 'error' | 'info'

interface Toast {
    id: number
    message: string
    type: ToastType
}

const toasts = reactive<Toast[]>([])

export function useToast() {
    const showToast = (message: string, type: ToastType = 'info') => {
        const id = Date.now()
        toasts.push({ id, message, type })

        // Удаляем через 3 секунды (2.8с показ + 0.2с анимация ухода)
        setTimeout(() => {
            const index = toasts.findIndex(t => t.id === id)
            if (index !== -1) toasts.splice(index, 1)
        }, 3000)
    }

    return { toasts, showToast }
}