create database hospitalfinaldatabaseinacorparationintheworldhimisgoodandgod;
use hospitalfinaldatabaseinacorparationintheworldhimisgoodandgod;

create table WorkPosition(
    work_id int primary key auto_increment,
    work_type varchar(60),
    work_name varchar(200)
);

create table patient(
    card_number int primary key auto_increment,
    surname varchar(50),
    name varchar(50),
    lastname varchar(50),
    pass varchar(11),
    gender enum('М','Ж'),
    birthday date,
    phone int,
    photo varchar(250),
    email varchar(100),
    card_get_date date,
    work_pos_id int references WorkPosition(work_id)
);

create table Polis(
    polis_id int primary key auto_increment,
    card_number int references patient(card_number),
    polis_number int,
    get_date date,
    ending_date date,
    company varchar(60)
);

create table Specialization(
    specialization_id int primary key auto_increment,
    specialization_name varchar(60)
);

create table Doctors(
    doctor_id int primary key auto_increment,
    surname varchar(80),
    name varchar(80),
    lastname varchar(80),
    specialization int references Specialization(specialization_id)
);

create table Sendings(
    sending_id int primary key auto_increment,
    patient_card int references patient(card_number),
    doctor_id int references Doctors(doctor_id),
    date_sending date,
    message varchar(400),
    status enum('Активный','Просрочен')
);

create table DiagnosCategory(
    category_id int primary key auto_increment,
    category_name varchar(100)
);

create table Diagnosis(
    diagnose_id int primary key auto_increment,
    name varchar(80),
    category_id int references DiagnosCategory(category_id)
);

create table SendingToDiagnos(
    send_diagnose_id int primary key auto_increment,
    send_id int references Sendings(sending_id),
    diagnose_id int references Diagnosis(diagnose_id)
)