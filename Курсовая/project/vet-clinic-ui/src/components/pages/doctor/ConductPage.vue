<script setup>
import { ref, computed, onMounted, watch } from "vue";
import BaseModal from "../../ui/BaseModal.vue";
import PetCombobox from "../../ui/PetCombobox.vue";

const emit = defineEmits(["navigate"]);

const selectedPet = ref(null);
const allServices = ref([]);

const anamnesis = ref("");
const diagnosis = ref("");
const assignments = ref([]);

const isServiceModalOpen = ref(false);
const selectedServiceId = ref(null);
const loading = ref(false);
const toastMessage = ref("");
const toastVisible = ref(false);

// Логируем изменение выбранного питомца
watch(selectedPet, (newPet) => {
    console.log("[ConductVisit] Выбран питомец:", newPet);
});

const totalCost = computed(() => {
    const total = assignments.value.reduce(
        (sum, item) => sum + item.price * item.qty,
        0,
    );
    console.log("[ConductVisit] Пересчет итоговой стоимости:", total);
    return total;
});

function showToast(message) {
    console.log("[ConductVisit] Показ уведомления:", message);
    toastMessage.value = message;
    toastVisible.value = true;
    setTimeout(() => (toastVisible.value = false), 3000);
}

async function loadServices() {
    console.log("[ConductVisit] Загрузка списка услуг...");
    const token = localStorage.getItem("token");

    try {
        const response = await fetch(`/api/services`, {
            headers: { Authorization: `Bearer ${token}` },
        });
        if (response.ok) {
            const data = await response.json();
            allServices.value = data ?? [];
            console.log(
                "[ConductVisit] Услуги успешно загружены. Количество:",
                allServices.value.length,
            );
        } else {
            console.error(
                "[ConductVisit] Ошибка при загрузке услуг. Статус:",
                response.status,
            );
        }
    } catch (e) {
        console.error("[ConductVisit] Критическая ошибка loadServices:", e);
    }
}

function removeItem(id) {
    console.log("[ConductVisit] Удаление услуги с ID:", id);
    assignments.value = assignments.value.filter((i) => i.id !== id);
}

function addService() {
    console.log(
        "[ConductVisit] Попытка добавить услугу. ID выбранной услуги:",
        selectedServiceId.value,
    );

    if (!selectedServiceId.value) {
        console.warn("[ConductVisit] Услуга не выбрана");
        showToast("Выберите услугу");
        return;
    }

    const service = allServices.value.find(
        (s) => s.id === selectedServiceId.value,
    );

    if (service) {
        const existing = assignments.value.find((a) => a.id === service.id);
        if (existing) {
            existing.qty++;
            console.log(
                `[ConductVisit] Услуга "${service.name}" уже есть, увеличиваем количество до ${existing.qty}`,
            );
        } else {
            assignments.value.push({
                id: service.id,
                name: service.name,
                type: "Услуга",
                price: service.price,
                qty: 1,
            });
            console.log(
                `[ConductVisit] Добавлена новая услуга: ${service.name}`,
            );
        }
        isServiceModalOpen.value = false;
        selectedServiceId.value = null;
        showToast("Услуга добавлена");
    } else {
        console.error(
            "[ConductVisit] Выбранная услуга не найдена в общем списке allServices",
        );
    }
}

async function saveVisit() {
    console.log("[ConductVisit] Начало сохранения...");

    if (!selectedPet.value) {
        showToast("Выберите питомца");
        return;
    }
    if (!diagnosis.value.trim()) {
        showToast("Введите диагноз");
        return;
    }

    if (assignments.value.length === 0) {
        showToast("Добавьте хотя бы одну услугу");
        return;
    }

    loading.value = true;
    const token = localStorage.getItem("token");

    // Формируем payload строго под твой Go DTO
    const payload = {
        selectedPet: {
            // ВАЖНО: Твой бэкенд использует это как AppointmentID.
            // Если мы идем через поиск (без записи), это может вызвать ошибку в базе.
            id: Number(selectedPet.value.pet_id || selectedPet.value.id),
            name: selectedPet.value.pet_name || selectedPet.value.name,
            avatar: "🐾",
            breed: selectedPet.value.breed || "",
            owner: selectedPet.value.owner_name || "",
        },
        anamnesis: anamnesis.value,
        diagnosis: diagnosis.value,
        assignments: assignments.value.map((a) => ({
            id: Number(a.id),
            name: a.name,
            type: "service", // БЫЛО "Услуга", ДОЛЖНО БЫТЬ "service" (для БД)
            price: Number(a.price),
            qty: Number(a.qty),
        })),
        totalCost: Math.round(totalCost.value), // Go ждет int64
    };

    console.log(
        "[ConductVisit] Итоговый Payload для отправки:",
        JSON.stringify(payload, null, 2),
    );

    try {
        const response = await fetch(`/api/visits`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                Authorization: `Bearer ${token}`,
            },
            body: JSON.stringify(payload),
        });

        if (response.ok) {
            console.log("[ConductVisit] Успешно сохранено!");
            showToast("Приём сохранён успешно");
            // Очистка
            selectedPet.value = null;
            anamnesis.value = "";
            diagnosis.value = "";
            assignments.value = [];
            setTimeout(() => emit("navigate", "today"), 1500);
        } else {
            const errorData = await response.json();
            console.error("[ConductVisit] Ошибка сервера:", errorData);
            showToast(errorData.error || "Ошибка сервера");
        }
    } catch (e) {
        console.error("[ConductVisit] Ошибка запроса:", e);
        showToast("Сетевая ошибка: " + e.message);
    } finally {
        loading.value = false;
    }
}

