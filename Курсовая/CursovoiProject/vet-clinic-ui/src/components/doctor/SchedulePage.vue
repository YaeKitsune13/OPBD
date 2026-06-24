<script setup lang="ts">
import { ref, reactive, computed, onMounted } from "vue";
import BaseModal from "../elements/BaseModal.vue";
import { useToast } from "../../utils/useToast";
const { showToast } = useToast();

// --- СОСТОЯНИЕ ---
const appointments = ref<any[]>([]);
const diagnosesList = ref<{ id: number; name: string }[]>([]);
const loading = ref(false);
const isReportModalOpen = ref(false);
const selectedApp = ref<any>(null);
const isDiagnosisPickerOpen = ref(false);
const availableDiagnoses = computed(() =>
    diagnosesList.value.filter((d) => !reportForm.diagnosisIds.includes(d.id)),
);

const reportForm = reactive({
    weight: "",
    diagnosisIds: [] as number[],
    treatment: "",
    medications: "",
});

const statusMap: Record<string, { label: string; cls: string }> = {
    waiting: { label: "Ожидает", cls: "badge-waiting" },
    confirmed: { label: "Принята", cls: "badge-confirmed" },
    rejected: { label: "Отклонена", cls: "badge-rejected" },
    done: { label: "Завершена", cls: "badge-info" },
};

const token = localStorage.getItem("token");
const doctorId = JSON.parse(localStorage.getItem("user") || "{}")?.id;

async function loadSchedule() {
    loading.value = true;
    try {
        const res = await fetch(`/api/doctor/schedule?doctor_id=${doctorId}`, {
            headers: {
                Authorization: `Bearer ${token}`,
            },
        });
        if (res.ok) {
            appointments.value = await res.json();
        } else {
            appointments.value = [];
        }
    } catch (err) {
        console.error("Ошибка загрузки расписания:", err);
        appointments.value = [];
    } finally {
        loading.value = false;
    }
}

async function loadDiagnoses() {
    try {
        const res = await fetch(`/api/diagnoses`, {
            headers: { Authorization: `Bearer ${token}` },
        });
        if (res.ok) {
            diagnosesList.value = await res.json();
        }
    } catch (err) {
        console.error("Ошибка загрузки справочника диагнозов:", err);
    }
}

function openCompleteModal(app: any) {
    selectedApp.value = app;
    Object.assign(reportForm, {
        weight: "",
        diagnosisIds: [],
        treatment: "",
        medications: "",
    });
    isDiagnosisPickerOpen.value = false;
    if (!diagnosesList.value.length) {
        loadDiagnoses();
    }
    isReportModalOpen.value = true;
}
async function saveProtocol() {
    if (!reportForm.weight || !reportForm.diagnosisIds.length)
        return showToast(
            "Поле 'Вес' и хотя бы один диагноз обязательны для заполнения",
            "info",
        );

    try {
        const res = await fetch(
            `/api/appointments/${selectedApp.value.id}/complete`,
            {
                method: "PATCH",
                headers: {
                    "Content-Type": "application/json",
                    Authorization: `Bearer ${token}`,
                },
                body: JSON.stringify({
                    weight: Number(reportForm.weight),
                    diagnosis_ids: reportForm.diagnosisIds,
                    treatment: reportForm.treatment,
                    medications: reportForm.medications,
                }),
            },
        );

        if (res.ok) {
            isReportModalOpen.value = false;
            await loadSchedule();
        } else {
            const err = await res.json();
            alert(err.message || "Ошибка при сохранении протокола");
        }
    } catch (err) {
        alert("Ошибка связи с сервером");
    }
}

async function updateStatus(appId: number, status: string) {
    try {
        const res = await fetch(`/api/appointments/${appId}/status`, {
            method: "PATCH",
            headers: {
                "Content-Type": "application/json",
                Authorization: `Bearer ${token}`,
            },
            body: JSON.stringify({ status }),
        });
        if (res.ok) {
            await loadSchedule();
        }
    } catch (err) {
        alert("Ошибка при обновлении статуса");
    }
}

