<script setup>
import { ref, onMounted } from "vue";

const emit = defineEmits(["navigate"]);
const appointments = ref([]);
const doctorName = ref("");
const today = ref("");
const loading = ref(false);
const toastMessage = ref("");
const toastVisible = ref(false);

function showToast(message) {
    toastMessage.value = message;
    toastVisible.value = true;
    setTimeout(() => (toastVisible.value = false), 3000);
}

function formatToday() {
    return new Date().toLocaleDateString("ru-RU", {
        day: "numeric",
        month: "long",
        year: "numeric",
    });
}

function badgeClass(status) {
    switch (status) {
        case "confirmed":
            return "badge-confirmed";
        case "waiting":
            return "badge-waiting";
        case "rejected":
            return "badge-cancelled";
        default:
            return "badge-pending";
    }
}

async function loadSchedule() {
    loading.value = true;
    const token = localStorage.getItem("token");
    try {
        const meRes = await fetch("/api/doctors/me", {
            headers: { Authorization: `Bearer ${token}` },
        });
        if (meRes.ok) {
            const me = await meRes.json();
            doctorName.value =
                `${me.user.last_name} ${me.user.first_name} ${me.user.middle_name ?? ""}`.trim();
        }
        const schedRes = await fetch("/api/doctors/me/schedule", {
            headers: { Authorization: `Bearer ${token}` },
        });
        if (schedRes.ok) appointments.value = (await schedRes.json()) ?? [];
    } catch (e) {
        showToast("Ошибка загрузки");
    } finally {
        loading.value = false;
    }
}

async function updateStatus(id, status) {
    const token = localStorage.getItem("token");
    try {
        const res = await fetch(`/api/appointments/${id}/status`, {
            method: "PUT",
            headers: {
                Authorization: `Bearer ${token}`,
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ status }),
        });
        if (res.ok) {
            const item = appointments.value.find((a) => a.appointmentId === id);
            if (item) item.status = status;
        }
    } catch (e) {
        showToast("Ошибка");
    }
}

onMounted(() => {
    today.value = formatToday();
    loadSchedule();
});
</script>

<template>
    <div class="page">
        <div class="page-header">
            <div>
                <div class="page-title">Расписание на сегодня</div>
                <div class="page-sub">{{ today }} · {{ doctorName }}</div>
            </div>
            <div class="row">
                <span class="badge badge-info"
                    >{{ appointments.length }} записи</span
                >
            </div>
        </div>

        <div class="card">
            <div v-if="loading" class="empty">⏳ Загрузка расписания...</div>
            <div v-else-if="appointments.length === 0" class="empty">
                🎉 Сегодня приёмов нет
            </div>

            <div v-else class="schedule-list">
                <div
                    v-for="app in appointments"
                    :key="app.appointmentId"
                    class="schedule-item"
                >
                    <div class="schedule-time">{{ app.time }}</div>
                    <div class="schedule-info">
                        <div class="schedule-name">
                            {{ app.petLabel }} — {{ app.ownerName }}
                        </div>
                        <div class="schedule-sub">
                            {{ app.breed }}
                            <span v-if="app.reason" class="reason-text"
                                >· {{ app.reason }}</span
                            >
                        </div>
                    </div>

                    <span class="badge" :class="badgeClass(app.status)">{{
                        app.status === "confirmed" ? "Подтвержден" : app.status
                    }}</span>

                    <div class="schedule-actions">
                        <!-- РАЦИОНАЛЬНО: Кнопка "Принять" переводит на ведение приёма -->
                        <template v-if="app.status === 'confirmed'">
                            <button
                                class="btn btn-primary btn-sm"
                                @click="emit('navigate', 'conduct', app)"
                            >
                                Принять
                            </button>
                        </template>
                        <template v-else-if="app.status === 'waiting'">
                            <button
                                class="btn btn-sm btn-ghost"
                                @click="
                                    updateStatus(app.appointmentId, 'confirmed')
                                "
                            >
                                ✓
                            </button>
                            <button
                                class="btn btn-sm btn-ghost text-red"
                                @click="
                                    updateStatus(app.appointmentId, 'rejected')
                                "
                            >
                                ✕
                            </button>
                        </template>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped>
.schedule-list {
    padding: 4px 8px;
}
.schedule-item {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 14px 8px;
    border-bottom: 1px solid var(--border);
}
.schedule-time {
    font-family: monospace;
    font-size: 15px;
    font-weight: 700;
    color: var(--accent);
    min-width: 50px;
}
.schedule-info {
    flex: 1;
    min-width: 0;
}
.schedule-name {
    font-weight: 600;
    font-size: 14px;
}
.schedule-sub {
    font-size: 12px;
    color: var(--text3);
}
.reason-text {
    color: var(--accent);
    font-style: italic;
}
.schedule-actions {
    display: flex;
    gap: 6px;
}
.text-red {
    color: #ff4d4f;
}
</style>
