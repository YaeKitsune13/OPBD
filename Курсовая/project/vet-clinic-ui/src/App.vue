<script setup lang="ts">
import { ref, onMounted } from "vue";
import Auth from "./components/screens/Auth.vue";
import Main from "./components/Main.vue";
import ToastContainer from './components/ui/ToastContainer.vue'

// --- СОСТОЯНИЕ ---
const isAuthenticated = ref(false);
const userRole = ref("");
const userName = ref("");
const userId = ref<number | null>(null);

// --- ЛОГИКА ПРИ ЗАГРУЗКЕ ---
// Проверяем, есть ли сохраненный токен и данные пользователя в браузере
onMounted(() => {
    const token = localStorage.getItem("token");
    const savedUser = localStorage.getItem("user");

    if (token && savedUser) {
        try {
            const user = JSON.parse(savedUser);
            userRole.value = user.role;
            userName.value = user.name;
            userId.value = user.id;
            isAuthenticated.value = true;
        } catch (e) {
            // Если данные в localStorage повреждены — очищаем всё
            handleLogout();
        }
    }
});

// --- ФУНКЦИИ ---

// Вызывается, когда Auth.vue делает emit("login-success", role)
function handleLogin(role: string) {
    // После вызова этой функции сработает onMounted (если обновлять)
    // или мы можем обновить состояние напрямую из localStorage
    const savedUser = localStorage.getItem("user");
    if (savedUser) {
        const user = JSON.parse(savedUser);
        userRole.value = user.role;
        userName.value = user.name;
        userId.value = user.id;
    } else {
        userRole.value = role;
    }

    isAuthenticated.value = true;
}

// Функция для выхода
function handleLogout() {
    // 1. Очищаем хранилище браузера
    localStorage.removeItem("token");
    localStorage.removeItem("user");

    // 2. Сбрасываем состояние в коде
    isAuthenticated.value = false;
    userRole.value = "";
    userName.value = "";
    userId.value = null;
}
</script>

<template>
    <!-- Если НЕ авторизован — показываем экран входа -->
    <!-- v-model тут связывается с isVisible внутри Auth (если ты его оставил) -->
    <Auth
        v-if="!isAuthenticated"
        :modelValue="true"
        @login-success="handleLogin"
    />

    <!-- Если авторизован — показываем основной интерфейс -->
    <Main
        v-else
        :user-role="userRole"
        :user-name="userName"
        @logout="handleLogout"
    />
    <ToastContainer />
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
    font-family: "Onest", sans-serif;
    margin: 0;
}
</style>
