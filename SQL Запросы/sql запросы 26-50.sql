use books_deal;

-- 26 задание --
select * from publishing_house where publish in
(
	select publish
    from books
    where title_book
    like '%Горе%'
)
and City = 'Москва';

-- 27.	Вывести список названий компаний-поставщиков (поле Name_company)
-- и названия книг (поле Title_book),
-- которые они поставили в период с 01.01.2002 по 31.12.2003 (условие по полю Date_order).
select distinct d.Name_company, b.Title_book
from purshases p
join books b on p.code_book = b.code_book
join deliveries d on p.code_delivery = d.code_delivery
where p.Date_order between '2002-01-01' and '2003-12-31';

-- 28.	Вывести список авторов (поле Name_author),
-- книги которых были выпущены в издательстве ‘Мир’ (условие по полю Publish).
select a.name_autor, b.title_book
from books b
join publishing_house p on p.code_publish = b.code_publish
join authors a on b.code_autor = a.code_autor
where p.publish = 'Мир';

-- 29.	Вывести список поставщиков (поле Name_company),
-- которые поставляют книги издательства ‘Питер’ (условие по полю Publish).
select distinct d.Name_company
from books b
join purshases p on b.code_book = p.code_book
join deliveries d on p.code_delivery = d.code_delivery
join publishing_house ph on b.code_publish = ph.code_publish
where ph.publish like '%Питер%';

-- 31.	Вывести суммарную стоимость партии одноименных книг (использовать поля Amount и Cost)
-- и название книги (поле Title_book) в каждой поставке.
select
    b.Title_book,
    p.purshases_id,
    sum(p.Cost) as Total_Cost
from purshases p
join books b on p.code_book = b.code_book
group by b.Title_book, p.purshases_id;

-- 32.	Вывести стоимость одной печатной страницы каждой книги (использовать поля Cost и Pages)
-- названия соответствующих книг (поле Title_book)
select
    b.Title_book,
    round((p.Cost / b.Pages), 4) as Cost_per_Page
from purshases p
join books b on p.code_book = b.code_book;

-- 33.	Вывести количество лет с момента рождения авторов (использовать поле Birthday)
-- и имена соответствующих авторов (поле Name_author).
select
    a.Name_autor,
    year(current_date) - year(a.Birthday) -
    case
        when month(current_date) < month(a.Birthday)
             or (month(current_date) = month(a.Birthday) and day(current_date) < day(a.Birthday))
        then 1
        else 0
    end as Age
from authors a;

-- 34.	Вывести общую сумму поставок книг (использовать поле Cost),
-- выполненных ‘ЗАО Оптторг’ (условие по полю Name_company).
select
    d.Name_company,
    sum(p.Cost) as Total_Supply_Cost
from purshases p
join deliveries d on p.code_delivery = d.code_delivery
where d.Name_company = 'ЗАО Оптторг'
group by d.Name_company;

-- 35.	Вывести общее количество всех поставок (использовать любое по-ле из таблицы Purchases),
-- выполненных в период с 01.01.2003 по 01.02.2003 (условие по полю Date_order).
select
    count(*) as Total_Supplies
from purshases p
where p.Date_order between '2003-01-01' and '2003-02-01';

-- 36.	Вывести среднюю стоимость (использовать поле Cost)
-- и среднее количество экземпляров книг (использовать поле Amount) в одной поставке,
-- где автором книги является ‘Акунин’ (условие по полю Name_author).
select
    round(avg(p.Cost), 2) as Average_Cost
from purshases p
join books b on p.code_book = b.code_book
join authors a on b.code_autor = a.code_autor
where a.Name_autor like '%Акунин%';

-- 37.	Вывести все сведения о поставке (все поля таблицы Purchases),
-- а также название книги (поле Title_book) с минимальной общей стоимостью (использовать поля Cost и Amount).
select
    b.Title_book,
    (p.Cost) as Minimal_cost
from purshases p
join books b on p.code_book = b.code_book
where (p.Cost) = (
    select min(p2.Cost)
    from purshases p2
);

