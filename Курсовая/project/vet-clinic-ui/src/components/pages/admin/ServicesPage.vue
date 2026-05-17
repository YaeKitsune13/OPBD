<script setup>
import { ref, computed, reactive, onMounted } from "vue";
import BaseModal from "../../ui/BaseModal.vue";

const searchQuery = ref("");
const isAddModalOpen = ref(false);
const isEditModalOpen = ref(false);
const services = ref([]);

const toastMessage = ref("");
const toastVisible = ref(false);

const newService = reactive({
    name: "",
    desc: "",
    price: null,
});

const editingService = reactive({
    id: null,
    name: "",
    desc: "",
    price: null,
});

const filteredServices = computed(() => {
    if (!Array.isArray(services.value)) {
        return [];
    }
    const query = searchQuery.value.toLowerCase();
    return services.value.filter(
        (s) =>
            s.name.toLowerCase().includes(query) ||
            String(s.id).includes(query),
    );
});

function showToast(message) {
    toastMessage.value = message;
    toastVisible.value = true;
    setTimeout(() => (toastVisible.value = false), 3000);
}

async function loadServices() {
    const token = localStorage.getItem("token");
    try {
        const response = await fetch("/api/admin/services", {
            headers: { Authorization: `Bearer ${token}` },
        });
        if (response.ok) {
            services.value = await response.json();
        } else {
            showToast("Ошибка загрузки сервисов");
        }
    } catch (e) {
        console.error(e);
        showToast("Ошибка загрузки сервисов");
    }
}

onMounted(() => {
    loadServices();
});

async function saveService() {
    if (!newService.name || !newService.price) {
        showToast("Пожалуйста, заполните название и цену");
        return;
    }

    const token = localStorage.getItem("token");
    try {
        const response = await fetch("/api/admin/services", {
            method: "POST",
            headers: {
                Authorization: `Bearer ${token}`,
                "Content-Type": "application/json",
            },
            body: JSON.stringify(newService),
        });

        if (response.ok) {
            isAddModalOpen.value = false;
            newService.name = "";
            newService.desc = "";
            newService.price = null;

            await loadServices(); // ← перезагружаем список, т.к. бэкенд не возвращает созданный объект
            showToast("Услуга добавлена");
        } else {
            const err = await response.json().catch(() => ({}));
            showToast(err.error || "Ошибка сохранения сервиса");
        }
    } catch (e) {
        console.error(e);
        showToast("Ошибка сохранения сервиса");
    }
}

function openEditModal(service) {
    editingService.id = service.id;
    editingService.name = service.name;
    editingService.desc = service.desc;
    editingService.price = service.price;
    isEditModalOpen.value = true;
}

