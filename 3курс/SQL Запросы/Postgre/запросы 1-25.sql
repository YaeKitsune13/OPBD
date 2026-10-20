-- Сортировка

-- Задание 1
-- Выбрать все сведения о книгах из таблицы Books,
-- и отсортировать результат по коду книги (поле Code_book).
-- 1
select *
from books
order by code_book;

-- Задание 2
-- Выбрать из таблицы Books коды книг, названия и количество страниц
-- (поля Code_book, Title_book и Pages),
-- отсортировать результат по названиям книг (поле Title_book по возрастанию),
-- и по полю Pages (по убыванию).
-- 2
select b.code_book, b.title_book, b.pages
from books b
order by b.title_book asc, b.pages desc;

-- Задание 3
-- Выбрать из таблицы Deliveries список поставщиков
-- (поля Name_delivery, Phone и INN),
-- отсортировать результат по полю INN (по убыванию).
-- 3
select d.name_company, d.phone, d.inn
from deliveries d
order by d.inn desc;


-- Изменение порядка следования полей

-- Задание 4
-- Выбрать все поля из таблицы Deliveries таким образом,
-- чтобы в результате порядок столбцов был следующим:
-- Name_delivery, INN, Phone, Address, Code_delivery.
-- 4
select d.name_company, d.inn, d.phone, d.address, d.code_delivery
from deliveries d;


-- Задание 5
-- Выбрать все поля из таблицы Publishing_house таким образом,
-- чтобы в результате порядок столбцов был следующим:
-- Publish, City, Code_publish.
-- 5
select ph.publish, ph.city, ph.code_publish
from publishing_houses ph;

-- Выбор некоторых полей из двух таблиц

-- Задание 6
-- Выбрать из таблицы Books названия книг и количество страниц
-- (поля Title_book и Pages),
-- а из таблицы Authors выбрать имя соответствующего автора книги
-- (поле Name_author).
-- 6
select b.title_book, a.name, b.pages
from books b
join authors a on b.code_author = a.code_author;

-- Задание 7
-- Выбрать из таблицы Books названия книг и количество страниц
-- (поля Title_book и Pages),
-- а из таблицы Deliveries выбрать имя соответствующего поставщика книги
-- (поле Name_delivery).
-- 7
select d.name_company, b.title_book, b.pages
from books b
join deliveries d on b.code_book = d.code_delivery;

-- Задание 8
-- Выбрать из таблицы Books названия книг и количество страниц
-- (поля Title_book и Pages),
-- а из таблицы Publishing_house выбрать название соответствующего издательства
-- и место издания (поля Publish и City).
-- 8
select ph.city, ph.publish, b.title_book, b.pages
from books b
join publishing_houses ph on b.code_publish = ph.code_publish;

-- Условие неточного совпадения

-- Задание 9
-- Выбрать из справочника поставщиков (таблица Deliveries)
-- названия компаний, телефоны и ИНН
-- (поля Name_company, Phone и INN),
-- у которых название компании (поле Name_company) начинается с 'ОАО'.
-- 9
select d.name_company, d.phone, d.inn
from deliveries d
where d.name_company ilike 'оао%';

-- Задание 10
-- Выбрать из таблицы Books названия книг и количество страниц
-- (поля Title_book и Pages),
-- а из таблицы Authors выбрать имя соответствующего автора книги
-- (поле Name_author),
-- у которых название книги начинается со слова 'Мемуары'.
-- 10
select b.title_book, a.name, b.pages
from books b
join authors a on b.code_author = a.code_author
where b.title_book ilike 'мемуары%'

-- Условие неточного совпадения

-- Задание 11
-- Выбрать из таблицы Authors фамилии, имена, отчества авторов
-- (поле Name_author),
-- значения которых начинаются с 'Иванов'.
-- 11
select a.name
from authors a
where a.name ilike 'Иванов%';

-- Точное несовпадение значений одного из полей

-- Задание 12
-- Вывести список названий издательств (поле Publish)
-- из таблицы Publishing_house,
-- которые не находятся в городе 'Москва' (условие по полю City).
-- 12
select ph.publish, ph.city
from publishing_houses ph
where ph.city not ilike 'Москва'
order by ph.city desc;

-- Задание 13
-- Вывести список названий книг (поле Title_book)
-- из таблицы Books,
-- которые выпущены любыми издательствами, кроме издательства 'Питер-Софт'
-- (поле Publish из таблицы Publishing_house).
-- 13
select b.title_book, ph.publish
from books b
join publishing_houses ph on ph.code_publish = b.code_publish
where ph.publish not ilike 'Питер-Софт'
order by ph.publish desc;

-- Выбор записей по диапазону значений (Between)

