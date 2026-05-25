<script setup lang="ts">
import { ref, reactive, watch } from "vue";
import { useToast } from "../../utils/useToast";
const { showToast } = useToast();

const props = defineProps({ isOpen: Boolean });
const emit = defineEmits(["close", "logout"]);

const activeSection = ref("view");

const loading = ref(false);
const stats = reactive({
    petsCount: 0,
    visitsCount: 0,
});
const userData = ref(JSON.parse(localStorage.getItem("user") || "{}"));
const editForm = reactive({
    lastName: "",
    firstName: "",
    phone: "",
});

const passwordForm = reactive({
    current: "",
    next: "",
    confirm: "",
});

const token = localStorage.getItem("token");

async function fetchProfileData() {
    if (!props.isOpen) return;
    loading.value = true;
    try {
        const res = await fetch(`/api/users/${userData.value.id}`, {
            headers: { Authorization: `Bearer ${token}` },
        });
        if (res.ok) {
            const freshData = await res.json();
            userData.value = freshData;
            editForm.lastName = freshData.lastName;
            editForm.firstName = freshData.firstName;
            editForm.phone = freshData.phone;

            stats.petsCount = freshData.petsCount || 0;
            stats.visitsCount = freshData.visitsCount || 0;

            localStorage.setItem("user", JSON.stringify(freshData));
        }
    } catch (e) {
        console.error("Ошибка обновления профиля", e);
    } finally {
        loading.value = false;
    }
}

watch(
    () => props.isOpen,
    (val) => {
        if (val) {
            setSection("view");
            fetchProfileData();
        }
    },
);

function setSection(section: string) {
    activeSection.value = section;
    passwordForm.current = "";
    passwordForm.next = "";
    passwordForm.confirm = "";
}

async function saveProfile() {
    loading.value = true;
    try {
        const res = await fetch(`/api/users/${userData.value.id}`, {
            method: "PUT",
            headers: {
                "Content-Type": "application/json",
                Authorization: `Bearer ${token}`,
            },
            body: JSON.stringify(editForm),
        });
        if (res.ok) {
            showToast("Данные успешно обновлены");
            await fetchProfileData();
            setSection("view");
        }
    } catch (e) {
        showToast("Ошибка при сохранении");
    } finally {
        loading.value = false;
    }
}

async function changePassword() {
    if (!passwordForm.current || !passwordForm.next)
        return showToast("Заполните поля", "info");
    if (passwordForm.next !== passwordForm.confirm)
        return showToast("Пароли не совпадают", "info");

    loading.value = true;
    try {
        const res = await fetch(`/api/users/${userData.value.id}/password`, {
            method: "PUT",
            headers: {
                "Content-Type": "application/json",
                Authorization: `Bearer ${token}`,
            },
            body: JSON.stringify({
                current: passwordForm.current,
                next: passwordForm.next,
            }),
        });
        if (res.ok) {
            showToast("Пароль успешно изменен", "success");
            setSection("view");
        } else {
            const err = await res.json();
            showToast(err.error || "Неверный текущий пароль", "info");
        }
    } catch (e) {
        showToast("Ошибка сервера", "error");
    } finally {
        loading.value = false;
    }
}
</script>

