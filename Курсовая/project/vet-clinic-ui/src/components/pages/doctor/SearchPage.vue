<script setup>
import { ref, watch, onMounted } from "vue";

const searchQuery = ref("");
const results = ref([]);
const selectedPet = ref(null);
const visits = ref([]);
const loading = ref(false);

const toastMessage = ref("");
const toastVisible = ref(false);

function showToast(message) {
    toastMessage.value = message;
    toastVisible.value = true;
    setTimeout(() => (toastVisible.value = false), 3000);
}

function petEmoji(species) {
    const s = species?.toLowerCase();
    if (s?.includes("кош") || s?.includes("кот") || s === "cat") return "🐱";
    if (s?.includes("соб") || s === "dog") return "🐶";
    if (s?.includes("птиц") || s === "bird") return "🐦";
    if (s?.includes("кроли") || s === "rabbit") return "🐰";
    return "🐾";
}

function formatDate(dateStr) {
    if (!dateStr) return "—";
    return new Date(dateStr).toLocaleDateString("ru-RU");
}

function badgeClass(status) {
    switch (status) {
        case "completed":
            return "badge-confirmed";
        case "cancelled":
            return "badge-cancelled";
        default:
            return "badge-pending";
    }
}

function statusLabel(status) {
    switch (status) {
        case "completed":
            return "Завершён";
        case "cancelled":
            return "Отменён";
        default:
            return "Запланирован";
    }
}

async function doSearch(query = "") {
    loading.value = true;
    selectedPet.value = null;
    visits.value = [];

    const token = localStorage.getItem("token");
    try {
        const q = query.trim();
        const url =
            q.length >= 2
                ? `/api/search?q=${encodeURIComponent(q)}`
                : `/api/search?q=`; // пустой запрос — вернёт всех
        const response = await fetch(url, {
            headers: { Authorization: `Bearer ${token}` },
        });
        if (response.ok) {
            results.value = (await response.json()) ?? [];
        } else {
            showToast("Ошибка поиска");
            results.value = [];
        }
    } catch (e) {
        console.error(e);
        showToast("Ошибка поиска");
        results.value = [];
    } finally {
        loading.value = false;
    }
}

// Debounce чтобы не долбить сервер на каждый символ
let debounceTimer = null;
watch(searchQuery, (val) => {
    clearTimeout(debounceTimer);
    debounceTimer = setTimeout(() => doSearch(val), 300);
});

// Загрузить всех при входе
onMounted(() => doSearch());

async function selectPet(pet) {
    selectedPet.value = pet;
    visits.value = [];

    const token = localStorage.getItem("token");
    try {
        const response = await fetch(`/api/visits/pet/${pet.pet_id}`, {
            headers: { Authorization: `Bearer ${token}` },
        });
        if (response.ok) {
            visits.value = (await response.json()) ?? [];
        } else {
            showToast("Не удалось загрузить историю визитов");
        }
    } catch (e) {
        console.error(e);
        showToast("Ошибка загрузки визитов");
    }
}
</script>

<template>
    <div class="page">
        <div class="page-header">
            <div>
                <div class="page-title">Поиск пациента</div>
                <div class="page-sub">
                    История болезней любого питомца по кличке или ФИО
                </div>
            </div>
        </div>

        <!-- Строка поиска -->
        <div class="card">
            <div class="card-body">
                <input
                    type="text"
                    v-model="searchQuery"
                    placeholder="Кличка питомца или ФИО владельца..."
                    style="width: 100%"
                />
            </div>
        </div>

        <!-- Список результатов -->
        <div v-if="!selectedPet">
            <div v-if="loading" class="empty">
                <div class="empty-icon">⏳</div>
                <div>Поиск...</div>
            </div>

            <div v-else-if="results.length === 0" class="empty">
                <div class="empty-icon">😔</div>
                <div>Ничего не найдено</div>
            </div>

            <div v-else class="card">
                <div class="card-header">
                    <span class="card-title"
                        >Пациентов: {{ results.length }}</span
                    >
                </div>
                <div
                    v-for="pet in results"
                    :key="pet.pet_id"
                    class="search-result-row"
                    @click="selectPet(pet)"
                >
                    <div class="row" style="gap: 12px; align-items: center">
                        <div class="pet-avatar">
                            {{ petEmoji(pet.species) }}
                        </div>
                        <div>
                            <div style="font-weight: 600">
                                {{ pet.pet_name }}
                            </div>
                            <div class="text-muted" style="font-size: 12px">
                                {{ pet.species }} · {{ pet.breed }} ·
                                {{ pet.weight }} кг · {{ pet.age }} лет
                            </div>
                            <div class="text-muted" style="font-size: 12px">
                                Владелец: {{ pet.owner_name }} ·
                                {{ pet.owner_phone }}
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <!-- Детали питомца + история визитов -->
        <div v-if="selectedPet" class="card">
            <div class="card-header">
                <button
                    class="btn btn-ghost btn-sm"
                    @click="selectedPet = null"
                >
                    ← Назад
                </button>
                <span class="card-title" style="margin-left: 8px">
                    {{ selectedPet.pet_name }} · {{ selectedPet.owner_name }}
                </span>
            </div>

            <div style="padding: 14px; border-bottom: 1px solid var(--border)">
                <div class="row" style="gap: 12px">
                    <div class="pet-avatar">
                        {{ petEmoji(selectedPet.species) }}
                    </div>
                    <div>
                        <div style="font-weight: 600">
                            {{ selectedPet.pet_name }}
                        </div>
                        <div class="text-muted" style="font-size: 12px">
                            {{ selectedPet.species }} ·
                            {{ selectedPet.breed }} ·
                            {{ selectedPet.weight }} кг ·
                            {{ selectedPet.age }} лет
                        </div>
                        <div class="text-muted" style="font-size: 12px">
                            Владелец: {{ selectedPet.owner_name }} ·
                            {{ selectedPet.owner_phone }}
                        </div>
                    </div>
                </div>
            </div>

            <div style="padding: 14px">
                <div
                    v-if="visits.length === 0"
                    class="text-muted"
                    style="font-size: 13px"
                >
                    Визитов не найдено
                </div>
                <div
                    v-for="visit in visits"
                    :key="visit.id"
                    class="visit-item"
                    style="
                        border: 1px solid var(--border2);
                        border-radius: 4px;
                        padding: 10px;
                        margin-bottom: 8px;
                    "
                >
                    <div class="row-between">
                        <span class="mono" style="font-size: 12px">
                            {{ formatDate(visit.date) }}
                        </span>
                        <span class="badge" :class="badgeClass(visit.status)">
                            {{ statusLabel(visit.status) }}
                        </span>
                    </div>
                    <div style="margin-top: 5px; font-weight: 500">
                        🔬 {{ visit.diagnosis }}
                    </div>
                    <div
                        v-if="visit.notes"
                        class="text-muted"
                        style="font-size: 12px; margin-top: 4px"
                    >
                        {{ visit.notes }}
                    </div>
                </div>
            </div>
        </div>

        <div v-if="toastVisible" class="toast">{{ toastMessage }}</div>
    </div>
</template>

<style scoped>
.search-result-row {
    padding: 14px;
    border-bottom: 1px solid var(--border);
    cursor: pointer;
    transition: background 0.15s;
}

.search-result-row:hover {
    background: var(--bg2);
}

.search-result-row:last-child {
    border-bottom: none;
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