onMounted(() => {
    console.log("[ConductVisit] Компонент смонтирован");
    loadServices();
});
</script>

<template>
    <div class="page">
        <div class="page-header conduct-header">
            <div class="title-with-select">
                <div class="page-title">Ведение приёма:</div>
                <!-- isDoctor=true → загружает всех питомцев с поиском -->
                <PetCombobox v-model="selectedPet" :is-doctor="true" />
            </div>
        </div>

        <div class="grid-2">
            <!-- Медкарта -->
            <div class="card">
                <div class="card-header">
                    <span class="card-title">Медкарта</span>
                </div>
                <div class="card-body">
                    <div class="form-group">
                        <label>АНАМНЕЗ</label>
                        <textarea
                            v-model="anamnesis"
                            placeholder="Жалобы пациента..."
                        ></textarea>
                    </div>
                    <div class="form-group mt-12">
                        <label>ДИАГНОЗ</label>
                        <input
                            v-model="diagnosis"
                            type="text"
                            placeholder="Введите диагноз"
                        />
                    </div>
                </div>
            </div>

            <!-- Назначения -->
            <div class="card">
                <div class="card-header">
                    <span class="card-title">Назначения</span>
                    <button
                        class="btn btn-ghost btn-sm"
                        @click="isServiceModalOpen = true"
                    >
                        + Добавить услугу
                    </button>
                </div>
                <div class="table-wrap">
                    <table>
                        <thead>
                            <tr>
                                <th>НАИМЕНОВАНИЕ</th>
                                <th>ЦЕНА</th>
                                <th></th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-for="item in assignments" :key="item.id">
                                <td class="td-main">{{ item.name }}</td>
                                <td class="mono">{{ item.price }} р.</td>
                                <td>
                                    <button
                                        class="btn btn-ghost btn-sm"
                                        @click="removeItem(item.id)"
                                    >
                                        ✕
                                    </button>
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>
                <div class="card-body row-between total-row">
                    <span class="text-muted">Итого:</span>
                    <span class="text-accent mono total-sum"
                        >{{ totalCost }} руб.</span
                    >
                </div>
            </div>
        </div>

        <div class="row conduct-footer">
            <button class="btn btn-ghost" @click="emit('navigate', 'today')">
                Отмена
            </button>
            <button
                class="btn btn-primary"
                @click="saveVisit"
                :disabled="loading"
            >
                {{ loading ? "Сохраняю..." : "Сохранить прием" }}
            </button>
        </div>

        <!-- Модалка добавления услуги -->
        <BaseModal
            :show="isServiceModalOpen"
            title="Добавить услугу"
            @close="isServiceModalOpen = false"
        >
            <div class="form-group">
                <label>Выберите услугу</label>
                <select v-model.number="selectedServiceId" class="mt-4">
                    <option :value="null">-- Выберите --</option>
                    <option
                        v-for="service in allServices"
                        :key="service.id"
                        :value="service.id"
                    >
                        {{ service.name }} — {{ service.price }} руб.
                    </option>
                </select>
            </div>
            <template #footer>
                <button class="btn btn-primary" @click="addService">
                    Добавить
                </button>
            </template>
        </BaseModal>

        <div v-if="toastVisible" class="toast">{{ toastMessage }}</div>
    </div>
</template>

<style scoped>
.conduct-header {
    margin-bottom: 20px;
}
.title-with-select {
    display: flex;
    align-items: center;
    gap: 12px;
    width: 100%;
}

.total-row {
    border-top: 1px solid var(--border);
}
.total-sum {
    font-size: 18px;
    font-weight: 700;
}

.conduct-footer {
    justify-content: flex-end;
    gap: 10px;
    margin-top: 24px;
}

textarea {
    min-height: 120px;
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
