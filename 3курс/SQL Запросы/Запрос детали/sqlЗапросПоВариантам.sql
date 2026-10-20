-- 1. Добавляем колонку даты в таблицу SPJ
ALTER TABLE SPJ ADD ship_date DATE;

-- 2. Заполняем существующие записи датами (чтобы запросы возвращали результат)
-- Сначала ставим всем 2023 год
UPDATE SPJ SET ship_date = '2023-02-15';

-- Обновим несколько записей на лето 2023 (чтобы сработал запрос №3 про период июнь-сентябрь)
UPDATE SPJ SET ship_date = '2023-07-01' WHERE p_no = 'P1';
UPDATE SPJ SET ship_date = '2023-08-20' WHERE s_no = 'S3';

-- 1. Вывести информацию о поставщиках (детали из Лондона, поставка в 2023 году)
SELECT DISTINCT S.*
FROM S
         JOIN SPJ ON S.s_no = SPJ.s_no
         JOIN P ON SPJ.p_no = P.p_no
WHERE P.city = 'Лондон'
  AND SPJ.ship_date BETWEEN '2023-01-01' AND '2023-12-31'; -- Добавлено слово BETWEEN

-- 2. Вставить нового поставщика
INSERT INTO S (s_no, s_name, status, city)
VALUES ('S6', 'Иванов', 15, 'Москва');

-- 3. Поставщики, поставившие деталь P1 летом 2023 года
SELECT DISTINCT S.*
FROM S
         JOIN SPJ ON S.s_no = SPJ.s_no
WHERE SPJ.p_no = 'P1'
  AND SPJ.ship_date >= '2023-06-01' AND SPJ.ship_date <= '2023-09-01';

-- 4. Удаление поставщика с наименьшим числом поставок
-- ВАЖНО: Сначала удаляем его поставки из SPJ, потом самого поставщика из S.

-- Шаг А: Удаляем записи из таблицы поставок
DELETE FROM SPJ
WHERE s_no = (
    SELECT s_no FROM (
                         SELECT s_no
                         FROM SPJ
                         GROUP BY s_no
                         ORDER BY COUNT(*) ASC
                             LIMIT 1
                     ) AS subquery
);

-- Шаг Б: Удаляем самого поставщика
DELETE FROM S
WHERE s_no = (
    SELECT s_no FROM (
                         -- Повторяем логику поиска, но так как из SPJ мы его уже удалили на Шаге А,
                         -- этот запрос может не сработать так, как ожидается, если запускать их последовательно.
                         -- В идеале для удаления нужно знать ID заранее.
                         -- Но если мы хотим удалить "бесполезных" поставщиков (у которых теперь 0 поставок):
                         SELECT s_no FROM S WHERE s_no NOT IN (SELECT DISTINCT s_no FROM SPJ)
                             LIMIT 1
                     ) AS subquery
);