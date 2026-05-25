<script setup lang="ts">
import { ref, reactive, computed, onMounted, watch } from "vue";
import { useToast } from "../../utils/useToast";

const emit = defineEmits<{ navigate: [page: string] }>();
const { showToast } = useToast();

interface Pet {
    id: number;
    name: string;
    avatar: string;
    breed?: string;
    weight?: number;
    dob?: string;
}
interface Doctor {
    doctorId: number;
    lastName: string;
    firstName: string;
    specialization: string;
}
interface Service {
    id: number;
    name: string;
    price: number;
}

const pets = ref<Pet[]>([]);
const doctors = ref<Doctor[]>([]);
const services = ref<Service[]>([]);
const busySlots = ref<string[]>([]);
const loading = ref(false);

const TIME_SLOTS = [
    "09:00",
    "10:00",
    "11:00",
    "13:00",
    "14:00",
    "15:00",
    "16:00",
];

const form = reactive({
    petId: null as number | null,
    serviceId: null as number | null,
    specialization: "",
    doctorId: null as number | null,
    date: "",
    time: "",
    comment: "",
});

const token = localStorage.getItem("token");
const user = JSON.parse(localStorage.getItem("user") || "{}");

async function loadData() {
    console.group("🚀 loadData: Инициализация данных");
    console.log("User ID:", user.id);
    loading.value = true;
    try {
        const res = await fetch(`/api/book/init/${user.id}`, {
            headers: { Authorization: `Bearer ${token}` },
        });
        console.log("Статус ответа init:", res.status);

        if (res.ok) {
            const data = await res.json();
            console.log("Данные успешно получены:", data);
            pets.value = data.pets || [];
            doctors.value = data.doctors || [];
            services.value = data.services || [];
        } else {
            console.error("Ошибка при загрузке init данных");
            showToast("Не удалось загрузить данные для записи", "error");
        }
    } catch (err) {
        console.error("Критическая ошибка loadData:", err);
        showToast("Ошибка подключения к серверу", "error");
    } finally {
        loading.value = false;
        console.groupEnd();
    }
}

async function fetchBusySlots() {
    console.group("🕒 fetchBusySlots: Проверка занятых слотов");
    console.log("Параметры запроса:", {
        doctorId: form.doctorId,
        date: form.date,
    });

    if (!form.doctorId || !form.date) {
        console.log("Запрос отменен: не выбран врач или дата");
        console.groupEnd();
        return;
    }

    try {
        const res = await fetch(
            `/api/appointments/busy-slots?doctor_id=${form.doctorId}&date=${form.date}`,
            { headers: { Authorization: `Bearer ${token}` } },
        );

        if (res.ok) {
            const data = await res.json();
            console.log("Занятые слоты получены:", data);
            busySlots.value = data;
        } else {
            console.warn("Сервер вернул ошибку при получении слотов");
            busySlots.value = [];
        }
    } catch (err) {
        console.error("Ошибка при fetchBusySlots:", err);
        busySlots.value = [];
    } finally {
        console.groupEnd();
    }
}

watch(
    () => [form.doctorId, form.date],
    () => {
        console.log("👀 Watcher (doctorId/date): Изменение параметров времени");
        fetchBusySlots();
    },
);

watch(
    () => form.specialization,
    (newVal) => {
        console.log(
            "👀 Watcher (specialization): Смена направления на",
            newVal,
        );
        form.doctorId = null;
        form.time = "";
    },
);

const specializations = computed(() => {
    const specs = [...new Set(doctors.value.map((d) => d.specialization))];
    return specs;
});

const filteredDoctors = computed(() => {
    const filtered = doctors.value.filter(
        (d) => d.specialization === form.specialization,
    );
    return filtered;
});

const selectedService = computed(() =>
    services.value.find((s) => s.id === form.serviceId),
);
const selectedDoctor = computed(() =>
    doctors.value.find((d) => d.doctorId === form.doctorId),
);
const minDate = new Date().toISOString().split("T")[0];

async function sendRequest() {
    console.group("📤 sendRequest: Отправка формы");

    if (!form.petId || !form.doctorId || !form.date || !form.time) {
        console.warn("Валидация не прошла");
        console.groupEnd();
        return showToast("Пожалуйста, заполните все обязательные поля", "info");
    }

    const payload = {
        pet_id: form.petId,
        doctor_id: form.doctorId,
        service_id: form.serviceId,
        scheduled_at: `${form.date}T${form.time}:00`,
        comment: form.comment,
    };

    loading.value = true;
    try {
        const res = await fetch("/api/appointments", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                Authorization: `Bearer ${token}`,
            },
            body: JSON.stringify(payload),
        });

        if (res.ok) {
            showToast("Вы успешно записаны на приём!", "success");

            Object.assign(form, {
                petId: null,
                serviceId: null,
                specialization: "",
                doctorId: null,
                date: "",
                time: "",
                comment: "",
            });

            emit("navigate", "appointments");
        } else {
            const errData = await res.json();
            showToast(errData.error || "Ошибка при создании записи", "error");
        }
    } catch (err) {
        console.error("❌ Ошибка:", err);
        showToast("Не удалось отправить запрос", "error");
    } finally {
        loading.value = false;
        console.groupEnd();
    }
}