-- Задание 14
-- Вывести фамилии, имена, отчества авторов (поле Name_author)
-- из таблицы Authors,
-- у которых дата рождения (поле Birthday) находится в диапазоне
-- 01.01.1840 – 01.06.1860.
-- 14
select a.name, a.birthday
from authors a
where a.birthday between '1840-01-01' and '1860.06.01'
order by a.birthday desc;

-- Задание 15
-- Вывести список названий книг (поле Title_book из таблицы Books)
-- и количество экземпляров (поле Amount из таблицы Purchases),
-- которые были закуплены в период с 12.03.2003 по 15.06.2003
-- (условие по полю Date_order из таблицы Purchases).
-- 15
select b.title_book,ph.date_order
from books b
join purshases ph on ph.code_book = b.code_book
where ph.date_order between '2003-03-12' and '2017-06-15';

-- Задание 16
-- Вывести список названий книг (поле Title_book)
-- и количество страниц (поле Pages) из таблицы Books,
-- у которых объем в страницах укладывается в диапазон 200 – 300
-- (условие по полю Pages).
-- 16
select b.title_book, b.pages
from books b
where b.pages > 200 and b.pages < 400
order by b.pages asc;

-- Задание 17
-- Вывести список фамилий, имен, отчеств авторов (поле Name_author)
-- из таблицы Authors,
-- у которых фамилия начинается на одну из букв диапазона 'В' – 'Г'
-- (условие по полю Name_author).
-- 17
select a.name
from authors a
where left(a.name,1) between 'В' and 'Г';

-- Выбор записей по диапазону значений (In)

-- Задание 18
-- Вывести список названий книг (поле Title_book из таблицы Books)
-- и количество (поле Amount из таблицы Purchases),
-- которые были поставлены поставщиками с кодами 3, 7, 9, 11
-- (условие по полю Code_delivery из таблицы Purchases).
-- 18
select b.title_book, p.amount, p.code_delivery
from books b
join purshases p on b.code_book = p.code_book
where p.code_delivery in (3,7,9,11)
order by p.code_delivery;
-- Задание 19
-- Вывести список названий книг (поле Title_book) из таблицы Books,
-- которые выпущены следующими издательствами:
-- 'Питер-Софт', 'Альфа', 'Наука'
-- (условие по полю Publish из таблицы Publishing_house).
-- 19
select b.title_book, ph.publish
from books b
join publishing_houses ph on ph.code_publish = b.code_publish
where ph.publish in ('Питер-Софт','Водолей','Владос');

-- Задание 20
-- Вывести список названий книг (поле Title_book) из таблицы Books,
-- которые написаны следующими авторами:
-- 'Толстой Л.Н.', 'Достоевский Ф.М.', 'Пушкин А.С.'
-- (условие по полю Name_author из таблицы Authors).
-- 20
select b.title_book, a.name
from books b
join authors a on b.code_author = a.code_author
where a.name in (' Иван Гончаров', 'Михаил Лермонтов', ' Владимир Высоцкий')

-- Выбор записей с использованием Like

-- Задание 21
-- Вывести список авторов (поле Name_author) из таблицы Authors,
-- которые начинаются на букву 'К'.
-- 21
select a.name
from authors a
where a.name ilike 'К%'

-- Задание 22
-- Вывести названия издательств (поле Publish)
-- из таблицы Publishing_house,
-- которые содержат в названии сочетание 'софт'.
-- 22
select ph.publish
from publishing_houses ph
where ph.publish ilike '%мир%'

-- Задание 23
-- Выбрать названия компаний (поле Name_company)
-- из таблицы Deliveries,
-- у которых значение оканчивается на 'ский'.
-- 23
select d.name_company
from deliveries d
where d.name_company ilike '%»';

-- Выбор записей по нескольким условиям

-- Задание 24
-- Выбрать коды поставщиков (поле Code_delivery),
-- даты заказов (поле Date_order) и названия книг (поле Title_book),
-- если количество книг (поле Amount) в заказе больше 100,
-- или цена (поле Cost) за книгу находится в диапазоне от 200 до 500.
-- 24
select b.code_book, p.amount, p.cost
from books b
join purshases p on p.code_book = b.code_book
join deliveries d on d.code_delivery = p.code_delivery
where p.amount > 100 and (p.cost between 200 and 500)
order by p.amount, p.cost;

-- Задание 25
-- Выбрать коды авторов (поле Code_author),
-- имена авторов (поле Name_author),
-- названия соответствующих книг (поле Title_book),
-- если код издательства (поле Code_Publish) находится в диапазоне от 10 до 25,
-- и количество страниц (поле Pages) в книге больше 120.
-- 25
select b.code_publish, a.name, b.title_book, b.pages
from authors a
join books b on b.code_author = a.code_author
where (b.code_publish between 10 and 2500) and b.pages > 120
