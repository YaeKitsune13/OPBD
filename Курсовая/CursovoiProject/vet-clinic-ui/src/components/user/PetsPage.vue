<script setup lang="ts">
import { ref, reactive, onMounted } from "vue";
import BaseModal from "../elements/BaseModal.vue";
import { PET_SPECIES } from "../../utils/constants";

interface Pet {
    id: number;
    name: string;
    species: string;
    breed: string;
    dob: string;
    weight: number;
    avatar: string;
}

interface PetForm {
    id: number;
    name: string;
    species: string;
    breed: string;
    dob: string;
    weight: string;
    avatar: string;
}

const pets = ref<Pet[]>([]);
const loading = ref(true);
const isAddModalOpen = ref(false);
const isEditModalOpen = ref(false);
const isDeleteModalOpen = ref(false);
const petToDelete = ref<Pet | null>(null);

const newPet = reactive({
    name: "",
    species: "Кошка",
    breed: "",
    dob: "",
    weight: "",
});
const editingPet = reactive<PetForm>({
    id: 0,
    name: "",
    species: "",
    breed: "",
    dob: "",
    weight: "",
    avatar: "",
});

function normalizePet(raw: any): Pet {
    return {
        id: raw.ID ?? raw.petId ?? raw.id ?? 0,
        name: raw.name ?? "",
        species: raw.species ?? "",
        breed: raw.breed ?? "",
        dob: (raw.dob ?? "").split("T")[0],
        weight: Number(raw.weight ?? 0),
        avatar: raw.avatar ?? "",
    };
}

async function loadPets() {
    loading.value = true;
    try {
        const token = localStorage.getItem("token");
        const userId = JSON.parse(localStorage.getItem("user") || "{}")?.id;

        console.group("🐾 loadPets");
        console.log("userId:", userId);

        const res = await fetch(`/api/pets/owner/${userId}`, {
            headers: { Authorization: `Bearer ${token}` },
        });

        console.log("status:", res.status);

        if (res.ok) {
            const raw = await res.json();
            console.log("raw response:", raw);
            pets.value = Array.isArray(raw) ? raw.map(normalizePet) : [];
            console.log("normalized:", pets.value);
        } else {
            console.warn("Ошибка ответа, pets пустые");
            pets.value = [];
        }
    } catch (e) {
        console.error("loadPets error:", e);
        pets.value = [];
    } finally {
        loading.value = false;
        console.groupEnd();
    }
}

async function savePet() {
    if (!newPet.name || !newPet.dob || !newPet.species)
        return alert("Заполните обязательные поля: Кличка, Вид, Дата рождения");

    const token = localStorage.getItem("token");
    const userId = JSON.parse(localStorage.getItem("user") || "{}")?.id;

    const payload = {
        name: newPet.name,
        species: newPet.species,
        breed: newPet.breed,
        dob: newPet.dob,
        weight: Number(newPet.weight),
        avatar: "",
    };

    console.group("🐾 savePet");
    console.log("payload:", payload);

    try {
        const res = await fetch(`/api/pets?ownerId=${userId}`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                Authorization: `Bearer ${token}`,
            },
            body: JSON.stringify(payload),
        });

        console.log("status:", res.status);
        const result = await res.json().catch(() => ({}));
        console.log("response:", result);

        if (res.ok) {
            await loadPets();
            isAddModalOpen.value = false;
            resetNewPet();
        } else {
            alert(`Ошибка при сохранении: ${result?.error || res.status}`);
        }
    } catch (e) {
        console.error("savePet error:", e);
        alert("Ошибка соединения");
    } finally {
        console.groupEnd();
    }
}

function openEditModal(pet: Pet) {
    editingPet.id = pet.id;
    editingPet.name = pet.name;
    editingPet.species = pet.species;
    editingPet.breed = pet.breed;
    editingPet.dob = pet.dob;
    editingPet.weight = String(pet.weight);
    editingPet.avatar = pet.avatar;
    isEditModalOpen.value = true;
}

