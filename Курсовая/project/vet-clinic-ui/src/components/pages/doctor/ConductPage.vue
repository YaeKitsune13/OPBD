<script setup>
import { ref, computed } from 'vue'
import BaseModal from '../../ui/BaseModal.vue'
import PetCombobox from '../../ui/PetCombobox.vue'

const emit = defineEmits(['navigate'])

// База питомцев для выбора
const allPets = [
  { id: 1, name: 'Барсик', avatar: '🐱', breed: 'Британец', owner: 'Иванов И.И.' },
  { id: 2, name: 'Рыжик', avatar: '🐶', breed: 'Лабрадор', owner: 'Петров С.Р.' },
  { id: 3, name: 'Снежок', avatar: '🐇', breed: 'Вислоухий', owner: 'Сидорова А.А.' },
]

const selectedPet = ref(allPets[0]) // По умолчанию выбран Барсик

const anamnesis = ref('Владелец отмечает чихание и снижение аппетита.')
const diagnosis = ref('ОРВИ кошек')
const assignments = ref([{ id: 1, name: 'Первичный осмотр', type: 'Услуга', price: 850, qty: 1 }])

// Модалки
const isServiceModalOpen = ref(false)

const totalCost = computed(() => {
  return assignments.value.reduce((sum, item) => sum + item.price * item.qty, 0)
})

function removeItem(id) {
  assignments.value = assignments.value.filter((i) => i.id !== id)
}
</script>

<template>
  <div class="page">
    <div class="page-header conduct-header">
      <div class="title-with-select">
        <div class="page-title">Ведение приёма:</div>
        <!-- НАШ НОВЫЙ COMBOBOX -->
        <PetCombobox v-model="selectedPet" :pets="allPets" />
      </div>
    </div>

    <div class="grid-2">
      <!-- Медкарта -->
      <div class="card">
        <div class="card-header"><span class="card-title">Медкарта</span></div>
        <div class="card-body">
          <div class="form-group">
            <label>АНАМНЕЗ</label>
            <textarea v-model="anamnesis" placeholder="Жалобы пациента..."></textarea>
          </div>
          <div class="form-group mt-12">
            <label>ДИАГНОЗ</label>
            <input v-model="diagnosis" type="text" placeholder="Введите диагноз" />
          </div>
        </div>
      </div>

      <!-- Назначения -->
      <div class="card">
        <div class="card-header">
          <span class="card-title">Назначения</span>
          <button class="btn btn-ghost btn-sm" @click="isServiceModalOpen = true">
            + Добавить услугу
          </button>
        </div>
        <div class="table-wrap">
          <table>
            <thead>
              <tr>
                <th>НАИМЕНОВАНИЕ</th>
                <th>ЦЕНА</th>
                <th></th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="item in assignments" :key="item.id">
                <td class="td-main">{{ item.name }}</td>
                <td class="mono">{{ item.price }} р.</td>
                <td>
                  <button class="btn btn-ghost btn-sm" @click="removeItem(item.id)">✕</button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
        <div class="card-body row-between total-row">
          <span class="text-muted">Итого:</span>
          <span class="text-accent mono total-sum">{{ totalCost }} руб.</span>
        </div>
      </div>
    </div>

    <div class="row conduct-footer">
      <button class="btn btn-ghost" @click="emit('navigate', 'today')">Отмена</button>
      <button class="btn btn-primary" @click="emit('navigate', 'today')">Сохранить прием</button>
    </div>

    <!-- Модалка (как была раньше) -->
    <BaseModal
      :show="isServiceModalOpen"
      title="Добавить услугу"
      @close="isServiceModalOpen = false"
    >
      <div class="form-group">
        <label>Выберите услугу</label>
        <select class="mt-4">
          <option>УЗИ брюшной полости — 1500 руб.</option>
          <option>Вакцинация — 850 руб.</option>
        </select>
      </div>
      <template #footer>
        <button class="btn btn-primary" @click="isServiceModalOpen = false">Добавить</button>
      </template>
    </BaseModal>
  </div>
</template>

<style scoped>
.conduct-header {
  margin-bottom: 20px;
}
.title-with-select {
  display: flex;
  align-items: center;
  gap: 12px;
  width: 100%;
}

.total-row {
  border-top: 1px solid var(--border);
}
.total-sum {
  font-size: 18px;
  font-weight: 700;
}

.conduct-footer {
  justify-content: flex-end;
  gap: 10px;
  margin-top: 24px;
}

textarea {
  min-height: 120px;
}
</style>
