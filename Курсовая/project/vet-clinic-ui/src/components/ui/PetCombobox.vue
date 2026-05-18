<script setup>
import { ref, onMounted, onUnmounted, watch } from "vue";

const props = defineProps({
    modelValue: Object,
    isDoctor: Boolean, // Теперь мы используем этот проп
    placeholder: { type: String, default: "Выберите питомца" },
});

const emit = defineEmits(["update:modelValue"]);

const searchInput = ref(null);
const isOpen = ref(false);
const isLoading = ref(false);
const searchQuery = ref("");
const pets = ref([]); // Список питомцев теперь хранится тут
const containerRef = ref(null);

// Функция загрузки/поиска питомцев
async function fetchPets(query = "") {
    isLoading.value = true;
    const token = localStorage.getItem("token");

    try {
        const url = query
            ? `/api/search?query=${encodeURIComponent(query)}`
            : `/api/search`;
        const response = await fetch(url, {
            headers: { Authorization: `Bearer ${token}` },
        });

        if (response.ok) {
            const data = await response.json();
            console.log("[PetCombobox] Пришло с сервера:", data);

            pets.value = (data || []).map((p) => ({
                // Используем ключи, которые мы увидели в консоли:
                id: p.pet_id,
                name: p.pet_name,
                breed: p.breed,
                owner: p.owner_name,
                avatar: getEmoji(p.species),
            }));
        }
    } catch (e) {
        console.error("[PetCombobox] Ошибка:", e);
    } finally {
        isLoading.value = false;
    }
}

function getEmoji(species) {
    const s = species?.toLowerCase() || "";
    if (s.includes("кот") || s.includes("кош") || s.includes("cat"))
        return "🐱";
    if (s.includes("соб") || s.includes("пес") || s.includes("dog"))
        return "🐶";
    return "🐾";
}

// Следим за вводом текста (поиск с задержкой)
let timeout = null;
watch(searchQuery, (newQuery) => {
    clearTimeout(timeout);
    timeout = setTimeout(() => {
        fetchPets(newQuery);
    }, 300);
});

function selectPet(pet) {
    emit("update:modelValue", pet);
    isOpen.value = false;
    searchQuery.value = ""; // Очищаем поиск при выборе
}

function toggleDropdown() {
    isOpen.value = !isOpen.value;
    if (isOpen.value && pets.value.length === 0) {
        fetchPets(); // Загружаем начальный список при открытии
    }
}

const clickOutside = (e) => {
    if (containerRef.value && !containerRef.value.contains(e.target)) {
        isOpen.value = false;
    }
};

onMounted(() => document.addEventListener("click", clickOutside));
onUnmounted(() => document.removeEventListener("click", clickOutside));

watch(isOpen, async (newVal) => {
    if (newVal) {
        // Ждем отрисовки и ставим фокус
        setTimeout(() => {
            const el = document.querySelector(".search-input");
            if (el) el.focus();
        }, 50);
    }
});
</script>

<template>
    <div class="pet-select-custom" ref="containerRef">
        <div
            class="select-trigger"
            :class="{ 'is-open': isOpen }"
            @click="toggleDropdown"
        >
            <div class="selected-content">
                <template v-if="modelValue && !isOpen">
                    <span class="selected-avatar">{{ modelValue.avatar }}</span>
                    <span class="selected-name">{{ modelValue.name }}</span>
                </template>
                <template v-else-if="isOpen">
                    <!-- Инпут поиска внутри селекта -->
                    <input
                        v-model="searchQuery"
                        class="search-input"
                        placeholder="Поиск по кличке или владельцу..."
                        @click.stop
                        autofocus
                    />
                </template>
                <span v-else class="placeholder">{{ placeholder }}</span>
            </div>
            <div class="arrow" :class="{ rotated: isOpen }">▼</div>
        </div>

        <Transition name="slide-fade">
            <div v-if="isOpen" class="dropdown-list">
                <div v-if="isLoading" class="info-item">Загрузка...</div>
                <div v-else-if="pets.length === 0" class="info-item">
                    Никого не нашли
                </div>

                <div
                    v-for="pet in pets"
                    :key="pet.id"
                    class="pet-option"
                    :class="{
                        'is-active': pet.id && modelValue?.id === pet.id,
                    }"
                    @click="selectPet(pet)"
                >
                    <div class="option-avatar">{{ pet.avatar }}</div>
                    <div class="option-info">
                        <div class="option-name">{{ pet.name }}</div>
                        <div class="option-sub">
                            {{ pet.breed }} • {{ pet.owner }}
                        </div>
                    </div>
                    <div v-if="modelValue?.id === pet.id" class="check-mark">
                        ✓
                    </div>
                </div>
            </div>
        </Transition>
    </div>
</template>

<style scoped>
/* Твои старые стили + новые для поиска */
.pet-select-custom {
    position: relative;
    width: 100%;
    max-width: 400px;
    user-select: none;
}

.search-input {
    background: transparent;
    border: none;
    color: var(--text);
    outline: none;
    width: 100%;
    font-size: 14px;
}

.info-item {
    padding: 12px;
    text-align: center;
    color: var(--text3);
    font-size: 13px;
}

.select-trigger {
    display: flex;
    align-items: center;
    justify-content: space-between;
    background: var(--surface2);
    border: 1px solid var(--border);
    border-radius: 8px;
    padding: 8px 12px;
    min-height: 40px;
    cursor: pointer;
}

.select-trigger.is-open {
    border-color: var(--accent);
}

.selected-content {
    display: flex;
    align-items: center;
    gap: 10px;
    flex: 1;
}

.dropdown-list {
    position: absolute;
    top: calc(100% + 6px);
    left: 0;
    right: 0;
    background: #1a1a1a; /* Тёмный фон под твой скриншот */
    border: 1px solid #333;
    border-radius: 8px;
    max-height: 300px;
    overflow-y: auto;
    z-index: 9999;
    box-shadow: 0 10px 30px rgba(0, 0, 0, 0.5);
}

.pet-option {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 10px 14px;
    cursor: pointer;
}

.pet-option:hover {
    background: #333;
}

.option-name {
    font-size: 14px;
    font-weight: 600;
    color: #fff;
}

.option-sub {
    font-size: 11px;
    color: #888;
}

.arrow {
    font-size: 9px;
    transition: transform 0.3s;
}

.arrow.rotated {
    transform: rotate(180deg);
}
</style>
