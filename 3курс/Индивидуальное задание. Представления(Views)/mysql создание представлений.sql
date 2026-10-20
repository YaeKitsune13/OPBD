-- =====================================================
-- ПРЕДСТАВЛЕНИЯ ДЛЯ БАЗЫ ДАННЫХ INSURANCE
-- =====================================================

-- 1. Представление "Кратковременное страхование"
-- (с продолжительностью не более 2-х недель)
CREATE OR REPLACE VIEW short_term_insurance AS
SELECT 
    p.policy_number,
    p.full_name AS policyholder_name,
    p.passport,
    it.name AS insurance_type,
    p.contract_date,
    p.end_date,
    DATEDIFF(p.end_date, p.contract_date) AS duration_days,
    p.premium_amount,
    p.policy_cost,
    e.full_name AS employee_name
FROM policyholders p
JOIN insurance_types it ON it.insurance_type_id = p.insurance_type_id
JOIN employees e ON e.employee_id = p.employee_id
WHERE DATEDIFF(p.end_date, p.contract_date) <= 14
ORDER BY p.contract_date DESC;

-- =====================================================

-- 2. Представление "Доходы по видам страхования"
CREATE OR REPLACE VIEW insurance_revenue_by_type AS
SELECT 
    it.name AS insurance_type,
    YEAR(p.contract_date) AS year,
    COUNT(DISTINCT p.policy_number) AS policy_count,
    SUM(p.policy_cost) AS total_policy_cost,
    COALESCE(SUM(c.payout), 0) AS total_payouts,
    SUM(p.policy_cost) - COALESCE(SUM(c.payout), 0) AS revenue
FROM insurance_types it
LEFT JOIN policyholders p ON p.insurance_type_id = it.insurance_type_id
LEFT JOIN claims c ON c.policy_number = p.policy_number
WHERE p.policy_number IS NOT NULL
GROUP BY it. insurance_type_id, it.name, YEAR(p.contract_date)
ORDER BY year DESC, insurance_type;

-- =====================================================

-- 3. Представление "Страховые выплаты"
CREATE OR REPLACE VIEW insurance_payouts AS
SELECT 
    p.policy_number,
    it.name AS insurance_type,
    p.premium_amount,
    p.policy_cost,
    COALESCE(SUM(c. payout), 0) AS total_payouts,
    p.policy_cost - COALESCE(SUM(c.payout), 0) AS difference
FROM policyholders p
JOIN insurance_types it ON it. insurance_type_id = p. insurance_type_id
LEFT JOIN claims c ON c.policy_number = p.policy_number
GROUP BY 
    p.policy_number, 
    it.name, 
    p.premium_amount, 
    p.policy_cost
ORDER BY p.policy_number;

-- =====================================================
-- ПРИМЕРЫ ИСПОЛЬЗОВАНИЯ ПРЕДСТАВЛЕНИЙ
-- =====================================================

-- Проверка представления "Кратковременное страхование"
SELECT * FROM short_term_insurance;

-- Проверка представления "Доходы по видам страхования"
SELECT * FROM insurance_revenue_by_type;

-- Проверка представления "Страховые выплаты"
SELECT * FROM insurance_payouts;

-- =====================================================
-- ПРИМЕРЫ INSERT
-- =====================================================

-- Пример 1: Добавление новых видов страхования
INSERT INTO insurance_types (insurance_type_id, name, description, annual_cost) 
VALUES
    (4, 'здоровье', 'дмс и медицинское страхование', 25000.00),
    (5, 'жизнь', 'страхование жизни', 30000.00),
    (6, 'спорт', 'краткосрочное спортивное страхование', 5000.00);

-- Пример 2: Добавление новых сотрудников
INSERT INTO employees (employee_id, full_name, passport, position) 
VALUES
    (1005, 'Васильева Мария Ивановна', '4004 567890', 'агент'),
    (1006, 'Козлов Дмитрий Андреевич', '4005 678901', 'старший агент'),
    (1007, 'Морозова Елена Викторовна', '4006 789012', 'менеджер');

-- Пример 3: Добавление кратковременных полисов (для демонстрации представления 1)
INSERT INTO policyholders (policy_number, passport, full_name, birth_date, insurance_type_id, employee_id, contract_date, end_date, premium_amount, policy_cost) 
VALUES
    -- Полис на 7 дней
    ('PL00000009', '4500 999999', 'Туристов Олег', DATE '1987-04-10', 2, 1001, 
     CURDATE(), DATE_ADD(CURDATE(), INTERVAL 7 DAY), 3000, 5000),
    
    -- Полис на 10 дней
    ('PL00000010', '4501 111111', 'Спортсменова Анна', DATE '1994-08-22', 6, 1005, 
     CURDATE(), DATE_ADD(CURDATE(), INTERVAL 10 DAY), 2500, 5000),
    
    -- Полис на 14 дней (ровно 2 недели)
    ('PL00000011', '4501 222222', 'Путешественников Сергей', DATE '1990-11-30', 2, 1002, 
     DATE_SUB(CURDATE(), INTERVAL 2 DAY), DATE_ADD(CURDATE(), INTERVAL 12 DAY), 7000, 15000),
    
    -- Обычный долгосрочный полис для сравнения
    ('PL00000012', '4501 333333', 'Надёжнов Пётр', DATE '1986-02-14', 4, 1006, 
     CURDATE(), DATE_ADD(CURDATE(), INTERVAL 1 YEAR), 20000, 25000);

