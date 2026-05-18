<script setup>
import { ref, computed, onMounted, watch } from "vue";
import BaseModal from "../../ui/BaseModal.vue";
import PetCombobox from "../../ui/PetCombobox.vue";

const props = defineProps({
    initialData: Object,
});

const emit = defineEmits(["navigate"]);

// Состояние приёма
const selectedPet = ref(null);
const allServices = ref([]);
const anamnesis = ref("");
const diagnosis = ref("");
const assignments = ref([]);

// UI Состояние
const isRegModalOpen = ref(false);
const isServiceModalOpen = ref(false);
const selectedServiceId = ref(null);
const loading = ref(false);
const toastMessage = ref("");
const toastVisible = ref(false);

// Данные для быстрого создания питомца (Владелец теперь всегда ID=1 на бэкенде)
const regData = ref({
    pet_name: "",
    species: "Кот",
    breed: "",
});

function showToast(message) {
    toastMessage.value = message;
    toastVisible.value = true;
    setTimeout(() => (toastVisible.value = false), 3000);
}

onMounted(async () => {
    await loadServices();
    if (props.initialData) {
        // Если переход из расписания
        if (props.initialData.appointmentId) {
            selectedPet.value = {
                id: props.initialData.appointmentId,
                pet_id: props.initialData.petId || 0,
                name: props.initialData.petLabel,
                owner_name: props.initialData.ownerName,
            };
            anamnesis.value = props.initialData.reason || "";
        }
        // Если переход из поиска
        else if (props.initialData.pet_id) {
            selectedPet.value = {
                id: 0,
                pet_id: props.initialData.pet_id,
                name: props.initialData.pet_name,
                owner_name: props.initialData.owner_name,
            };
        }
    }
});

// РАЦИОНАЛЬНО: Создание питомца (обычное или полностью анонимное)
async function quickRegister(isAnon = false) {
    loading.value = true;

    // Подготавливаем данные: если анонимно, шлем пустые строки
    const payload = isAnon
        ? {
              pet_name: "",
              species: "Другое",
              breed: "Аноним",
              is_anonymous: true,
          }
        : { ...regData.value, is_anonymous: true };

    try {
        const response = await fetch("/api/pets/quick-register", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                Authorization: `Bearer ${localStorage.getItem("token")}`,
            },
            body: JSON.stringify(payload),
        });

        if (response.ok) {
            const data = await response.json();
            // Выбираем созданного питомца
            selectedPet.value = {
                id: 0, // Записи в расписании нет
                pet_id: data.pet_id,
                name: data.pet_name,
                owner_name: data.owner_name,
            };
            isRegModalOpen.value = false;
            showToast(isAnon ? "Анонимный приём начат" : "Пациент добавлен");
        } else {
            showToast("Ошибка при создании записи");
        }
    } catch (e) {
        showToast("Ошибка сети");
    } finally {
        loading.value = false;
    }
}

async function loadServices() {
    try {
        const res = await fetch(`/api/services`, {
            headers: {
                Authorization: `Bearer ${localStorage.getItem("token")}`,
            },
        });
        if (res.ok) allServices.value = (await res.json()) ?? [];
    } catch (e) {
        console.error(e);
    }
}

function addService() {
    if (!selectedServiceId.value) return;
    const s = allServices.value.find((s) => s.id === selectedServiceId.value);
    if (s) {
        const existing = assignments.value.find((a) => a.id === s.id);
        if (existing) existing.qty++;
        else
            assignments.value.push({
                id: s.id,
                name: s.name,
                price: s.price,
                qty: 1,
            });
        isServiceModalOpen.value = false;
        selectedServiceId.value = null;
    }
}

