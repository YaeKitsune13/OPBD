<script setup>
import { ref, reactive } from 'vue'

const props = defineProps({
  // Можно передать извне, показывать ли экран (например, если токен протух)
  modelValue: Boolean,
})

const emit = defineEmits(['update:modelValue', 'login-success'])

// --- СОСТОЯНИЕ ---
const activeTab = ref('login')

const loginData = reactive({
  email: 'ivanov@mail.ru',
  pass: 'password123',
})

const regData = reactive({
  lastName: '',
  firstName: '',
  phone: '',
  email: '',
  pass: '',
})

// --- ЛОГИКА ---
function doLogin() {
  if (!loginData.email || !loginData.pass) {
    alert('Введите почту и пароль') // Позже заменишь на красивый Toast
    return
  }

  // Твоя логика имитации ролей
  let role = 'client'
  if (loginData.email.includes('doctor')) role = 'doctor'
  if (loginData.email.includes('admin')) role = 'admin'

  finalizeLogin(role)
}

function doRegister() {
  if (!regData.firstName || !regData.email) {
    alert('Заполните обязательные поля')
    return
  }
  finalizeLogin('client')
}

function quickLogin(role) {
  finalizeLogin(role)
}

function finalizeLogin(role) {
  // 1. Сообщаем родителю роль
  emit('login-success', role)
  // 2. Закрываем окно (через v-model если используешь в App.vue)
  emit('update:modelValue', false)

  // Очистка паролей для безопасности
  loginData.pass = ''
  regData.pass = ''
}
</script>

<template>
  <!-- Используем props.modelValue для контроля видимости из App.vue -->
  <div v-if="modelValue" class="auth-screen">
    <div class="auth-logo">
      <div class="auth-logo-icon">🐾</div>
      <span class="auth-logo-text">ВетКлиника</span>
    </div>

    <div class="auth-box">
      <div class="auth-tabs">
        <div
          class="auth-tab"
          :class="{ active: activeTab === 'login' }"
          @click="activeTab = 'login'"
        >
          Вход
        </div>
        <div class="auth-tab" :class="{ active: activeTab === 'reg' }" @click="activeTab = 'reg'">
          Регистрация
        </div>
      </div>

      <!-- ФОРМА ВХОДА -->
      <div v-if="activeTab === 'login'" class="auth-form">
        <div class="form-group">
          <label>Email</label>
          <input v-model="loginData.email" type="email" placeholder="example@mail.ru" />
        </div>
        <div class="form-group">
          <label>Пароль</label>
          <input
            v-model="loginData.pass"
            type="password"
            placeholder="Минимум 8 символов"
            @keyup.enter="doLogin"
          />
        </div>
        <button class="auth-submit" @click="doLogin">Войти</button>
        <div class="auth-footer-text">
          Нет аккаунта?
          <span class="auth-link" @click="activeTab = 'reg'">Зарегистрироваться</span>
        </div>
      </div>

      <!-- ФОРМА РЕГИСТРАЦИИ -->
      <div v-else class="auth-form">
        <div class="form-row">
          <div class="form-group">
            <label>Фамилия</label>
            <input v-model="regData.lastName" type="text" placeholder="Иванов" />
          </div>
          <div class="form-group">
            <label>Имя</label>
            <input v-model="regData.firstName" type="text" placeholder="Иван" />
          </div>
        </div>
        <div class="form-group">
          <label>Телефон</label>
          <input v-model="regData.phone" type="text" placeholder="+79001234567" />
        </div>
        <div class="form-group">
          <label>Email</label>
          <input v-model="regData.email" type="email" placeholder="example@mail.ru" />
        </div>
        <div class="form-group">
          <label>Пароль</label>
          <input v-model="regData.pass" type="password" placeholder="Минимум 8 символов" />
        </div>
        <button class="auth-submit" @click="doRegister">Зарегистрироваться</button>
        <div class="auth-footer-text">
          Уже есть аккаунт?
          <span class="auth-link" @click="activeTab = 'login'">Войти</span>
        </div>
      </div>
    </div>

    <!-- БАР ДЛЯ РАЗРАБОТКИ -->
    <div class="auth-devbar">
      <span class="auth-devbar-label">dev / быстрый вход:</span>
      <button class="dev-btn client" @click="quickLogin('client')">👤 Клиент</button>
      <button class="dev-btn doctor" @click="quickLogin('doctor')">🩺 Врач</button>
      <button class="dev-btn admin" @click="quickLogin('admin')">🔑 Администратор</button>
    </div>
  </div>
