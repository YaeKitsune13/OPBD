<script setup lang="ts">
import { ref, computed, onMounted, watch } from "vue";
import { Line } from "vue-chartjs";
import {
    Chart as ChartJS,
    CategoryScale,
    LinearScale,
    PointElement,
    LineElement,
    Title,
    Tooltip,
    Filler,
} from "chart.js";

ChartJS.register(
    CategoryScale,
    LinearScale,
    PointElement,
    LineElement,
    Title,
    Tooltip,
    Filler,
);

const userId = JSON.parse(localStorage.getItem("user") || "{}")?.id;
const token = localStorage.getItem("token");

const pets = ref<any[]>([]);
const selectedPetId = ref<number | null>(null);
const history = ref<any[]>([]);
const loading = ref(false);

async function loadPets() {
    try {
        const res = await fetch(`/api/stats/pets/${userId}`, {
            headers: { Authorization: `Bearer ${token}` },
        });
        if (res.ok) {
            const data = await res.json();
            pets.value = data;
            if (pets.value.length > 0) {
                selectedPetId.value = pets.value[0].id;
            }
        }
    } catch (err) {
        console.error("Ошибка при загрузке питомцев:", err);
    }
}

async function loadWeightData() {
    if (!selectedPetId.value) return;
    loading.value = true;
    try {
        const res = await fetch(`/api/stats/weight/${selectedPetId.value}`, {
            headers: { Authorization: `Bearer ${token}` },
        });
        if (res.ok) {
            history.value = await res.json();
        } else {
            history.value = [];
        }
    } catch (err) {
        console.error("Ошибка при загрузке данных веса:", err);
        history.value = [];
    } finally {
        loading.value = false;
    }
}

watch(selectedPetId, loadWeightData);

onMounted(loadPets);

const chartData = computed(() => ({
    labels: history.value.map((h) => h.date),
    datasets: [
        {
            label: "Вес (кг)",
            data: history.value.map((h) => h.weight),
            borderColor: "#10b981",
            backgroundColor: "rgba(16, 185, 129, 0.1)",
            fill: true,
            tension: 0.4,
            pointBackgroundColor: "#fff",
            pointBorderColor: "#10b981",
            pointRadius: 4,
        },
    ],
}));

const chartOptions = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
        legend: { display: false },
        tooltip: {
            backgroundColor: "#1a1c21",
            padding: 12,
            titleFont: { size: 14, weight: "bold" as const },
            bodyFont: { size: 13 },
            displayColors: false,
        },
    },
    scales: {
        y: {
            beginAtZero: false,
            grid: { color: "rgba(0,0,0,0.05)" },
        },
        x: {
            grid: { display: false },
        },
    },
};

function handlePrint() {
    window.print();
}
</script>

<template>
    <div class="page">
        <div class="page-header no-print">
            <div>
                <h1 class="page-title">Аналитика здоровья</h1>
                <p class="page-sub">
                    Динамика веса питомца по результатам осмотров
                </p>
            </div>
            <div class="header-actions">
                <select
                    v-model="selectedPetId"
                    class="select-custom"
                    :disabled="loading"
                >
                    <option v-if="!pets.length" value="">Загрузка...</option>
                    <option v-for="p in pets" :key="p.id" :value="p.id">
                        {{ p.name }}
                    </option>
                </select>
            </div>
        </div>

        <div class="card mt-20">
            <div class="card-header no-print">
                <span class="card-title">График изменения массы тела</span>
            </div>
            <div class="card-body chart-container">
                <div v-if="loading" class="loader">
                    <span class="text-muted">Получение данных...</span>
                </div>

                <template v-else>
                    <Line
                        v-if="history.length"
                        :data="chartData"
                        :options="chartOptions"
                    />
                    <div v-else class="empty-state">
                        <div class="empty-icon">📊</div>
                        <p class="text-muted">
                            Нет данных о взвешиваниях для этого питомца
                        </p>
                        <small
                            >Данные появляются автоматически после визитов к
                            врачу</small
                        >
                    </div>
                </template>
            </div>
        </div>

        <div class="print-only-info">
            <hr />
            <p>Отчет сформирован: {{ new Date().toLocaleDateString() }}</p>
            <p>Клиент ID: {{ userId }} | Питомец ID: {{ selectedPetId }}</p>
        </div>
    </div>
</template>

<style scoped>
.header-actions {
    display: flex;
    gap: 12px;
    align-items: center;
}

.select-custom {
    padding: 10px 16px;
    border-radius: var(--radius);
    border: 1px solid var(--border);
    background: var(--surface);
    font-weight: 600;
    cursor: pointer;
    min-width: 160px;
}

.chart-container {
    height: 400px;
    position: relative;
}

.loader,
.empty-state {
    height: 100%;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    text-align: center;
}

.empty-icon {
    font-size: 48px;
    margin-bottom: 12px;
    opacity: 0.5;
}

.empty-state small {
    margin-top: 8px;
    color: var(--text3);
}

.print-only-info {
    display: none;
    margin-top: 40px;
    font-size: 12px;
    color: var(--text3);
}

@media print {
    .no-print {
        display: none !important;
    }
    .print-only-info {
        display: block;
    }
    .card {
        border: none;
        box-shadow: none;
    }
    .page {
        padding: 0;
    }
}
</style>