async function saveVisit() {
    if (
        !selectedPet.value ||
        !diagnosis.value.trim() ||
        assignments.value.length === 0
    ) {
        return showToast("Заполните все данные приёма");
    }
    loading.value = true;
    const payload = {
        selectedPet: {
            id: Number(selectedPet.value.id || 0), // AppointmentID
            pet_id: Number(selectedPet.value.pet_id), // PetID
            name: selectedPet.value.name,
            owner: selectedPet.value.owner_name,
        },
        anamnesis: anamnesis.value,
        diagnosis: diagnosis.value,
        assignments: assignments.value.map((a) => ({ ...a, type: "service" })),
        totalCost: Math.round(totalCost.value),
    };
    try {
        const res = await fetch(`/api/visits`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                Authorization: `Bearer ${localStorage.getItem("token")}`,
            },
            body: JSON.stringify(payload),
        });
        if (res.ok) {
            showToast("Приём сохранён успешно");
            setTimeout(() => emit("navigate", "today"), 1000);
        }
    } catch (e) {
        showToast("Ошибка связи");
    } finally {
        loading.value = false;
    }
}

const totalCost = computed(() =>
    assignments.value.reduce((s, i) => s + i.price * i.qty, 0),
);
</script>

<template>
    <div class="page-container">
        <!-- Header -->
        <div class="page-header">
            <div class="header-main">
                <h1 class="page-title">Приём пациента</h1>

                <div v-if="selectedPet" class="patient-badge">
                    <div class="patient-info">
                        <span class="p-name">{{ selectedPet.name }}</span>
                        <span class="p-owner">{{
                            selectedPet.owner_name
                        }}</span>
                    </div>
                    <button class="btn-icon-close" @click="selectedPet = null">
                        ✕
                    </button>
                </div>

                <template v-else>
                    <PetCombobox v-model="selectedPet" :is-doctor="true" />
                    <div class="quick-actions">
                        <button
                            class="btn btn-secondary"
                            @click="isRegModalOpen = true"
                        >
                            + Питомец
                        </button>
                        <button
                            class="btn btn-secondary"
                            @click="quickRegister(true)"
                        >
                            👤 Анонимно
                        </button>
                    </div>
                </template>
            </div>
        </div>

        <div class="conduct-grid">
            <!-- Left: Card Data -->
            <div class="card">
                <div class="card-header">
                    <h3 class="card-title">Медицинская карта</h3>
                </div>
                <div class="card-body">
                    <div class="field">
                        <label>Жалобы и анамнез</label>
                        <textarea
                            v-model="anamnesis"
                            placeholder="Опишите состояние пациента..."
                        ></textarea>
                    </div>
                    <div class="field mt-16">
                        <label>Диагноз</label>
                        <input
                            v-model="diagnosis"
                            type="text"
                            placeholder="Введите диагноз"
                        />
                    </div>
                </div>
            </div>

            <!-- Right: Bill -->
            <div class="card">
                <div class="card-header">
                    <h3 class="card-title">Оказанные услуги</h3>
                    <button
                        class="btn-text-accent"
                        @click="isServiceModalOpen = true"
                    >
                        + Услуга
                    </button>
                </div>
                <div class="table-container">
                    <table class="conduct-table">
                        <thead>
                            <tr>
                                <th>Наименование</th>
                                <th class="text-right">Всего</th>
                                <th style="width: 32px"></th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-for="item in assignments" :key="item.id">
                                <td class="td-name">
                                    {{ item.name }}
                                    <span class="text-dim"
                                        >x{{ item.qty }}</span
                                    >
                                </td>
                                <td class="td-price">
                                    {{ item.price * item.qty }} ₽
                                </td>
                                <td class="text-right">
                                    <button
                                        class="btn-delete"
                                        @click="
                                            assignments = assignments.filter(
                                                (i) => i.id !== item.id,
                                            )
                                        "
                                    >
                                        ✕
                                    </button>
                                </td>
                            </tr>
                            <tr v-if="assignments.length === 0">
                                <td colspan="3" class="empty-row">
                                    Добавьте услуги из прейскуранта
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>
                <div class="card-footer-total">
                    <span class="total-label">Итого:</span>
                    <span class="total-value">{{ totalCost }} ₽</span>
                </div>
            </div>
        </div>

        <div class="footer-actions">
            <button class="btn btn-ghost" @click="emit('navigate', 'today')">
                Отмена
            </button>
            <button
                class="btn btn-primary"
                @click="saveVisit"
                :disabled="loading || !selectedPet"
            >
                {{ loading ? "Загрузка..." : "Сохранить приём" }}
            </button>
        </div>

        <!-- Registration Modal (Only Pet info needed now) -->
        <BaseModal
            :show="isRegModalOpen"
            title="Быстрая регистрация питомца"
            @close="isRegModalOpen = false"
        >
            <div class="modal-form">
                <div class="field">
                    <label>Кличка животного</label>
                    <input
                        v-model="regData.pet_name"
                        type="text"
                        placeholder="Барсик, Рекс..."
                    />
                </div>
                <div class="grid-2 mt-12">
                    <div class="field">
                        <label>Вид</label>
                        <select v-model="regData.species">
                            <option>Кот</option>
                            <option>Собака</option>
                            <option>Птица</option>
                            <option>Экзот</option>
                        </select>
                    </div>
                    <div class="field">
                        <label>Порода</label>
                        <input
                            v-model="regData.breed"
                            type="text"
                            placeholder="Сиамская..."
                        />
                    </div>
                </div>
                <p class="hint-text">
                    Пациент будет привязан к системному аккаунту клиники.
                </p>
            </div>
            <template #footer>
                <button
                    class="btn btn-primary w-full"
                    @click="quickRegister(false)"
                    :disabled="loading"
                >
                    Добавить и выбрать
                </button>
            </template>
        </BaseModal>

        <!-- Service Modal -->
        <BaseModal
            :show="isServiceModalOpen"
            title="Выбор услуги"
            @close="isServiceModalOpen = false"
        >
            <div class="field">
                <label>Выберите позицию из списка</label>
                <select v-model.number="selectedServiceId" class="select-full">
                    <option :value="null">-- Список услуг --</option>
                    <option v-for="s in allServices" :key="s.id" :value="s.id">
                        {{ s.name }} ({{ s.price }} ₽)
                    </option>
                </select>
            </div>
            <template #footer
                ><button class="btn btn-primary w-full" @click="addService">
                    Добавить в счёт
                </button></template
            >
        </BaseModal>

        <div v-if="toastVisible" class="toast-fixed">{{ toastMessage }}</div>
    </div>
