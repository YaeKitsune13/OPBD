<script setup>
import { ref, computed, reactive } from 'vue'
import BaseModal from '../../ui/BaseModal.vue'

// Состояние
const searchQuery = ref('')
const isAddModalOpen = ref(false)

// Данные таблицы
const services = ref([
  { id: '001', name: 'Первичный осмотр', desc: 'Общий клинический осмотр животного', price: 850 },
  { id: '002', name: 'Вакцинация', desc: 'Введение вакцины, комплексная защита', price: 850 },
  {
    id: '003',
    name: 'УЗИ брюшной полости',
    desc: 'Ультразвуковая диагностика органов',
    price: 1500,
  },
])

// Поля новой услуги
const newService = reactive({
  name: '',
  desc: '',
  price: null,
})

// Фильтрация для поиска
const filteredServices = computed(() => {
  const query = searchQuery.value.toLowerCase()
  return services.value.filter((s) => s.name.toLowerCase().includes(query) || s.id.includes(query))
})

// Логика сохранения
function saveService() {
  if (!newService.name || !newService.price) {
    alert('Пожалуйста, заполните название и цену')
    return
  }

  const id = String(services.value.length + 1).padStart(3, '0')

  services.value.push({
    id,
    name: newService.name,
    desc: newService.desc,
    price: newService.price,
  })

  // Закрываем и сбрасываем форму
  isAddModalOpen.value = false
  newService.name = ''
  newService.desc = ''
  newService.price = null
}

function deleteService(id) {
  if (confirm('Удалить эту услугу из справочника?')) {
    services.value = services.value.filter((s) => s.id !== id)
  }
}
</script>

<template>
  <div class="page">
    <div class="page-header">
      <div>
        <div class="page-title">Справочник услуг</div>
        <div class="page-sub">Управление каталогом процедур и ценами</div>
      </div>
      <button class="btn btn-primary" @click="isAddModalOpen = true">✚ Добавить услугу</button>
    </div>

    <div class="card">
      <div class="card-header">
        <input
          type="text"
          v-model="searchQuery"
          placeholder="Поиск по названию или ID..."
          style="max-width: 300px"
        />
      </div>
      <div class="table-wrap">
        <table>
          <thead>
            <tr>
              <th style="width: 80px">#</th>
              <th>Наименование</th>
              <th>Стоимость</th>
              <th style="width: 120px">Действия</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="s in filteredServices" :key="s.id">
              <td class="mono text-muted">{{ s.id }}</td>
              <td class="td-main">
                {{ s.name }}
                <div class="service-desc">{{ s.desc }}</div>
              </td>
              <td class="mono">{{ s.price }} руб.</td>
              <td>
                <div class="row">
                  <button class="btn btn-ghost btn-sm" title="Редактировать">✏</button>
                  <button
                    class="btn btn-danger btn-sm"
                    @click="deleteService(s.id)"
                    title="Удалить"
                  >
                    🗑
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- МОДАЛКА ДОБАВЛЕНИЯ -->
    <BaseModal :show="isAddModalOpen" title="Добавить новую услугу" @close="isAddModalOpen = false">
      <div class="form-grid">
        <div class="form-group full">
          <label>Наименование услуги *</label>
          <input v-model="newService.name" type="text" placeholder="Напр. Чистка зубов" />
        </div>

        <div class="form-group full">
          <label>Описание</label>
          <textarea
            v-model="newService.desc"
            placeholder="Краткое описание процедуры..."
          ></textarea>
        </div>

        <div class="form-group full">
          <label>Стоимость (руб.) *</label>
          <input v-model="newService.price" type="number" placeholder="0" />
        </div>
      </div>

      <template #footer>
        <button class="btn btn-ghost" @click="isAddModalOpen = false">Отмена</button>
        <button class="btn btn-primary" @click="saveService">Сохранить</button>
      </template>
    </BaseModal>
  </div>
</template>

<style scoped>
.service-desc {
  font-size: 11px;
  font-weight: 400;
  color: var(--text3);
  margin-top: 2px;
}

textarea {
  min-height: 100px;
  resize: vertical;
}

/* Стилизация таблицы для соответствия скриншоту */
table {
  border-spacing: 0;
  width: 100%;
}

tr:last-child td {
  border-bottom: none;
}
</style>
