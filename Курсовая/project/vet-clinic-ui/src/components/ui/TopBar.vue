<script setup>
// Определяем, какие данные компонент принимает (пропсы)
defineProps({
  userName: String,
  userAvatar: String,
  currentRole: String,
})

// Определяем события, которые компонент отправляет наверх (в App.vue)
const emit = defineEmits(['update-role', 'toggle-sidebar', 'open-profile'])

const roles = [
  { id: 'client', label: 'Клиент' },
  { id: 'doctor', label: 'Врач' },
  { id: 'admin', label: 'Администрация' },
]

function handleRoleChange(roleId) {
  // Мы не меняем DOM здесь, мы просто говорим родителю: "Эй, роль изменилась!"
  emit('update-role', roleId)
}
</script>

<template>
  <div class="topbar">
    <!-- Кнопка гамбургера для мобилок -->
    <button class="hamburger" @click="emit('toggle-sidebar')">☰</button>

    <div class="logo">
      <div class="logo-icon">🐾</div>
      <span>ВетКлиника</span>
    </div>

    <div class="topbar-sep"></div>

    <!-- Переключатель ролей -->
    <div class="role-switcher">
      <button
        v-for="role in roles"
        :key="role.id"
        class="role-btn"
        :class="{ active: currentRole === role.id }"
        @click="handleRoleChange(role.id)"
      >
        {{ role.label }}
      </button>
    </div>

    <!-- Блок пользователя -->
    <div class="topbar-user" @click="emit('open-profile')">
      <div class="avatar">{{ userAvatar }}</div>
      <span class="topbar-user-name">{{ userName }}</span>
    </div>
  </div>
</template>

<style scoped>
/* Добавил scoped, чтобы стили не конфликтовали с другими компонентами */
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

.role-switcher {
  display: flex;
  gap: 2px;
  background: var(--surface2);
  border: 1px solid var(--border);
  border-radius: var(--radius);
  padding: 3px;
}

.role-btn {
  padding: 5px 14px;
  border: none;
  background: transparent;
  color: var(--text2);
  font-family: var(--font);
  font-size: 12px;
  font-weight: 500;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.15s;
}

.role-btn.active {
  background: var(--surface3);
  color: var(--text);
  border: 1px solid var(--border2);
}

.role-btn:hover:not(.active) {
  color: var(--text);
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
