<script setup lang="ts">
const props = defineProps<{
    currentRole: "client" | "doctor";
    activePage: string;
    isOpen: boolean;
}>();

const emit = defineEmits<{
    navigate: [page: string];
    close: [];
}>();

function goTo(page: string) {
    emit("navigate", page);
    emit("close");
}
</script>

<template>
    <aside class="sidebar" :class="{ open: isOpen }">
        <div v-if="currentRole === 'client'" class="nav-root">
            <div class="nav-section">
                <div class="nav-section-label">Главное</div>
                <button
                    class="nav-item"
                    :class="{ active: activePage === 'dashboard' }"
                    @click="goTo('dashboard')"
                >
                    <span class="nav-dot" />
                    Обзор
                </button>
                <button
                    class="nav-item"
                    :class="{ active: activePage === 'pets' }"
                    @click="goTo('pets')"
                >
                    <span class="nav-dot" />
                    Мои питомцы
                </button>
            </div>

            <div class="nav-section">
                <div class="nav-section-label">Услуги</div>
                <button
                    class="nav-item"
                    :class="{ active: activePage === 'appointments' }"
                    @click="goTo('appointments')"
                >
                    <span class="nav-dot" />
                    Мои записи
                </button>
                <button
                    class="nav-item"
                    :class="{ active: activePage === 'book' }"
                    @click="goTo('book')"
                >
                    <span class="nav-dot" />
                    Записаться
                </button>
                <button
                    class="nav-item"
                    :class="{ active: activePage === 'shop' }"
                    @click="goTo('shop')"
                >
                    <span class="nav-dot" />
                    Медикаменты
                </button>
                <button
                    class="nav-item"
                    :class="{ active: activePage === 'orders' }"
                    @click="goTo('orders')"
                >
                    <span class="nav-dot" />
                    Мои заказы
                </button>
            </div>

            <div class="nav-section">
                <div class="nav-section-label">Питомец</div>
                <button
                    class="nav-item"
                    :class="{ active: activePage === 'stats' }"
                    @click="goTo('stats')"
                >
                    <span class="nav-dot" />
                    Статистика
                </button>
            </div>
        </div>

        <div v-else-if="currentRole === 'doctor'" class="nav-root">
            <div class="nav-section">
                <div class="nav-section-label">Приём</div>
                <button
                    class="nav-item"
                    :class="{ active: activePage === 'schedule' }"
                    @click="goTo('schedule')"
                >
                    <span class="nav-dot" />
                    Расписание
                </button>
            </div>
            <div class="nav-section">
                <div class="nav-section-label">Пациенты</div>
                <button
                    class="nav-item"
                    :class="{ active: activePage === 'patients' }"
                    @click="goTo('patients')"
                >
                    <span class="nav-dot" />
                    Пациенты
                </button>
            </div>
        </div>
    </aside>

    <div v-if="isOpen" class="sidebar-overlay" @click="emit('close')" />
</template>

<style scoped>
.sidebar {
    width: var(--sidebar-w);
    background: var(--surface);
    border-right: 1px solid var(--border);
    display: flex;
    flex-direction: column;
    overflow-y: auto;
    flex-shrink: 0;
    transition: transform 0.25s ease;
}

.nav-root {
    padding: 8px 0 16px;
}

.nav-section {
    padding: 16px 0 4px;
}

.nav-section-label {
    padding: 0 20px 8px;
    font-size: 9px;
    font-weight: 700;
    letter-spacing: 0.12em;
    text-transform: uppercase;
    color: var(--text3);
}

.nav-item {
    display: flex;
    align-items: center;
    gap: 10px;
    width: 100%;
    padding: 8px 20px;
    cursor: pointer;
    color: var(--text2);
    font-size: 13px;
    font-weight: 500;
    font-family: var(--font);
    background: none;
    border: none;
    border-left: 2px solid transparent;
    text-align: left;
    transition:
        color 0.12s,
        background 0.12s,
        border-color 0.12s;
    position: relative;
}

.nav-item:hover {
    background: var(--surface2);
    color: var(--text);
}

.nav-item.active {
    background: var(--accent-dim);
    color: var(--accent);
    border-left-color: var(--accent);
    font-weight: 600;
}

/* Точка-индикатор */
.nav-dot {
    width: 5px;
    height: 5px;
    border-radius: 50%;
    background: var(--border2);
    flex-shrink: 0;
    transition:
        background 0.12s,
        transform 0.12s;
}

.nav-item.active .nav-dot {
    background: var(--accent);
    transform: scale(1.3);
}

.nav-item:hover .nav-dot {
    background: var(--text3);
}

/* Разделитель между секциями */
.nav-section + .nav-section {
    border-top: 1px solid var(--border);
    margin-top: 4px;
}

.sidebar-overlay {
    display: none;
    position: fixed;
    inset: 0;
    background: rgba(0, 0, 0, 0.5);
    z-index: 99;
}

@media (max-width: 768px) {
    .sidebar {
        position: fixed;
        left: 0;
        top: var(--topbar-h);
        bottom: 0;
        z-index: 100;
        transform: translateX(-100%);
    }
    .sidebar.open {
        transform: translateX(0);
    }
    .sidebar-overlay {
        display: block;
    }
}
</style>
