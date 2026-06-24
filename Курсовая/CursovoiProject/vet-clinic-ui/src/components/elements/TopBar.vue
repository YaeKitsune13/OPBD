<script setup lang="ts">
defineProps<{
    userName: string;
    userAvatar: string;
    currentRole: "client" | "doctor";
}>();

const emit = defineEmits<{
    "toggle-sidebar": [];
    "open-profile": [];
}>();
</script>

<template>
    <div class="topbar">
        <button class="hamburger" @click="emit('toggle-sidebar')">☰</button>

        <div class="logo">
            <div class="logo-icon">🐾</div>
            <span>ВетКлиника</span>
        </div>

        <div class="topbar-sep" />

        <div class="topbar-user" @click="emit('open-profile')">
            <div class="avatar">{{ userAvatar }}</div>
            <span class="topbar-user-name">{{ userName }}</span>
        </div>
    </div>
</template>

<style scoped>
.topbar {
    height: var(--topbar-h);
    background: var(--surface);
    border-bottom: 1px solid var(--border);
    display: flex;
    align-items: center;
    padding: 0 20px;
    gap: 16px;
    flex-shrink: 0;
    z-index: 100;
}

.logo {
    display: flex;
    align-items: center;
    gap: 10px;
    font-size: 15px;
    font-weight: 600;
    letter-spacing: -0.3px;
    color: var(--text);
    white-space: nowrap;
}
.logo-icon {
    width: 28px;
    height: 28px;
    background: var(--accent);
    border-radius: 6px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 14px;
}

.topbar-sep {
    flex: 1;
}

.topbar-user {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 4px 10px;
    border-radius: var(--radius);
    cursor: pointer;
    transition: background 0.15s;
}
.topbar-user:hover {
    background: var(--surface2);
}

.avatar {
    width: 28px;
    height: 28px;
    border-radius: 50%;
    background: var(--accent-dim);
    border: 1px solid var(--accent);
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 11px;
    font-weight: 600;
    color: var(--accent);
}
.topbar-user-name {
    font-size: 13px;
    font-weight: 500;
}

.hamburger {
    display: none;
    background: none;
    border: none;
    color: var(--text2);
    cursor: pointer;
    padding: 4px;
    font-size: 18px;
}
@media (max-width: 768px) {
    .hamburger {
        display: flex;
    }
}
</style>
