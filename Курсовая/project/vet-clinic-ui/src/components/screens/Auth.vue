<script setup lang="ts">
import { ref, reactive } from "vue";
import { useToast } from '../../utils/useToast'
const { showToast } = useToast()

const props = defineProps({
    modelValue: Boolean,
});

const emit = defineEmits(["update:modelValue", "login-success"]);

// --- СОСТОЯНИЕ ---
const activeTab = ref("login");

const loginData = reactive({
    email: "",
    pass: "",
});

const regData = reactive({
    lastName: "",
    firstName: "",
    phone: "",
    email: "",
    pass: "",
});

// --- ВСПОМОГАТЕЛЬНЫЕ ФУНКЦИИ ---

function clearData() {
    loginData.email = "";
    loginData.pass = "";

    regData.email = "";
    regData.firstName = "";
    regData.lastName = "";
    regData.pass = "";
    regData.phone = "";
}

/**
 * Приводит телефон к формату +79991234567
 */
function normalizePhone(phone: string): string {
    let digits = phone.replace(/\D/g, "");

    if (digits.startsWith("8")) {
        digits = "7" + digits.slice(1);
    }

    if (digits.length === 10) {
        digits = "7" + digits;
    }

    if (digits.length === 11) {
        return "+" + digits;
    }

    return phone;
}

// --- ОСНОВНАЯ ЛОГИКА ---

async function doLogin() {
    if (!loginData.email || !loginData.pass) {
        showToast("Введите почту и пароль", "info");
        return;
    }

    try {
        const response = await fetch("/api/auth/login", {
            method: "POST",
            headers: {
                Accept: "application/json",
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                email: loginData.email,
                password: loginData.pass,
            }),
        });

        const result = await response.json();

        if (result.error) {
            alert(result.error);
            return;
        }

        if (result.token) {
            // 1. Сохраняем сессию
            localStorage.setItem("token", result.token);
            localStorage.setItem(
                "user",
                JSON.stringify({
                    id: result.userId,
                    role: result.role,
                    name: result.userName,
                }),
            );

            // 2. Закрываем окно и уведомляем App.vue
            finalizeLogin(result.role);
            clearData();
        }
    } catch (e) {
        showToast("Ошибка соединения с сервером","error");
    }
}

async function doRegister() {
    // 1. Обязательные поля
    if (
        !regData.firstName ||
        !regData.lastName ||
        !regData.email ||
        !regData.pass
    ) {
        showToast("Заполните все обязательные поля (Имя, Фамилия, Email, Пароль)","info");
        return;
    }

    // 2. Валидация пароля (>= 8 символов)
    if (regData.pass.length < 8) {
        showToast("Безопасность: Пароль должен быть не менее 8 символов","info");
        return;
    }

    // 3. Валидация Email
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailRegex.test(regData.email)) {
        showToast("Введите корректный Email адрес","info");
        return;
    }

    // 4. Валидация телефона (через регулярку разрешаем ввод со скобками/тире)
    const phoneRegex =
        /^(\+7|7|8)?[\s\-]?\(?[489][0-9]{2}\)?[\s\-]?[0-9]{3}[\s\-]?[0-9]{2}[\s\-]?[0-9]{2}$/;
    if (regData.phone && !phoneRegex.test(regData.phone)) {
        showToast("Введите корректный номер телефона (например: +79001234567)","info");
        return;
    }

    try {
        const response = await fetch("/api/auth/register", {
            method: "POST",
            headers: {
                Accept: "application/json",
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                email: regData.email,
                firstName: regData.firstName,
                lastName: regData.lastName,
                password: regData.pass,
                phone: normalizePhone(regData.phone), // Отправляем чистые данные
            }),
        });

        const result = await response.json();

        if (result.error) {
            showToast(result.error, "error");
            return;
        }

        showToast("Регистрация успешна! Теперь вы можете войти.","success");
        activeTab.value = "login";
        clearData();
    } catch (e) {
        showToast("Ошибка при регистрации","error");
    }
}

function finalizeLogin(role: string) {
    emit("login-success", role);
    emit("update:modelValue", false);

    loginData.pass = "";
    regData.pass = "";
}
</script>

