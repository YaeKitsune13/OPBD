<script setup lang="ts">
import { ref, computed, onMounted } from "vue";
import BaseModal from "../elements/BaseModal.vue";
import { useToast } from "../../utils/useToast";
const { showToast } = useToast();
interface Product {
    id?: number;
    ID?: number;
    name: string;
    description: string;
    price: number;
    category: string;
}
interface CartItem {
    id: number;
    productId: number;
    name: string;
    price: number;
    quantity: number;
}

const search = ref("");
const medicationsList = ref<Product[]>([]);
const cart = ref<CartItem[]>([]);
const isCartOpen = ref(false);
const loading = ref(false);
const isSubmittingOrder = ref(false);

const userRaw = localStorage.getItem("user");
const user = JSON.parse(userRaw || "{}");
const userId = user.id;
const token = localStorage.getItem("token");

console.log("[INIT] ShopPage загружен. Данные из LocalStorage:", {
    userId,
    token: token ? "есть" : "отсутствует",
});

const getHeaders = () => {
    return {
        Authorization: `Bearer ${token}`,
        "Content-Type": "application/json",
    };
};

async function loadData() {
    if (!userId) {
        console.error("[LOAD] Ошибка: userId не найден в localStorage");
        return;
    }
    console.group("[API] Загрузка данных магазина...");
    loading.value = true;
    try {
        const [resMeds, resCart] = await Promise.all([
            fetch("/api/medications", { headers: getHeaders() }),
            fetch(`/api/cart/${userId}`, { headers: getHeaders() }),
        ]);

        console.log("Статус медикаментов:", resMeds.status);
        console.log("Статус корзины:", resCart.status);

        if (resMeds.ok) {
            const medsData = await resMeds.json();
            medicationsList.value = medsData || [];
            console.log(
                "Список товаров загружен:",
                medicationsList.value.length,
            );
        } else {
            console.error("Ошибка загрузки товаров:", resMeds.statusText);
        }

        if (resCart.ok) {
            const cartData = await resCart.json();
            cart.value = cartData || [];
            console.log("Текущая корзина:", cart.value);
        } else {
            console.error("Ошибка загрузки корзины:", resCart.statusText);
        }
    } catch (err) {
        console.error("Критическая ошибка fetch:", err);
    } finally {
        loading.value = false;
        console.groupEnd();
    }
}

async function addToCart(product: Product) {
    const pid = product.id || product.ID;

    console.group(`[ACTION] Добавление в корзину: ${product.name}`);
    console.log("Распознанный ID:", pid);

    if (!pid) {
        console.error("ID продукта не найден! Объект продукта:", product);
        showToast("Ошибка: Не удалось определить ID товара", "error");
        console.groupEnd();
        return;
    }

    const payload = {
        productId: Number(pid),
        quantity: 1,
    };

    console.log("Payload для отправки:", payload);

    try {
        const res = await fetch(`/api/cart/${userId}`, {
            method: "POST",
            headers: getHeaders(),
            body: JSON.stringify(payload),
        });

        if (res.ok) {
            console.log("Успешно добавлено.");
            await loadData();
        } else {
            const errorText = await res.text();
            console.error("Сервер вернул 400/500:", errorText);
            showToast(`Ошибка добавления: ${errorText}`, "error");
        }
    } catch (err) {
        console.error("Ошибка сетевого запроса:", err);
    } finally {
        console.groupEnd();
    }
}

async function updateQty(itemId: number, delta: number) {
    const item = cart.value.find((i) => i.id === itemId);
    if (!item) return;

    const newQty = item.quantity + delta;
    if (newQty <= 0) return deleteFromCart(itemId);

    console.log(
        `[ACTION] Изменение кол-ва itemID ${itemId}: ${item.quantity} -> ${newQty}`,
    );

    try {
        const res = await fetch(`/api/cart/${itemId}`, {
            method: "PUT",
            headers: getHeaders(),
            body: JSON.stringify({ quantity: Number(newQty) }),
        });
        if (res.ok) {
            await loadData();
        }
    } catch (err) {
        console.error("Ошибка при обновлении кол-ва:", err);
    }
}

async function deleteFromCart(itemId: number) {
    console.log(`[ACTION] Удаление товара itemID: ${itemId}`);
    try {
        const res = await fetch(`/api/cart/${itemId}`, {
            method: "DELETE",
            headers: getHeaders(),
        });
        if (res.ok) await loadData();
    } catch (err) {
        console.error(err);
    }
}

async function placeOrder() {
    console.group("[ACTION] Оформление заказа");
    isSubmittingOrder.value = true;
    try {
        const res = await fetch("/api/orders", {
            method: "POST",
            headers: getHeaders(),
            body: JSON.stringify({
                userId: Number(userId),
                totalAmount: Number(cartTotal.value),
            }),
        });

        if (res.ok) {
            showToast("Заказ успешно оформлен!", "success");
            isCartOpen.value = false;
            await loadData();
        } else {
            const errText = await res.text();
            showToast(`Ошибка оформления: ${errText}`, "error");
        }
    } catch (err) {
        console.error(err);
    } finally {
        isSubmittingOrder.value = false;
        console.groupEnd();
    }
}

const filteredMeds = computed(() => {
    return medicationsList.value.filter((m) =>
        m.name.toLowerCase().includes(search.value.toLowerCase()),
    );
});

const cartTotal = computed(() => {
    if (!cart.value) return 0;
    return cart.value.reduce(
        (sum, item) => sum + item.price * item.quantity,
        0,
    );
});

onMounted(() => {
    loadData();
});
</script>

