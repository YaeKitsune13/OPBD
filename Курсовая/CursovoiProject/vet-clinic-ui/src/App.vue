<script setup lang="ts">
import { ref, onMounted } from "vue";
import Auth from "./components/elements/Auth.vue";
import Main from "./components/Main.vue";
import ToastContainer from "./components/elements/ToastContainer.vue";

const isAuthenticated = ref(false);
const userRole = ref<"client" | "doctor" | null>(null);

onMounted(() => {
    const token = localStorage.getItem("token");
    const savedUser = localStorage.getItem("user");

    if (token && savedUser) {
        try {
            const user = JSON.parse(savedUser);
            userRole.value = user.role as "client" | "doctor";
            isAuthenticated.value = true;
        } catch {
            handleLogout();
        }
    }
});

function handleLogin(role: string) {
    userRole.value = role as "client" | "doctor";
    isAuthenticated.value = true;
}

function handleLogout() {
    localStorage.removeItem("token");
    localStorage.removeItem("user");
    isAuthenticated.value = false;
    userRole.value = null;
}
</script>

<template>
    <Auth v-if="!isAuthenticated" @login-success="handleLogin" />
    <Main v-else :user-role="userRole!" @logout="handleLogout" />

    <ToastContainer />
</template>
