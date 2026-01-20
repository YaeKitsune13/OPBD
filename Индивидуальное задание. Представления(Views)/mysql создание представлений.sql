-- ============================================
-- СОЗДАНИЕ ПРЕДСТАВЛЕНИЙ
-- ============================================

CREATE VIEW v_policies_full AS
SELECT
    p.policy_number,
    p.full_name AS policyholder_name,
    p.passport,
    p.birth_date,
    TIMESTAMPDIFF(YEAR, p.birth_date, CURDATE()) AS age,
    it.name AS insurance_type,
    it.description AS insurance_description,
    e.full_name AS agent_name,
    e.position AS agent_position,
    p.contract_date,
    p.end_date,
    DATEDIFF(p.end_date, CURDATE()) AS days_until_expiration,
    p.premium_amount,
    p.policy_cost
FROM policyholders p
JOIN insurance_types it ON p.insurance_type_id = it.insurance_type_id
JOIN employees e ON p.employee_id = e.employee_id;

CREATE VIEW v_claims_summary AS
SELECT
    p.policy_number,
    p.full_name AS policyholder_name,
    it.name AS insurance_type,
    p.premium_amount,
    COUNT(c.claim_id) AS total_claims,
    COALESCE(SUM(c. payout), 0) AS total_payout,
    p. premium_amount - COALESCE(SUM(c.payout), 0) AS remaining_coverage
FROM policyholders p
LEFT JOIN claims c ON p.policy_number = c.policy_number
JOIN insurance_types it ON p.insurance_type_id = it.insurance_type_id
GROUP BY p.policy_number, p.full_name, it.name, p.premium_amount;

CREATE VIEW v_policy_violations AS
-- Полисы с превышением выплат
SELECT
    'Превышение выплат' AS violation_type,
    p.policy_number,
    p.full_name,
    p.premium_amount,
    SUM(c.payout) AS total_payout
FROM policyholders p
JOIN claims c ON p.policy_number = c.policy_number
GROUP BY p.policy_number, p. full_name, p.premium_amount
HAVING SUM(c.payout) > p.premium_amount

UNION ALL

-- Полисы с ранними выплатами (до 6 дней)
SELECT
    'Выплата раньше 6 дней' AS violation_type,
    p.policy_number,
    p.full_name,
    p.premium_amount,
    c.payout
FROM policyholders p
JOIN claims c ON p.policy_number = c.policy_number
WHERE c.event_date < DATE_ADD(p.contract_date, INTERVAL 6 DAY);

CREATE VIEW v_employee_performance AS
SELECT
    e.employee_id,
    e.full_name AS agent_name,
    e.position,
    COUNT(p.policy_number) AS total_policies,
    SUM(p.premium_amount) AS total_premiums,
    COUNT(CASE WHEN p.contract_date >= DATE_SUB(CURDATE(), INTERVAL 30 DAY)
               THEN 1 END) AS policies_last_30_days,
    SUM(CASE WHEN p.contract_date >= DATE_SUB(CURDATE(), INTERVAL 30 DAY)
             THEN p.premium_amount ELSE 0 END) AS revenue_last_30_days
FROM employees e
LEFT JOIN policyholders p ON e.employee_id = p.employee_id
GROUP BY e.employee_id, e.full_name, e.position;

CREATE VIEW v_expiring_soon AS
SELECT
    policy_number,
    full_name,
    end_date,
    DATEDIFF(end_date, CURDATE()) AS days_left
FROM policyholders
WHERE end_date BETWEEN CURDATE() AND DATE_ADD(CURDATE(), INTERVAL 30 DAY)
ORDER BY end_date;

-- ============================================
-- ПРОСТОЕ ОБНОВЛЯЕМОЕ ПРЕДСТАВЛЕНИЕ
-- ============================================

-- Представление только для сотрудников (простое, обновляемое)
CREATE VIEW v_employees_simple AS
SELECT
    employee_id,
    full_name,
    position,
    phone,
    email
FROM employees;

-- Представление только для держателей полисов (обновляемое)
CREATE VIEW v_policyholders_simple AS
SELECT
    policy_number,
    full_name,
    passport,
    birth_date,
    phone,
    email,
    contract_date,
    end_date,
    premium_amount,
    policy_cost,
    insurance_type_id,
    employee_id
