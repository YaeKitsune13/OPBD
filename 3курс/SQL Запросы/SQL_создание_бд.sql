create database if not exists books_deal;

use books_deal;


create table if not exists authors(
	code_autor int primary key auto_increment,
    name_autor varchar(100),
    birthday date
);
create table if not exists deliveries(
	code_delivery int primary key auto_increment,
    name_company varchar(80),
    adress varchar(120),
    phone bigint,
    INN bigint
);
create table if not exists publishing_house(
	code_publish int primary key auto_increment,
    publish varchar(120),
    City varchar(40)
);

create table if not exists books(
	code_book int primary key auto_increment,
    title_book varchar(40),
    code_autor int,
    pages int,
    code_publish int,
    constraint foreign key (code_autor) references authors (code_autor),
    constraint foreign key (code_publish) references publishing_house (code_publish)
);
create table if not exists purshases(
	purshases_id int primary key auto_increment,
	code_book int,
    date_order date,
    code_delivery int,
    cost int,
    code_purchaase int,
    constraint foreign key (code_delivery) references deliveries (code_delivery),
    constraint foreign key (code_book) references books (code_book)
);