<template>
    <div class="page">
        <div class="page-header" style="min-height: 80px; align-items: center">
            <div>
                <h1 class="page-title">Аптека</h1>
                <p class="page-sub">Товары для здоровья ваших питомцев</p>
            </div>
            <button class="cart-btn" @click="isCartOpen = true">
                <span class="cart-label">Корзина</span>
                <span class="cart-count" v-if="cart && cart.length">{{
                    cart.length
                }}</span>
            </button>
        </div>

        <div class="search-wrap">
            <input
                v-model="search"
                type="text"
                placeholder="Поиск медикаментов..."
                class="styled-search"
            />
        </div>

        <div
            v-if="loading && medicationsList.length === 0"
            class="text-muted py-40 center"
        >
            Загрузка каталога...
        </div>

        <div v-else class="med-grid">
            <div
                v-for="med in filteredMeds"
                :key="med.id || med.ID"
                class="med-card"
            >
                <div class="med-content">
                    <span class="med-cat">{{ med.category }}</span>
                    <h3 class="med-title">{{ med.name }}</h3>
                    <p class="med-text">{{ med.description }}</p>
                </div>
                <div class="med-footer">
                    <span class="med-price"
                        >{{ med.price }} <small>₽</small></span
                    >
                    <button class="add-btn" @click="addToCart(med)">
                        В корзину
                    </button>
                </div>
            </div>
        </div>

        <BaseModal
            :show="isCartOpen"
            title="Ваша корзина"
            maxWidth="480px"
            @close="isCartOpen = false"
        >
            <div v-if="!cart || cart.length === 0" class="empty-cart">
                <p class="text-muted">Корзина пуста</p>
            </div>
            <div v-else class="cart-items-list">
                <div v-for="item in cart" :key="item.id" class="cart-item">
                    <div class="cart-item-info">
                        <div class="cart-item-name">{{ item.name }}</div>
                        <div class="cart-item-price">
                            {{ item.price }} ₽ / шт.
                        </div>
                    </div>
                    <div class="cart-item-actions">
                        <div class="qty-control">
                            <button @click="updateQty(item.id, -1)">-</button>
                            <span class="qty-num">{{ item.quantity }}</span>
                            <button @click="updateQty(item.id, 1)">+</button>
                            <button
                                class="item-del-btn"
                                @click="deleteFromCart(item.id)"
                            >
                                x
                            </button>
                        </div>
                    </div>
                </div>
                <div class="cart-summary">
                    <div class="summary-line">
                        <span>Итого к оплате:</span>
                        <span class="total-price">{{ cartTotal }} ₽</span>
                    </div>
                </div>
            </div>
            <template #footer v-if="cart && cart.length > 0">
                <button
                    class="btn btn-primary full"
                    :disabled="isSubmittingOrder"
                    @click="placeOrder"
                >
                    {{
                        isSubmittingOrder
                            ? "Оформление..."
                            : `Оформить заказ на ${cartTotal} ₽`
                    }}
                </button>
            </template>
        </BaseModal>
    </div>
</template>

<style scoped>
.page-header {
    display: flex;
    justify-content: space-between;
}
.cart-btn {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 10px 16px;
    background: var(--accent);
    color: white;
    border: none;
    border-radius: 12px;
    font-weight: 600;
    cursor: pointer;
}
.cart-count {
    background: rgba(0, 0, 0, 0.2);
    padding: 2px 8px;
    border-radius: 6px;
    font-size: 12px;
}
.med-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
    gap: 20px;
}
.med-card {
    background: var(--surface);
    border: 1px solid var(--border);
    border-radius: 16px;
    padding: 20px;
    display: flex;
    flex-direction: column;
    justify-content: space-between;
}
.med-cat {
    font-size: 10px;
    text-transform: uppercase;
    color: var(--accent);
    font-weight: 700;
}
.med-title {
    font-size: 16px;
    margin: 8px 0;
    font-weight: 700;
}
.med-text {
    font-size: 13px;
    color: var(--text2);
    line-height: 1.4;
    margin-bottom: 16px;
}
.med-footer {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding-top: 12px;
    border-top: 1px solid var(--surface2);
}
.med-price {
    font-weight: 800;
    font-size: 18px;
}
.add-btn {
    padding: 8px 16px;
    background: var(--surface2);
    border: none;
    border-radius: 8px;
    cursor: pointer;
    font-weight: 600;
    transition: 0.2s;
}
.add-btn:hover {
    background: var(--accent);
    color: white;
}
.styled-search {
    width: 100%;
    padding: 12px 16px;
    border-radius: 10px;
    border: 1px solid var(--border);
    margin-bottom: 24px;
    font-size: 14px;
    outline: none;
}
.styled-search:focus {
    border-color: var(--accent);
}
.cart-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 12px 0;
    border-bottom: 1px solid var(--border);
}
.qty-control {
    display: flex;
    align-items: center;
    gap: 12px;
    background: var(--surface2);
    padding: 4px 12px;
    border-radius: 8px;
    border: 1px solid var(--border);
}
.qty-control button {
    border: none;
    background: none;
    cursor: pointer;
    font-weight: bold;
    color: var(--accent);
    font-size: 16px;
}
.qty-num {
    font-weight: 700;
    min-width: 20px;
    text-align: center;
}
.total-price {
    font-size: 20px;
    font-weight: 800;
    color: var(--accent);
}
.cart-summary {
    margin-top: 20px;
    padding-top: 16px;
    border-top: 2px dashed var(--border);
}
.summary-line {
    display: flex;
    justify-content: space-between;
    align-items: center;
}
.full {
    width: 100%;
}
.center {
    text-align: center;
}
.py-40 {
    padding: 40px 0;
}
.empty-cart {
    padding: 40px 0;
    text-align: center;
}
.item-del-btn {
    background: none;
    border: none;
    color: var(--red);
    cursor: pointer;
    font-size: 18px;
    padding: 4px;
}
</style>
