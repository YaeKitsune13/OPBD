<script setup lang="ts">
import { ref, onMounted, nextTick } from "vue";
import BaseModal from "../elements/BaseModal.vue";
import { useToast } from "../../utils/useToast";
import html2pdf from "html2pdf.js";

const { showToast } = useToast();

const appointments = ref<any[]>([]);
const loading = ref(true);
const isProtocolOpen = ref(false);
const selectedProtocol = ref<any>(null);
const protocolRef = ref<HTMLElement | null>(null);

const statusMap: Record<string, { label: string; cls: string }> = {
    waiting: { label: "Ожидает", cls: "badge-waiting" },
    confirmed: { label: "Принята", cls: "badge-confirmed" },
    rejected: { label: "Отклонена", cls: "badge-rejected" },
    done: { label: "Завершена", cls: "badge-info" },
};

async function loadApps() {
    loading.value = true;
    try {
        const user = JSON.parse(localStorage.getItem("user") || "{}");
        const token = localStorage.getItem("token");

        const res = await fetch(`/api/appointments/client/${user.id}`, {
            headers: { Authorization: `Bearer ${token}` },
        });

        if (res.ok) {
            appointments.value = (await res.json()) || [];
        } else {
            showToast("Ошибка при загрузке данных", "error");
        }
    } catch (error) {
        console.error(error);
        showToast("Ошибка соединения с сервером", "error");
    } finally {
        loading.value = false;
    }
}

function showProtocol(app: any) {
    selectedProtocol.value = app;
    isProtocolOpen.value = true;
}
function diagnosisNames(protocol: any) {
    return (protocol?.diagnoses || []).map((d: any) => d.name).join(", ");
}
async function printProtocol() {
    showToast("Генерация PDF документа...", "info");

    await nextTick();
    const element = protocolRef.value;

    if (!element) return;

    const opt = {
        margin: [10, 10, 10, 10],
        filename: `Медицинское_заключение_${selectedProtocol.value?.petName || "пациент"}.pdf`,
        image: { type: "jpeg", quality: 0.98 },
        html2canvas: { scale: 3, useCORS: true, letterRendering: true },
        jsPDF: { unit: "mm", format: "a4", orientation: "portrait" },
    };

    try {
        const pdf = await html2pdf().set(opt).from(element).toPdf().get("pdf");
        window.open(pdf.output("bloburl"), "_blank");
    } catch (err) {
        showToast("Ошибка при создании PDF", "error");
    }
}

onMounted(loadApps);
</script>

