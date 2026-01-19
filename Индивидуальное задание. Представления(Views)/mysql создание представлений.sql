-- Создание Представлений

CREATE VIEW v_policies_full AS
SELECT
    p. policy_number,
    p. full_name AS policyholder_name,
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
    p.premium_amount - COALESCE(SUM(c.payout), 0) AS remaining_coverage
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
GROUP BY p.policy_number, p.full_name, p.premium_amount
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
    SUM(CASE WHEN p. contract_date >= DATE_SUB(CURDATE(), INTERVAL 30 DAY)
             THEN p.premium_amount ELSE 0 END) AS revenue_last_30_days
FROM employees e
LEFT JOIN policyholders p ON e.employee_id = p.employee_id
GROUP BY e.employee_id, e.full_name, e.position;

-- Представление для мониторинга истекающих полисов
CREATE VIEW v_expiring_soon AS
SELECT
    policy_number,
    full_name,
    end_date,
    DATEDIFF(end_date, CURDATE()) AS days_left
FROM policyholders
WHERE end_date BETWEEN CURDATE() AND DATE_ADD(CURDATE(), INTERVAL 30 DAY)
ORDER BY end_date;

-- Использование

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

-- Список всех представлений в БД
SELECT TABLE_NAME, VIEW_DEFINITION
FROM INFORMATION_SCHEMA.VIEWS
WHERE TABLE_SCHEMA = 'insurance';
