# 🐾 Vet Clinic — Информационная система учёта посещений

Веб-приложение для автоматизации работы ветеринарной клиники.  
Стек: **Vue 3 + TypeScript** (frontend) · **NestJS + TypeORM** (backend)

---

## ER-диаграмма базы данных

```mermaid
erDiagram
    users {
        int id PK
        string email
        string password_hash
        enum role "vet | client"
        datetime created_at
    }

    owners {
        int id PK
        int user_id FK
        string full_name
        string phone
        string email
        string address
    }

    vets {
        int id PK
        int user_id FK
        string full_name
        string phone
        string specialization
    }

    pets {
        int id PK
        int owner_id FK
        string name
        string species
        string breed
        date birth_date
        float weight
        string photo_url
    }

    appointments {
        int id PK
        int pet_id FK
        int vet_id FK
        datetime datetime
        enum status "pending | confirmed | rejected"
        datetime created_at
    }

    appointment_status_history {
        int id PK
        int appointment_id FK
        int changed_by_user_id FK
        enum status "pending | confirmed | rejected"
        datetime changed_at
    }

    visits {
        int id PK
        int appointment_id FK
        int pet_id FK
        int vet_id FK
        datetime date_time
        text anamnesis
        text diagnosis
        text treatment
        float total_cost
    }

    services {
        int id PK
        string name
        text description
        float price
    }

    visit_services {
        int id PK
        int visit_id FK
        int service_id FK
    }

    medications {
        int id PK
        string name
        text description
        float price
    }

    visit_medications {
        int id PK
        int visit_id FK
        int medication_id FK
        int quantity
    }

    pet_weights {
        int id PK
        int pet_id FK
        float weight
        datetime measured_at
    }

    users ||--o| owners : "client"
    users ||--o| vets : "vet"
    owners ||--o{ pets : "owns"
    pets ||--o{ appointments : "has"
    vets ||--o{ appointments : "conducts"
    appointments ||--o{ appointment_status_history : "tracks"
    users ||--o{ appointment_status_history : "changes"
    appointments ||--o| visits : "becomes"
    pets ||--o{ visits : "attended"
    vets ||--o{ visits : "conducts"
    visits ||--o{ visit_services : "includes"
    services ||--o{ visit_services : "used in"
    visits ||--o{ visit_medications : "prescribes"
    medications ||--o{ visit_medications : "used in"
    pets ||--o{ pet_weights : "tracked"
```

---

## Роли

| Роль | Описание |
|---|---|
| `client` | Регистрируется сам, управляет питомцами, записывается на приём |
| `vet` | Ведёт приём, заполняет диагнозы, назначает лечение |

---

## Модули

| Модуль | Описание |
|---|---|
| `auth` | Регистрация, логин, JWT |
| `users` | Управление пользователями |
| `owners` | Владельцы животных |
| `vets` | Ветеринары |
| `pets` | Животные, фото, вес |
| `appointments` | Записи на приём, статусы |
| `visits` | Визиты, диагнозы, лечение |
| `services` | Справочник услуг |
| `medications` | Справочник лекарств |
| `stats` | Статистика и графики |

---

## Запуск

```bash
# Установка зависимостей
cd frontend && npm install
cd ../backend && npm install

# Запуск (из корня проекта)
startServer.bat
```

Swagger UI доступен по адресу: `http://localhost:3000/api/docs`