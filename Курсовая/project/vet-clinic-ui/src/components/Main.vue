<script setup>
import { ref, computed, defineAsyncComponent, watch, nextTick, onMounted } from 'vue'
import TopBar from './ui/TopBar.vue'
import SideBar from './ui/SideBar.vue'
import PrintZone from './ui/PrintZone.vue'
import ProfilePanel from './ui/ProfilePanel.vue'

// Принимаем данные из App.vue
const props = defineProps({
  userRole: String,
  userName: String
})
const emit = defineEmits(['logout'])

// --- СОСТОЯНИЕ ---
const role = ref(props.userRole)
const currentPage = ref('') // Инициализируем пустотой, настроим в onMounted
const mobileMenu = ref(false)
const isProfileOpen = ref(false)

// --- ДАННЫЕ ПОЛЬЗОВАТЕЛЯ (Динамические) ---
const currentUser = computed(() => {
  // Названия ролей для отображения в профиле
  const roleLabels = {
    admin: 'Администратор',
    doctor: 'Врач',
    client: 'Клиент'
  }

  // Генерация аватара из инициалов (Иван Иванов -> ИИ)
  const getInitials = (name) => {
    if (!name) return '??'
    const parts = name.split(' ')
    if (parts.length >= 2) return (parts[0][0] + parts[1][0]).toUpperCase()
    return name[0].toUpperCase()
  }

  return {
    name: props.userName || 'Пользователь',
    avatar: getInitials(props.userName),
    roleName: roleLabels[role.value] || 'Пользователь',
    // Эти данные можно будет тоже прокинуть через props, если добавишь в БД
    email: '---', 
    phone: '---',
  }
})

// Установка начальной страницы при загрузке
onMounted(() => {
  setInitialPage(props.userRole)
})

// Следим за изменением роли (например, если админ переключил режим)
watch(() => props.userRole, (newRole) => {
  if (newRole) {
    role.value = newRole
    setInitialPage(newRole)
  }
})

function setInitialPage(currentRole) {
  const firstPages = { 
    client: 'dashboard', 
    doctor: 'today', 
    admin: 'analytics' 
  }
  currentPage.value = firstPages[currentRole] || 'dashboard'
}

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
const ClientsPage = defineAsyncComponent(() => import('./pages/admin/ClientsPage.vue'))

// Метод смены роли (для тестов или админа)
function updateRole(newRole) {
  role.value = newRole
  setInitialPage(newRole)
}

// Печать
const printData = ref({ pet: null, visits: [], type: 'history' })
function triggerPrint(data) {
  printData.value = data
  nextTick(() => { window.print() })
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
        <ClientsPage v-if="currentPage === 'clients'"/>
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