FROM policyholders;

-- Представление для типов страхования
CREATE VIEW v_insurance_types_simple AS
SELECT
    insurance_type_id,
    name,
    description
FROM insurance_types;

-- ============================================
-- ИСПОЛЬЗОВАНИЕ:  SELECT
-- ============================================

-- Просто выбираем из представления
SELECT * FROM v_policies_full
WHERE days_until_expiration < 30
ORDER BY days_until_expiration;

-- Найти полисы с превышением выплат
SELECT * FROM v_claims_summary
WHERE remaining_coverage < 0;

-- Статистика по типам страхования
SELECT
    insurance_type,
    COUNT(*) AS policies_count,
    SUM(total_payout) AS total_paid
FROM v_claims_summary
GROUP BY insurance_type;

-- ============================================
-- ИСПОЛЬЗОВАНИЕ:  INSERT
-- ============================================

-- INSERT через простое представление (работает!)
INSERT INTO v_employees_simple (full_name, position, phone, email)
VALUES ('Петров Петр Петрович', 'Агент', '+7-900-123-45-67', 'petrov@insurance.com');

-- INSERT для типа страхования
INSERT INTO v_insurance_types_simple (name, description)
VALUES ('Страхование имущества', 'Защита от повреждения или утраты имущества');

-- INSERT для держателя полиса
INSERT INTO v_policyholders_simple (
    policy_number,
    full_name,
    passport,
    birth_date,
    phone,
    email,
    contract_date,
    end_date,
    premium_amount,
    policy_cost,
    insurance_type_id,
    employee_id
)
VALUES (
    'POL-2026-00100',
    'Сидоров Сидор Сидорович',
    '4567 123456',
    '1985-03-15',
    '+7-900-555-11-22',
    'sidorov@email.com',
    '2026-01-20',
    '2027-01-20',
    50000.00,
    55000.00,
    1,
    1
);

-- ПРИМЕЧАНИЕ: INSERT НЕ работает для сложных представлений с JOIN или GROUP BY
-- Пример НЕРАБОТАЮЩЕГО INSERT:
-- INSERT INTO v_policies_full (...) VALUES (...);  -- ОШИБКА!

-- ============================================
-- ИСПОЛЬЗОВАНИЕ: UPDATE
-- ============================================

-- UPDATE через простое представление
UPDATE v_employees_simple
SET phone = '+7-900-999-88-77'
WHERE employee_id = 1;

-- UPDATE для типа страхования
UPDATE v_insurance_types_simple
SET description = 'Комплексное медицинское страхование'
WHERE name = 'Медицинское страхование';

-- UPDATE для держателя полиса
UPDATE v_policyholders_simple
SET premium_amount = 60000.00,
    policy_cost = 66000.00
WHERE policy_number = 'POL-2026-00100';

-- UPDATE с условием
UPDATE v_policyholders_simple
SET end_date = DATE_ADD(end_date, INTERVAL 1 YEAR)
WHERE policy_number IN (
    SELECT policy_number FROM v_expiring_soon
    WHERE days_left < 10
);

-- ПРИМЕЧАНИЕ: UPDATE НЕ работает для представлений с агрегацией
-- Пример НЕРАБОТАЮЩЕГО UPDATE:
-- UPDATE v_claims_summary SET total_claims = 5 WHERE policy_number = 'POL-001';  -- ОШИБКА!

-- ============================================
-- ИСПОЛЬЗОВАНИЕ: DELETE
-- ============================================

-- DELETE через простое представление
DELETE FROM v_employees_simple
WHERE employee_id = 999;

-- DELETE для типа страхования
DELETE FROM v_insurance_types_simple
WHERE insurance_type_id = 10;

-- DELETE для держателя полиса
DELETE FROM v_policyholders_simple
WHERE policy_number = 'POL-2026-00100';

-- DELETE с условием
DELETE FROM v_policyholders_simple
WHERE end_date < DATE_SUB(CURDATE(), INTERVAL 5 YEAR)
  AND policy_number NOT IN (SELECT policy_number FROM claims);