<template>
    <div class="page">
        <div class="page-header">
            <div>
                <h1 class="page-title">Мои визиты</h1>
                <p class="page-sub">
                    История посещений и медицинские заключения
                </p>
            </div>
            <button
                class="btn btn-ghost btn-sm"
                @click="loadApps"
                :disabled="loading"
            >
                🔄 Обновить
            </button>
        </div>

        <div v-if="loading" class="text-muted py-40 center">Загрузка...</div>
        <div v-else-if="!appointments.length" class="card py-40 center">
            <p class="text-muted">Записей не найдено</p>
        </div>
        <div v-else class="table-wrap">
            <table>
                <thead>
                    <tr>
                        <th>Питомец</th>
                        <th>Дата и время</th>
                        <th>Врач</th>
                        <th>Статус</th>
                        <th></th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for="app in appointments" :key="app.id">
                        <td class="td-main">{{ app.petName }}</td>
                        <td>
                            <div class="mono" style="font-size: 13px">
                                {{ app.date }}
                            </div>
                            <div class="text-muted" style="font-size: 11px">
                                {{ app.time }}
                            </div>
                        </td>
                        <td>{{ app.doctorName }}</td>
                        <td>
                            <span
                                v-if="statusMap[app.status]"
                                :class="['badge', statusMap[app.status].cls]"
                            >
                                {{ statusMap[app.status].label }}
                            </span>
                        </td>
                        <td style="text-align: right">
                            <button
                                v-if="app.status === 'done'"
                                class="btn btn-primary btn-sm"
                                @click="showProtocol(app)"
                            >
                                Протокол
                            </button>
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>

        <BaseModal
            :show="isProtocolOpen"
            title="Медицинское заключение"
            @close="isProtocolOpen = false"
        >
            <div class="protocol-view" v-if="selectedProtocol">
                <template v-if="selectedProtocol.protocol">
                    <div class="p-summary">
                        <div class="p-summary-item">
                            <label>Пациент</label>
                            <span>{{ selectedProtocol.petName }}</span>
                        </div>
                        <div class="p-summary-item">
                            <label>Вес</label>
                            <span
                                >{{
                                    selectedProtocol.protocol.weight ||
                                    selectedProtocol.protocol.weight_at_visit
                                }}
                                кг</span
                            >
                        </div>
                    </div>
                    <div class="p-section">
                        <label>Поставленный диагноз</label>
                        <div class="p-val diagnosis-tags">
                            <span
                                v-for="d in selectedProtocol.protocol.diagnoses"
                                :key="d.id"
                                class="diagnosis-tag-readonly"
                            >
                                {{ d.name }}
                            </span>
                            <span
                                v-if="
                                    !selectedProtocol.protocol.diagnoses?.length
                                "
                                class="text-muted"
                            >
                                Диагноз не указан
                            </span>
                        </div>
                    </div>
                    <div class="p-section">
                        <label>Рекомендации и лечение</label>
                        <div class="p-val">
                            {{ selectedProtocol.protocol.treatment }}
                        </div>
                    </div>
                    <button
                        class="btn btn-accent full mt-20"
                        @click="printProtocol"
                    >
                        🖨 Сохранить в PDF
                    </button>
                </template>
                <div v-else class="center py-20">
                    <p class="text-muted">
                        Протокол приема еще не заполнен врачом.
                    </p>
                </div>
            </div>
        </BaseModal>

        <div style="position: absolute; left: -9999px; top: 0">
            <div ref="protocolRef" class="pdf-document">
                <div class="pdf-header">
                    <div class="pdf-logo">🐾</div>
                    <div class="pdf-clinic-info">
                        <h2>ВЕТЕРИНАРНЫЙ ЦЕНТР "ВетКлиника"</h2>
                        <p>Лицензия № ВЕТ-77-01-000842 | petcare-clinic.ru</p>
                        <p>
                            г. Москва, ул. Ветеринарная, 12 | +7 (495) 000-00-00
                        </p>
                    </div>
                </div>

                <h1 class="pdf-title">
                    МЕДИЦИНСКОЕ ЗАКЛЮЧЕНИЕ №{{ selectedProtocol?.id }}
                </h1>

                <div class="pdf-meta-grid">
                    <div class="pdf-meta-item">
                        <strong>Пациент:</strong>
                        {{ selectedProtocol?.petName }}
                    </div>
                    <div class="pdf-meta-item">
                        <strong>Вид животного:</strong>
                        {{
                            selectedProtocol?.petSpecies || "Домашнее животное"
                        }}
                    </div>
                    <div class="pdf-meta-item">
                        <strong>Дата приема:</strong>
                        {{ selectedProtocol?.date }}
                    </div>
                    <div class="pdf-meta-item">
                        <strong>Вес пациента:</strong>
                        {{
                            selectedProtocol?.protocol?.weight ||
                            selectedProtocol?.protocol?.weight_at_visit
                        }}
                        кг
                    </div>
                    <div class="pdf-meta-item">
                        <strong>Лечащий врач:</strong>
                        {{ selectedProtocol?.doctorName }}
                    </div>
                </div>

                <div class="pdf-content-block">
                    <h3 class="pdf-section-h">1. Анамнез и диагноз</h3>
                    <div class="pdf-text-box">
                        {{ diagnosisNames(selectedProtocol?.protocol) }}
                    </div>
                </div>

                <div class="pdf-content-block">
                    <h3 class="pdf-section-h">
                        2. Назначенное лечение и рекомендации
                    </h3>
                    <div class="pdf-text-box">
                        {{ selectedProtocol?.protocol?.treatment }}
                    </div>
                </div>

                <div
                    class="pdf-content-block"
                    v-if="selectedProtocol?.protocol?.medications"
                >
                    <h3 class="pdf-section-h">3. Медикаментозные назначения</h3>
                    <div class="pdf-text-box highlight">
                        {{ selectedProtocol?.protocol?.medications }}
                    </div>
                </div>

                <div class="pdf-footer">
                    <div class="pdf-signature">
                        <div class="pdf-sig-line"></div>
                        <p>
                            Подпись врача: ____________________ /
                            {{ selectedProtocol?.doctorName }} /
                        </p>
                    </div>
                    <div class="pdf-stamp">ДЛЯ ДОКУМЕНТОВ</div>
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped>
.center {
    text-align: center;
}
.py-40 {
    padding: 40px 0;
}
.protocol-view {
    display: flex;
    flex-direction: column;
    gap: 16px;
}
.p-summary {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 12px;
    background: var(--surface2);
    padding: 14px;
    border-radius: 10px;
}
.p-summary-item label {
    display: block;
    font-size: 10px;
    text-transform: uppercase;
    color: var(--text3);
    font-weight: 800;
}
.p-summary-item span {
    font-weight: 700;
    color: var(--text);
}
.p-section label {
    font-size: 11px;
    font-weight: 700;
    color: var(--accent);
    text-transform: uppercase;
    display: block;
    margin-bottom: 6px;
}
.p-val {
    background: var(--surface3);
    padding: 12px;
    border-radius: 8px;
    font-size: 14px;
    line-height: 1.6;
}
.full {
    width: 100%;
}