async function updatePet() {
    if (!editingPet.name || !editingPet.dob)
        return alert("Кличка и дата рождения обязательны");

    const token = localStorage.getItem("token");

    const payload = {
        name: editingPet.name,
        species: editingPet.species,
        breed: editingPet.breed,
        dob: editingPet.dob,
        weight: Number(editingPet.weight),
        avatar: editingPet.avatar,
    };

    console.group("🐾 updatePet");
    console.log("id:", editingPet.id);
    console.log("payload:", payload);

    try {
        const res = await fetch(`/api/pets/${editingPet.id}`, {
            method: "PUT",
            headers: {
                "Content-Type": "application/json",
                Authorization: `Bearer ${token}`,
            },
            body: JSON.stringify(payload),
        });

        console.log("status:", res.status);
        const result = await res.json().catch(() => ({}));
        console.log("response:", result);

        if (res.ok) {
            await loadPets();
            isEditModalOpen.value = false;
        } else {
            alert(`Ошибка при обновлении: ${result?.error || res.status}`);
        }
    } catch (e) {
        console.error("updatePet error:", e);
        alert("Ошибка соединения");
    } finally {
        console.groupEnd();
    }
}

function confirmDelete(pet: Pet) {
    petToDelete.value = pet;
    isDeleteModalOpen.value = true;
}

async function deletePet() {
    if (!petToDelete.value) return;

    const token = localStorage.getItem("token");

    console.group("🐾 deletePet");
    console.log("id:", petToDelete.value.id);

    try {
        const res = await fetch(`/api/pets/${petToDelete.value.id}`, {
            method: "DELETE",
            headers: { Authorization: `Bearer ${token}` },
        });

        console.log("status:", res.status);

        if (res.ok) {
            await loadPets();
        } else {
            const result = await res.json().catch(() => ({}));
            console.warn("delete error:", result);
            alert(`Ошибка при удалении: ${result?.error || res.status}`);
        }
    } catch (e) {
        console.error("deletePet error:", e);
        alert("Ошибка соединения");
    } finally {
        isDeleteModalOpen.value = false;
        petToDelete.value = null;
        console.groupEnd();
    }
}

function resetNewPet() {
    newPet.name = "";
    newPet.species = "Кошка";
    newPet.breed = "";
    newPet.dob = "";
    newPet.weight = "";
}

function calcAge(dob: string) {
    if (!dob) return "—";
    const diff = Date.now() - new Date(dob).getTime();
    const y = Math.floor(diff / (1000 * 60 * 60 * 24 * 365.25));
    return y === 1 ? "1 год" : `${y} лет`;
}

function formatDob(dob: string) {
    if (!dob) return "—";
    return new Date(dob).toLocaleDateString("ru-RU", {
        day: "2-digit",
        month: "long",
        year: "numeric",
    });
}

onMounted(loadPets);
</script>

