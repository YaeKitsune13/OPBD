<script setup lang="ts">
import { onMounted, onUnmounted, ref, reactive, watch } from 'vue'
const emit = defineEmits(['navigate'])
const userRaw = localStorage.getItem('user')
const user = userRaw ? JSON.parse(userRaw) : null

const DashboardForm = reactive({
  petsCount: 0,
  appointmentsCount: 0,
  visits: 0,
  appointments: [],
  pets: [],
})

async function loadData() {
  const token = localStorage.getItem('token')
  const userRaw = localStorage.getItem('user')
  if (!userRaw) return;

  const userData = JSON.parse(userRaw)
  const userId = userData.user_id || userData.id

  try {
    const petsRes = await fetch(`/api/pets/owner/${userId}`, {
      headers: { "Authorization": `Bearer ${token}` }
    })
    if (petsRes.ok) DashboardForm.pets = await petsRes.json();

    const appointmentsRes = await fetch(`/api/appointments/owner/${userId}`, {
      headers: { "Authorization": `Bearer ${token}` }
    })
    if (appointmentsRes.ok) DashboardForm.appointments = await appointmentsRes.json();
  } catch (e) {
    showToast("Ошибка при загрузке данных", "error")
  }

  DashboardForm.petsCount = DashboardForm.pets.length
  DashboardForm.appointmentsCount = DashboardForm.appointments.length
}

function formatDate(dateStr) {
  if (!dateStr) return '—'
  const [day, month, year] = dateStr.split('.')
  const d = new Date(`${year}-${month}-${day}`)
  if (isNaN(d)) return dateStr
  return d.toLocaleDateString('ru-RU', { day: '2-digit', month: 'short' })
}

function calcAge(dob) {
  if (!dob) return '—'
  const [day, month, year] = dob.split('.')
  const birth = new Date(`${year}-${month}-${day}`)
  const diff = Date.now() - birth.getTime()
  return Math.floor(diff / (1000 * 60 * 60 * 24 * 365.25))
}

function getNextAppointment() {
  const now = new Date()
  const future = DashboardForm.appointments
    .filter(app => {
      const [day, month, year] = app.scheduledDate.split('.')
      return new Date(`${year}-${month}-${day}`) >= now
    })
    .sort((a, b) => {
      const [ad, am, ay] = a.scheduledDate.split('.')
      const [bd, bm, by] = b.scheduledDate.split('.')
      return new Date(`${ay}-${am}-${ad}`) - new Date(`${by}-${bm}-${bd}`)
    })
  return future[0] || null
}
loadData()
</script>

<template>
  <div class="page-content">
    <div class="page-header">
      <div>
        <div class="page-title">Добро пожаловать, {{ user.lastName }} {{ user.firstName }} 👋</div>
        <div class="page-sub">Сегодня, {{ new Date().toLocaleString() }}</div>
      </div>
      <button class="btn btn-primary" @click="emit('navigate', 'book')">✚ Записаться</button>
    </div>

    <div class="stats-grid">
      <div class="stat-card">
        <div class="stat-label">Питомцы</div>
        <div class="stat-value text-accent"> {{ DashboardForm.petsCount }}</div>
        <div class="stat-change">активных карточки</div>
      </div>
      <div class="stat-card">
        <div class="stat-label">Записи</div>
        <div class="stat-value">{{ DashboardForm.appointmentsCount }}</div>
        <div class="stat-change up">↑ ожидают подтверждения</div>
      </div>
      <div class="stat-card">
        <div class="stat-label">Визитов</div>
        <div class="stat-value">{{ DashboardForm.visits }}</div>
        <div class="stat-change">за всё время</div>
      </div>
      <div class="stat-card">
        <div class="stat-label">Следующий приём</div>
        <template v-if="getNextAppointment()">
          <div class="stat-value mono" style="font-size: 18px">
            {{ formatDate(getNextAppointment().scheduledDate) }}
          </div>
          <div class="stat-change">
            {{ getNextAppointment().petLabel }} · {{ getNextAppointment().scheduledTime }}
          </div>
        </template>
        <template v-else>
          <div class="stat-value mono" style="font-size: 18px">—</div>
          <div class="stat-change">нет записей</div>
        </template>
      </div>
    </div>

    <div class="grid-2">
      <!-- Последние записи -->
      <div class="card">
        <div class="card-header">
          <span class="card-title">Последние записи</span>
          <button class="btn btn-ghost btn-sm" @click="emit('navigate', 'appointments')">
            Все →
          </button>
        </div>
        <div class="schedule-list">
          <div class="schedule-item" v-for="app in DashboardForm.appointments" :key="app.id">
            <div class="schedule-time">{{ formatDate(app.scheduledDate) }}</div>
            <div class="schedule-info">
              <div class="schedule-name">{{ app.petLabel }} — {{ app.doctorName }}</div>
              <div class="schedule-sub">{{ app.scheduledTime }} · {{ app.specialty }}</div>
            </div>
            <span class="badge badge-waiting">Ожидание</span>
          </div>
        </div>
      </div>

      <!-- Питомцы кратко -->
      <div class="card">
        <div class="card-header">
          <span class="card-title">Мои питомцы</span>
          <button class="btn btn-ghost btn-sm" @click="emit('navigate', 'pets')">Все →</button>
        </div>
        <div class="card-body pet-list-mini">
          <div class="row" v-for="pet in DashboardForm.pets" :key="pet.petId">
            <div class="pet-avatar-mini">{{ pet.avatar || '🐾' }}</div>
            <div class="pet-info-mini">
              <div class="pet-name-mini">{{ pet.name }}</div>
              <div class="pet-meta-mini">{{ pet.species }} · {{ pet.breed }} · {{ pet.weight }} кг</div>
            </div>
            <span class="tag">{{ calcAge(pet.dob) }} лет</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* Сюда копируем стили .stats-grid, .stat-card, .schedule-list и т.д. */
.pet-list-mini {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.pet-avatar-mini {
  width: 36px;
  height: 36px;
  font-size: 18px;
  background: var(--surface3);
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.pet-name-mini {
  font-size: 13px;
  font-weight: 600;
}

.pet-meta-mini {
  font-size: 12px;
  color: var(--text2);
}
</style>