async function saveEditedService() {
    if (!editingService.name || !editingService.price) {
        showToast("Пожалуйста, заполните название и цену");
        return;
    }

    const token = localStorage.getItem("token");
    try {
        const response = await fetch(
            `/api/admin/services/${editingService.id}`,
            {
                method: "PUT",
                headers: {
                    Authorization: `Bearer ${token}`,
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({
                    name: editingService.name,
                    desc: editingService.desc,
                    price: editingService.price,
                }),
            },
        );

        if (response.ok) {
            isEditModalOpen.value = false;
            await loadServices(); // ← перезагружаем список вместо ручного обновления по индексу
            showToast("Услуга обновлена");
        } else {
            const err = await response.json().catch(() => ({}));
            showToast(err.error || "Ошибка обновления сервиса");
        }
    } catch (e) {
        console.error(e);
        showToast("Ошибка обновления сервиса");
    }
}

async function deleteService(id) {
    if (!confirm("Удалить эту услугу из справочника?")) {
        return;
    }

    const token = localStorage.getItem("token");
    try {
        const response = await fetch(`/api/admin/services/${id}`, {
            method: "DELETE",
            headers: { Authorization: `Bearer ${token}` },
        });

        if (response.ok) {
            services.value = services.value.filter((s) => s.id !== id);
            showToast("Услуга удалена");
        } else {
            const err = await response.json().catch(() => ({}));
            showToast(err.error || "Ошибка удаления услуги");
        }
    } catch (e) {
        console.error(e);
        showToast("Ошибка удаления сервиса");
    }
}
</script>

<template>
    <div class="page">
        <div class="page-header">
            <div>
                <div class="page-title">Справочник услуг</div>
                <div class="page-sub">
                    Управление каталогом процедур и ценами
                </div>
            </div>
            <button class="btn btn-primary" @click="isAddModalOpen = true">
                ✚ Добавить услугу
            </button>
        </div>

        <div class="card">
            <div class="card-header">
                <input
                    type="text"
                    v-model="searchQuery"
                    placeholder="Поиск по названию или ID..."
                    style="max-width: 300px"
                />
            </div>
            <div class="table-wrap">
                <table>
                    <thead>
                        <tr>
                            <th style="width: 80px">#</th>
                            <th>Наименование</th>
                            <th>Стоимость</th>
                            <th style="width: 120px">Действия</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="s in filteredServices" :key="s.id">
                            <td class="mono text-muted">{{ s.id }}</td>
                            <td class="td-main">
                                {{ s.name }}
                                <div class="service-desc">{{ s.desc }}</div>
                            </td>
                            <td class="mono">{{ s.price }} руб.</td>
                            <td>
                                <div class="row">
                                    <button
                                        class="btn btn-ghost btn-sm"
                                        @click="openEditModal(s)"
                                        title="Редактировать"
                                    >
                                        ✏
                                    </button>
                                    <button
                                        class="btn btn-danger btn-sm"
                                        @click="deleteService(s.id)"
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
            title="Добавить новую услугу"
            @close="isAddModalOpen = false"
        >
            <div class="form-grid">
                <div class="form-group full">
                    <label>Наименование услуги *</label>
                    <input
                        v-model="newService.name"
                        type="text"
                        placeholder="Напр. Чистка зубов"
                    />
                </div>
                <div class="form-group full">
                    <label>Описание</label>
                    <textarea
                        v-model="newService.desc"
                        placeholder="Краткое описание процедуры..."
                    ></textarea>
                </div>
                <div class="form-group full">
                    <label>Стоимость (руб.) *</label>
                    <input
                        v-model="newService.price"
                        type="number"
                        placeholder="0"
                    />
                </div>
            </div>
            <template #footer>
                <button class="btn btn-ghost" @click="isAddModalOpen = false">
                    Отмена
                </button>
                <button class="btn btn-primary" @click="saveService">
                    Сохранить
                </button>
            </template>
        </BaseModal>

        <!-- МОДАЛКА РЕДАКТИРОВАНИЯ -->
        <BaseModal
            :show="isEditModalOpen"
            title="Редактировать услугу"
            @close="isEditModalOpen = false"
        >
            <div class="form-grid">
                <div class="form-group full">
                    <label>Наименование услуги *</label>
                    <input
                        v-model="editingService.name"
                        type="text"
                        placeholder="Напр. Чистка зубов"
                    />
                </div>
                <div class="form-group full">
                    <label>Описание</label>
                    <textarea
                        v-model="editingService.desc"
                        placeholder="Краткое описание процедуры..."
                    ></textarea>
                </div>
                <div class="form-group full">
                    <label>Стоимость (руб.) *</label>
                    <input
                        v-model="editingService.price"
                        type="number"
                        placeholder="0"
                    />
                </div>
            </div>
            <template #footer>
                <button class="btn btn-ghost" @click="isEditModalOpen = false">
                    Отмена
                </button>
                <button class="btn btn-primary" @click="saveEditedService">
                    Сохранить изменения
                </button>
            </template>
        </BaseModal>

        <!-- ТОСТ -->
        <div v-if="toastVisible" class="toast">{{ toastMessage }}</div>
    </div>
</template>

<style scoped>
.service-desc {
    font-size: 11px;
    font-weight: 400;
    color: var(--text3);
    margin-top: 2px;
}

textarea {
    min-height: 100px;
    resize: vertical;
}

table {
    border-spacing: 0;
    width: 100%;
}

tr:last-child td {
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