-- Пример 4: Добавление страховых случаев
INSERT INTO claims (policy_number, description, event_date, payout) 
VALUES
    ('PL00000009', 'потеря багажа', DATE_ADD(CURDATE(), INTERVAL 4 DAY), 1500.00),
    ('PL00000010', 'спортивная травма', DATE_ADD(CURDATE(), INTERVAL 7 DAY), 2000.00),
    ('PL00000011', 'медицинская помощь за границей', DATE_ADD(CURDATE(), INTERVAL 5 DAY), 8000.00),
    ('PL00000012', 'плановое лечение', DATE_ADD(CURDATE(), INTERVAL 30 DAY), 15000.00),
    ('PL00000003', 'пожар на кухне', DATE_ADD(CURDATE(), INTERVAL 10 DAY), 18000.00);

-- =====================================================
-- ПРИМЕРЫ UPDATE
-- =====================================================

-- Пример 1: Повышение сотрудников в должности
UPDATE employees
SET position = 'старший агент'
WHERE employee_id = 1001 AND position = 'агент';

UPDATE employees
SET position = 'руководитель отдела'
WHERE employee_id = 1003 AND position = 'старший агент';

-- Пример 2: Корректировка стоимости полисов
UPDATE policyholders
SET policy_cost = 12000,
    premium_amount = 10000
WHERE policy_number = 'PL00000009';

-- Пример 3: Изменение годовой стоимости видов страхования
UPDATE insurance_types
SET annual_cost = 11000.00
WHERE insurance_type_id = 1;

UPDATE insurance_types
SET annual_cost = 16000.00,
    description = 'страхование выезжающих за рубеж и внутренние поездки'
WHERE insurance_type_id = 2;

-- Пример 4: Продление срока действия полиса
UPDATE policyholders
SET end_date = DATE_ADD(end_date, INTERVAL 30 DAY)
WHERE policy_number = 'PL00000010';

-- Пример 5: Корректировка суммы выплаты по страховому случаю
UPDATE claims
SET payout = 3500.00,
    description = 'дтп.  повреждение бампера (пересчёт после экспертизы)'
WHERE claim_id = 1;

-- Пример 6: Массовое обновление - увеличение всех премий на 5%
UPDATE policyholders
SET premium_amount = ROUND(premium_amount * 1.05, 0)
WHERE insurance_type_id = 2;

-- =====================================================
-- ПРИМЕРЫ DELETE
-- =====================================================

-- Пример 1: Удаление страхового случая (сначала удаляем зависимые записи)
DELETE FROM claims
WHERE claim_id = 5 AND policy_number = 'PL00000004';

-- Пример 2: Удаление полиса (сначала нужно удалить связанные страховые случаи)
-- Сначала удаляем все страховые случаи по полису
DELETE FROM claims
WHERE policy_number = 'PL00000010';

-- Затем удаляем сам полис
DELETE FROM policyholders
WHERE policy_number = 'PL00000010';

-- Пример 3: Удаление сотрудника, у которого нет полисов
DELETE FROM employees
WHERE employee_id = 1007 
  AND NOT EXISTS (
      SELECT 1 FROM policyholders WHERE employee_id = 1007
  );

-- Пример 4: Удаление вида страхования (сначала удалить все связанные данные)
-- Проверяем, есть ли полисы по этому виду
SELECT COUNT(*) FROM policyholders WHERE insurance_type_id = 6;

-- Если есть, сначала удаляем страховые случаи
DELETE FROM claims
WHERE policy_number IN (
    SELECT policy_number FROM policyholders WHERE insurance_type_id = 6
);

-- Затем удаляем полисы
DELETE FROM policyholders
WHERE insurance_type_id = 6;

-- И наконец удаляем вид страхования
DELETE FROM insurance_types
WHERE insurance_type_id = 6;

-- Пример 5: Массовое удаление - удаление старых страховых случаев (старше 2 лет)
DELETE FROM claims
WHERE event_date < DATE_SUB(CURDATE(), INTERVAL 2 YEAR);

-- Пример 6: Удаление полисов с истёкшим сроком действия (старше 1 года от текущей даты)
-- Сначала удаляем связанные страховые случаи
DELETE FROM claims
WHERE policy_number IN (
    SELECT policy_number FROM policyholders 
    WHERE end_date < DATE_SUB(CURDATE(), INTERVAL 1 YEAR)
);

-- Затем удаляем сами полисы
DELETE FROM policyholders
WHERE end_date < DATE_SUB(CURDATE(), INTERVAL 1 YEAR);

-- =====================================================
-- ПРОВЕРКА РЕЗУЛЬТАТОВ ПОСЛЕ ОПЕРАЦИЙ
-- =====================================================

-- Проверяем кратковременное страхование
SELECT 'Кратковременное страхование: ' AS title;
SELECT * FROM short_term_insurance;

-- Проверяем доходы по видам страхования
SELECT 'Доходы по видам страхования:' AS title;
SELECT * FROM insurance_revenue_by_type;

-- Проверяем страховые выплаты
SELECT 'Страховые выплаты:' AS title;
SELECT * FROM insurance_payouts;

-- Общая статистика
SELECT 'Общая статистика:' AS title;
SELECT 
    (SELECT COUNT(*) FROM insurance_types) AS total_insurance_types,
    (SELECT COUNT(*) FROM employees) AS total_employees,
    (SELECT COUNT(*) FROM policyholders) AS total_policies,
    (SELECT COUNT(*) FROM claims) AS total_claims;