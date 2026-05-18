<script setup>
import { ref, reactive, computed, onMounted } from "vue";
import BaseModal from "../../ui/BaseModal.vue";

const meds = ref([]);
const searchQuery = ref("");
const isAddModalOpen = ref(false);
const isEditModalOpen = ref(false);

const toastMessage = ref("");
const toastVisible = ref(false);

// Генерация списков для селектов
const months = [
    { value: "01", label: "Январь" },
    { value: "02", label: "Февраль" },
    { value: "03", label: "Март" },
    { value: "04", label: "Апрель" },
    { value: "05", label: "Май" },
    { value: "06", label: "Июнь" },
    { value: "07", label: "Июль" },
    { value: "08", label: "Август" },
    { value: "09", label: "Сентябрь" },
    { value: "10", label: "Октябрь" },
    { value: "11", label: "Ноябрь" },
    { value: "12", label: "Декабрь" },
];

const currentYear = new Date().getFullYear();
const years = Array.from({ length: 10 }, (_, i) => currentYear + i);

// Вспомогательные реактивные объекты для селектов
const newMedDate = reactive({ month: "01", year: String(currentYear) });
const editMedDate = reactive({ month: "01", year: String(currentYear) });

function buildExpiry(month, year) {
    return `${year}-${month}`;
}

function parseExpiry(expiry) {
    if (!expiry) return { month: "01", year: String(currentYear) };
    const [year, month] = expiry.split("-");
    return { month: month || "01", year: year || String(currentYear) };
}

const newMed = reactive({
    name: "",
    desc: "",
    price: null,
});

const editingMed = reactive({
    id: "",
    name: "",
    desc: "",
    price: null,
});

function showToast(message) {
    toastMessage.value = message;
    toastVisible.value = true;
    setTimeout(() => (toastVisible.value = false), 3000);
}

const filteredMeds = computed(() => {
    if (!Array.isArray(meds.value)) return [];
    const query = searchQuery.value.toLowerCase();
    return meds.value.filter(
        (m) =>
            m.name.toLowerCase().includes(query) ||
            String(m.id).toLowerCase().includes(query),
    );
});

async function loadMeds() {
    const token = localStorage.getItem("token");
    try {
        const response = await fetch("/api/meds", {
            headers: { Authorization: `Bearer ${token}` },
        });
        if (response.ok) {
            meds.value = await response.json();
        } else {
            showToast("Ошибка загрузки медикаментов");
        }
    } catch (e) {
        console.error(e);
        showToast("Ошибка загрузки медикаментов");
    }
}

onMounted(() => {
    loadMeds();
});

async function saveMed() {
    if (!newMed.name || !newMed.price) {
        showToast("Заполните обязательные поля");
        return;
    }

    const token = localStorage.getItem("token");
    try {
        const response = await fetch("/api/admin/meds", {
            method: "POST",
            headers: {
                Authorization: `Bearer ${token}`,
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                ...newMed,
                expiry: buildExpiry(newMedDate.month, newMedDate.year),
            }),
        });

        if (response.ok) {
            isAddModalOpen.value = false;
            Object.assign(newMed, { name: "", desc: "", price: null });
            Object.assign(newMedDate, {
                month: "01",
                year: String(currentYear),
            });
            await loadMeds();
            showToast("Медикамент добавлен");
        } else {
            const err = await response.json().catch(() => ({}));
            showToast(err.error || "Ошибка сохранения медикамента");
        }
    } catch (e) {
        console.error(e);
        showToast("Ошибка сохранения медикамента");
    }
}

function openEditModal(med) {
    editingMed.id = med.id;
    editingMed.name = med.name;
    editingMed.desc = med.desc;
    editingMed.price = med.price;
    const parsed = parseExpiry(med.expiry);
    editMedDate.month = parsed.month;
    editMedDate.year = parsed.year;
    isEditModalOpen.value = true;
}