</template>

<style scoped>
.page-container {
    font-family: var(--font);
    color: var(--text);
}
.page-header {
    margin-bottom: 24px;
}
.header-main {
    display: flex;
    align-items: center;
    gap: 16px;
}
.page-title {
    font-size: 24px;
    font-weight: 700;
    margin: 0;
}
.quick-actions {
    display: flex;
    gap: 8px;
}

/* Patient Badge */
.patient-badge {
    display: flex;
    align-items: center;
    gap: 12px;
    background: var(--surface2);
    border: 1px solid var(--accent);
    padding: 6px 14px;
    border-radius: var(--radius);
    animation: fadeIn 0.2s ease;
}
.patient-info {
    display: flex;
    flex-direction: column;
}
.p-name {
    font-weight: 700;
    color: var(--accent);
    font-size: 14px;
}
.p-owner {
    font-size: 11px;
    color: var(--text3);
}
.btn-icon-close {
    background: none;
    border: none;
    color: var(--text3);
    cursor: pointer;
    padding: 4px;
}
.btn-icon-close:hover {
    color: var(--red);
}

/* Layout */
.conduct-grid {
    display: grid;
    grid-template-columns: 1fr 400px;
    gap: 20px;
}
.card {
    background: var(--surface);
    border: 1px solid var(--border);
    border-radius: var(--radius);
    overflow: hidden;
}
.card-header {
    padding: 14px 16px;
    border-bottom: 1px solid var(--border);
    display: flex;
    justify-content: space-between;
    align-items: center;
}
.card-title {
    font-size: 14px;
    font-weight: 600;
    margin: 0;
    text-transform: uppercase;
    color: var(--text2);
    letter-spacing: 0.5px;
}
.card-body {
    padding: 16px;
}

