<script setup>
import { ref, computed } from 'vue'

const emit = defineEmits(['navigate'])

// Данные формы
const anamnesis = ref('Владелец отмечает снижение аппетита в течение 3 дней, чихание.')
const diagnosis = ref('ОРВИ, ринотрахеит кошек')
const currentWeight = ref(4.2)

// Список услуг и лекарств в текущем приеме
const assignments = ref([
  { id: 1, name: 'Первичный осмотр', type: 'Услуга', price: 850, qty: 1 },
  { id: 2, name: 'Амоксициллин 50мг', type: 'Медикамент', price: 320, qty: 1 },
])

// Автоматический расчет итога
const totalCost = computed(() => {
  return assignments.value.reduce((sum, item) => sum + item.price * item.qty, 0)
})

function removeItem(id) {
  assignments.value = assignments.value.filter((i) => i.id !== id)
}

function saveVisit() {
  alert('Медкарта сохранена!')
  emit('navigate', 'today')
}
</script>

<template>
  <div class="page">
    <div class="page-header">
      <div>
        <div class="page-title">Ведение приёма</div>
        <div class="page-sub">🐱 Барсик · Иванов И.И. · 28.04.2026 09:00</div>
      </div>
    </div>

    <div class="grid-2">
      <!-- Инфо о питомце -->
      <div class="card">
        <div class="card-header"><span class="card-title">Информация</span></div>
        <div class="card-body" style="display: flex; flex-direction: column; gap: 8px">
          <div class="row-between"><span class="text-muted">Кличка</span><span>Барсик</span></div>
          <div class="row-between">
            <span class="text-muted">Последний вес</span><span class="text-accent">4.20 кг</span>
          </div>
          <hr class="sep" />
          <div class="form-group">
            <label>Вес на приёме (кг)</label>
            <input type="number" step="0.1" v-model="currentWeight" />
          </div>
        </div>
      </div>

      <!-- Медкарта -->
      <div class="card">
        <div class="card-header"><span class="card-title">Медицинская карта</span></div>
        <div class="card-body" style="display: flex; flex-direction: column; gap: 12px">
          <div class="form-group">
            <label>Анамнез <span class="text-red">*</span></label>
            <textarea v-model="anamnesis" rows="3"></textarea>
          </div>
          <div class="form-group">
            <label>Диагноз <span class="text-red">*</span></label>
            <input type="text" v-model="diagnosis" />
          </div>
        </div>
      </div>
    </div>

    <!-- Назначения -->
    <div class="card">
      <div class="card-header">
        <span class="card-title">Назначения</span>
        <div class="row">
          <button class="btn btn-ghost btn-sm">+ Услуга</button>
          <button class="btn btn-ghost btn-sm">+ Медикамент</button>
        </div>
      </div>
      <div class="table-wrap">
        <table>
          <thead>
            <tr>
              <th>Наименование</th>
              <th>Тип</th>
              <th>Кол-во</th>
              <th>Стоимость</th>
              <th></th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="item in assignments" :key="item.id">
              <td class="td-main">{{ item.name }}</td>
              <td>
                <span
                  class="badge"
                  :class="item.type === 'Услуга' ? 'badge-info' : 'badge-waiting'"
                  >{{ item.type }}</span
                >
              </td>
              <td class="mono">{{ item.qty }}</td>
              <td class="mono">{{ item.price }} руб.</td>
              <td>
                <button class="btn btn-ghost btn-sm" @click="removeItem(item.id)">✕</button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      <div class="card-footer-cost">
        <span class="text-muted" style="font-size: 12px">Итоговая стоимость (автоматически)</span>
        <span class="total-price">{{ totalCost }} руб.</span>
      </div>
    </div>

    <div class="row" style="justify-content: flex-end; gap: 8px">
      <button class="btn btn-ghost" @click="emit('navigate', 'today')">Отмена</button>
      <button class="btn btn-primary" @click="saveVisit">Сохранить карту</button>
    </div>
  </div>
</template>

<style scoped>
.card-footer-cost {
  padding: 12px 14px;
  border-top: 1px solid var(--border);
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.total-price {
  font-size: 16px;
  font-weight: 700;
  font-family: var(--mono);
  color: var(--accent);
}
.sep {
  border: none;
  border-top: 1px solid var(--border);
  margin: 4px 0;
}
</style>
