<script setup>
import { ref, reactive } from 'vue'

const props = defineProps({
  isOpen: Boolean,
  user: Object, // { name, avatar, role, email, phone }
})

const emit = defineEmits(['close', 'logout'])

// Состояния внутренних секций
const activeSection = ref('view') // 'view', 'edit', 'password'

// Данные формы редактирования (копия данных пользователя)
const editForm = reactive({
  lastName: 'Иванов',
  firstName: 'Иван',
  middleName: 'Иванович',
  phone: '+79001234567',
  email: 'ivanov@mail.ru',
  address: 'г. Москва, ул. Ленина, д. 1',
})

function saveProfile() {
  alert('Профиль обновлён!')
  activeSection.value = 'view'
}

function changePassword() {
  alert('Пароль изменён!')
  activeSection.value = 'view'
}

function handleLogout() {
  emit('logout')
}
</script>

<template>
  <div class="profile-panel" :class="{ open: isOpen }">
    <!-- Шапка панели -->
    <div class="profile-header">
      <span class="panel-title">Профиль</span>
      <button class="profile-close" @click="emit('close')">✕</button>
    </div>

    <!-- Основная инфо (Аватар и Имя) -->
    <div class="profile-avatar-wrap">
      <div class="profile-avatar-big">{{ user.avatar }}</div>
      <div class="profile-name">{{ user.name }}</div>
      <div class="profile-role-badge">{{ user.roleName }}</div>
    </div>

    <div class="profile-body">
      <!-- РЕЖИМ ПРОСМОТРА -->
      <template v-if="activeSection === 'view'">
        <div>
          <div class="profile-section-title">Контакты</div>
          <div class="profile-info-row">
            <span class="profile-info-label">📧 Email</span>
            <span class="profile-info-val">{{ user.email }}</span>
          </div>
          <div class="profile-info-row">
            <span class="profile-info-label">📱 Телефон</span>
            <span class="profile-info-val">{{ user.phone }}</span>
          </div>
        </div>

        <div>
          <div class="profile-section-title">Статистика</div>
          <div class="profile-stat-grid">
            <div class="profile-stat">
              <div class="profile-stat-val">12</div>
              <div class="profile-stat-label">Визитов</div>
            </div>
            <div class="profile-stat">
              <div class="profile-stat-val">3</div>
              <div class="profile-stat-label">Питомца</div>
            </div>
          </div>
        </div>

        <button class="profile-logout" @click="handleLogout">↩ Выйти из системы</button>
      </template>

      <!-- РЕЖИМ РЕДАКТИРОВАНИЯ -->
      <div v-else-if="activeSection === 'edit'" class="edit-section">
        <div class="profile-section-title">Редактировать профиль</div>
        <div class="form-group-mini">
          <label>Фамилия</label>
          <input v-model="editForm.lastName" type="text" />
        </div>
        <div class="form-group-mini">
          <label>Имя</label>
          <input v-model="editForm.firstName" type="text" />
        </div>
        <div class="form-group-mini">
          <label>Телефон</label>
          <input v-model="editForm.phone" type="text" />
        </div>
        <div class="btn-group">
          <button class="btn btn-primary btn-sm" @click="saveProfile">Сохранить</button>
          <button class="btn btn-ghost btn-sm" @click="activeSection = 'view'">Отмена</button>
        </div>
      </div>

      <!-- РЕЖИМ СМЕНЫ ПАРОЛЯ -->
      <div v-else-if="activeSection === 'password'" class="edit-section">
        <div class="profile-section-title">Смена пароля</div>
        <div class="form-group-mini">
          <label>Текущий пароль</label>
          <input type="password" />
        </div>
        <div class="form-group-mini">
          <label>Новый пароль</label>
          <input type="password" />
        </div>
        <div class="btn-group">
          <button class="btn btn-primary btn-sm" @click="changePassword">Обновить</button>
          <button class="btn btn-ghost btn-sm" @click="activeSection = 'view'">Отмена</button>
        </div>
      </div>
    </div>

    <!-- Кнопки действий в футере (видны только в режиме просмотра) -->
    <div class="profile-actions" v-if="activeSection === 'view'">
      <button class="btn btn-ghost btn-sm" @click="activeSection = 'edit'">✏ Редактировать</button>
      <button class="btn btn-ghost btn-sm" @click="activeSection = 'password'">🔒 Пароль</button>
    </div>
  </div>

  <!-- Оверлей для закрытия панели кликом вне её -->
  <div v-if="isOpen" class="profile-overlay" @click="emit('close')"></div>
