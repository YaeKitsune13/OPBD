<script setup>
import { ref, reactive, computed } from 'vue'
import BaseModal from '../../ui/BaseModal.vue'

// Данные склада
const meds = ref([
  {
    id: 'M001',
    name: 'Амоксициллин 50мг',
    desc: 'Антибиотик широкого спектра',
    price: 320,
    expiry: '2027-12',
    status: 'ok',
  },
  {
    id: 'M002',
    name: 'Отипакс',
    desc: 'Ушные капли для животных',
    price: 480,
    expiry: '2026-06',
    status: 'ok',
  },
  {
    id: 'M003',
    name: 'Нафтизин',
    desc: 'Сосудосуживающее средство',
    price: 120,
    expiry: '2026-01',
    status: 'expired',
  },
])

// Состояние UI
const searchQuery = ref('')
const isAddModalOpen = ref(false)

// Форма нового препарата
const newMed = reactive({
  name: '',
  desc: '',
  price: null,
  expiry: '',
})

// Поиск по складу
const filteredMeds = computed(() => {
  const query = searchQuery.value.toLowerCase()
  return meds.value.filter(
    (m) => m.name.toLowerCase().includes(query) || m.id.toLowerCase().includes(query),
  )
})

// Проверка на просрочку (для стиля)
function isExpired(expiryDate) {
  const now = new Date()
  const exp = new Date(expiryDate)
  return exp < now
}

// Сохранение
function saveMed() {
  if (!newMed.name || !newMed.price) return alert('Заполните обязательные поля')

  const id = 'M' + String(meds.value.length + 1).padStart(3, '0')

  meds.value.push({
    id,
    name: newMed.name,
    desc: newMed.desc,
    price: newMed.price,
    expiry: newMed.expiry,
    status: isExpired(newMed.expiry) ? 'expired' : 'ok',
  })

  // Закрываем и очищаем
  isAddModalOpen.value = false
  Object.assign(newMed, { name: '', desc: '', price: null, expiry: '' })
}

function deleteMed(id) {
  if (confirm('Удалить препарат из базы данных?')) {
    meds.value = meds.value.filter((m) => m.id !== id)
  }
}
</script>

<template>
  <div class="page">
    <div class="page-header">
      <div>
        <div class="page-title">Справочник медикаментов</div>
        <div class="page-sub">Управление препаратами и контроль сроков годности</div>
      </div>
      <button class="btn btn-primary" @click="isAddModalOpen = true">✚ Добавить медикамент</button>
    </div>

    <div class="card">
      <div class="card-header">
        <input
          type="text"
          v-model="searchQuery"
          placeholder="Поиск по названию или артикулу..."
          style="max-width: 300px"
        />
        <button class="btn btn-danger btn-sm" @click="alert('Списание просрочки...')">
          Списать просроченные
        </button>
      </div>
      <div class="table-wrap">
        <table>
          <thead>
            <tr>
              <th style="width: 80px">#</th>
              <th>Наименование</th>
              <th>Цена / ед.</th>
              <th>Срок годности</th>
              <th style="width: 120px">Действия</th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="m in filteredMeds"
              :key="m.id"
              :class="{ 'expired-row': m.status === 'expired' }"
            >
              <td class="mono text-muted">{{ m.id }}</td>
              <td class="td-main">
                <span :class="{ 'text-red': m.status === 'expired' }">{{ m.name }}</span>
                <div class="med-desc">{{ m.desc }}</div>
              </td>
              <td class="mono">{{ m.price }} руб.</td>
              <td class="mono" :class="{ 'text-red': m.status === 'expired' }">
                {{ m.expiry }}
                <span v-if="m.status === 'expired'" title="Срок годности истек"> ⚠️</span>
              </td>
              <td>
                <div class="row">
                  <button class="btn btn-ghost btn-sm">✏</button>
                  <button class="btn btn-danger btn-sm" @click="deleteMed(m.id)">🗑</button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- МОДАЛКА ИЗ ОРИГИНАЛЬНОГО UI -->
    <BaseModal :show="isAddModalOpen" title="Добавить медикамент" @close="isAddModalOpen = false">
      <div class="form-grid">
        <div class="form-group full">
          <label>Название препарата *</label>
          <input v-model="newMed.name" type="text" placeholder="Напр. Амоксициллин" />
        </div>

        <div class="form-group full">
          <label>Описание / Показания</label>
          <textarea v-model="newMed.desc" placeholder="Описание препарата..."></textarea>
        </div>

        <div class="form-group">
          <label>Цена за ед. (руб.) *</label>
          <input v-model="newMed.price" type="number" placeholder="0" />
        </div>

        <div class="form-group">
          <label>Срок годности</label>
          <!-- Используем type="month" как в оригинальном HTML -->
          <input v-model="newMed.expiry" type="month" />
        </div>
      </div>

      <template #footer>
        <button class="btn btn-ghost" @click="isAddModalOpen = false">Отмена</button>
        <button class="btn btn-primary" @click="saveMed">Сохранить</button>
      </template>
    </BaseModal>
  </div>
</template>

<style scoped>
.med-desc {
  font-size: 11px;
  color: var(--text3);
  margin-top: 2px;
}

.expired-row {
  background: var(--red-dim);
}

textarea {
  min-height: 80px;
}

input[type='month'] {
  color-scheme: dark; /* Делает календарь темным в хроме */
}
</style>