async function saveEditedMed() {
    if (!editingMed.name || !editingMed.price) {
        showToast("Заполните обязательные поля");
        return;
    }

    const token = localStorage.getItem("token");
    try {
        const response = await fetch(`/api/admin/meds/${editingMed.id}`, {
            method: "PUT",
            headers: {
                Authorization: `Bearer ${token}`,
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                ...editingMed,
                expiry: buildExpiry(editMedDate.month, editMedDate.year),
            }),
        });

        if (response.ok) {
            isEditModalOpen.value = false;
            await loadMeds();
            showToast("Медикамент обновлён");
        } else {
            const err = await response.json().catch(() => ({}));
            showToast(err.error || "Ошибка обновления медикамента");
        }
    } catch (e) {
        console.error(e);
        showToast("Ошибка обновления медикамента");
    }
}

async function deleteMed(id) {
    if (!confirm("Удалить препарат из базы данных?")) return;

    const token = localStorage.getItem("token");
    try {
        const response = await fetch(`/api/admin/meds/${id}`, {
            method: "DELETE",
            headers: { Authorization: `Bearer ${token}` },
        });

        if (response.ok) {
            meds.value = meds.value.filter((m) => m.id !== id);
            showToast("Медикамент удалён");
        } else {
            const err = await response.json().catch(() => ({}));
            showToast(err.error || "Ошибка удаления медикамента");
        }
    } catch (e) {
        console.error(e);
        showToast("Ошибка удаления медикамента");
    }
}

async function writeOffExpired() {
    if (!confirm("Списать все просроченные препараты?")) return;

    const token = localStorage.getItem("token");
    const expired = meds.value.filter((m) => m.status === "expired");

    try {
        await Promise.all(
            expired.map((m) =>
                fetch(`/api/admin/meds/${m.id}`, {
                    method: "DELETE",
                    headers: { Authorization: `Bearer ${token}` },
                }),
            ),
        );
        await loadMeds();
        showToast(`Списано препаратов: ${expired.length}`);
    } catch (e) {
        console.error(e);
        showToast("Ошибка при списании");
    }
}
</script>

