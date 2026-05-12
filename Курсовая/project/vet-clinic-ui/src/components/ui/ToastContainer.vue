<script setup lang="ts">
import { useToast } from '../../utils/useToast'
const { toasts } = useToast()

const icons = { success: "✓", error: "✕", info: "ℹ" };
</script>

<template>
    <div class="toast-area">
        <TransitionGroup name="toast">
            <div v-for="toast in toasts" :key="toast.id" :class="['toast', toast.type]">
                <span class="toast-icon">{{ icons[toast.type] }}</span>
                <span>{{ toast.message }}</span>
            </div>
        </TransitionGroup>
    </div>
</template>

<style scoped>
.toast-area {
    position: fixed;
    bottom: 20px;
    right: 20px;
    display: flex;
    flex-direction: column;
    gap: 10px;
    z-index: 9999;
}

.toast {
    background: var(--surface2); /* Используем твои переменные */
    border: 1px solid var(--border);
    border-radius: var(--radius);
    padding: 12px 16px;
    font-size: 13px;
    display: flex;
    align-items: center;
    gap: 10px;
    color: var(--text);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
    min-width: 200px;
}

.toast-icon {
    font-weight: bold;
}

/* Твои стили типов */
.toast.success { border-left: 3px solid var(--accent); }
.toast.error { border-left: 3px solid #ef4444; } /* Можно заменить на var(--red) */
.toast.info { border-left: 3px solid #3b82f6; } /* Можно заменить на var(--blue) */

/* Анимация (аналог твоего toastIn) */
.toast-enter-active {
    animation: toastIn 0.2s ease;
}
.toast-leave-active {
    transition: all 0.2s ease;
}
.toast-leave-to {
    opacity: 0;
    transform: translateX(20px);
}

@keyframes toastIn {
    from {
        opacity: 0;
        transform: translateX(20px);
    }
    to {
        opacity: 1;
        transform: translateX(0);
    }
}
</style>