<template>
    <div class="profile-drawer" :class="{ open: isOpen }">
        <div class="drawer-header">
            <h3>
                {{
                    activeSection === "view"
                        ? "Личный кабинет"
                        : activeSection === "edit"
                          ? "Редактирование"
                          : "Безопасность"
                }}
            </h3>
            <button class="close-btn" @click="emit('close')">✕</button>
        </div>

        <div class="drawer-content">
            <div class="user-hero" :class="{ mini: activeSection !== 'view' }">
                <div class="avatar-box">
                    {{ userData.firstName?.[0] }}{{ userData.lastName?.[0] }}
                </div>
                <div class="hero-text">
                    <h2>{{ userData.lastName }} {{ userData.firstName }}</h2>
                    <span class="role-tag">{{
                        userData.role === "doctor" ? "Врач" : "Клиент"
                    }}</span>
                </div>
            </div>

            <div v-if="activeSection === 'view'" class="section-fade">
                <div class="info-list">
                    <div class="info-box">
                        <label>Электронная почта</label>
                        <p>{{ userData.email }}</p>
                    </div>
                    <div class="info-box">
                        <label>Номер телефона</label>
                        <p>{{ userData.phone || "Не указан" }}</p>
                    </div>
                </div>

                <div class="stats-row">
                    <div class="stat-item">
                        <span class="val">{{ stats.petsCount }}</span>
                        <span class="lab">Питомца</span>
                    </div>
                    <div class="stat-item">
                        <span class="val">{{ stats.visitsCount }}</span>
                        <span class="lab">Визитов</span>
                    </div>
                </div>

                <div class="menu-list">
                    <button class="menu-btn" @click="setSection('edit')">
                        Изменить данные
                    </button>
                    <button class="menu-btn" @click="setSection('password')">
                        Сменить пароль
                    </button>
                </div>

                <button class="logout-btn" @click="emit('logout')">
                    Выйти из системы
                </button>
            </div>

            <div v-else-if="activeSection === 'edit'" class="section-fade p-24">
                <div class="form-grid-mini">
                    <div class="form-group">
                        <label>Фамилия</label>
                        <input
                            v-model="editForm.lastName"
                            type="text"
                            :disabled="loading"
                        />
                    </div>
                    <div class="form-group">
                        <label>Имя</label>
                        <input
                            v-model="editForm.firstName"
                            type="text"
                            :disabled="loading"
                        />
                    </div>
                    <div class="form-group">
                        <label>Телефон</label>
                        <input
                            v-model="editForm.phone"
                            type="text"
                            placeholder="+7..."
                            :disabled="loading"
                        />
                    </div>
                </div>
                <div class="btn-footer">
                    <button
                        class="btn btn-primary full"
                        @click="saveProfile"
                        :disabled="loading"
                    >
                        {{ loading ? "Сохранение..." : "Сохранить изменения" }}
                    </button>
                    <button
                        class="btn btn-ghost full"
                        @click="setSection('view')"
                    >
                        Отмена
                    </button>
                </div>
            </div>

            <div
                v-else-if="activeSection === 'password'"
                class="section-fade p-24"
            >
                <div class="form-grid-mini">
                    <div class="form-group">
                        <label>Текущий пароль</label>
                        <input
                            v-model="passwordForm.current"
                            type="password"
                            :disabled="loading"
                        />
                    </div>
                    <div class="form-group">
                        <label>Новый пароль</label>
                        <input
                            v-model="passwordForm.next"
                            type="password"
                            :disabled="loading"
                        />
                    </div>
                    <div class="form-group">
                        <label>Подтвердите пароль</label>
                        <input
                            v-model="passwordForm.confirm"
                            type="password"
                            :disabled="loading"
                        />
                    </div>
                </div>
                <div class="btn-footer">
                    <button
                        class="btn btn-primary full"
                        @click="changePassword"
                        :disabled="loading"
                    >
                        {{ loading ? "Обновление..." : "Обновить пароль" }}
                    </button>
                    <button
                        class="btn btn-ghost full"
                        @click="setSection('view')"
                    >
                        Отмена
                    </button>
                </div>
            </div>
        </div>
    </div>

    <div v-if="isOpen" class="drawer-overlay" @click="emit('close')"></div>
</template>