function getDiagnosisName(id: number) {
    return diagnosesList.value.find((d) => d.id === id)?.name || "";
}

function toggleDiagnosisPicker() {
    isDiagnosisPickerOpen.value = !isDiagnosisPickerOpen.value;
}

function addDiagnosis(id: number) {
    reportForm.diagnosisIds.push(id);
    isDiagnosisPickerOpen.value = false;
}

function removeDiagnosis(id: number) {
    reportForm.diagnosisIds = reportForm.diagnosisIds.filter((x) => x !== id);
}

onMounted(loadSchedule);
</script>

<template>
    <div class="page">
        <div class="page-header">
            <div>
                <h1 class="page-title">Приём пациентов</h1>
                <p class="page-sub">
                    Управление текущими визитами и заполнение мед. карт
                </p>
            </div>
            <button
                class="btn btn-ghost"
                @click="loadSchedule"
                :disabled="loading"
            >
                {{ loading ? "Обновление..." : "🔄 Обновить" }}
            </button>
        </div>

        <div
            v-if="loading && !appointments.length"
            class="text-muted py-40 center"
        >
            Загрузка расписания...
        </div>

        <div v-else-if="!appointments.length" class="card py-40 center">
            <p class="text-muted">На сегодня записей нет</p>
        </div>

        <div v-else class="schedule-list">
            <div
                v-for="app in appointments"
                :key="app.id"
                class="card app-card-doctor"
            >
                <div class="row-between">
                    <div class="app-main-info">
                        <span class="app-time-tag">{{ app.time }}</span>
                        <span class="app-pet-name">{{ app.petName }}</span>
                        <span class="app-owner-name"
                            >({{ app.ownerName }})</span
                        >
                    </div>

                    <div class="app-actions">
                        <span
                            v-if="statusMap[app.status]"
                            :class="['badge', statusMap[app.status].cls]"
                        >
                            {{ statusMap[app.status].label }}
                        </span>

                        <!-- Действия для новых записей -->
                        <div
                            v-if="app.status === 'waiting'"
                            class="btn-group-mini"
                        >
                            <button
                                class="btn btn-primary btn-sm"
                                @click="updateStatus(app.id, 'confirmed')"
                            >
                                Принять
                            </button>
                            <button
                                class="btn btn-ghost btn-sm text-red"
                                @click="updateStatus(app.id, 'rejected')"
                            >
                                Отклонить
                            </button>
                        </div>

                        <!-- Действие для подтвержденных (завершение приема) -->
                        <button
                            v-if="app.status === 'confirmed'"
                            class="btn btn-primary btn-sm"
                            @click="openCompleteModal(app)"
                        >
                            Завершить приём
                        </button>
                    </div>
                </div>
            </div>
        </div>

        <!-- Модалка Протокола -->
        <BaseModal
            :show="isReportModalOpen"
            title="Медицинский протокол"
            @close="isReportModalOpen = false"
        >
            <div class="form-grid" v-if="selectedApp">
                <div class="form-group">
                    <label>Вес питомца (кг) *</label>
                    <input
                        v-model="reportForm.weight"
                        type="number"
                        step="0.1"
                        placeholder="Например: 5.4"
                    />
                </div>
                <div class="form-group full">
                    <label>Диагноз *</label>
                    <div class="diagnosis-box">
                        <span
                            v-for="id in reportForm.diagnosisIds"
                            :key="id"
                            class="diagnosis-tag"
                        >
                            {{ getDiagnosisName(id) }}
                            <button
                                type="button"
                                class="tag-remove"
                                @click="removeDiagnosis(id)"
                            >
                                ×
                            </button>
                        </span>

                        <div class="diagnosis-add">
                            <button
                                type="button"
                                class="tag-add-btn"
                                @click="toggleDiagnosisPicker"
                            >
                                +
                            </button>

                            <div
                                v-if="isDiagnosisPickerOpen"
                                class="diagnosis-dropdown"
                            >
                                <div
                                    v-for="d in availableDiagnoses"
                                    :key="d.id"
                                    class="diagnosis-option"
                                    @click="addDiagnosis(d.id)"
                                >
                                    {{ d.name }}
                                </div>
                                <div
                                    v-if="!availableDiagnoses.length"
                                    class="diagnosis-empty"
                                >
                                    {{
                                        diagnosesList.length
                                            ? "Все диагнозы добавлены"
                                            : "Загрузка..."
                                    }}
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
                <div class="form-group full">
                    <label>Рекомендации и лечение</label>
                    <textarea
                        v-model="reportForm.treatment"
                        rows="2"
                        placeholder="Манипуляции, советы по уходу..."
                    ></textarea>
                </div>
                <div class="form-group full">
                    <label>Назначенные препараты</label>
                    <input
                        v-model="reportForm.medications"
                        type="text"
                        placeholder="Лекарства, дозировка, длительность"
                    />
                </div>
            </div>
            <template #footer>
                <button class="btn btn-primary" @click="saveProtocol">
                    Подтвердить и закрыть визит
                </button>
            </template>
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
.text-red {
    color: var(--red);
}
.btn-group-mini {
    display: flex;
    gap: 8px;
}

