<script setup>
// Если нужно будет переходить на другие страницы из кнопок внутри Dashboard
const emit = defineEmits(['navigate'])
</script>

<template>
  <div class="page-content">
    <div class="page-header">
      <div>
        <div class="page-title">Добрый день, Иван 👋</div>
        <div class="page-sub">Сегодня, 28 апреля 2026</div>
      </div>
      <button class="btn btn-primary" @click="emit('navigate', 'book')">✚ Записаться</button>
    </div>

    <div class="stats-grid">
      <div class="stat-card">
        <div class="stat-label">Питомцы</div>
        <div class="stat-value text-accent">3</div>
        <div class="stat-change">активных карточки</div>
      </div>
      <div class="stat-card">
        <div class="stat-label">Записи</div>
        <div class="stat-value">2</div>
        <div class="stat-change up">↑ ожидают подтверждения</div>
      </div>
      <div class="stat-card">
        <div class="stat-label">Визитов</div>
        <div class="stat-value">12</div>
        <div class="stat-change">за всё время</div>
      </div>
      <div class="stat-card">
        <div class="stat-label">Следующий приём</div>
        <div class="stat-value mono" style="font-size: 18px">02 мая</div>
        <div class="stat-change">Барсик · 10:30</div>
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
          <div class="schedule-item">
            <div class="schedule-time">02 мая</div>
            <div class="schedule-info">
              <div class="schedule-name">Барсик — Терапевт Кузнецов</div>
              <div class="schedule-sub">10:30 · Плановый осмотр</div>
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
          <div class="row">
            <div class="pet-avatar-mini">🐱</div>
            <div class="pet-info-mini">
              <div class="pet-name-mini">Барсик</div>
              <div class="pet-meta-mini">Кот · Британец · 4.2 кг</div>
            </div>
            <span class="tag">5 лет</span>
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
