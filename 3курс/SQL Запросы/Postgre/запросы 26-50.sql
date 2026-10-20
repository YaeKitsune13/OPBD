\c books_deals
-- Задание 26
-- Вывести список издательств (поле Publish) из таблицы Publishing_house,
-- в которых выпущены книги, названия которых (поле Title_book)
-- начинаются со слова 'Труды',
-- и город издания (поле City) — 'Новосибирск'.
select ph.publish, b.title_book, ph.city
from publishing_houses ph
join books b on ph.code_publish = b.code_publish
where b.title_book like 'А%' and ph.city = 'Москва';

-- Многотабличные запросы (JOIN)

-- Задание 27
-- Вывести список названий компаний-поставщиков (поле Name_company)
-- и названия книг (поле Title_book),
-- которые они поставили в период с 01.01.2002 по 31.12.2003
-- (условие по полю Date_order).
select d.name_company, b.title_book, p.date_order, p.amount
from deliveries d
join purshases p on p.code_delivery = d.code_delivery
join books b on b.code_book = p.code_book
where p.date_order between '2002-01-01' and '2018-12-31';

-- Задание 28
-- Вывести список авторов (поле Name_author),
-- книги которых были выпущены в издательстве 'Мир'
-- (условие по полю Publish).

select a.name, b.title_book, ph.publish
from authors a
join books b on a.code_author = b.code_author
join publishing_houses ph on b.code_publish = ph.code_publish
where ph.publish ilike 'Аграф';

-- Задание 29
-- Вывести список поставщиков (поле Name_company),
-- которые поставляют книги издательства 'Питер'
-- (условие по полю Publish).
select d.name_company, ph.publish, b.title_book
from deliveries d
join purshases p on p.code_delivery = d.code_delivery
join books b on b.code_book = p.code_book
join publishing_houses ph on b.code_publish = ph.code_publish
where ph.publish ilike 'Memories';

-- Задание 30
-- Вывести список авторов (поле Name_author)
-- и названия книг (поле Title_book),
-- которые были поставлены поставщиком 'ОАО Книготорг'
-- (условие по полю Name_company).
select a.name, b.title_book, d.name_company
from authors a
join books b on a.code_author = b.code_author
join purshases p on b.code_book = p.code_book
join deliveries d on p.code_delivery = d.code_delivery
where d.name_company ilike 'ИП «Примак Ю.П.»';

-- Вычисления

-- Задание 31
-- Вывести суммарную стоимость партии одноименных книг
-- (использовать поля Amount и Cost)
-- и название книги (поле Title_book) в каждой поставке.
select title_book,amount, cost, (p.amount * p.cost) as total_cost
from books b
join purshases p on b.code_book = p.code_book
order by total_cost desc;

-- Задание 32
-- Вывести стоимость одной печатной страницы каждой книги
-- (использовать поля Cost и Pages)
-- и названия соответствующих книг (поле Title_book).
select title_book, cost, pages, round(cast(cost as numeric)/ nullif(pages,0),2) as cost_per_page
from books
join purshases p on books.code_book = p.code_book
order by cost_per_page desc;

-- Задание 33
-- Вывести количество лет с момента рождения авторов
-- (использовать поле Birthday)
-- и имена соответствующих авторов (поле Name_author).
select name, date(birthday), date_part('year',age(current_date,birthday)) as age
from authors
order by age desc;

-- Вычисление итоговых значений с использованием агрегатных функций

-- Задание 34
-- Вывести общую сумму поставок книг (использовать поле Cost),
-- выполненных 'ЗАО Оптторг' (условие по полю Name_company).
select name_company, sum(cost * amount) as total_sum
from deliveries d
join purshases p on d.code_delivery = p.code_delivery
group by name_company
order by total_sum desc;

-- Задание 35
-- Вывести общее количество всех поставок
-- (использовать любое поле из таблицы Purchases),
-- выполненных в период с 01.01.2003 по 01.02.2003
-- (условие по полю Date_order).
select count(*) as total_orders
from publishing_houses
join books b on publishing_houses.code_publish = b.code_publish
join purshases p on b.code_book = p.code_book
where p.date_order between '2003-01-01' and '2018-02-01';

-- Задание 36
-- Вывести среднюю стоимость (использовать поле Cost)
-- и среднее количество экземпляров книг (использовать поле Amount)
-- в одной поставке, где автором книги является 'Акунин'
-- (условие по полю Name_author).
select title_book, avg(cost) as avg_cost, avg(amount) as avg_amount
from authors a
join books b on a.code_author = b.code_author
join purshases p on b.code_book = p.code_book
-- where a.name ilike 'Акунин'
group by title_book;

-- Задание 37
-- Вывести все сведения о поставке (все поля таблицы Purchases),
-- а также название книги (поле Title_book)
-- с минимальной общей стоимостью
-- (использовать поля Cost и Amount).
select p.*, b.title_book, (p.cost * p.amount) as total_cost
from purshases p
join books b on p.code_book = b.code_book
where (p.cost * p.amount) = (
    select min(cost * amount)
    from purshases
);

-- Задание 38
-- Вывести все сведения о поставке (все поля таблицы Purchases),
-- а также название книги (поле Title_book)
-- с максимальной общей стоимостью
-- (использовать поля Cost и Amount).
select p.*, b.title_book, (p.cost * p.amount) as total_cost
from purshases p
join books b on p.code_book = b.code_book
where (p.cost * p.amount) = (
    select max(cost * amount)
    from purshases
);

