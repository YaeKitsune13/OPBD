-- Создание таблицы Поставщиков (S)
CREATE TABLE S (
                   s_no CHAR(6) NOT NULL PRIMARY KEY,
                   s_name VARCHAR(20) NOT NULL,
                   status INT,
                   city VARCHAR(20)
);

-- Создание таблицы Деталей (P)
CREATE TABLE P (
                   p_no CHAR(6) NOT NULL PRIMARY KEY,
                   p_name VARCHAR(20) NOT NULL,
                   color VARCHAR(20),
                   weight INT,
                   city VARCHAR(20)
);

-- Создание таблицы Изделий (J)
CREATE TABLE J (
                   j_no CHAR(6) NOT NULL PRIMARY KEY,
                   j_name VARCHAR(20) NOT NULL,
                   city VARCHAR(20)
);

-- Создание таблицы Поставок (SPJ)
CREATE TABLE SPJ (
                     s_no CHAR(6) NOT NULL,
                     p_no CHAR(6) NOT NULL,
                     j_no CHAR(6) NOT NULL,
                     qty INT,
                     PRIMARY KEY (s_no, p_no, j_no),
                     FOREIGN KEY (s_no) REFERENCES S(s_no),
                     FOREIGN KEY (p_no) REFERENCES P(p_no),
                     FOREIGN KEY (j_no) REFERENCES J(j_no)
);

-- Заполнение таблицы S (Поставщики)
INSERT INTO S (s_no, s_name, status, city) VALUES
                                               ('S1', 'Смит', 20, 'Лондон'),
                                               ('S2', 'Джонс', 10, 'Париж'),
                                               ('S3', 'Блейк', 30, 'Париж'),
                                               ('S4', 'Кларк', 20, 'Лондон'),
                                               ('S5', 'Адамс', 30, 'Афины');

-- Заполнение таблицы P (Детали)
INSERT INTO P (p_no, p_name, color, weight, city) VALUES
                                                      ('P1', 'Гайка', 'Красный', 12, 'Лондон'),
                                                      ('P2', 'Болт', 'Зеленый', 17, 'Париж'),
                                                      ('P3', 'Винт', 'Голубой', 17, 'Рим'),
                                                      ('P4', 'Винт', 'Красный', 14, 'Лондон'),
                                                      ('P5', 'Кулачок', 'Голубой', 12, 'Париж'),
                                                      ('P6', 'Блюм', 'Красный', 19, 'Лондон');

-- Заполнение таблицы J (Изделия)
INSERT INTO J (j_no, j_name, city) VALUES
                                       ('J1', 'Жесткий диск', 'Париж'),
                                       ('J2', 'Перфоратор', 'Рим'),
                                       ('J3', 'Считыватель', 'Афины'),
                                       ('J4', 'Принтер', 'Афины'),
                                       ('J5', 'Флоппи-диск', 'Лондон'),
                                       ('J6', 'Терминал', 'Осло'),
                                       ('J7', 'Лента', 'Лондон');

-- Заполнение таблицы SPJ (Поставки)
INSERT INTO SPJ (s_no, p_no, j_no, qty) VALUES
                                            ('S1', 'P1', 'J1', 200),
                                            ('S1', 'P1', 'J4', 700),
                                            ('S2', 'P3', 'J1', 400),
                                            ('S2', 'P3', 'J2', 200),
                                            ('S2', 'P3', 'J3', 200),
                                            ('S2', 'P3', 'J4', 500),
                                            ('S2', 'P3', 'J5', 600),
                                            ('S2', 'P3', 'J6', 400),
                                            ('S2', 'P3', 'J7', 800),
                                            ('S2', 'P5', 'J2', 100),
                                            ('S3', 'P3', 'J1', 200),
                                            ('S3', 'P4', 'J2', 500),
                                            ('S4', 'P6', 'J3', 300),
                                            ('S4', 'P6', 'J7', 300),
                                            ('S5', 'P2', 'J2', 200),
                                            ('S5', 'P2', 'J4', 100),
                                            ('S5', 'P5', 'J5', 500),
                                            ('S5', 'P5', 'J7', 100),
                                            ('S5', 'P6', 'J2', 200),
                                            ('S5', 'P1', 'J4', 100),
                                            ('S5', 'P3', 'J4', 200),
                                            ('S5', 'P4', 'J4', 800),
                                            ('S5', 'P5', 'J4', 400),
                                            ('S5', 'P6', 'J4', 500);