</template>

<style scoped>
.profile-panel {
  position: fixed;
  top: var(--topbar-h);
  right: 0;
  width: 320px;
  height: calc(100vh - var(--topbar-h));
  background: var(--surface);
  border-left: 1px solid var(--border);
  z-index: 5000;
  transform: translateX(100%);
  transition: transform 0.3s ease;
  display: flex;
  flex-direction: column;
  overflow-y: auto;
}

.profile-panel.open {
  transform: translateX(0);
  box-shadow: -10px 0 30px rgba(0, 0, 0, 0.3);
}

.profile-header {
  padding: 16px 20px;
  border-bottom: 1px solid var(--border);
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.panel-title {
  font-size: 13px;
  font-weight: 600;
}
.profile-close {
  background: none;
  border: none;
  color: var(--text3);
  cursor: pointer;
  font-size: 16px;
}

.profile-avatar-wrap {
  padding: 24px 20px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  border-bottom: 1px solid var(--border);
}
.profile-avatar-big {
  width: 64px;
  height: 64px;
  border-radius: 50%;
  background: var(--accent-dim);
  border: 2px solid var(--accent);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 22px;
  font-weight: 700;
  color: var(--accent);
}
.profile-name {
  font-size: 16px;
  font-weight: 600;
  text-align: center;
}
.profile-role-badge {
  padding: 2px 10px;
  background: var(--surface2);
  border-radius: 20px;
  font-size: 10px;
  font-weight: 700;
  color: var(--text3);
  text-transform: uppercase;
}

.profile-body {
  padding: 20px;
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 24px;
}
.profile-section-title {
  font-size: 10px;
  font-weight: 700;
  text-transform: uppercase;
  color: var(--text3);
  letter-spacing: 0.1em;
  margin-bottom: 8px;
}
.profile-info-row {
  display: flex;
  justify-content: space-between;
  padding: 8px 0;
  border-bottom: 1px solid var(--border);
  font-size: 13px;
}
.profile-info-label {
  color: var(--text2);
}

.profile-stat-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 10px;
}
.profile-stat {
  background: var(--surface2);
  padding: 12px;
  border-radius: 8px;
  text-align: center;
}
.profile-stat-val {
  font-size: 20px;
  font-weight: 700;
  color: var(--accent);
  font-family: var(--mono);
}
.profile-stat-label {
  font-size: 11px;
  color: var(--text3);
}

.profile-logout {
  margin-top: 20px;
  width: 100%;
  padding: 10px;
  background: var(--red-dim);
  color: var(--red);
  border: 1px solid var(--red);
  border-radius: 6px;
  cursor: pointer;
  font-weight: 600;
  font-size: 13px;
  transition: 0.2s;
}
.profile-logout:hover {
  background: var(--red);
  color: white;
}

.profile-actions {
  padding: 16px 20px;
  border-top: 1px solid var(--border);
  display: flex;
  gap: 8px;
}
.profile-actions button {
  flex: 1;
  font-size: 11px;
}

.edit-section {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
.form-group-mini {
  display: flex;
  flex-direction: column;
  gap: 4px;
}
.form-group-mini label {
  font-size: 11px;
  color: var(--text3);
}
.form-group-mini input {
  padding: 6px 10px;
  font-size: 13px;
}
.btn-group {
  display: flex;
  gap: 8px;
  margin-top: 10px;
}

.profile-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.4);
  z-index: 4999;
}
</style>