<template>
    <!-- Используем props.modelValue для контроля видимости из App.vue -->
    <div v-if="modelValue" class="auth-screen">
        <div class="auth-logo">
            <div class="auth-logo-icon">🐾</div>
            <span class="auth-logo-text">ВетКлиника</span>
        </div>

        <div class="auth-box">
            <div class="auth-tabs">
                <div
                    class="auth-tab"
                    :class="{ active: activeTab === 'login' }"
                    @click="activeTab = 'login'"
                >
                    Вход
                </div>
                <div
                    class="auth-tab"
                    :class="{ active: activeTab === 'reg' }"
                    @click="activeTab = 'reg'"
                >
                    Регистрация
                </div>
            </div>

            <!-- ФОРМА ВХОДА -->
            <div v-if="activeTab === 'login'" class="auth-form">
                <div class="form-group">
                    <label>Email</label>
                    <input
                        v-model="loginData.email"
                        type="email"
                        placeholder="example@mail.ru"
                    />
                </div>
                <div class="form-group">
                    <label>Пароль</label>
                    <input
                        v-model="loginData.pass"
                        type="password"
                        placeholder="Минимум 8 символов"
                        @keyup.enter="doLogin"
                    />
                </div>
                <button class="auth-submit" @click="doLogin">Войти</button>
                <div class="auth-footer-text">
                    Нет аккаунта?
                    <span class="auth-link" @click="activeTab = 'reg'"
                        >Зарегистрироваться</span
                    >
                </div>
            </div>

            <!-- ФОРМА РЕГИСТРАЦИИ -->
            <div v-else class="auth-form">
                <div class="form-row">
                    <div class="form-group">
                        <label>Фамилия</label>
                        <input
                            v-model="regData.lastName"
                            type="text"
                            placeholder="Иванов"
                        />
                    </div>
                    <div class="form-group">
                        <label>Имя</label>
                        <input
                            v-model="regData.firstName"
                            type="text"
                            placeholder="Иван"
                        />
                    </div>
                </div>
                <div class="form-group">
                    <label>Телефон</label>
                    <input
                        v-model="regData.phone"
                        type="text"
                        placeholder="+79001234567"
                    />
                </div>
                <div class="form-group">
                    <label>Email</label>
                    <input
                        v-model="regData.email"
                        type="email"
                        placeholder="example@mail.ru"
                    />
                </div>
                <div class="form-group">
                    <label>Пароль</label>
                    <input
                        v-model="regData.pass"
                        type="password"
                        placeholder="Минимум 8 символов"
                    />
                </div>
                <button class="auth-submit" @click="doRegister">
                    Зарегистрироваться
                </button>
                <div class="auth-footer-text">
                    Уже есть аккаунт?
                    <span class="auth-link" @click="activeTab = 'login'"
                        >Войти</span
                    >
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped>
.auth-footer-text {
    text-align: center;
    font-size: 12px;
    color: var(--text3);
}
.auth-link {
    color: var(--accent);
    cursor: pointer;
}
.form-row {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 10px;
}
/* ── AUTH SCREEN ── */
.auth-screen {
    position: fixed;
    inset: 0;
    background: var(--bg);
    z-index: 9000;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 20px;
}

.auth-screen.hidden {
    display: none;
}

.auth-logo {
    display: flex;
    align-items: center;
    gap: 10px;
    margin-bottom: 32px;
}

.auth-logo-icon {
    width: 40px;
    height: 40px;
    background: var(--accent);
    border-radius: 10px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 20px;
}

.auth-logo-text {
    font-size: 20px;
    font-weight: 700;
    letter-spacing: -0.5px;
}

.auth-box {
    width: 100%;
    max-width: 380px;
    background: var(--surface);
    border: 1px solid var(--border);
    border-radius: 10px;
    overflow: hidden;
}

.auth-tabs {
    display: flex;
    border-bottom: 1px solid var(--border);
}

.auth-tab {
    flex: 1;
    padding: 14px;
    text-align: center;
    font-size: 13px;
    font-weight: 600;
    color: var(--text2);
    cursor: pointer;
    transition: all 0.15s;
    border-bottom: 2px solid transparent;
    margin-bottom: -1px;
}

.auth-tab.active {
    color: var(--accent);
    border-bottom-color: var(--accent);
}
.auth-tab:hover:not(.active) {
    color: var(--text);
}

.auth-form {
    padding: 24px;
    display: flex;
    flex-direction: column;
    gap: 14px;
}

.auth-form .form-group {
    display: flex;
    flex-direction: column;
    gap: 5px;
}

.auth-form input {
    background: var(--surface2);
    border: 1px solid var(--border);
    border-radius: var(--radius);
    color: var(--text);
    font-family: var(--font);
    font-size: 13px;
    padding: 10px 12px;
    outline: none;
    transition: border-color 0.15s;
    width: 100%;
}

.auth-form input:focus {
    border-color: var(--accent);
}

.auth-form label {
    font-size: 11px;
    font-weight: 600;
    color: var(--text3);
    text-transform: uppercase;
    letter-spacing: 0.05em;
}

.auth-submit {
    width: 100%;
    padding: 11px;
    background: var(--accent);
    color: #000;
    font-family: var(--font);
    font-size: 14px;
    font-weight: 600;
    border: none;
    border-radius: var(--radius);
    cursor: pointer;
    transition: background 0.15s;
    margin-top: 4px;
}

.auth-submit:hover {
    background: #22c55e;
}

/* dev quick-login buttons at bottom */
.auth-devbar {
    position: fixed;
    bottom: 0;
    left: 0;
    right: 0;
    background: var(--surface);
    border-top: 1px solid var(--border);
    padding: 10px 20px;
    display: flex;
    align-items: center;
    gap: 10px;
    z-index: 9100;
    flex-wrap: wrap;
}

.auth-devbar-label {
    font-size: 11px;
    color: var(--text3);
    font-family: var(--mono);
    white-space: nowrap;
}
</style>