onMounted(() => {
    loadData();
});
</script>

<template>
    <div class="page">
        <div class="page-header">
            <div>
                <h1 class="page-title">Записаться на визит</h1>
                <p class="page-sub">Выберите время для посещения клиники</p>
            </div>
        </div>

        <div class="card booking-card">
            <div class="card-body">
                <div class="booking-row">
                    <div class="form-group">
                        <label>Питомец *</label>
                        <select v-model.number="form.petId">
                            <option :value="null">— Выбрать —</option>
                            <option v-for="p in pets" :key="p.id" :value="p.id">
                                {{ p.avatar }} {{ p.name }}
                            </option>
                        </select>
                    </div>
                    <div class="form-group">
                        <label>Цель визита</label>
                        <select v-model.number="form.serviceId">
                            <option :value="null">— Осмотр —</option>
                            <option
                                v-for="s in services"
                                :key="s.id"
                                :value="s.id"
                            >
                                {{ s.name }} ({{ s.price }}₽)
                            </option>
                        </select>
                    </div>
                    <div class="form-group">
                        <label>Направление *</label>
                        <select v-model="form.specialization">
                            <option value="">— Профиль —</option>
                            <option
                                v-for="s in specializations"
                                :key="s"
                                :value="s"
                            >
                                {{ s }}
                            </option>
                        </select>
                    </div>
                    <div class="form-group">
                        <label>Врач *</label>
                        <select
                            v-model.number="form.doctorId"
                            :disabled="!form.specialization"
                        >
                            <option :value="null">— Выберите —</option>
                            <option
                                v-for="d in filteredDoctors"
                                :key="d.doctorId"
                                :value="d.doctorId"
                            >
                                {{ d.lastName }} {{ d.firstName }}
                            </option>
                        </select>
                    </div>
                    <div class="form-group">
                        <label>Дата *</label>
                        <input v-model="form.date" type="date" :min="minDate" />
                    </div>
                    <div class="form-group">
                        <label>Время *</label>
                        <select
                            v-model="form.time"
                            :disabled="!form.date || !form.doctorId"
                        >
                            <option value="">--:--</option>
                            <option
                                v-for="s in TIME_SLOTS"
                                :key="s"
                                :value="s"
                                :disabled="busySlots.includes(s)"
                            >
                                {{ s }}
                                {{ busySlots.includes(s) ? "(занято)" : "" }}
                            </option>
                        </select>
                    </div>
                </div>

                <div class="form-group mt-16">
                    <label>Комментарий</label>
                    <textarea
                        v-model="form.comment"
                        rows="2"
                        placeholder="Опишите жалобы или причину обращения..."
                    ></textarea>
                </div>

                <div
                    v-if="form.date && form.time && form.doctorId"
                    class="summary-plate"
                >
                    <div class="summary-info">
                        <div class="summary-date">
                            {{ form.date.split("-").reverse().join(".") }} в
                            {{ form.time }}
                        </div>
                        <div class="summary-details">
                            Врач: {{ selectedDoctor?.lastName }}
                            {{ selectedDoctor?.firstName }}
                            <span v-if="selectedService">
                                • {{ selectedService.name }}</span
                            >
                        </div>
                    </div>
                    <div class="summary-price" v-if="selectedService">
                        {{ selectedService.price }} ₽
                    </div>
                </div>

                <div class="row mt-20" style="gap: 12px">
                    <button
                        class="btn btn-primary"
                        :disabled="loading"
                        @click="sendRequest"
                    >
                        {{ loading ? "Отправка..." : "✓ Подтвердить запись" }}
                    </button>
                    <button
                        class="btn btn-ghost"
                        @click="emit('navigate', 'dashboard')"
                    >
                        Отмена
                    </button>
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped>
.booking-row {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
    gap: 15px;
}
.booking-card {
    border-radius: 16px;
    border: 1px solid var(--border);
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.03);
}
label {
    font-size: 10px;
    text-transform: uppercase;
    color: var(--text3);
    font-weight: 700;
    margin-bottom: 5px;
    display: block;
}
select,
input,
textarea {
    width: 100%;
    padding: 10px;
    border-radius: 8px;
    border: 1px solid var(--border);
    background: var(--surface2);
    font-family: inherit;
    font-size: 14px;
    transition: all 0.2s;
}
select:focus,
input:focus,
textarea:focus {
    outline: none;
    border-color: var(--accent);
    background: var(--surface);
}

.summary-plate {
    margin-top: 24px;
    padding: 16px 20px;
    background: var(--accent-dim);
    border: 1px solid var(--accent);
    border-radius: 12px;
    display: flex;
    justify-content: space-between;
    align-items: center;
}
.summary-date {
    font-weight: 800;
    font-size: 17px;
    color: var(--text);
}
.summary-details {
    font-size: 13px;
    color: var(--text2);
    margin-top: 4px;
}
.summary-price {
    font-size: 22px;
    font-weight: 900;
    color: var(--accent);
}

.mt-16 {
    margin-top: 16px;
}
.mt-20 {
    margin-top: 20px;
}
</style>
