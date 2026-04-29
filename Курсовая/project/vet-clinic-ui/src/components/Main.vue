<script setup>
import { ref, computed, defineAsyncComponent, watch, nextTick } from 'vue'
import TopBar from './ui/TopBar.vue'
import SideBar from './ui/SideBar.vue'
import PrintZone from './ui/PrintZone.vue'
import ProfilePanel from './ui/ProfilePanel.vue'

const props = defineProps(['initialRole'])
const emit = defineEmits(['logout'])

// --- СОСТОЯНИЕ ---
const role = ref(props.initialRole)
const currentPage = ref('dashboard')
const mobileMenu = ref(false)
const isProfileOpen = ref(false)

// Данные для печати
const printData = ref({
  pet: null,
  visits: [],
  type: 'history',
})

// --- ДАННЫЕ ПОЛЬЗОВАТЕЛЯ (Добавлено) ---
// Вычисляем данные профиля на основе текущей роли
const currentUser = computed(() => {
  const configs = {
    client: {
      name: 'Иванов Иван Иванович',
      avatar: 'ИВ',
      roleName: 'Клиент',
      email: 'ivanov@mail.ru',
      phone: '+7 (900) 123-45-67',
    },
    doctor: {
      name: 'Кузнецов Андрей Владимирович',
      avatar: 'КА',
      roleName: 'Ветеринарный врач',
      email: 'dr.kuznetsov@vet.ru',
      phone: '+7 (900) 777-88-99',
    },
    admin: {
      name: 'Администратор системы',
      avatar: 'АД',
      roleName: 'Администратор',
      email: 'admin@vet.ru',
      phone: '+7 (900) 000-00-01',
    },
  }
  return configs[role.value] || configs.client
})

// Следим за изменением роли извне
watch(
  () => props.initialRole,
  (newRole) => {
    if (newRole) updateRole(newRole)
  },
)

// --- ЛЕНИВАЯ ЗАГРУЗКА ---
const DashboardPage = defineAsyncComponent(() => import('./pages/client/DashboardPage.vue'))
const PetsPage = defineAsyncComponent(() => import('./pages/client/PetsPage.vue'))
const AppointmentsPage = defineAsyncComponent(() => import('./pages/client/AppointmentsPage.vue'))
const BookPage = defineAsyncComponent(() => import('./pages/client/BookPage.vue'))
const HistoryPage = defineAsyncComponent(() => import('./pages/client/HistoryPage.vue'))
const WeightPage = defineAsyncComponent(() => import('./pages/client/WeightPage.vue'))

const TodayPage = defineAsyncComponent(() => import('./pages/doctor/TodayPage.vue'))
const ConductPage = defineAsyncComponent(() => import('./pages/doctor/ConductPage.vue'))
const SearchPage = defineAsyncComponent(() => import('./pages/doctor/SearchPage.vue'))

const AnalyticsPage = defineAsyncComponent(() => import('./pages/admin/AnalyticsPage.vue'))
const RevenuePage = defineAsyncComponent(() => import('./pages/admin/RevenuePage.vue'))
const ServicesPage = defineAsyncComponent(() => import('./pages/admin/ServicesPage.vue'))
const MedsPage = defineAsyncComponent(() => import('./pages/admin/MedsPage.vue'))

// Метод смены роли
function updateRole(newRole) {
  role.value = newRole
  const firstPages = { client: 'dashboard', doctor: 'today', admin: 'analytics' }
  currentPage.value = firstPages[newRole]
}

function triggerPrint(data) {
  printData.value = data
  nextTick(() => {
    window.print()
  })
}
</script>

<template>
  <TopBar
    :current-role="role"
    :user-name="currentUser.name"
    :user-avatar="currentUser.avatar"
    @update-role="updateRole"
    @toggle-sidebar="mobileMenu = !mobileMenu"
    @open-profile="isProfileOpen = true"
    @logout="emit('logout')"
  />

  <div class="layout">
    <SideBar
      :current-role="role"
      :active-page="currentPage"
      :is-open="mobileMenu"
      @navigate="(page) => (currentPage = page)"
      @close="mobileMenu = false"
    />

    <main class="main" id="main">
      <div class="page-container">
        <!-- Client -->
        <DashboardPage v-if="currentPage === 'dashboard'" @navigate="(p) => (currentPage = p)" />
        <PetsPage v-if="currentPage === 'pets'" />
        <AppointmentsPage
          v-if="currentPage === 'appointments'"
          @navigate="(p) => (currentPage = p)"
        />
        <BookPage v-if="currentPage === 'book'" @navigate="(p) => (currentPage = p)" />
        <HistoryPage v-if="currentPage === 'history'" @print="triggerPrint" />
        <WeightPage v-if="currentPage === 'weight'" />

        <!-- Doctor -->
        <TodayPage v-if="currentPage === 'today'" @navigate="(p) => (currentPage = p)" />
        <ConductPage v-if="currentPage === 'conduct'" @navigate="(p) => (currentPage = p)" />
        <SearchPage v-if="currentPage === 'search'" />

        <!-- Admin -->
        <AnalyticsPage v-if="currentPage === 'analytics'" />
        <RevenuePage v-if="currentPage === 'revenue'" />
        <ServicesPage v-if="currentPage === 'services'" />
        <MedsPage v-if="currentPage === 'meds'" />
      </div>
    </main>
  </div>

  <!-- Профиль -->
  <ProfilePanel
    :is-open="isProfileOpen"
    :user="currentUser"
    @close="isProfileOpen = false"
    @logout="emit('logout')"
  />

  <!-- Печать -->
  <Teleport to="body">
    <PrintZone
      id="printRoot"
      :pet="printData.pet"
      :visits="printData.visits"
      :type="printData.type"
    />
  </Teleport>
</template>

<style scoped>
.layout {
  display: flex;
  height: calc(100vh - var(--topbar-h));
  overflow: hidden;
}
.main {
  flex: 1;
  overflow-y: auto;
  background: var(--bg);
  scroll-behavior: smooth;
}
.page-container {
  padding: 24px;
  max-width: 1400px;
  margin: 0 auto;
  width: 100%;
}
@media (max-width: 768px) {
  .page-container {
    padding: 16px;
  }
}
</style>
