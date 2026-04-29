<script setup>
import { ref } from 'vue'

const meds = ref([
  { id: 'M001', name: 'Амоксициллин 50мг', price: 320, expiry: '12.2027', status: 'ok' },
  { id: 'M003', name: 'Нафтизин', price: 120, expiry: '01.2026', status: 'expired' },
])
</script>

<template>
  <div class="page">
    <div class="page-header">
      <div>
        <div class="page-title">Медикаменты</div>
        <div class="page-sub">Склад и сроки годности</div>
      </div>
      <button class="btn btn-danger btn-sm">Списать просроченные</button>
    </div>

    <div class="card">
      <div class="table-wrap">
        <table>
          <thead>
            <tr>
              <th>#</th>
              <th>Название</th>
              <th>Цена</th>
              <th>Срок</th>
              <th></th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="m in meds"
              :key="m.id"
              :style="m.status === 'expired' ? 'background: var(--red-dim)' : ''"
            >
              <td class="mono text-muted">{{ m.id }}</td>
              <td class="td-main" :class="{ 'text-red': m.status === 'expired' }">{{ m.name }}</td>
              <td class="mono">{{ m.price }} руб.</td>
              <td class="mono" :class="{ 'text-red': m.status === 'expired' }">
                {{ m.expiry }} {{ m.status === 'expired' ? '⚠' : '' }}
              </td>
              <td>
                <button v-if="m.status === 'expired'" class="btn btn-danger btn-sm">Списать</button>
                <button v-else class="btn btn-ghost btn-sm">✏</button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>
