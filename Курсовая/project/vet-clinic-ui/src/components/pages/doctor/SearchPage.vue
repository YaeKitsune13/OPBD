<script setup>
import { ref, watch, onMounted, computed } from "vue";

const emit = defineEmits(["navigate"]);
const searchQuery = ref("");
const results = ref([]);
const selectedPet = ref(null);
const visits = ref([]);
const loading = ref(false);

// РАЦИОНАЛЬНО: Фильтруем результаты, чтобы скрыть "Анонимов" и системных "Гостей"
const filteredResults = computed(() => {
    return results.value.filter((pet) => {
        // Скрываем, если имя владельца содержит "Анонимный" или телефон системный "0000"
        const isAnon =
            pet.owner_name?.toLowerCase().includes("аноним") ||
            pet.owner_phone === "0000";
        return !isAnon;
    });
});

function petEmoji(species) {
    const s = species?.toLowerCase();
    if (s?.includes("кош") || s?.includes("кот")) return "🐱";
    if (s?.includes("соб") || s?.includes("пес")) return "🐶";
    if (s?.includes("птиц")) return "🐦";
    return "🐾";
}

function formatDate(dateStr) {
    if (!dateStr) return "—";
    return new Date(dateStr).toLocaleDateString("ru-RU");
}

async function doSearch(query = "") {
    loading.value = true;
    const token = localStorage.getItem("token");
    try {
        const q = query.trim();
        const url =
            q.length >= 2
                ? `/api/search?q=${encodeURIComponent(q)}`
                : `/api/search?q=`;

        const response = await fetch(url, {
            headers: { Authorization: `Bearer ${token}` },
        });
        if (response.ok) {
            results.value = (await response.json()) ?? [];
        }
    } catch (e) {
        console.error("Ошибка поиска:", e);
    } finally {
        loading.value = false;
    }
}

let debounceTimer = null;
watch(searchQuery, (val) => {
    clearTimeout(debounceTimer);
    debounceTimer = setTimeout(() => doSearch(val), 300);
});

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
        }
    } catch (e) {
        console.error("Ошибка загрузки истории:", e);
    }
}
</script>

<template>
    <div class="page-container">
        <div class="page-header">
            <h1 class="page-title">Поиск пациента</h1>
            <p class="page-sub">
                Поиск по кличке или ФИО владельца (скрыты анонимные записи)
            </p>
        </div>

        <!-- Поле поиска -->
        <div class="search-card">
            <div class="search-wrapper">
                <span class="search-icon">🔍</span>
                <input
                    type="text"
                    v-model="searchQuery"
                    placeholder="Например: Барсик или Иванов..."
                    class="search-input"
                />
            </div>
        </div>

        <!-- Список результатов -->
        <div v-if="!selectedPet">
            <div v-if="loading" class="state-msg">
                ⏳ Поиск в базе данных...
            </div>

            <div v-else-if="filteredResults.length === 0" class="state-msg">
                Ничего не найдено (или запись является анонимной)
            </div>

            <div v-else class="results-grid">
                <div
                    v-for="pet in filteredResults"
                    :key="pet.pet_id"
                    class="pet-row"
                    @click="selectPet(pet)"
                >
                    <div class="pet-main-info">
                        <div class="pet-avatar-circle">
                            {{ petEmoji(pet.species) }}
                        </div>
                        <div class="pet-text">
                            <div class="pet-nickname">{{ pet.pet_name }}</div>
                            <div class="pet-meta">
                                {{ pet.breed }} · {{ pet.owner_name }}
                            </div>
                        </div>
                    </div>

                    <div class="row-actions">
                        <button
                            class="btn-action"
                            @click.stop="emit('navigate', 'conduct', pet)"
                        >
                            Принять
                        </button>
                    </div>
                </div>
            </div>
        </div>

        <!-- Детали питомца и история -->
        <div v-if="selectedPet" class="details-view">
            <div class="details-header">
                <button class="btn-back" @click="selectedPet = null">
                    ← К списку
                </button>
                <button
                    class="btn-start"
                    @click="emit('navigate', 'conduct', selectedPet)"
                >
                    Начать приём
                </button>
            </div>

            <div class="pet-profile-card">
                <div class="profile-main">
                    <span class="profile-emoji">{{
                        petEmoji(selectedPet.species)
                    }}</span>
                    <div class="profile-info">
                        <h2>{{ selectedPet.pet_name }}</h2>
                        <p>
                            {{ selectedPet.species }} ·
                            {{ selectedPet.breed }} ·
                            {{ selectedPet.weight }} кг
                        </p>
                        <p class="owner-link">
                            Владелец: {{ selectedPet.owner_name }} ({{
                                selectedPet.owner_phone
                            }})
                        </p>
                    </div>
                </div>
            </div>

            <div class="history-section">
                <h3 class="section-title">История визитов</h3>
                <div v-if="visits.length === 0" class="history-empty">
                    Записей о лечении не найдено
                </div>
                <div
                    v-for="visit in visits"
                    :key="visit.visitId"
                    class="visit-card"
                >
                    <div class="visit-top">
                        <span class="visit-date"
                            >{{ visit.date }} {{ visit.time }}</span
                        >
                        <span class="visit-price">{{ visit.price }}</span>
                    </div>
                    <div class="visit-diag">🔬 {{ visit.diagnosis }}</div>
                    <div class="visit-desc">{{ visit.details }}</div>
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped>
/* Используем переменные твоей Dark Theme */
.page-container {
    color: var(--text);
    font-family: var(--font);
}
.page-header {
    margin-bottom: 24px;
}
.page-title {
    font-size: 24px;
    font-weight: 700;
    margin: 0;
}
.page-sub {
    color: var(--text3);
    font-size: 13px;
    margin-top: 4px;
}

