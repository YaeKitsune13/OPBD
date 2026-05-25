# Документация API — Ветеринарная клиника

## 1. Авторизация и профиль
| Маршрут | Метод | Описание | Принимает (JSON) | Вернёт (JSON) |
| :--- | :---: | :--- | :--- | :--- |
| `/api/auth/login` | POST | Вход в систему | `{email, password}` | `{token, userId, role, userName, lastName, phone}` |
| `/api/auth/register` | POST | Регистрация клиента | `{email, firstName, lastName, password, phone}` | `{message, userId}` |
| `/api/users/:id` | GET | Данные профиля + статистика | — | `{id, firstName, lastName, email, phone, role, petsCount, visitsCount}` |
| `/api/users/:id` | PUT | Обновление ФИО и телефона | `{firstName, lastName, phone}` | Обновленный объект пользователя |
| `/api/users/:id/password` | PUT | Смена пароля | `{current, next}` | `{status: "ok"}` или `{error: "текст"}` |

## 2. Главная страница (Dashboard)
| Маршрут | Метод | Описание | Принимает | Вернёт (JSON) |
| :--- | :---: | :--- | :--- | :--- |
| `/api/dashboard/:userId` | GET | Сводная информация для клиента | — | `{petsCount, appointmentsCount, visitsCount, nextAppointment: {date, time, petName}, recentAppointments: [], pets: []}` |

## 3. Питомцы (Pets)
| Маршрут | Метод | Описание | Принимает (JSON) | Вернёт (JSON) |
| :--- | :---: | :--- | :--- | :--- |
| `/api/pets/owner/:userId` | GET | Список всех питомцев владельца | — | `Array<Pet>` |
| `/api/pets?ownerId=:userId` | POST | Добавить нового питомца | `{name, species, breed, dob, weight, avatar}` | Созданный объект `Pet` |
| `/api/pets/:petId` | PUT | Редактировать данные питомца | `{name, species, breed, dob, weight, avatar}` | Обновленный объект `Pet` |
| `/api/pets/:petId` | DELETE | Удалить карточку питомца | — | `{success: true}` |

## 4. Запись на прием и визиты
| Маршрут | Метод | Описание | Параметры | Вернёт (JSON) |
| :--- | :---: | :--- | :--- | :--- |
| `/api/book/init/:userId` | GET | Данные для формы записи | — | `{pets: [], doctors: [], services: []}` |
| `/api/appointments/busy-slots` | GET | Занятые часы врача | `?doctor_id=1&date=2024-05-20` | `Array<string>` (напр. `["09:00", "14:00"]`) |
| `/api/appointments` | POST | Создать новую запись | `{pet_id, doctor_id, service_id, scheduled_at, comment}` | `{success: true, appointmentId}` |
| `/api/appointments/client/:userId` | GET | История визитов клиента | — | `Array<{id, petName, date, time, status, doctorName, protocol: {...}}>` |

## 5. Аптека, корзина и заказы
| Маршрут | Метод | Описание | Принимает (JSON) | Вернёт (JSON) |
| :--- | :---: | :--- | :--- | :--- |
| `/api/medications` | GET | Каталог медикаментов | — | `Array<Product>` |
| `/api/cart/:userId` | GET | Получить корзину юзера | — | `Array<CartItem>` |
| `/api/cart/:userId` | POST | Добавить товар в корзину | `{productId, quantity}` | Созданный `CartItem` |
| `/api/cart/:itemId` | PUT | Изменить количество товара | `{quantity}` | Обновленный `CartItem` |
| `/api/cart/:itemId` | DELETE | Удалить позицию из корзины | — | `{success: true}` |
| `/api/orders` | POST | Оформить заказ (очистка корзины) | `{userId, totalAmount, items: []}` | `{success: true, orderId}` |

## 6. Аналитика и статистика
| Маршрут | Метод | Описание | Принимает | Вернёт (JSON) |
| :--- | :---: | :--- | :--- | :--- |
| `/api/stats/pets/:userId` | GET | Список питомцев для графиков | — | `Array<{id, name}>` |
| `/api/stats/weight/:petId` | GET | История веса для Chart.js | — | `Array<{date, weight}>` (сорт. по дате) |

## 7. Модуль Врача
| Маршрут | Метод | Описание | Параметры / Body | Вернёт (JSON) |
| :--- | :---: | :--- | :--- | :--- |
| `/api/doctor/schedule` | GET | Расписание (все или по врачу) | `?doctor_id=1` | `Array<{id, time, petName, ownerName, status}>` |
| `/api/appointments/:id/status` | PATCH | Принять или отклонить запись | `{status: "confirmed" \| "rejected"}` | `{success: true}` |
| `/api/appointments/:id/complete` | PATCH | Завершить прием + Протокол | `{weight, diagnosis, treatment, medications}` | `{success: true}` |
| `/api/doctor/patients` | GET | Поиск по базе владельцев | `?search=Иванов` | `Array<{id, fullName, phone, email, petsCount}>` |
| `/api/doctor/patients/:id/history` | GET | Медкарта всех питомцев юзера | — | `Array<{petName, petIcon, visits: [...]}>` |
