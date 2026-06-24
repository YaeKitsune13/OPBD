<script setup lang="ts">
import { ref, onMounted } from "vue";
import BaseModal from "../elements/BaseModal.vue";

// --- ТИПЫ ---
interface Patient {
    id: number;
    fullName: string;
    phone: string;
    email: string;
    petsCount: number;
    lastVisit: string;
}
interface VisitHistory {
    petName: string;
    petIcon: string;
    visits: {
        id: number;
        date: string;
        diagnoses: { id: number; name: string }[];
        treatment: string;
    }[];
}

const search = ref("");
const patients = ref<Patient[]>([]);
const loading = ref(false);
const historyLoading = ref(false);

const selectedOwner = ref<Patient | null>(null);
const ownerHistory = ref<VisitHistory[]>([]);
const isHistoryOpen = ref(false);

const token = localStorage.getItem("token");

async function loadPatients() {
    loading.value = true;
    try {
        const res = await fetch(
            `/api/doctor/patients?search=${encodeURIComponent(search.value)}`,
            {
                headers: { Authorization: `Bearer ${token}` },
            },
        );
        if (res.ok) {
            patients.value = await res.json();
        } else {
            patients.value = [];
        }
    } catch (err) {
        console.error("Failed to load patients:", err);
        patients.value = [];
    } finally {
        loading.value = false;
    }
}

async function openHistory(owner: Patient) {
    selectedOwner.value = owner;
    isHistoryOpen.value = true;
    historyLoading.value = true;
    ownerHistory.value = [];

    try {
        const res = await fetch(`/api/doctor/patients/${owner.id}/history`, {
            headers: { Authorization: `Bearer ${token}` },
        });
        if (res.ok) {
            ownerHistory.value = await res.json();
        }
    } catch (err) {
        console.error("Failed to load history:", err);
    } finally {
        historyLoading.value = false;
    }
}

onMounted(loadPatients);
</script>

<template>
    <div class="page">
        <div class="page-header">
            <div>
                <h1 class="page-title">Картотека пациентов</h1>
                <p class="page-sub">
                    Поиск владельцев и медицинская история питомцев
                </p>
            </div>
        </div>

        <!-- Поиск -->
        <div class="search-panel card">
            <input
                v-model="search"
                type="text"
                placeholder="Поиск по ФИО или номеру телефона..."
                @input="loadPatients"
            />
        </div>

        <!-- Таблица -->
        <div class="table-wrap mt-20">
            <div v-if="loading" class="text-muted py-40 center">
                Поиск в базе данных...
            </div>

            <table v-else>
                <thead>
                    <tr>
                        <th>Владелец</th>
                        <th>Контакты</th>
                        <th>Питомцев</th>
                        <th>Последний визит</th>
                        <th></th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-if="patients.length === 0">
                        <td colspan="5" class="center py-40 text-muted">
                            Совпадений не найдено
                        </td>
                    </tr>
                    <tr v-for="p in patients" :key="p.id">
                        <td>
                            <div class="td-main">{{ p.fullName }}</div>
                        </td>
                        <td>
                            <div class="td-phone">{{ p.phone }}</div>
                            <div class="td-sub">{{ p.email }}</div>
                        </td>
                        <td>
                            <span class="tag">{{ p.petsCount }}</span>
                        </td>
                        <td>{{ p.lastVisit || "—" }}</td>
                        <td style="text-align: right">
                            <button
                                class="btn btn-ghost btn-sm"
                                @click="openHistory(p)"
                            >
                                Медкарта
                            </button>
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>

        <!-- Модалка истории (Медкарта) -->
        <BaseModal
            :show="isHistoryOpen"
            title="Медицинская история"
            @close="isHistoryOpen = false"
        >
            <div v-if="selectedOwner" class="history-content">
                <div class="owner-summary mb-20">
                    <h3>{{ selectedOwner.fullName }}</h3>
                    <p class="text-muted">{{ selectedOwner.phone }}</p>
                </div>

                <div v-if="historyLoading" class="center py-20 text-muted">
                    Загрузка медицинских записей...
                </div>

                <div
                    v-else-if="ownerHistory.length === 0"
                    class="center py-20 text-muted"
                >
                    История визитов пуста
                </div>

                <div
                    v-for="pet in ownerHistory"
                    :key="pet.petName"
                    class="pet-history-item"
                >
                    <div class="pet-header">
                        <span class="pet-icon">{{ pet.petIcon || "🐾" }}</span>
                        <span class="pet-name">{{ pet.petName }}</span>
                    </div>

                    <div class="history-timeline">
                        <div
                            v-for="visit in pet.visits"
                            :key="visit.id"
                            class="timeline-point"
                        >
                            <div class="t-date">{{ visit.date }}</div>
                            <div class="t-info">
                                <div class="t-diagnoses">
                                    <span
                                        v-for="d in visit.diagnoses"
                                        :key="d.id"
                                        class="t-diagnosis-tag"
                                    >
                                        {{ d.name }}
                                    </span>
                                    <span
                                        v-if="!visit.diagnoses.length"
                                        class="text-muted"
                                    >
                                        Диагноз не указан
                                    </span>
                                </div>
                                <p>{{ visit.treatment }}</p>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </BaseModal>
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

.search-panel {
    padding: 8px;
    margin-bottom: 24px;
}
.search-panel input {
    width: 100%;
    padding: 12px 16px;
    border: none;
    background: transparent;
    font-size: 15px;
    outline: none;
}

.td-main {
    font-weight: 700;
    color: var(--text);
}
.td-phone {
    font-weight: 600;
    font-family: var(--mono);
    font-size: 13px;
}
.td-sub {
    font-size: 12px;
    color: var(--text3);
}

.history-content {
    max-height: 70vh;
    overflow-y: auto;
    padding-right: 4px;
}

.pet-history-item {
    margin-bottom: 24px;
}

.pet-header {
    display: flex;
    align-items: center;
    gap: 10px;
    margin-bottom: 16px;
    background: var(--surface2);
    padding: 10px;
    border-radius: 8px;
}
.pet-icon {
    font-size: 20px;
}
.pet-name {
    font-weight: 800;
    font-size: 16px;
}

.history-timeline {
    padding-left: 8px;
}

.timeline-point {
    display: flex;
    gap: 16px;
    padding-left: 16px;
    border-left: 2px solid var(--border);
    margin-bottom: 20px;
    position: relative;
}
.timeline-point::before {
    content: "";
    position: absolute;
    left: -6px;
    top: 0;
    width: 10px;
    height: 10px;
    background: var(--accent);
    border-radius: 50%;
    border: 2px solid var(--surface);
}

.t-date {
    font-size: 12px;
    font-weight: 700;
    color: var(--text3);
    min-width: 85px;
    font-family: var(--mono);
}
.t-diagnoses {
    display: flex;
    flex-wrap: wrap;
    gap: 6px;
    margin-bottom: 4px;
}
.t-diagnosis-tag {
    font-size: 13px;
    font-weight: 700;
    color: var(--text);
    background: var(--accent-soft, #e6f4ea);
    color: var(--accent, #1a8f4c);
    padding: 2px 8px;
    border-radius: 999px;
}
.t-info p {
    font-size: 13px;
    color: var(--text2);
    line-height: 1.4;
}

.mb-20 {
    margin-bottom: 20px;
}
.mt-20 {
    margin-top: 20px;
}
</style>
