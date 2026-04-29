<script setup>
import { ref, computed } from 'vue'

const emit = defineEmits(['navigate'])

// Логика выбора докторов (как была в JS, но на Vue)
const doctors = {
  therapist: ['Кузнецов А.В.', 'Ломов Д.Е.'],
  surgeon: ['Попова М.С.'],
  cardio: ['Ломов Д.Е.'],
  neuro: ['Крылова Н.А.'],
  derm: ['Крылова Н.А.', 'Кузнецов А.В.'],
}

const selectedSpec = ref('')
const selectedDoctor = ref('')

// Автоматически обновляемый список врачей в зависимости от специальности
const availableDoctors = computed(() => {
  return doctors[selectedSpec.value] || []
})

function sendRequest() {
  alert('Заявка отправлена! Ожидайте подтверждения.')
  emit('navigate', 'appointments') // После записи перекидываем на страницу со списком
}
</script>

<template>
  <div class="page">
    <div class="page-header">
      <div>
        <div class="page-title">Записаться на приём</div>
        <div class="page-sub">Выберите питомца, врача и удобное время</div>
      </div>
    </div>

    <div class="card">
      <div class="card-body">
        <div class="form-grid">
          <div class="form-group">
            <label>Питомец</label>
            <select>
              <option>🐱 Барсик</option>
              <option>🐶 Рыжик</option>
            </select>
          </div>

          <div class="form-group">
            <label>Специализация</label>
            <select v-model="selectedSpec" @change="selectedDoctor = ''">
              <option value="">— Выберите —</option>
              <option value="therapist">Терапевт</option>
              <option value="surgeon">Хирург</option>
              <option value="cardio">Кардиолог</option>
            </select>
          </div>

          <div class="form-group">
            <label>Врач</label>
            <select v-model="selectedDoctor" :disabled="!selectedSpec">
              <option value="">— Выберите врача —</option>
              <option v-for="doc in availableDoctors" :key="doc" :value="doc">{{ doc }}</option>
            </select>
          </div>

          <div class="form-group">
            <label>Дата приёма</label>
            <input type="date" />
          </div>

          <div class="form-group full">
            <label>Комментарий</label>
            <textarea placeholder="Опишите проблему..."></textarea>
          </div>
        </div>

        <div class="mt-12">
          <button class="btn btn-primary" @click="sendRequest">Отправить заявку</button>
        </div>
      </div>
    </div>
  </div>
</template>