.pdf-document {
    width: 180mm;
    padding: 0;
    color: #1a1a1a;
    font-family: "DejaVu Sans", "Arial", sans-serif;
    line-height: 1.5;
}

.pdf-header {
    display: flex;
    align-items: center;
    gap: 20px;
    border-bottom: 3px solid #10b981;
    padding-bottom: 15px;
    margin-bottom: 25px;
}

.pdf-logo {
    font-size: 50px;
}
.pdf-clinic-info h2 {
    margin: 0;
    font-size: 18px;
    color: #10b981;
    letter-spacing: 1px;
}
.pdf-clinic-info p {
    margin: 2px 0;
    font-size: 11px;
    color: #4b5563;
}

.pdf-title {
    text-align: center;
    font-size: 20px;
    font-weight: 900;
    margin-bottom: 30px;
    text-decoration: underline;
}

.pdf-meta-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 10px 20px;
    margin-bottom: 30px;
    background: #f9fafb;
    padding: 15px;
    border-radius: 8px;
}

.pdf-meta-item {
    font-size: 13px;
}

.pdf-content-block {
    margin-bottom: 25px;
}
.pdf-section-h {
    font-size: 13px;
    text-transform: uppercase;
    color: #374151;
    border-bottom: 1px solid #e5e7eb;
    padding-bottom: 5px;
    margin-bottom: 10px;
}

.pdf-text-box {
    padding: 10px;
    min-height: 40px;
    font-size: 14px;
    text-align: justify;
}

.pdf-text-box.highlight {
    background: #f0fdf4;
    border: 1px solid #10b981;
    border-radius: 6px;
}

.pdf-footer {
    margin-top: 50px;
    display: flex;
    justify-content: space-between;
    align-items: flex-end;
}

.pdf-signature {
    font-size: 12px;
}
.pdf-sig-line {
    width: 250px;
    border-bottom: 1px solid #000;
    margin-bottom: 8px;
}

.pdf-stamp {
    width: 100px;
    height: 100px;
    border: 3px double #10b981;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    text-align: center;
    font-size: 9px;
    font-weight: 900;
    color: #10b981;
    transform: rotate(-15deg);
    opacity: 0.6;
}
.diagnosis-tags {
    display: flex;
    flex-wrap: wrap;
    gap: 6px;
}
.diagnosis-tag-readonly {
    font-size: 12px;
    font-weight: 700;
    padding: 3px 9px;
    border-radius: 999px;
    background: var(--accent-soft, #e6f4ea);
    color: var(--accent, #1a8f4c);
}
</style>
