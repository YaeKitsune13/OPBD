<script setup lang="ts">
import { ref, onMounted, nextTick } from "vue";
import html2pdf from "html2pdf.js";

const orders = ref<any[]>([]);
const loading = ref(true);
const receiptRef = ref<HTMLElement | null>(null);
const selectedOrder = ref<any>(null);

const statusMap: Record<string, { label: string; cls: string }> = {
    paid: { label: "Оплачен", cls: "badge-confirmed" },
    confirmed: { label: "Подтвержден", cls: "badge-info" },
    delivered: { label: "Доставлен", cls: "badge-waiting" },
};

async function loadOrders() {
    loading.value = true;
    try {
        const user = JSON.parse(localStorage.getItem("user") || "{}");
        const token = localStorage.getItem("token");
        const res = await fetch(`/api/orders/${user.id}`, {
            headers: { Authorization: `Bearer ${token}` },
        });
        if (res.ok) {
            orders.value = await res.json();
        }
    } catch (e) {
        console.error(e);
    } finally {
        loading.value = false;
    }
}

async function printReceipt(order: any) {
    selectedOrder.value = order;
    await nextTick();

    const element = receiptRef.value;
    const opt = {
        margin: 10,
        filename: `order_receipt_${order.id}.pdf`,
        image: { type: "jpeg", quality: 1 },
        html2canvas: { scale: 3, logging: false, useCORS: true },
        jsPDF: { unit: "mm", format: "a4", orientation: "portrait" },
    };

    html2pdf()
        .set(opt)
        .from(element)
        .toPdf()
        .get("pdf")
        .then((pdf: any) => {
            window.open(pdf.output("bloburl"), "_blank");
        });
}

onMounted(loadOrders);
</script>

<template>
    <div class="page">
        <div class="page-header">
            <div>
                <h1 class="page-title">Мои заказы</h1>
                <p class="page-sub">История покупок в аптеке</p>
            </div>
        </div>

        <div v-if="loading" class="text-muted py-40 center">Загрузка...</div>

        <div
            v-else-if="!orders || orders.length === 0"
            class="card py-40 center"
        >
            <p class="text-muted">История заказов пуста</p>
        </div>

        <div v-else class="table-wrap">
            <table>
                <thead>
                    <tr>
                        <th>№ Заказа</th>
                        <th>Дата</th>
                        <th>Сумма</th>
                        <th>Статус</th>
                        <th style="text-align: right">Действие</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for="order in orders" :key="order.id">
                        <td class="td-main">#{{ order.id }}</td>
                        <td class="mono">{{ order.date }}</td>
                        <td class="text-accent" style="font-weight: 700">
                            {{ order.totalAmount }} ₽
                        </td>
                        <td>
                            <span
                                v-if="statusMap[order.status]"
                                :class="['badge', statusMap[order.status].cls]"
                            >
                                {{ statusMap[order.status].label }}
                            </span>
                        </td>
                        <td style="text-align: right">
                            <button
                                class="btn btn-ghost btn-sm"
                                @click="printReceipt(order)"
                            >
                                📄 Чек
                            </button>
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>

        <div style="position: absolute; left: -9999px; top: 0">
            <div ref="receiptRef" class="receipt-container">
                <div class="receipt-card">
                    <div class="receipt-brand">
                        <div class="receipt-logo">🐾</div>
                        <h2>ВЕТКЛИНИКА "ВетКлиника"</h2>
                        <p>Лицензия №77-14-001234</p>
                    </div>

                    <div class="receipt-meta">
                        <div>
                            <strong>Чек №:</strong> {{ selectedOrder?.id }}
                        </div>
                        <div>
                            <strong>Дата:</strong> {{ selectedOrder?.date }}
                        </div>
                        <div>
                            <strong>Статус:</strong>
                            {{ statusMap[selectedOrder?.status]?.label }}
                        </div>
                    </div>

                    <div class="receipt-divider"></div>

                    <table class="items-table">
                        <thead>
                            <tr>
                                <th>Наименование</th>
                                <th style="text-align: center">Кол-во</th>
                                <th style="text-align: right">Цена</th>
                                <th style="text-align: right">Сумма</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr
                                v-for="item in selectedOrder?.items"
                                :key="item.name"
                            >
                                <td>{{ item.name }}</td>
                                <td style="text-align: center">
                                    {{ item.quantity }}
                                </td>
                                <td style="text-align: right">
                                    {{ item.price }} ₽
                                </td>
                                <td style="text-align: right">
                                    {{ item.price * item.quantity }} ₽
                                </td>
                            </tr>
                        </tbody>
                    </table>

                    <div class="receipt-divider"></div>

                    <div class="receipt-total">
                        <div class="total-row">
                            <span>Итого к оплате:</span>
                            <span class="total-val"
                                >{{ selectedOrder?.totalAmount }} ₽</span
                            >
                        </div>
                    </div>

                    <div class="receipt-footer">
                        <p>Электронный чек сформирован в системе PetCare</p>
                        <p>
                            Спасибо, что доверяете нам здоровье ваших питомцев!
                        </p>
                        <div class="qr-placeholder">QR-КОД ОПЛАТЫ</div>
                    </div>
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

.receipt-container {
    width: 190mm;
    background: white;
    padding: 0;
    margin: 0;
}

.receipt-card {
    padding: 40px;
    color: #1a1c21;
}

.receipt-brand {
    text-align: center;
    margin-bottom: 30px;
}

.receipt-logo {
    font-size: 40px;
    margin-bottom: 10px;
}

.receipt-brand h2 {
    margin: 0;
    font-size: 24px;
    letter-spacing: 2px;
}

.receipt-brand p {
    color: #64748b;
    font-size: 12px;
}

.receipt-meta {
    display: flex;
    justify-content: space-between;
    font-size: 14px;
    margin-bottom: 20px;
}

.receipt-divider {
    border-top: 2px dashed #e2e8f0;
    margin: 20px 0;
}

.items-table {
    width: 100%;
    border-collapse: collapse;
}

.items-table th {
    text-align: left;
    font-size: 12px;
    text-transform: uppercase;
    color: #64748b;
    padding-bottom: 10px;
}

.items-table td {
    padding: 12px 0;
    font-size: 14px;
    border-bottom: 1px solid #f1f5f9;
}

.receipt-total {
    display: flex;
    justify-content: flex-end;
    margin-top: 20px;
}

.total-row {
    display: flex;
    gap: 20px;
    align-items: center;
}

.total-row span:first-child {
    font-size: 16px;
    font-weight: 600;
}

.total-val {
    font-size: 28px;
    font-weight: 800;
    color: #10b981;
}

.receipt-footer {
    margin-top: 50px;
    text-align: center;
    font-size: 12px;
    color: #94a3b8;
}

.qr-placeholder {
    margin: 20px auto 0;
    width: 80px;
    height: 80px;
    border: 1px solid #e2e8f0;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 8px;
    color: #cbd5e1;
}
</style>
