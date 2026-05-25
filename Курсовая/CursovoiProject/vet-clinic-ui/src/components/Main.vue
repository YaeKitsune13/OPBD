<script setup lang="ts">
import { ref, computed } from "vue";
import TopBar from "./elements/TopBar.vue";
import SideBar from "./elements/SideBar.vue";
import Profile from "./elements/Profile.vue";

import DashboardPage from "./user/DashboardPage.vue";
import PetsPage from "./user/PetsPage.vue";
import BookPage from "./user/BookPage.vue";
import AppointmentsPage from "./user/AppointmentsPage.vue";
import StatsPage from "./user/StatsPage.vue";
import ShopPage from "./user/ShopPage.vue";

import SchedulePage from "./doctor/SchedulePage.vue";
import PatientsPage from "./doctor/PatientsPage.vue";
import OrdersPage from "./user/OrdersPage.vue";

const props = defineProps<{ userRole: "client" | "doctor" }>();
const emit = defineEmits<{ logout: [] }>();

const mobileMenu = ref(false);
const isProfileOpen = ref(false);
const currentPage = ref(props.userRole === "client" ? "dashboard" : "schedule");

const currentUser = computed(() => {
    try {
        const u = JSON.parse(localStorage.getItem("user") || "{}");
        const l = u.lastName?.[0] || "";
        const f = u.firstName?.[0] || "";
        return {
            name:
                `${u.lastName || ""} ${u.firstName || ""}`.trim() ||
                "Пользователь",
            avatar: (l + f).toUpperCase() || "??",
        };
    } catch {
        return { name: "Пользователь", avatar: "??" };
    }
});
</script>

<template>
    <TopBar
        :user-name="currentUser.name"
        :user-avatar="currentUser.avatar"
        :current-role="userRole"
        @toggle-sidebar="mobileMenu = !mobileMenu"
        @open-profile="isProfileOpen = true"
    />

    <div class="layout">
        <SideBar
            :current-role="userRole"
            :active-page="currentPage"
            :is-open="mobileMenu"
            @navigate="(p) => (currentPage = p)"
            @close="mobileMenu = false"
        />

        <main class="main">
            <div class="page-container">
                <template v-if="userRole === 'client'">
                    <DashboardPage
                        v-if="currentPage === 'dashboard'"
                        @navigate="(p) => (currentPage = p)"
                    />
                    <PetsPage v-else-if="currentPage === 'pets'" />
                    <BookPage v-else-if="currentPage === 'book'" />
                    <AppointmentsPage
                        v-else-if="currentPage === 'appointments'"
                    />
                    <StatsPage v-else-if="currentPage === 'stats'" />
                    <ShopPage v-else-if="currentPage === 'shop'" />
                    <OrdersPage v-else-if="currentPage === 'orders'" />
                </template>

                <template v-else-if="userRole === 'doctor'">
                    <SchedulePage v-if="currentPage === 'schedule'" />
                    <PatientsPage v-else-if="currentPage === 'patients'" />
                </template>

                <div v-else class="page">
                    <div class="page-header">
                        <div>
                            <div class="page-title">{{ currentPage }}</div>
                            <div class="page-sub">
                                Страница не найдена или недоступна
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </main>
    </div>

    <Profile
        :is-open="isProfileOpen"
        @close="isProfileOpen = false"
        @logout="emit('logout')"
    />
</template>
