<script setup lang="ts">
import { ref, reactive } from "vue";
import { useToast } from "../../utils/useToast";
const { showToast } = useToast();

const emit = defineEmits<{ "login-success": [role: string] }>();

const activeTab = ref<"login" | "reg">("login");
const loading = ref(false);

const loginData = reactive({ email: "", pass: "" });
const regData = reactive({
    lastName: "",
    firstName: "",
    phone: "",
    email: "",
    pass: "",
});

function clearForms() {
    loginData.email = "";
    loginData.pass = "";
    regData.lastName = "";
    regData.firstName = "";
    regData.phone = "";
    regData.email = "";
    regData.pass = "";
}

function normalizePhone(phone: string): string {
    let d = phone.replace(/\D/g, "");
    if (d.startsWith("8")) d = "7" + d.slice(1);
    if (d.length === 10) d = "7" + d;
    return d.length === 11 ? "+" + d : phone;
}

async function doLogin() {
    if (!loginData.email || !loginData.pass)
        return showToast("Введите почту и пароль", "info");
    loading.value = true;
    try {
        const res = await fetch("/api/auth/login", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({
                email: loginData.email,
                password: loginData.pass,
            }),
        });
        const data = await res.json();
        if (!res.ok) throw new Error(data.error || "Ошибка входа");

        localStorage.setItem("token", data.token);
        localStorage.setItem(
            "user",
            JSON.stringify({
                id: data.userId,
                role: data.role,
                firstName: data.userName,
                lastName: data.lastName,
                email: loginData.email,
                phone: data.phone || "",
            }),
        );
        clearForms();
        emit("login-success", data.role);
    } catch (err: any) {
        showToast(err.message, "error");
    } finally {
        loading.value = false;
    }
}

async function doRegister() {
    if (
        !regData.firstName ||
        !regData.lastName ||
        !regData.email ||
        !regData.pass
    )
        return showToast("Заполните все обязательные поля", "info");

    loading.value = true;
    try {
        const res = await fetch("/api/auth/register", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({
                email: regData.email,
                firstName: regData.firstName,
                lastName: regData.lastName,
                password: regData.pass,
                phone: normalizePhone(regData.phone),
            }),
        });
        const data = await res.json();
        if (!res.ok) throw new Error(data.error || "Ошибка регистрации");

        showToast("Регистрация успешна! Войдите.");
        activeTab.value = "login";
        clearForms();
    } catch (err: any) {
        showToast(err.message);
    } finally {
        loading.value = false;
    }
}
</script>

<template>
    <div class="auth-screen">
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

            <div v-if="activeTab === 'login'" class="auth-form">
                <div class="form-group">
                    <label>Email</label>
                    <input
                        v-model="loginData.email"
                        type="email"
                        placeholder="example@mail.ru"
                        :disabled="loading"
                    />
                </div>
                <div class="form-group">
                    <label>Пароль</label>
                    <input
                        v-model="loginData.pass"
                        type="password"
                        placeholder="Минимум 8 символов"
                        :disabled="loading"
                        @keyup.enter="doLogin"
                    />
                </div>
                <button
                    class="auth-submit"
                    @click="doLogin"
                    :disabled="loading"
                >
                    {{ loading ? "Вход..." : "Войти" }}
                </button>
                <p class="auth-footer-text">
                    Нет аккаунта?
                    <span class="auth-link" @click="activeTab = 'reg'"
                        >Зарегистрироваться</span
                    >
                </p>
            </div>

            <div v-else class="auth-form">
                <div class="form-row">
                    <div class="form-group">
                        <label>Фамилия</label>
                        <input
                            v-model="regData.lastName"
                            type="text"
                            placeholder="Иванов"
                            :disabled="loading"
                        />
                    </div>
                    <div class="form-group">
                        <label>Имя</label>
                        <input
                            v-model="regData.firstName"
                            type="text"
                            placeholder="Иван"
                            :disabled="loading"
                        />
                    </div>
                </div>
                <div class="form-group">
                    <label>Телефон</label>
                    <input
                        v-model="regData.phone"
                        type="text"
                        placeholder="+79001234567"
                        :disabled="loading"
                    />
                </div>
                <div class="form-group">
                    <label>Email</label>
                    <input
                        v-model="regData.email"
                        type="email"
                        placeholder="example@mail.ru"
                        :disabled="loading"
                    />
                </div>
                <div class="form-group">
                    <label>Пароль</label>
                    <input
                        v-model="regData.pass"
                        type="password"
                        placeholder="Минимум 8 символов"
                        :disabled="loading"
                        @keyup.enter="doRegister"
                    />
                </div>
                <button
                    class="auth-submit"
                    @click="doRegister"
                    :disabled="loading"
                >
                    {{ loading ? "Регистрация..." : "Зарегистрироваться" }}
                </button>
                <p class="auth-footer-text">
                    Уже есть аккаунт?
                    <span class="auth-link" @click="activeTab = 'login'"
                        >Войти</span
                    >
                </p>
            </div>
        </div>
    </div>
</template>

<style scoped>
.auth-screen {
    position: fixed;
    inset: 0;
    background: var(--bg);
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 20px;
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
    border-bottom: 2px solid transparent;
    margin-bottom: -1px;
    transition: all 0.15s;
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

.form-row {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 10px;
}

.auth-submit {
    width: 100%;
    padding: 11px;
    background: var(--accent);
    color: #000; /* Вернул черный цвет текста */
    font-family: var(--font);
    font-size: 14px;
    font-weight: 600;
    border: none;
    border-radius: var(--radius);
    cursor: pointer;
    transition: background 0.15s;
    margin-top: 4px;
}
.auth-submit:hover:not(:disabled) {
    background: #22c55e;
}
.auth-submit:disabled {
    opacity: 0.6;
    cursor: not-allowed;
}

.auth-footer-text {
    text-align: center;
    font-size: 12px;
    color: var(--text3);
}
.auth-link {
    color: var(--accent);
    cursor: pointer;
}

label {
    display: block;
    font-size: 13px;
    color: var(--text);
    margin-bottom: 4px;
}

input {
    width: 100%;
    padding: 10px;
    border: 1px solid var(--border);
    border-radius: var(--radius);
    background: var(--surface2);
    font-size: 14px;
}
input:focus {
    outline: none;
    border-color: var(--accent);
}
</style>