-- 38.	Вывести все сведения о поставке (все поля таблицы Purchases),
-- а также название книги (поле Title_book) с максимальной общей стоимостью (использовать поля Cost и Amount).
select
    b.Title_book,
    (p.Cost) as Maximal_cost
from purshases p
join books b on p.code_book = b.code_book
where (p.Cost) = (
    select max(p2.Cost)
    from purshases p2
);

-- 39.	Вывести название книги (поле Title_book),
-- суммарную стоимость партии одноименных книг (использовать поля Amount и Cost),
-- поместив в результат в поле с названием Itogo,
-- в поставках за период с 01.01.2002 по 01.06.2002 (условие по полю Date_order).
select
    b.Title_book,
    sum(p.Cost) as Itogo
from purshases p
join books b on p.code_book = b.code_book
where p.Date_order between '2002-01-01' and '2002-06-01'
group by b.Title_book;

-- 40.	Вывести стоимость одной печатной страницы каждой книги (использовать поля Cost и Pages),
-- поместив результат в поле с названием One_page, и названия соответствующих книг (поле Title_book).
select
    b.Title_book,
    round((p.Cost / b.Pages), 4) as One_page
from purshases p
join books b on p.code_book = b.code_book;

-- 41.	Вывести общую сумму поставок книг (использовать поле Cost)
 -- и поместить результат в поле с названием Sum_cost, выполненных ‘ОАО Луч ’ (условие по полю Name_company).
select
    d.Name_company,
    sum(p.Cost) as Sum_cost
from purshases p
join deliveries d on p.code_delivery = d.code_delivery
where d.Name_company = 'ОАО Луч'
group by d.Name_company;

-- 42.	Изменить в таблице Books содержимое поля Pages на 300, если код автора
-- (поле Code_author) =56 и название книги (поле Title_book) =’Мемуары’.
update books
set Pages = 300
where code_autor = 56 and title_book = 'Мемуары';

-- 43.	Изменить в таблице Deliveries содержимое поля Address на ‘нет сведений’,
-- если значение поля является пустым.
update deliveries
set adress = 'нет сведений'
where adress is null or adress = '';

-- 44.	Увеличить в таблице Purchases цену (поле Cost)
-- на 20 процентов, если заказы были оформлены в течение последнего месяца (условие по полю Date_order).
update purshases
set Cost = Cost * 1.2
where Date_order >= date_sub(current_date, interval 1 month);

-- 45.	Добавить в таблицу Purchases новую запись, причем так, чтобы код покупки (поле Code_purchase)
-- был автоматически увеличен на единицу, а в тип закупки (поле Type_purchase) внести значение ‘опт’.
insert into purshases (code_delivery, code_book, Date_order,cost,code_purchaase)
values (1, 1, current_date,100,2);

-- 46.	Добавить в таблицу Books новую запись, причем вместо ключевого поля поставить код (поле Code_book),
-- автоматически увеличенный на единицу от максимального кода в таблице,
-- вместо названия книги (поле Title_book) написать ‘Наука. Техника. Инновации’.
insert into books (code_publish, code_autor, title_book, pages)
values (221, 1, 'Наука. Техника. Инновации', 250);

-- 47.	Добавить в таблицу Publish_house новую запись, причем вместо ключевого поля поставить код (поле Code_publish),
-- автоматически увеличенный на единицу от максимального кода в таблице,
-- Вместо названия города – ‘Москва’ (поле City), вместо издательства – ‘Наука’ (поле Publish).
insert into publishing_house (city, publish)
values ('Москва', 'Наука');

-- 48.	Удалить из таблицы Purchases все записи, у которых количество книг в заказе (поле Amount) = 0.
delete from purshases
where cost = 0;

-- 49.	Удалить из таблицы Authors все записи, у которых нет имени автора в поле Name_Author
delete from authors
where Name_autor is null or Name_autor = '';

-- 50.	Удалить из таблицы Deliveries все записи, у которых не указан ИНН (поле INN пустое).
delete from deliveries
where INN is null or INN = '';
