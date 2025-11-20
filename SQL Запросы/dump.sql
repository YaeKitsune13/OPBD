create database HOMIQ;

use HOMIQ;
-- Таблица с адресами всех пользователей
create table Address(

)
-- Таблица для пользователей
create table Uses(
    user_Id int primary key auto_increment,
    login varchar(60),
    email varchar(100),
    password varchar(20)
);
-- Таблица для админов
create table Admins(
    admin_Id int primary key auto_increment,
    login varchar(60),
    email varchar(100),
    password varchar(20)
);
-- Таблица для рабочих
create table Employee(
    employee_Id int primary key auto_increment,
    login varchar(60),
    email varchar(100),
    password varchar(20)
);


