<script setup>
import { ref, defineAsyncComponent, watch, nextTick } from 'vue' // Добавили nextTick сюда
import TopBar from './ui/TopBar.vue'
import SideBar from './ui/SideBar.vue'
import PrintZone from './ui/PrintZone.vue'

const props = defineProps(['initialRole'])
const emit = defineEmits(['logout'])

// --- СОСТОЯНИЕ ---
const role = ref(props.initialRole)
const currentPage = ref('dashboard')
const mobileMenu = ref(false)

const printData = ref({
  pet: null,
  visits: [],
  type: 'history',
})

// Следим за изменением роли извне (если нужно)
watch(
  () => props.initialRole,
  (newRole) => {
    if (newRole) updateRole(newRole)
  },
)

// --- ЛЕНИВАЯ ЗАГРУЗКА СТРАНИЦ ---
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

// Метод смены роли и сброса страницы
function updateRole(newRole) {
  role.value = newRole
  const firstPages = {
    client: 'dashboard',
    doctor: 'today',
    admin: 'analytics',
  }
  currentPage.value = firstPages[newRole]
}

function triggerPrint(data) {
  printData.value = data // Записываем данные в реактивный объект для PrintZone
  nextTick(() => {
    window.print() // Запускаем системное окно печати
  })
}
</script>

<template>
  <TopBar
    :current-role="role"
    @update-role="updateRole"
    @toggle-sidebar="mobileMenu = !mobileMenu"
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
        <!-- Client Pages -->
        <DashboardPage v-if="currentPage === 'dashboard'" @navigate="(p) => (currentPage = p)" />
        <PetsPage v-if="currentPage === 'pets'" />
        <AppointmentsPage
          v-if="currentPage === 'appointments'"
          @navigate="(p) => (currentPage = p)"
        />
        <BookPage v-if="currentPage === 'book'" @navigate="(p) => (currentPage = p)" />
        <!-- <-- Добавлено -->
        <HistoryPage v-if="currentPage === 'history'" @print="triggerPrint" />
        <WeightPage v-if="currentPage === 'weight'" />

        <!-- Doctor Pages -->
        <TodayPage v-if="currentPage === 'today'" @navigate="(p) => (currentPage = p)" />
        <ConductPage v-if="currentPage === 'conduct'" />
        <SearchPage v-if="currentPage === 'search'" />

        <!-- Admin Pages -->
        <AnalyticsPage v-if="currentPage === 'analytics'" />
        <RevenuePage v-if="currentPage === 'revenue'" />
        <ServicesPage v-if="currentPage === 'services'" />
        <MedsPage v-if="currentPage === 'meds'" />
      </div>
    </main>
  </div>

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
  max-width: 1400px; /* Чтобы на очень широких мониторах не растягивалось слишком сильно */
  margin: 0 auto;
  width: 100%;
}

@media (max-width: 768px) {
  .page-container {
    padding: 16px;
  }
}

/* Анимация перехода между страницами */
.v-enter-active,
.v-leave-active {
  transition: opacity 0.2s ease;
}
.v-enter-from,
.v-leave-to {
  opacity: 0;
}
</style>
