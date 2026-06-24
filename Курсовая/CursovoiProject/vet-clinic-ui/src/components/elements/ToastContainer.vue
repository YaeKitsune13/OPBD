<script setup lang="ts">
import { useToast } from "../../utils/useToast";
const { toasts, removeToast } = useToast();
</script>

<template>
    <div class="toast-container">
        <TransitionGroup name="toast">
            <div
                v-for="toast in toasts"
                :key="toast.id"
                :class="['toast-item', toast.type]"
                @click="removeToast(toast.id)"
            >
                <span class="toast-icon">
                    {{
                        toast.type === "success"
                            ? ""
                            : toast.type === "error"
                              ? "x"
                              : "!"
                    }}
                </span>
                <span class="toast-msg">{{ toast.message }}</span>
            </div>
        </TransitionGroup>
    </div>
</template>

<style scoped>
.toast-container {
    position: fixed;
    bottom: 75px;
    right: 24px;
    z-index: 9999;
    display: flex;
    flex-direction: column;
    gap: 10px;
    pointer-events: none;
}

.toast-item {
    pointer-events: auto;
    min-width: 280px;
    max-width: 400px;
    padding: 14px 20px;
    background: var(--surface);
    border-radius: 12px;
    box-shadow: 0 10px 30px rgba(0, 0, 0, 0.15);
    border: 1px solid var(--border);
    display: flex;
    align-items: center;
    gap: 12px;
    cursor: pointer;
    transition: all 0.3s ease;
}

.toast-item.success {
    border-left: 4px solid var(--accent);
}
.toast-item.error {
    border-left: 4px solid var(--red);
}
.toast-item.info {
    border-left: 4px solid var(--blue);
}

.toast-icon {
    font-size: 18px;
}
.toast-msg {
    font-size: 14px;
    font-weight: 600;
    color: var(--text);
}

.toast-enter-from {
    opacity: 0;
    transform: translateX(100%);
}
.toast-leave-to {
    opacity: 0;
    transform: scale(0.9);
}
</style>
