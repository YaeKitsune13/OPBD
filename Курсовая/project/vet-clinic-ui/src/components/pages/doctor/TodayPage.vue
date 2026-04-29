<script setup>
import { ref } from 'vue'

const emit = defineEmits(['navigate'])

const appointments = ref([
  {
    id: 1,
    time: '09:00',
    pet: '🐱 Барсик',
    owner: 'Иванов И.И.',
    type: 'Британец',
    reason: 'Плановый осмотр',
    status: 'confirmed',
  },
  {
    id: 2,
    time: '10:30',
    pet: '🐶 Мухтар',
    owner: 'Петров С.Р.',
    type: 'Немецкая овчарка',
    reason: 'Вакцинация',
    status: 'waiting',
  },
  {
    id: 3,
    time: '12:00',
    pet: '🐇 Пуговка',
    owner: 'Сидорова О.А.',
    type: 'Декоративный',
    reason: 'Осмотр',
    status: 'confirmed',
  },
])

function startVisit(app) {
  // В реальном приложении мы бы передали ID визита в Store
  emit('navigate', 'conduct')
}

function confirm(id) {
  const item = appointments.value.find((a) => a.id === id)
  if (item) item.status = 'confirmed'
}
</script>

<template>
  <div class="page">
    <div class="page-header">
      <div>
        <div class="page-title">Расписание на сегодня</div>
        <div class="page-sub">28 апреля 2026 · Кузнецов Андрей Владимирович</div>
      </div>
      <div class="row">
        <span class="badge badge-info">{{ appointments.length }} записи</span>
        <button class="btn btn-ghost btn-sm">🖨️ Печать</button>
      </div>
    </div>

    <div class="card">
      <div class="schedule-list" style="padding: 8px">
        <div
          v-for="app in appointments"
          :key="app.id"
          class="schedule-item"
          :style="app.status === 'waiting' ? 'opacity: 0.9' : ''"
        >
          <div class="schedule-time">{{ app.time }}</div>
          <div class="schedule-info">
            <div class="schedule-name">{{ app.pet }} — {{ app.owner }}</div>
            <div class="schedule-sub">{{ app.type }} · {{ app.reason }}</div>
          </div>

          <span v-if="app.status === 'confirmed'" class="badge badge-confirmed">Подтверждено</span>
          <span v-else class="badge badge-waiting">Ожидание</span>

          <div class="row" style="margin-left: 10px">
            <template v-if="app.status === 'confirmed'">
              <button class="btn btn-primary btn-sm" @click="startVisit(app)">Принять</button>
            </template>
            <template v-else>
              <button
                class="btn btn-sm btn-ghost"
                @click="confirm(app.id)"
                style="color: var(--accent)"
              >
                ✓
              </button>
              <button class="btn btn-sm btn-ghost" style="color: var(--red)">✕</button>
            </template>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