-- Изменение наименований полей

-- Задание 39
-- Вывести название книги (поле Title_book),
-- суммарную стоимость партии одноименных книг
-- (использовать поля Amount и Cost),
-- поместив в результат в поле с названием Itogo,
-- в поставках за период с 01.01.2002 по 01.06.2002
-- (условие по полю Date_order).
select b.title_book,
       (p.amount * p.cost) as Itogo, date(date_order)
from books b
join purshases p on b.code_book = p.code_book
where p.date_order between '2002-01-01' and '2017-06-01'
order by Itogo desc;

-- Задание 40
-- Вывести стоимость одной печатной страницы каждой книги
-- (использовать поля Cost и Pages),
-- поместив результат в поле с названием One_page,
-- и названия соответствующих книг (поле Title_book).
select b.title_book,
       round(cast(p.cost as numeric)/ nullif(b.pages,0),2) as One_page
from books b
join purshases p on b.code_book = p.code_book
order by One_page desc;

-- Задание 41
-- Вывести общую сумму поставок книг (использовать поле Cost)
-- и поместить результат в поле с названием Sum_cost,
-- выполненных 'ОАО Луч' (условие по полю Name_company).
select d.name_company,
       sum(p.cost * p.amount) as Sum_cost
from deliveries d
join purshases p on d.code_delivery = p.code_delivery
-- where d.name_company ilike 'ОАО Луч'
group by d.name_company
order by Sum_cost desc;

-- Оператор обработки данных Update

-- Задание 42
-- Изменить в таблице Books содержимое поля Pages на 300,
-- если код автора (поле Code_author) = 56
-- и название книги (поле Title_book) = 'Мемуары'.
update books
set pages = 300
where code_author = 1 and title_book ilike '%М%';

select *
from books
where code_author = 1
  and title_book ilike '%М%';

-- Задание 43
-- Изменить в таблице Deliveries содержимое поля Address
-- на 'нет сведений',
-- если значение поля является пустым.
select address from deliveries;

update deliveries
set address = null
where address ilike 'Белгород, Б.Хмельницкого пр., д.135 корп.1 эт.';

update deliveries
set address = 'нет сведений'
where address is null;

select address from deliveries;

update deliveries
set address = 'Белгород, Б.Хмельницкого пр., д.135 корп.1 эт.'
where address ilike 'нет сведений';

select address from deliveries;
-- Задание 44
-- Увеличить в таблице Purchases цену (поле Cost) на 20 процентов,
-- если заказы были оформлены в течение последнего месяца
-- (условие по полю Date_order).

select * from purshases
where date_order >= current_date - interval '82 month';

update purshases
set cost = round(cost * 1.2, 2)
where date_order >= current_date - interval '82 month';

select * from purshases
where date_order >= current_date - interval '82 month';

-- Оператор обработки данных Insert

-- Задание 45
-- Добавить в таблицу Purchases новую запись,
-- причем так, чтобы код покупки (поле Code_purchase)
-- был автоматически увеличен на единицу,
-- а в тип закупки (поле Type_purchase) внести значение 'опт'.
insert into purshases (code_delivery, code_book, date_order, amount, cost, type_purchase)
values (3, 5, '2024-06-01', 50, 250.00, 'опт');

select * from purshases
where code_delivery = 3 and code_book = 5 and date_order = '2024-06-01';

-- Задание 46
-- Добавить в таблицу Books новую запись,
-- причем вместо ключевого поля поставить код (поле Code_book),
-- автоматически увеличенный на единицу от максимального кода в таблице,
-- вместо названия книги (поле Title_book) написать
-- 'Наука. Техника. Инновации'.
insert into books (code_author, code_publish, title_book, pages)
values (2, 3, 'Наука. Техника. Инновации', 400);

select * from books
where title_book ilike 'Наука. Техника. Инновации';

-- Задание 47
-- Добавить в таблицу Publishing_house новую запись,
-- причем вместо ключевого поля поставить код (поле Code_publish),
-- автоматически увеличенный на единицу от максимального кода в таблице,
-- вместо названия города — 'Москва' (поле City),
-- вместо издательства — 'Наука' (поле Publish).
insert into publishing_houses (city, publish)
values ('Москва', 'Наука');

select * from publishing_houses
where city ilike 'Москва' and publish ilike 'Наука';

-- Оператор обработки данных Delete

-- Задание 48
-- Удалить из таблицы Purchases все записи,
-- у которых количество книг в заказе (поле Amount) = 0.
insert into purshases (code_delivery, code_book, date_order, amount, cost, type_purchase)
values (2, 4, '2024-05-15', 0, 150.00, 'розница');

select * from purshases
where amount = 0;

delete from purshases
where amount = 0;

select * from purshases
where amount = 0;

-- Задание 49
-- Удалить из таблицы Authors все записи,
-- у которых нет имени автора в поле Name_author.
insert into authors (name, birthday)
values (null, '1970-01-01');

select * from authors
where name is null;

delete from authors
where name is null;

select * from authors
where name is null;

-- Задание 50
-- Удалить из таблицы Deliveries все записи,
-- у которых не указан ИНН (поле INN пустое).
insert into deliveries (name_company, address, inn)
values ('ИП Иванов И.И.', 'г. Казань, ул. Ленина, д.10', null);

select * from deliveries
where inn is null;

delete from deliveries
where inn is null;

select * from deliveries
where inn is null;