/* Forms */
label {
    display: block;
    font-size: 11px;
    font-weight: 700;
    color: var(--text3);
    text-transform: uppercase;
    margin-bottom: 6px;
}
input,
textarea,
select {
    width: 100%;
    background: var(--surface2);
    border: 1px solid var(--border2);
    border-radius: var(--radius);
    color: var(--text);
    padding: 10px 12px;
    font-family: var(--font);
    font-size: 14px;
    outline: none;
    transition: border 0.2s;
}
input:focus,
textarea:focus {
    border-color: var(--accent);
}
textarea {
    min-height: 180px;
    resize: none;
}

/* Buttons */
.btn {
    padding: 10px 16px;
    border-radius: var(--radius);
    cursor: pointer;
    font-weight: 600;
    font-size: 13px;
    transition: all 0.2s;
    border: none;
}
.btn-primary {
    background: var(--accent);
    color: var(--bg);
}
.btn-primary:hover {
    filter: brightness(1.1);
    box-shadow: 0 0 12px var(--accent-glow);
}
.btn-primary:disabled {
    opacity: 0.4;
    cursor: not-allowed;
}
.btn-secondary {
    background: var(--surface3);
    color: var(--text);
    border: 1px solid var(--border);
}
.btn-secondary:hover {
    background: var(--border);
}
.btn-ghost {
    color: var(--text2);
    background: none;
}
.btn-ghost:hover {
    color: var(--text);
}
.btn-text-accent {
    background: none;
    border: none;
    color: var(--accent);
    font-weight: 600;
    cursor: pointer;
    font-size: 12px;
}

/* Table */
.conduct-table {
    width: 100%;
    border-collapse: collapse;
}
.conduct-table th {
    text-align: left;
    padding: 12px 16px;
    color: var(--text3);
    font-size: 10px;
    text-transform: uppercase;
    border-bottom: 1px solid var(--border);
}
.conduct-table td {
    padding: 10px 16px;
    border-bottom: 1px solid var(--border);
    font-size: 13px;
}
.td-name {
    font-weight: 500;
}
.td-price {
    font-family: var(--mono);
    color: var(--accent);
    text-align: right;
    font-weight: 600;
}
.btn-delete {
    color: var(--text3);
    background: none;
    border: none;
    cursor: pointer;
}
.btn-delete:hover {
    color: var(--red);
}
.empty-row {
    text-align: center;
    color: var(--text3);
    padding: 40px !important;
    font-size: 12px;
}

/* Bill Footer */
.card-footer-total {
    padding: 16px;
    background: var(--surface2);
    display: flex;
    justify-content: space-between;
    align-items: center;
}
.total-label {
    font-weight: 600;
    color: var(--text2);
    font-size: 14px;
}
.total-value {
    font-family: var(--mono);
    font-size: 20px;
    font-weight: 700;
    color: var(--accent);
}

.footer-actions {
    margin-top: 24px;
    display: flex;
    justify-content: flex-end;
    gap: 12px;
    padding-bottom: 40px;
}
.hint-text {
    font-size: 11px;
    color: var(--text3);
    font-style: italic;
    margin-top: 12px;
    text-align: center;
}

.toast-fixed {
    position: fixed;
    bottom: 24px;
    right: 24px;
    background: var(--surface3);
    color: var(--text);
    padding: 12px 24px;
    border-radius: var(--radius);
    border-left: 4px solid var(--accent);
    box-shadow: 0 10px 30px rgba(0, 0, 0, 0.5);
    z-index: 9999;
}

.mt-16 {
    margin-top: 16px;
}
.mt-12 {
    margin-top: 12px;
}
.grid-2 {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 12px;
}
.w-full {
    width: 100%;
}
.text-dim {
    color: var(--text3);
    font-size: 11px;
    margin-left: 4px;
}

@keyframes fadeIn {
    from {
        opacity: 0;
        transform: translateY(-5px);
    }
    to {
        opacity: 1;
        transform: translateY(0);
    }
}
</style>
