<script setup>
import { ref, onMounted } from "vue";

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

function statusLabel(status) {
    switch (status) {
        case "confirmed":
            return "Подтверждено";
        case "waiting":
            return "Ожидание";
        case "rejected":
            return "Отклонено";
        default:
            return status;
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
        if (schedRes.ok) {
            appointments.value = (await schedRes.json()) ?? [];
        } else {
            showToast("Ошибка загрузки расписания");
        }
    } catch (e) {
        console.error(e);
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
            // ← было appointment_id, теперь appointmentId
            const item = appointments.value.find((a) => a.appointmentId === id);
            if (item) item.status = status;
        } else {
            showToast("Ошибка обновления статуса");
        }
    } catch (e) {
        console.error(e);
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
            <div v-if="loading" class="empty">
                <div class="empty-icon">⏳</div>
                <div>Загрузка...</div>
            </div>

            <div v-else-if="appointments.length === 0" class="empty">
                <div class="empty-icon">🎉</div>
                <div>На сегодня записей нет</div>
            </div>

            <div v-else class="schedule-list">
                <div
                    v-for="app in appointments"
                    :key="app.appointmentId"
                    class="schedule-item"
                    :class="{ 'item-waiting': app.status === 'waiting' }"
                >
                    <!-- Время -->
                    <div class="schedule-time">{{ app.time }}</div>

                    <!-- Инфо -->
                    <div class="schedule-info">
                        <div class="schedule-name">
                            {{ app.petLabel }} — {{ app.ownerName }}
                        </div>
                        <div class="schedule-sub">
                            {{ app.breed }}
                            <template v-if="app.reason"
                                >· {{ app.reason }}</template
                            >
                        </div>
                    </div>

                    <!-- Статус -->
                    <span class="badge" :class="badgeClass(app.status)">
                        {{ statusLabel(app.status) }}
                    </span>

                    <!-- Действия -->
                    <div class="schedule-actions">
                        <template v-if="app.status === 'confirmed'">
                            <button class="btn btn-primary btn-sm">
                                Принять
                            </button>
                        </template>
                        <template v-else-if="app.status === 'waiting'">
                            <button
                                class="btn btn-sm btn-ghost action-confirm"
                                title="Подтвердить"
                                @click="
                                    updateStatus(app.appointmentId, 'confirmed')
                                "
                            >
                                ✓
                            </button>
                            <button
                                class="btn btn-sm btn-ghost action-reject"
                                title="Отклонить"
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

        <div v-if="toastVisible" class="toast">{{ toastMessage }}</div>
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
    padding: 12px 8px;
    border-bottom: 1px solid var(--border);
    border-radius: 6px;
    transition: background 0.15s;
}

.schedule-item:last-child {
    border-bottom: none;
}

.schedule-item:hover {
    background: var(--bg2);
}

.item-waiting {
    opacity: 0.85;
}

.schedule-time {
    font-family: monospace;
    font-size: 15px;
    font-weight: 600;
    color: var(--accent);
    min-width: 48px;
}

.schedule-info {
    flex: 1;
    min-width: 0;
}

.schedule-name {
    font-weight: 600;
    font-size: 14px;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}

.schedule-sub {
    font-size: 12px;
    color: var(--text3);
    margin-top: 2px;
}

.schedule-actions {
    display: flex;
    gap: 4px;
    margin-left: 8px;
}

.action-confirm {
    color: var(--accent);
    font-size: 16px;
}

.action-reject {
    color: var(--red);
    font-size: 16px;
}

.toast {
    position: fixed;
    bottom: 24px;
    right: 24px;
    background: var(--text1, #222);
    color: #fff;
    padding: 12px 20px;
    border-radius: 8px;
    z-index: 9999;
    font-size: 14px;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}
</style>
