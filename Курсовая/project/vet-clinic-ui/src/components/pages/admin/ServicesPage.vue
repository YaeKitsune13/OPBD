<script setup>
import { ref, computed } from 'vue'

const searchQuery = ref('')
const services = ref([
  { id: '001', name: 'Первичный осмотр', desc: 'Общий клинический осмотр', price: 850 },
  { id: '002', name: 'Вакцинация', desc: 'Введение вакцины, комплексная защита', price: 850 },
  { id: '003', name: 'УЗИ брюшной полости', desc: 'Ультразвуковая диагностика', price: 1500 },
])

const filteredServices = computed(() => {
  return services.value.filter((s) =>
    s.name.toLowerCase().includes(searchQuery.value.toLowerCase()),
  )
})
</script>

<template>
  <div class="page">
    <div class="page-header">
      <div>
        <div class="page-title">Справочник услуг</div>
        <div class="page-sub">Управление ценами</div>
      </div>
      <button class="btn btn-primary" @click="$emit('open-modal', 'add-service')">
        ✚ Добавить
      </button>
    </div>

    <div class="card">
      <div class="card-header">
        <input
          type="text"
          v-model="searchQuery"
          placeholder="Поиск по названию..."
          style="max-width: 300px"
        />
      </div>
      <div class="table-wrap">
        <table>
          <thead>
            <tr>
              <th>#</th>
              <th>Наименование</th>
              <th>Стоимость</th>
              <th>Действия</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="s in filteredServices" :key="s.id">
              <td class="mono text-muted">{{ s.id }}</td>
              <td class="td-main">
                {{ s.name }}
                <div style="font-size: 11px; font-weight: 400; color: var(--text3)">
                  {{ s.desc }}
                </div>
              </td>
              <td class="mono">{{ s.price }} руб.</td>
              <td>
                <div class="row">
                  <button class="btn btn-ghost btn-sm">✏</button>
                  <button class="btn btn-danger btn-sm">🗑</button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>