/* Search Input */
.search-card {
    background: var(--surface);
    border: 1px solid var(--border);
    border-radius: var(--radius);
    padding: 16px;
    margin-bottom: 20px;
}
.search-wrapper {
    position: relative;
    display: flex;
    align-items: center;
}
.search-icon {
    position: absolute;
    left: 12px;
    color: var(--text3);
}
.search-input {
    width: 100%;
    background: var(--surface2);
    border: 1px solid var(--border2);
    color: var(--text);
    padding: 12px 12px 12px 40px;
    border-radius: var(--radius);
    outline: none;
    transition: border-color 0.2s;
}
.search-input:focus {
    border-color: var(--accent);
}

/* Results */
.results-grid {
    background: var(--surface);
    border: 1px solid var(--border);
    border-radius: var(--radius);
    overflow: hidden;
}
.pet-row {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 14px 16px;
    border-bottom: 1px solid var(--border);
    cursor: pointer;
    transition: background 0.2s;
}
.pet-row:hover {
    background: var(--surface2);
}
.pet-row:last-child {
    border-bottom: none;
}

.pet-main-info {
    display: flex;
    align-items: center;
    gap: 14px;
}
.pet-avatar-circle {
    width: 42px;
    height: 42px;
    background: var(--surface3);
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 20px;
}
.pet-nickname {
    font-weight: 600;
    font-size: 15px;
    color: var(--text);
}
.pet-meta {
    font-size: 12px;
    color: var(--text3);
    margin-top: 2px;
}

.btn-action {
    background: var(--accent-dim);
    color: var(--accent);
    border: 1px solid var(--accent);
    padding: 6px 16px;
    border-radius: var(--radius);
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s;
}
.btn-action:hover {
    background: var(--accent);
    color: var(--bg);
}

/* Details View */
.details-header {
    display: flex;
    justify-content: space-between;
    margin-bottom: 16px;
}
.btn-back {
    background: none;
    border: none;
    color: var(--text2);
    cursor: pointer;
    font-weight: 600;
}
.btn-start {
    background: var(--accent);
    color: var(--bg);
    border: none;
    padding: 8px 16px;
    border-radius: var(--radius);
    cursor: pointer;
    font-weight: 600;
}

.pet-profile-card {
    background: var(--surface2);
    border: 1px solid var(--border);
    padding: 20px;
    border-radius: var(--radius);
}
.profile-main {
    display: flex;
    gap: 20px;
    align-items: center;
}
.profile-emoji {
    font-size: 48px;
}
.profile-info h2 {
    margin: 0;
    color: var(--accent);
}
.profile-info p {
    margin: 4px 0;
    color: var(--text2);
    font-size: 14px;
}
.owner-link {
    color: var(--blue) !important;
}

/* History */
.history-section {
    margin-top: 24px;
}
.section-title {
    font-size: 16px;
    font-weight: 600;
    margin-bottom: 12px;
    color: var(--text2);
}
.visit-card {
    background: var(--surface);
    border: 1px solid var(--border);
    padding: 16px;
    border-radius: var(--radius);
    margin-bottom: 12px;
}
.visit-top {
    display: flex;
    justify-content: space-between;
    margin-bottom: 8px;
}
.visit-date {
    font-family: var(--mono);
    font-size: 12px;
    color: var(--text3);
}
.visit-price {
    color: var(--accent);
    font-weight: 600;
    font-family: var(--mono);
}
.visit-diag {
    font-weight: 600;
    color: var(--text);
    margin-bottom: 4px;
}
.visit-desc {
    font-size: 13px;
    color: var(--text3);
    line-height: 1.4;
}

.state-msg {
    padding: 40px;
    text-align: center;
    color: var(--text3);
}
</style>