<template>
    <div class="page">
        <div class="page-header">
            <div>
                <div class="page-title">Мои питомцы</div>
                <div class="page-sub">Карточки и медицинская информация</div>
            </div>
            <button class="btn btn-primary" @click="isAddModalOpen = true">
                ✚ Добавить питомца
            </button>
        </div>

        <div
            v-if="loading"
            class="text-muted"
            style="padding: 40px; text-align: center"
        >
            Загрузка...
        </div>

        <div v-else class="pet-grid">
            <div v-for="pet in pets" :key="pet.id" class="pet-card">
                <div class="pet-avatar">{{ pet.avatar || "🐾" }}</div>
                <div class="pet-name">{{ pet.name }}</div>
                <div class="pet-meta">
                    {{ pet.species }} · {{ pet.breed || "—" }}
                </div>
                <div class="pet-details">
                    <span class="tag">{{ calcAge(pet.dob) }}</span>
                    <span class="tag">⚖️ {{ pet.weight }} кг</span>
                </div>
                <div class="pet-dob text-muted">
                    {{ formatDob(pet.dob) }}
                </div>
                <div class="pet-actions">
                    <button
                        class="btn btn-ghost btn-sm"
                        style="flex: 1"
                        @click="openEditModal(pet)"
                    >
                        ✏ Изменить
                    </button>
                    <button
                        class="btn btn-danger btn-sm"
                        style="flex: 1"
                        @click="confirmDelete(pet)"
                    >
                        🗑 Удалить
                    </button>
                </div>
            </div>

            <div
                class="pet-card add-placeholder"
                @click="isAddModalOpen = true"
            >
                <div class="add-icon">+</div>
                <div style="font-size: 12px">Добавить питомца</div>
            </div>
        </div>

        <!-- МОДАЛКА: Добавить -->
        <BaseModal
            :show="isAddModalOpen"
            title="Добавить питомца"
            @close="isAddModalOpen = false"
        >
            <div class="form-grid">
                <div class="form-group">
                    <label>Кличка *</label>
                    <input
                        v-model="newPet.name"
                        type="text"
                        placeholder="Барсик"
                    />
                </div>
                <div class="form-group">
                    <label>Вид животного *</label>
                    <select v-model="newPet.species">
                        <option
                            v-for="s in PET_SPECIES"
                            :key="s.value"
                            :value="s.label"
                        >
                            {{ s.emoji }} {{ s.label }}
                        </option>
                    </select>
                </div>
                <div class="form-group">
                    <label>Порода</label>
                    <input
                        v-model="newPet.breed"
                        type="text"
                        placeholder="Британец"
                    />
                </div>
                <div class="form-group">
                    <label>Дата рождения *</label>
                    <input v-model="newPet.dob" type="date" />
                </div>
                <div class="form-group">
                    <label>Вес (кг)</label>
                    <input
                        v-model="newPet.weight"
                        type="number"
                        step="0.1"
                        placeholder="0.0"
                    />
                </div>
            </div>
            <template #footer>
                <button class="btn btn-ghost" @click="isAddModalOpen = false">
                    Отмена
                </button>
                <button class="btn btn-primary" @click="savePet">
                    Сохранить
                </button>
            </template>
        </BaseModal>

        <BaseModal
            :show="isEditModalOpen"
            title="Редактировать питомца"
            @close="isEditModalOpen = false"
        >
            <div class="form-grid">
                <div class="form-group">
                    <label>Кличка *</label>
                    <input v-model="editingPet.name" type="text" />
                </div>
                <div class="form-group">
                    <label>Вид животного</label>
                    <select v-model="editingPet.species">
                        <option
                            v-for="s in PET_SPECIES"
                            :key="s.value"
                            :value="s.label"
                        >
                            {{ s.emoji }} {{ s.label }}
                        </option>
                    </select>
                </div>
                <div class="form-group">
                    <label>Порода</label>
                    <input v-model="editingPet.breed" type="text" />
                </div>
                <div class="form-group">
                    <label>Дата рождения *</label>
                    <input v-model="editingPet.dob" type="date" />
                </div>
                <div class="form-group">
                    <label>Вес (кг)</label>
                    <input
                        v-model="editingPet.weight"
                        type="number"
                        step="0.1"
                    />
                </div>
            </div>
            <template #footer>
                <button class="btn btn-ghost" @click="isEditModalOpen = false">
                    Отмена
                </button>
                <button class="btn btn-primary" @click="updatePet">
                    Сохранить изменения
                </button>
            </template>
        </BaseModal>

        <!-- МОДАЛКА: Удалить -->
        <BaseModal
            :show="isDeleteModalOpen"
            title="Удалить питомца?"
            max-width="400px"
            @close="isDeleteModalOpen = false"
        >
            <p style="color: var(--text2); font-size: 14px">
                Вы уверены, что хотите удалить карточку питомца
                <b>{{ petToDelete?.name }}</b
                >? Это действие необратимо.
            </p>
            <template #footer>
                <button
                    class="btn btn-ghost"
                    @click="isDeleteModalOpen = false"
                >
                    Отмена
                </button>
                <button class="btn btn-danger" @click="deletePet">
                    Удалить
                </button>
            </template>
        </BaseModal>
    </div>
</template>

<style scoped>
.pet-name {
    font-size: 15px;
    font-weight: 600;
    margin-bottom: 2px;
}
.pet-meta {
    font-size: 12px;
    color: var(--text2);
    margin-bottom: 8px;
}
.pet-dob {
    font-size: 11px;
    margin-top: 6px;
}
.pet-details {
    display: flex;
    gap: 6px;
    flex-wrap: wrap;
}
.pet-actions {
    display: flex;
    gap: 6px;
    margin-top: 12px;
}

.add-placeholder {
    border-style: dashed !important;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    min-height: 180px;
    gap: 8px;
    color: var(--text3);
    cursor: pointer;
    transition: all 0.15s;
}
.add-placeholder:hover {
    color: var(--accent);
    border-color: var(--accent) !important;
}
.add-icon {
    font-size: 24px;
}
</style>
