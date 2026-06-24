<script setup lang="ts">
import { ref, onMounted } from "vue";
const emit = defineEmits<{ navigate: [page: string] }>();

interface DashboardData {
    petsCount: number;
    appointmentsCount: number;
    visitsCount: number;
    nextAppointment: { date: string; time: string; petName: string } | null;
    recentAppointments: {
        id: number;
        petName: string;
        service: string;
        date: string;
    }[];
    pets: { id: number; name: string; avatar: string; breed: string }[];
}

const loading = ref(true);
const data = ref<DashboardData | null>(null);
const error = ref(false);

const user = JSON.parse(localStorage.getItem("user") || "{}");
const userName =
    `${user.lastName || ""} ${user.firstName || ""}`.trim() || "Пользователь";

async function loadDashboard() {
    loading.value = true;
    error.value = false;
    try {
        const res = await fetch(`/api/dashboard/${user.id}`, {
            headers: {
                Authorization: `Bearer ${localStorage.getItem("token")}`,
                "Content-Type": "application/json",
            },
        });
        if (res.ok) data.value = await res.json();
        else error.value = true;
    } catch {
        error.value = true;
    } finally {
        loading.value = false;
    }
}

onMounted(loadDashboard);
</script>

<template>
    <div class="page">
        <div class="page-header">
            <div>
                <h1 class="page-title">Добрый день, {{ userName }}!</h1>
                <p class="page-sub">
                    {{
                        new Date().toLocaleDateString("ru-RU", {
                            weekday: "long",
                            day: "numeric",
                            month: "long",
                        })
                    }}
                </p>
            </div>
            <button class="btn btn-primary" @click="emit('navigate', 'book')">
                ✚ Записаться
            </button>
        </div>

        <div v-if="loading" class="text-muted py-40 center">
            Загрузка сводки...
        </div>
        <div v-else-if="error" class="card py-40 center">
            <p class="text-red">Ошибка загрузки данных</p>
            <button class="btn btn-ghost btn-sm mt-12" @click="loadDashboard">
                Повторить
            </button>
        </div>

        <template v-else-if="data">
            <div class="stats-grid">
                <div class="stat-card">
                    <div class="stat-label">Мои питомцы</div>
                    <div class="stat-value text-accent">
                        {{ data.petsCount }}
                    </div>
                </div>
                <div class="stat-card">
                    <div class="stat-label">Активные записи</div>
                    <div class="stat-value">{{ data.appointmentsCount }}</div>
                </div>
                <div class="stat-card">
                    <div class="stat-label">Всего визитов</div>
                    <div class="stat-value">{{ data.visitsCount }}</div>
                </div>
                <div class="stat-card">
                    <div class="stat-label">Ближайший прием</div>
                    <div class="stat-value mono" style="font-size: 15px">
                        {{
                            data.nextAppointment
                                ? `${data.nextAppointment.date} в ${data.nextAppointment.time}`
                                : "Нет записей"
                        }}
                    </div>
                </div>
            </div>

            <div class="grid-2 mt-20">
                <div class="card">
                    <div class="card-header">
                        <span class="card-title">Последние записи</span>
                        <button
                            class="btn btn-ghost btn-sm"
                            @click="emit('navigate', 'appointments')"
                        >
                            Все →
                        </button>
                    </div>
                    <div class="card-body">
                        <div
                            v-if="!data.recentAppointments.length"
                            class="text-muted py-20 center"
                        >
                            История пуста
                        </div>
                        <div
                            v-for="app in data.recentAppointments"
                            :key="app.id"
                            class="list-row"
                        >
                            <div>
                                <div class="row-name">{{ app.petName }}</div>
                                <div class="row-sub">{{ app.service }}</div>
                            </div>
                            <span class="tag mono">{{ app.date }}</span>
                        </div>
                    </div>
                </div>
                <div class="card">
                    <div class="card-header">
                        <span class="card-title">Ваши питомцы</span>
                        <button
                            class="btn btn-ghost btn-sm"
                            @click="emit('navigate', 'pets')"
                        >
                            Все →
                        </button>
                    </div>
                    <div class="card-body">
                        <div
                            v-if="!data.pets.length"
                            class="text-muted py-20 center"
                        >
                            Питомцев нет
                        </div>
                        <div
                            v-for="pet in data.pets"
                            :key="pet.id"
                            class="list-row"
                        >
                            <div class="row" style="gap: 12px">
                                <span class="pet-icon-sm">{{
                                    pet.avatar
                                }}</span>
                                <div>
                                    <div class="row-name">{{ pet.name }}</div>
                                    <div class="row-sub">{{ pet.breed }}</div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </template>
    </div>
</template>

<style scoped>
.center {
    text-align: center;
}
.py-40 {
    padding: 40px 0;
}
.py-20 {
    padding: 20px 0;
}
.list-row {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 12px 0;
    border-bottom: 1px solid var(--border);
}
.list-row:last-child {
    border-bottom: none;
}
.row-name {
    font-weight: 600;
    font-size: 14px;
}
.row-sub {
    font-size: 12px;
    color: var(--text3);
}
.pet-icon-sm {
    width: 32px;
    height: 32px;
    background: var(--surface2);
    border-radius: 8px;
    display: flex;
    align-items: center;
    justify-content: center;
}
</style>
