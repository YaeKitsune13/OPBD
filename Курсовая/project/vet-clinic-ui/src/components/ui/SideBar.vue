<script setup>
// Принимаем данные от родителя
const props = defineProps({
  currentRole: {
    type: String,
    default: 'client',
  },
  activePage: {
    type: String,
    default: 'dashboard',
  },
  isOpen: {
    // Состояние открытой шторки на мобилках
    type: Boolean,
    default: false,
  },
})

// Сообщаем родителю о действиях
const emit = defineEmits(['navigate', 'close'])

// Функция навигации
function goTo(pageId) {
  emit('navigate', pageId)
  emit('close') // Закрываем сайдбар после клика (важно для мобилок)
}
</script>

<template>
  <!-- Сайдбар: класс 'open' добавляется реактивно через пропс isOpen -->
  <aside class="sidebar" :class="{ open: isOpen }">
    <!-- Секция КЛИЕНТА -->
    <div v-if="currentRole === 'client'" id="nav-client">
      <div class="nav-section">
        <div class="nav-section-label">Главное</div>
        <div
          class="nav-item"
          :class="{ active: activePage === 'dashboard' }"
          @click="goTo('dashboard')"
        >
          <span class="nav-icon">⊞</span> Обзор
        </div>
        <div class="nav-item" :class="{ active: activePage === 'pets' }" @click="goTo('pets')">
          <span class="nav-icon">🐾</span> Мои питомцы
        </div>
      </div>
      <div class="nav-section">
        <div class="nav-section-label">Запись</div>
        <div
          class="nav-item"
          :class="{ active: activePage === 'appointments' }"
          @click="goTo('appointments')"
        >
          <span class="nav-icon">📅</span> Записи <span class="nav-badge">2</span>
        </div>
        <div class="nav-item" :class="{ active: activePage === 'book' }" @click="goTo('book')">
          <span class="nav-icon">✚</span> Записаться
        </div>
      </div>
      <div class="nav-section">
        <div class="nav-section-label">История</div>
        <div
          class="nav-item"
          :class="{ active: activePage === 'history' }"
          @click="goTo('history')"
        >
          <span class="nav-icon">📋</span> Журнал здоровья
        </div>
        <div class="nav-item" :class="{ active: activePage === 'weight' }" @click="goTo('weight')">
          <span class="nav-icon">📈</span> Динамика веса
        </div>
      </div>
    </div>

    <!-- Секция ВРАЧА -->
    <div v-else-if="currentRole === 'doctor'" id="nav-doctor">
      <div class="nav-section">
        <div class="nav-section-label">Приём</div>
        <div class="nav-item" :class="{ active: activePage === 'today' }" @click="goTo('today')">
          <span class="nav-icon">📋</span> Расписание
        </div>
        <div
          class="nav-item"
          :class="{ active: activePage === 'conduct' }"
          @click="goTo('conduct')"
        >
          <span class="nav-icon">✏️</span> Ведение приёма
        </div>
      </div>
      <div class="nav-section">
        <div class="nav-section-label">Пациенты</div>
        <div class="nav-item" :class="{ active: activePage === 'search' }" @click="goTo('search')">
          <span class="nav-icon">🔍</span> Поиск пациента
        </div>
      </div>
    </div>

    <!-- Секция АДМИНА -->
    <div v-else-if="currentRole === 'admin'" id="nav-admin">
      <div class="nav-section">
        <div class="nav-section-label">Аналитика</div>
        <div
          class="nav-item"
          :class="{ active: activePage === 'analytics' }"
          @click="goTo('analytics')"
        >
          <span class="nav-icon">📊</span> Сводка
        </div>
        <div
          class="nav-item"
          :class="{ active: activePage === 'revenue' }"
          @click="goTo('revenue')"
        >
          <span class="nav-icon">💰</span> Выручка
        </div>
      </div>
      <div class="nav-section">
        <div class="nav-section-label">Справочники</div>
        <div
          class="nav-item"
          :class="{ active: activePage === 'services' }"
          @click="goTo('services')"
        >
          <span class="nav-icon">🔧</span> Услуги
        </div>
        <!-- ДОБАВЛЕННЫЙ ПУНКТ -->
        <div class="nav-item" :class="{ active: activePage === 'meds' }" @click="goTo('meds')">
          <span class="nav-icon">💊</span> Медикаменты
        </div>
      </div>
    </div>
  </aside>

  <!-- Оверлей -->
  <div
    v-if="isOpen"
    class="sidebar-overlay"
    :class="{ active: isOpen }"
    @click="emit('close')"
  ></div>
</template>

<style scoped>
/* Стили остаются без изменений */
.sidebar {
  width: var(--sidebar-w);
  background: var(--surface);
  border-right: 1px solid var(--border);
  display: flex;
  flex-direction: column;
  overflow-y: auto;
  flex-shrink: 0;
  transition: transform 0.25s ease;
}

.nav-section {
  padding: 16px 0 4px;
}

.nav-section-label {
  padding: 0 16px 6px;
  font-size: 10px;
  font-weight: 600;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  color: var(--text3);
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 9px 16px;
  cursor: pointer;
  color: var(--text2);
  font-size: 13px;
  font-weight: 500;
  transition: all 0.12s;
  border-left: 2px solid transparent;
}

.nav-item:hover {
  background: var(--surface2);
  color: var(--text);
}

.nav-item.active {
  background: var(--accent-dim);
  color: var(--accent);
  border-left-color: var(--accent);
}

.nav-icon {
  font-size: 15px;
  width: 18px;
  text-align: center;
}

.nav-badge {
  margin-left: auto;
  background: var(--red);
  color: #fff;
  font-size: 10px;
  font-weight: 700;
  padding: 1px 5px;
  border-radius: 10px;
}

.sidebar-overlay {
  display: none;
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  z-index: 99;
}

@media (max-width: 768px) {
  .sidebar {
    position: fixed;
    left: 0;
    top: var(--topbar-h);
    bottom: 0;
    z-index: 100;
    transform: translateX(-100%);
  }
  .sidebar.open {
    transform: translateX(0);
  }
  .sidebar-overlay.active {
    display: block;
  }
}
</style>