.app-card-doctor {
    padding: 16px 20px;
    margin-bottom: 12px;
    border-radius: 12px;
    transition: transform 0.2s;
}
.app-card-doctor:hover {
    transform: translateX(4px);
}

.app-time-tag {
    font-family: var(--mono);
    font-weight: 800;
    color: var(--accent);
    margin-right: 15px;
    font-size: 16px;
}
.app-pet-name {
    font-weight: 700;
    font-size: 15px;
    margin-right: 5px;
}
.app-owner-name {
    color: var(--text3);
    font-size: 13px;
}
.app-actions {
    display: flex;
    align-items: center;
    gap: 15px;
}

.diagnosis-box {
    display: flex;
    flex-wrap: wrap;
    align-items: center;
    gap: 8px;
    min-height: 44px;
    padding: 8px;
    border: 1px solid var(--border, #e0e0e0);
    border-radius: 10px;
}

.diagnosis-tag {
    display: flex;
    align-items: center;
    gap: 6px;
    font-size: 13px;
    padding: 6px 10px;
    border-radius: 999px;
    background: var(--accent-soft, #e6f4ea);
    color: var(--accent, #1a8f4c);
    font-weight: 600;
}

.tag-remove {
    border: none;
    background: none;
    cursor: pointer;
    font-size: 14px;
    line-height: 1;
    color: inherit;
    opacity: 0.7;
    padding: 0;
}
.tag-remove:hover {
    opacity: 1;
}

.diagnosis-add {
    position: relative;
}

.tag-add-btn {
    width: 32px;
    height: 32px;
    border-radius: 50%;
    border: 2px solid var(--accent, #1a8f4c);
    color: var(--accent, #1a8f4c);
    background: transparent;
    font-size: 18px;
    line-height: 1;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
}
.tag-add-btn:hover {
    background: var(--accent-soft, #e6f4ea);
}

.diagnosis-dropdown {
    position: absolute;
    top: 40px;
    left: 0;
    z-index: 20;
    min-width: 220px;
    max-height: 220px;
    overflow-y: auto;
    background: var(--surface, #fff);
    border: 1px solid var(--border, #e0e0e0);
    border-radius: 10px;
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
    padding: 6px;
}

.diagnosis-option {
    padding: 8px 10px;
    border-radius: 6px;
    font-size: 13px;
    cursor: pointer;
}
.diagnosis-option:hover {
    background: var(--bg2, #f5f5f5);
}

.diagnosis-empty {
    padding: 8px 10px;
    font-size: 13px;
    color: var(--text3, #999);
}
</style>
