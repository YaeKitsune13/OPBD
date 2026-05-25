<script setup lang="ts">
defineProps<{ show: boolean; title: string; maxWidth?: string }>();
const emit = defineEmits<{ close: [] }>();
</script>

<template>
    <Teleport to="body">
        <div v-if="show" class="modal-overlay" @click.self="emit('close')">
            <div class="modal" :style="maxWidth ? `max-width:${maxWidth}` : ''">
                <div class="modal-header">
                    <span class="modal-title">{{ title }}</span>
                    <button class="modal-close" @click="emit('close')">
                        ✕
                    </button>
                </div>
                <div class="modal-body">
                    <slot />
                </div>
                <div class="modal-footer" v-if="$slots.footer">
                    <slot name="footer" />
                </div>
            </div>
        </div>
    </Teleport>
</template>

<style scoped>
.modal-overlay {
    position: fixed;
    inset: 0;
    background: rgba(0, 0, 0, 0.4);
    z-index: 1000;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 20px;
}
.modal {
    background: var(--surface);
    border: 1px solid var(--border);
    border-radius: 10px;
    width: 100%;
    max-width: 520px;
    overflow: hidden;
    animation: fadeIn 0.15s ease;
}
.modal-header {
    padding: 14px 16px;
    border-bottom: 1px solid var(--border);
    display: flex;
    align-items: center;
    justify-content: space-between;
}
.modal-title {
    font-size: 14px;
    font-weight: 600;
}
.modal-close {
    background: none;
    border: none;
    color: var(--text3);
    cursor: pointer;
    font-size: 14px;
    padding: 2px 6px;
    border-radius: 4px;
    transition: background 0.15s;
}
.modal-close:hover {
    background: var(--surface3);
    color: var(--text);
}
.modal-body {
    padding: 16px;
}
.modal-footer {
    padding: 12px 16px;
    border-top: 1px solid var(--border);
    display: flex;
    justify-content: flex-end;
    gap: 8px;
}
</style>
