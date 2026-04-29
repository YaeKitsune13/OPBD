<script setup>
import { ref } from 'vue'

const emit = defineEmits(['print'])

// Данные о питомце (нужны для формирования заголовка в печати)
const petInfo = {
  name: 'Барсик',
  type: 'Кот',
  breed: 'Шотландская вислоухая',
}

const visits = ref([
  {
    id: 1,
    date: '10.04.2026',
    time: '11:00',
    doctor: 'Кузнецов А.В.',
    diagnosis: '🔬 ОРВИ у кошек (ринотрахеит)',
    details: 'Назначено: Амоксициллин 250 мг × 1 уп., Осмотр первичный × 1',
    analysis: 'ОАК — лейкоциты 12,4 × 10⁹/л (↑), эритроциты 7,2 × 10¹²/л',
    recommendations: 'Ограничить контакт с другими животными на 14 дней',
    price: '1 850',
  },
  {
    id: 2,
    date: '15.02.2026',
    time: '09:30',
    doctor: 'Попова М.С.',
    diagnosis: '💉 Плановая вакцинация',
    details: 'Назначено: Вакцина Nobivac Tricat × 1, Осмотр × 1',
    analysis: null,
    recommendations: 'Следующая вакцинация: 15.02.2027. Вес: 4,1 кг',
    price: '850',
  },
])

// Отправляем все записи в родительский компонент Main.vue
function printAll() {
  emit('print', {
    pet: petInfo,
    visits: visits.value,
    type: 'history',
  })
}

// Отправляем только одну конкретную запись
function printSingle(visit) {
  emit('print', {
    pet: petInfo,
    visits: [visit],
    type: 'history',
  })
}
</script>

<template>
  <div class="page">
    <div class="page-header">
      <div>
        <div class="page-title">Журнал здоровья</div>
        <div class="page-sub">История визитов и диагнозов</div>
      </div>
      <div class="row">
        <!-- В будущем сюда можно добавить v-model для фильтрации по питомцу -->
        <select class="btn btn-ghost" style="width: auto">
          <option>🐱 Барсик</option>
          <option>🐶 Рыжик</option>
        </select>
        <button class="btn btn-ghost" @click="printAll">🖨️ Печать</button>
      </div>
    </div>

    <!-- Список визитов -->
    <div v-for="visit in visits" :key="visit.id" class="visit-item">
      <div class="visit-head">
        <div>
          <span class="visit-date">{{ visit.date }} · {{ visit.time }}</span>
          <span class="visit-doctor"> · {{ visit.doctor }}</span>
        </div>
        <div class="row">
          <span class="badge badge-confirmed">Завершён</span>
          <!-- Кнопка печати конкретного визита -->
          <button
            class="btn btn-ghost btn-sm"
            @click="printSingle(visit)"
            title="Распечатать выписку"
          >
            🖨️
          </button>
        </div>
      </div>
      <div class="visit-body">
        <div class="visit-diagnosis">{{ visit.diagnosis }}</div>
        <div class="visit-detail">{{ visit.details }}</div>
        <div v-if="visit.analysis" class="visit-detail mt-4">
          <b>Анализы:</b> {{ visit.analysis }}
        </div>
        <div v-if="visit.recommendations" class="visit-detail mt-4">
          <b>Рекомендации:</b> {{ visit.recommendations }}
        </div>
        <div class="mt-8 text-muted" style="font-size: 12px">
          Итого за приём: <b>{{ visit.price }} руб.</b>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* Твои стили верны, они соответствуют общему дизайну */
.visit-item {
  border: 1px solid var(--border);
  border-radius: var(--radius);
  overflow: hidden;
  margin-bottom: 12px;
  transition: border-color 0.2s;
}
.visit-item:hover {
  border-color: var(--border2);
}
.visit-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 11px 14px;
  background: var(--surface2);
}
.visit-date {
  font-size: 12px;
  font-family: var(--mono);
  color: var(--text2);
}
.visit-doctor {
  font-size: 12px;
  color: var(--text3);
}
.visit-body {
  padding: 14px;
}
.visit-diagnosis {
  font-size: 13px;
  font-weight: 600;
  margin-bottom: 8px;
  color: var(--text);
}
.visit-detail {
  font-size: 12px;
  color: var(--text2);
  line-height: 1.4;
}
</style>