<template>
    <div class="page">
        <div class="page-header">
            <div>
                <div class="page-title">Справочник медикаментов</div>
                <div class="page-sub">
                    Управление препаратами и контроль сроков годности
                </div>
            </div>
            <button class="btn btn-primary" @click="isAddModalOpen = true">
                ✚ Добавить медикамент
            </button>
        </div>

        <div class="card">
            <div class="card-header">
                <input
                    type="text"
                    v-model="searchQuery"
                    placeholder="Поиск по названию или артикулу..."
                    style="max-width: 300px"
                />
                <button class="btn btn-danger btn-sm" @click="writeOffExpired">
                    Списать просроченные
                </button>
            </div>
            <div class="table-wrap">
                <table>
                    <thead>
                        <tr>
                            <th style="width: 80px">#</th>
                            <th>Наименование</th>
                            <th>Цена / ед.</th>
                            <th>Срок годности</th>
                            <th style="width: 120px">Действия</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr
                            v-for="m in filteredMeds"
                            :key="m.id"
                            :class="{ 'expired-row': m.status === 'expired' }"
                        >
                            <td class="mono text-muted">{{ m.id }}</td>
                            <td class="td-main">
                                <span
                                    :class="{
                                        'text-red': m.status === 'expired',
                                    }"
                                    >{{ m.name }}</span
                                >
                                <div class="med-desc">{{ m.desc }}</div>
                            </td>
                            <td class="mono">{{ m.price }} руб.</td>
                            <td
                                class="mono"
                                :class="{ 'text-red': m.status === 'expired' }"
                            >
                                {{ m.expiry }}
                                <span
                                    v-if="m.status === 'expired'"
                                    title="Срок годности истек"
                                    >⚠️</span
                                >
                            </td>
                            <td>
                                <div class="row">
                                    <button
                                        class="btn btn-ghost btn-sm"
                                        @click="openEditModal(m)"
                                        title="Редактировать"
                                    >
                                        ✏
                                    </button>
                                    <button
                                        class="btn btn-danger btn-sm"
                                        @click="deleteMed(m.id)"
                                        title="Удалить"
                                    >
                                        🗑
                                    </button>
                                </div>
                            </td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </div>

        <!-- МОДАЛКА ДОБАВЛЕНИЯ -->
        <BaseModal
            :show="isAddModalOpen"
            title="Добавить медикамент"
            @close="isAddModalOpen = false"
        >
            <div class="form-grid">
                <div class="form-group full">
                    <label>Название препарата *</label>
                    <input
                        v-model="newMed.name"
                        type="text"
                        placeholder="Напр. Амоксициллин"
                    />
                </div>
                <div class="form-group full">
                    <label>Описание / Показания</label>
                    <textarea
                        v-model="newMed.desc"
                        placeholder="Описание препарата..."
                    ></textarea>
                </div>
                <div class="form-group">
                    <label>Цена за ед. (руб.) *</label>
                    <input
                        v-model="newMed.price"
                        type="number"
                        placeholder="0"
                    />
                </div>
                <div class="form-group">
                    <label>Срок годности</label>
                    <div class="date-selects">
                        <select v-model="newMedDate.month">
                            <option
                                v-for="m in months"
                                :key="m.value"
                                :value="m.value"
                            >
                                {{ m.label }}
                            </option>
                        </select>
                        <select v-model="newMedDate.year">
                            <option
                                v-for="y in years"
                                :key="y"
                                :value="String(y)"
                            >
                                {{ y }}
                            </option>
                        </select>
                    </div>
                </div>
            </div>
            <template #footer>
                <button class="btn btn-ghost" @click="isAddModalOpen = false">
                    Отмена
                </button>
                <button class="btn btn-primary" @click="saveMed">
                    Сохранить
                </button>
            </template>
        </BaseModal>

        <!-- МОДАЛКА РЕДАКТИРОВАНИЯ -->
        <BaseModal
            :show="isEditModalOpen"
            title="Редактировать медикамент"
            @close="isEditModalOpen = false"
        >
            <div class="form-grid">
                <div class="form-group full">
                    <label>Название препарата *</label>
                    <input
                        v-model="editingMed.name"
                        type="text"
                        placeholder="Напр. Амоксициллин"
                    />
                </div>
                <div class="form-group full">
                    <label>Описание / Показания</label>
                    <textarea
                        v-model="editingMed.desc"
                        placeholder="Описание препарата..."
                    ></textarea>
                </div>
                <div class="form-group">
                    <label>Цена за ед. (руб.) *</label>
                    <input
                        v-model="editingMed.price"
                        type="number"
                        placeholder="0"
                    />
                </div>
                <div class="form-group">
                    <label>Срок годности</label>
                    <div class="date-selects">
                        <select v-model="editMedDate.month">
                            <option
                                v-for="m in months"
                                :key="m.value"
                                :value="m.value"
                            >
                                {{ m.label }}
                            </option>
                        </select>
                        <select v-model="editMedDate.year">
                            <option
                                v-for="y in years"
                                :key="y"
                                :value="String(y)"
                            >
                                {{ y }}
                            </option>
                        </select>
                    </div>
                </div>
            </div>
            <template #footer>
                <button class="btn btn-ghost" @click="isEditModalOpen = false">
                    Отмена
                </button>
                <button class="btn btn-primary" @click="saveEditedMed">
                    Сохранить изменения
                </button>
            </template>
        </BaseModal>

        <!-- ТОСТ -->
        <div v-if="toastVisible" class="toast">{{ toastMessage }}</div>
    </div>
</template>

<style scoped>
.med-desc {
    font-size: 11px;
    color: var(--text3);
    margin-top: 2px;
}

.expired-row {
    background: var(--red-dim);
}

textarea {
    min-height: 80px;
}

.date-selects {
    display: flex;
    gap: 8px;
}

.date-selects select {
    flex: 1;
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
