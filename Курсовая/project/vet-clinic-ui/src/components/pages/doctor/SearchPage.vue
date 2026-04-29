<script setup>
import { ref } from 'vue'

const searchQuery = ref('')
const searchPerformed = ref(false)

function doSearch() {
  if (searchQuery.value.length > 2) {
    searchPerformed.value = true
  }
}
</script>

<template>
  <div class="page">
    <div class="page-header">
      <div>
        <div class="page-title">Поиск пациента</div>
        <div class="page-sub">История болезней любого питомца по кличке или ФИО</div>
      </div>
    </div>

    <div class="card">
      <div class="card-body">
        <div class="row" style="gap: 8px">
          <input
            type="text"
            v-model="searchQuery"
            placeholder="Кличка питомца или ФИО владельца..."
            style="flex: 1"
            @keyup.enter="doSearch"
          />
          <button class="btn btn-primary" @click="doSearch">Найти</button>
        </div>
      </div>
    </div>

    <!-- Результаты поиска (показываем только если поиск нажат) -->
    <div v-if="searchPerformed" class="card">
      <div class="card-header">
        <span class="card-title">Результат: Барсик · Иванов И.И.</span>
      </div>
      <div style="padding: 14px; border-bottom: 1px solid var(--border)">
        <div class="row" style="gap: 12px">
          <div class="pet-avatar">🐱</div>
          <div>
            <div style="font-weight: 600">Барсик</div>
            <div class="text-muted" style="font-size: 12px">Кот · Британец · 4.2 кг · 5 лет</div>
            <div class="text-muted" style="font-size: 12px">
              Владелец: Иванов Иван Иванович · +79001234567
            </div>
          </div>
        </div>
      </div>
      <div style="padding: 14px">
        <div
          class="visit-item"
          style="border: 1px solid var(--border2); border-radius: 4px; padding: 10px"
        >
          <div class="row-between">
            <span class="mono" style="font-size: 12px">10.04.2026</span>
            <span class="badge badge-confirmed">Завершён</span>
          </div>
          <div style="margin-top: 5px; font-weight: 500">🔬 ОРВИ (ринотрахеит)</div>
        </div>
      </div>
    </div>

    <div v-else class="empty">
      <div class="empty-icon">🔍</div>
      <div>Введите данные для поиска пациента</div>
    </div>
  </div>
</template>