<style scoped>
.profile-drawer {
    position: fixed;
    top: 0;
    right: 0;
    width: 360px;
    height: 100vh;
    background: var(--surface);
    z-index: 2000;
    transform: translateX(100%);
    transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    display: flex;
    flex-direction: column;
}
.profile-drawer.open {
    transform: translateX(0);
    box-shadow: -10px 0 30px rgba(0, 0, 0, 0.1);
}
.drawer-overlay {
    position: fixed;
    inset: 0;
    background: rgba(0, 0, 0, 0.4);
    z-index: 1999;
}
.drawer-header {
    padding: 20px 24px;
    display: flex;
    justify-content: space-between;
    align-items: center;
    border-bottom: 1px solid var(--border);
}
.drawer-header h3 {
    font-size: 15px;
    font-weight: 700;
    color: var(--text);
}
.close-btn {
    background: var(--surface2);
    border: none;
    width: 30px;
    height: 30px;
    border-radius: 50%;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
}
.drawer-content {
    flex: 1;
    overflow-y: auto;
    display: flex;
    flex-direction: column;
}
.user-hero {
    padding: 32px 24px;
    text-align: center;
    transition: all 0.3s ease;
}
.user-hero.mini {
    padding: 20px 24px;
    display: flex;
    align-items: center;
    text-align: left;
    gap: 16px;
    background: var(--surface2);
}
.avatar-box {
    width: 80px;
    height: 80px;
    background: var(--accent);
    color: white;
    border-radius: 24px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 24px;
    font-weight: 800;
    margin: 0 auto 16px;
    box-shadow: 0 8px 16px var(--accent-dim);
    transition: all 0.3s ease;
}
.mini .avatar-box {
    width: 50px;
    height: 50px;
    font-size: 16px;
    border-radius: 12px;
    margin: 0;
}
.mini h2 {
    font-size: 16px;
    margin: 0;
}
.user-hero h2 {
    font-size: 20px;
    font-weight: 700;
    margin-bottom: 4px;
}
.role-tag {
    font-size: 10px;
    font-weight: 700;
    color: var(--text3);
    text-transform: uppercase;
    background: var(--surface3);
    padding: 2px 10px;
    border-radius: 10px;
}
.section-fade {
    animation: fadeIn 0.3s ease;
    flex: 1;
    display: flex;
    flex-direction: column;
}
@keyframes fadeIn {
    from {
        opacity: 0;
        transform: translateY(5px);
    }
    to {
        opacity: 1;
        transform: translateY(0);
    }
}
.p-24 {
    padding: 24px;
}
.info-list {
    padding: 0 24px;
    margin-bottom: 24px;
}
.info-box {
    padding: 12px 0;
    border-bottom: 1px dotted var(--border);
}
.info-box label {
    font-size: 11px;
    color: var(--text3);
    font-weight: 600;
    display: block;
    margin-bottom: 4px;
}
.info-box p {
    font-size: 14px;
    font-weight: 600;
    margin: 0;
}
.stats-row {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 12px;
    padding: 0 24px;
    margin-bottom: 24px;
}
.stat-item {
    background: var(--surface2);
    padding: 16px;
    border-radius: 16px;
    text-align: center;
}
.stat-item .val {
    display: block;
    font-size: 20px;
    font-weight: 800;
    color: var(--accent);
}
.stat-item .lab {
    font-size: 11px;
    color: var(--text3);
    font-weight: 600;
}
.menu-list {
    padding: 0 24px;
    display: flex;
    flex-direction: column;
    gap: 8px;
}
.menu-btn {
    text-align: left;
    padding: 12px 16px;
    border: 1px solid var(--border);
    background: var(--surface);
    border-radius: 12px;
    font-size: 13px;
    font-weight: 600;
    cursor: pointer;
    transition: 0.2s;
}
.menu-btn:hover {
    border-color: var(--accent);
    background: var(--accent-dim);
}
.logout-btn {
    margin: auto 24px 24px;
    padding: 12px;
    color: var(--red);
    border: 1px solid var(--red);
    background: none;
    border-radius: 12px;
    font-weight: 700;
    cursor: pointer;
}
.logout-btn:hover {
    background: var(--red);
    color: white;
}
.form-grid-mini {
    display: flex;
    flex-direction: column;
    gap: 16px;
}
.form-group label {
    font-size: 11px;
    font-weight: 700;
    color: var(--text3);
    margin-bottom: 6px;
    display: block;
}
.form-group input {
    width: 100%;
    padding: 10px 14px;
    border: 1px solid var(--border);
    border-radius: 8px;
    font-size: 14px;
    background: var(--surface2);
    box-sizing: border-box;
}
.form-group input:focus {
    border-color: var(--accent);
    background: white;
    outline: none;
}
.btn-footer {
    margin-top: 32px;
    display: flex;
    flex-direction: column;
    gap: 10px;
}
.full {
    width: 100%;
}
</style>