</template>

<style scoped>
.auth-footer-text {
  text-align: center;
  font-size: 12px;
  color: var(--text3);
}
.auth-link {
  color: var(--accent);
  cursor: pointer;
}
.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 10px;
}
/* ── AUTH SCREEN ── */
.auth-screen {
  position: fixed;
  inset: 0;
  background: var(--bg);
  z-index: 9000;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 20px;
}

.auth-screen.hidden {
  display: none;
}

.auth-logo {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 32px;
}

.auth-logo-icon {
  width: 40px;
  height: 40px;
  background: var(--accent);
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
}

.auth-logo-text {
  font-size: 20px;
  font-weight: 700;
  letter-spacing: -0.5px;
}

.auth-box {
  width: 100%;
  max-width: 380px;
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: 10px;
  overflow: hidden;
}

.auth-tabs {
  display: flex;
  border-bottom: 1px solid var(--border);
}

.auth-tab {
  flex: 1;
  padding: 14px;
  text-align: center;
  font-size: 13px;
  font-weight: 600;
  color: var(--text2);
  cursor: pointer;
  transition: all 0.15s;
  border-bottom: 2px solid transparent;
  margin-bottom: -1px;
}

.auth-tab.active {
  color: var(--accent);
  border-bottom-color: var(--accent);
}
.auth-tab:hover:not(.active) {
  color: var(--text);
}

.auth-form {
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.auth-form .form-group {
  display: flex;
  flex-direction: column;
  gap: 5px;
}

.auth-form input {
  background: var(--surface2);
  border: 1px solid var(--border);
  border-radius: var(--radius);
  color: var(--text);
  font-family: var(--font);
  font-size: 13px;
  padding: 10px 12px;
  outline: none;
  transition: border-color 0.15s;
  width: 100%;
}

.auth-form input:focus {
  border-color: var(--accent);
}

.auth-form label {
  font-size: 11px;
  font-weight: 600;
  color: var(--text3);
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.auth-submit {
  width: 100%;
  padding: 11px;
  background: var(--accent);
  color: #000;
  font-family: var(--font);
  font-size: 14px;
  font-weight: 600;
  border: none;
  border-radius: var(--radius);
  cursor: pointer;
  transition: background 0.15s;
  margin-top: 4px;
}

.auth-submit:hover {
  background: #22c55e;
}

/* dev quick-login buttons at bottom */
.auth-devbar {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  background: var(--surface);
  border-top: 1px solid var(--border);
  padding: 10px 20px;
  display: flex;
  align-items: center;
  gap: 10px;
  z-index: 9100;
  flex-wrap: wrap;
}

.auth-devbar-label {
  font-size: 11px;
  color: var(--text3);
  font-family: var(--mono);
  white-space: nowrap;
}

.dev-btn {
  padding: 6px 14px;
  border-radius: var(--radius);
  font-family: var(--font);
  font-size: 12px;
  font-weight: 600;
  cursor: pointer;
  border: 1px solid var(--border2);
  background: var(--surface2);
  color: var(--text2);
  transition: all 0.15s;
}

.dev-btn:hover {
  color: var(--text);
  border-color: var(--accent);
}
.dev-btn.client {
  border-color: var(--accent);
  color: var(--accent);
}
.dev-btn.doctor {
  border-color: var(--blue);
  color: var(--blue);
}
.dev-btn.admin {
  border-color: var(--yellow);
  color: var(--yellow);
}
</style>