-- DELETE устаревших полисов
DELETE FROM v_policyholders_simple
WHERE policy_number IN (
    SELECT policy_number
    FROM policyholders
    WHERE end_date < '2020-01-01'
);

-- ПРИМЕЧАНИЕ: DELETE НЕ работает для сложных представлений
-- Пример НЕРАБО��АЮЩЕГО DELETE:
-- DELETE FROM v_policies_full WHERE days_until_expiration < 0;  -- ОШИБКА!

-- ============================================
-- СОЗДАНИЕ ОБНОВЛЯЕМОГО ПРЕДСТАВЛЕНИЯ С ОГРАНИЧЕНИЯМИ
-- ============================================

-- WITH CHECK OPTION гарантирует, что вставленные/обновленные строки
-- будут удовлетворять условию представления

CREATE VIEW v_active_policies AS
SELECT
    policy_number,
    full_name,
    passport,
    contract_date,
    end_date,
    premium_amount,
    insurance_type_id,
    employee_id
FROM policyholders
WHERE end_date >= CURDATE()
WITH CHECK OPTION;

-- Попытка вставить истекший полис вызовет ошибку:
-- INSERT INTO v_active_policies (policy_number, full_name, end_date, ...)
-- VALUES ('POL-OLD', 'Test User', '2020-01-01', ...);  -- ОШИБКА!

-- Это сработает:
INSERT INTO v_active_policies (
    policy_number,
    full_name,
    passport,
    contract_date,
    end_date,
    premium_amount,
    insurance_type_id,
    employee_id
)
VALUES (
    'POL-2026-00101',
    'Новый Клиент',
    '1234 567890',
    CURDATE(),
    DATE_ADD(CURDATE(), INTERVAL 1 YEAR),
    75000.00,
    1,
    1
);

-- ============================================
-- ПРЕДСТАВЛЕНИЕ ДЛЯ CLAIMS (обновляемое)
-- ============================================

CREATE VIEW v_claims_simple AS
SELECT
    claim_id,
    policy_number,
    event_date,
    event_description,
    payout,
    claim_date
FROM claims;

-- INSERT нового страхового случая
INSERT INTO v_claims_simple (policy_number, event_date, event_description, payout, claim_date)
VALUES ('POL-001', '2026-01-15', 'ДТП на перекрестке', 25000.00, CURDATE());

-- UPDATE страхового случая
UPDATE v_claims_simple
SET payout = 30000.00
WHERE claim_id = 1;

-- DELETE страхового случая
DELETE FROM v_claims_simple
WHERE claim_id = 100;

-- ============================================
-- ПРОВЕРКА ОБНОВЛЯЕМОСТИ ПРЕДСТАВЛЕНИЯ
-- ============================================

-- Список всех представлений в БД
SELECT
    TABLE_NAME,
    IS_UPDATABLE,
    CHECK_OPTION
FROM INFORMATION_SCHEMA. VIEWS
WHERE TABLE_SCHEMA = 'insurance';

-- Детальная информация о представлении
SELECT
    TABLE_NAME,
    IS_UPDATABLE,
    CHECK_OPTION,
    DEFINER,
    SECURITY_TYPE
FROM INFORMATION_SCHEMA. VIEWS
WHERE TABLE_SCHEMA = 'insurance'
  AND TABLE_NAME = 'v_employees_simple';

-- ============================================
-- УДАЛЕНИЕ ПРЕДСТАВЛЕНИЙ
-- ============================================

-- Удалить одно представление
DROP VIEW IF EXISTS v_employees_simple;

-- Удалить несколько представлений
DROP VIEW IF EXISTS
    v_policies_full,
    v_claims_summary,
    v_policy_violations,
    v_employee_performance,
    v_expiring_soon,
    v_employees_simple,
    v_policyholders_simple,
    v_insurance_types_simple,
    v_active_policies,
    v_claims_simple;

-- ============================================
-- ПРАВИЛА ОБНОВЛЯЕМОСТИ ПРЕДСТАВЛЕНИЙ В MySQL
-- ============================================

/*
Представление ОБНОВЛЯЕМО (можно INSERT/UPDATE/DELETE), если оно:

✅ МОЖНО обновлять:
1. Содержит простой SELECT без JOIN
2. Не содержит DISTINCT, GROUP BY, HAVING
