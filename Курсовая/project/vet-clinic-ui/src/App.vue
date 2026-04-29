<script setup>
import { ref } from 'vue'
import Auth from './components/screens/Auth.vue'
import Main from './components/Main.vue'

// Состояние: авторизован ли пользователь
const isAuthenticated = ref(false)
// Состояние: роль пользователя (получим из Auth)
const userRole = ref('')

// Функция, которая сработает при успешном входе
function handleLogin(role) {
  userRole.value = role
  isAuthenticated.value = true
}

// Функция для выхода (можно передать в Main -> TopBar)
function handleLogout() {
  isAuthenticated.value = false
  userRole.value = ''
}
</script>

<template>
  <!-- Если НЕ авторизован — показываем экран входа -->
  <!-- v-model тут связывается с isVisible внутри Auth (если ты его оставил) -->
  <Auth v-if="!isAuthenticated" :modelValue="true" @login-success="handleLogin" />

  <!-- Если авторизован — показываем основной интерфейс -->
  <Main v-else :initial-role="userRole" @logout="handleLogout" />
</template>

<style>
/* Сюда стоит перенести глобальные переменные :root из твоего HTML,
   чтобы они были доступны во всех компонентах */
:root {
  --bg: #0e0f11;
  --surface: #16181c;
  --accent: #4ade80;
  /* ... все остальные переменные ... */
  --topbar-h: 52px;
  --sidebar-w: 220px;
}

body {
  background: var(--bg);
  color: var(--text);
  font-family: 'Onest', sans-serif;
  margin: 0;
}
</style>
