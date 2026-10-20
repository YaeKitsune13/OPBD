-- drop database if exists insurance;
-- create database insurance;

-- \c insurance

-- 1. виды страхования
create table insurance_types (
                                 insurance_type_id smallint primary key,
                                 name varchar(100),
                                 description text,
                                 annual_cost decimal(12,2)
);

-- 2. сотрудники
create table employees (
                           employee_id int primary key,
                           full_name varchar(100),
                           passport varchar(50),
                           position varchar(60)
);

-- 3. страхователи (полисы)
create table policyholders (
                               policy_number char(10) primary key,
                               passport varchar(50),
                               full_name varchar(40),
                               birth_date date,
                               insurance_type_id smallint,
                               employee_id int,
                               contract_date date,
                               end_date date,
                               premium_amount decimal(8,0),
                               policy_cost decimal(8,0),
                               constraint fk_policyholders_type foreign key (insurance_type_id) references insurance_types(insurance_type_id),
                               constraint fk_policyholders_employee foreign key (employee_id) references employees(employee_id)
);

-- 4. страховые случаи
create table claims (
                        claim_id integer generated always as identity primary key,
                        policy_number char(10),
                        description text,
                        event_date date,
                        payout decimal(12,2),
                        constraint fk_claims_policy foreign key (policy_number) references policyholders(policy_number)
);

-- ТЕСТОВЫЕ ДАННЫЕ

insert into insurance_types (insurance_type_id, name, description, annual_cost) values
                                                                                    (1, 'автострахование', 'осаго/каско', 10000.00),
                                                                                    (2, 'путешествия', 'страхование выезжающих', 15000.00),
                                                                                    (3, 'имущество', 'страхование квартиры/дома', 20000.00);

insert into employees (employee_id, full_name, passport, position) values
                                                                       (1001, 'иванов иван иванович', '4000 123456', 'агент'),
                                                                       (1002, 'петров пётр петрович', '4001 234567', 'агент'),
                                                                       (1003, 'сидорова анна сергеевна', '4002 345678', 'старший агент');

-- Используем CURRENT_DATE
insert into policyholders (policy_number, passport, full_name, birth_date, insurance_type_id, employee_id, contract_date, end_date, premium_amount, policy_cost) values
                                                                                                                                                                     ('pl00000001', '4500 111111', 'кузнецов максим', date '1990-05-01', 1, 1001, current_date, current_date + interval '1 year', 12000, 10000),
                                                                                                                                                                     ('pl00000002', '4500 222222', 'смирнова ольга',  date '1988-03-12', 2, 1002, current_date - interval '10 days', current_date + interval '355 days', 15000, 16000),
                                                                                                                                                                     ('pl00000003', '4500 333333', 'попов алексей',   date '1985-11-21', 3, 1002, current_date - interval '30 days', current_date + interval '11 months', 5000, 20000),
                                                                                                                                                                     ('pl00000004', '4500 444444', 'соколова ирина',  date '1995-07-15', 1, 1003, current_date - interval '8 days', current_date + interval '1 year', 8000, 10000);

insert into claims (policy_number, description, event_date, payout) values
                                                                        ('pl00000001', 'дтп. повреждение бампера', (select contract_date + interval '7 days' from policyholders where policy_number='pl00000001'), 3000.00),
('pl00000001', 'замена стекла',            (select contract_date + interval '20 days' from policyholders where policy_number='pl00000001'), 5000.00),
('pl00000002', 'несчастный случай',        (select contract_date + interval '3 days' from policyholders where policy_number='pl00000002'), 6000.00),
('pl00000002', 'медицинские расходы',      (select contract_date + interval '9 days' from policyholders where policy_number='pl00000002'), 12000.00),
('pl00000004', 'дтп. сколы лкп',           (select contract_date + interval '6 days' from policyholders where policy_number='pl00000004'), 2000.00);

-- ЗАПРОСЫ

-- 1)
with payouts as (
    select policy_number, coalesce(sum(payout),0) as total_payout
    from claims
    group by policy_number
)
select
    p.policy_number,
    p.premium_amount,
    coalesce(py.total_payout, 0) as total_payout
from policyholders p
         left join payouts py on py.policy_number = p.policy_number
where coalesce(py.total_payout, 0) > p.premium_amount
order by p.policy_number;

-- 2)
with first_claim as (
    select policy_number, min(event_date) as first_event_date
    from claims
    group by policy_number
)
select
    p.policy_number,
    p.contract_date,
    fc.first_event_date,
    p.contract_date + interval '6 days' as must_not_be_earlier_than
from policyholders p
    join first_claim fc on fc.policy_number = p.policy_number
where fc.first_event_date < (p.contract_date + interval '6 days')
order by p.policy_number;

-- 3)
select
    it.insurance_type_id,
    it.name as insurance_type,
    p.policy_number,
    p.full_name as policyholder_name,
    p.contract_date,
    p.end_date
from policyholders p
         join insurance_types it on it.insurance_type_id = p.insurance_type_id
order by insurance_type, policyholder_name;

-- 4)
select e.*
from employees e
where not exists (
    select 1
    from policyholders p
    where p.employee_id = e.employee_id
      and p.contract_date = current_date
)
order by e.full_name;

-- 5)
select
    p.policy_number,
    it.name as insurance_type,
    p.policy_cost,
    it.annual_cost
from policyholders p
         join insurance_types it on it.insurance_type_id = p.insurance_type_id
where p.policy_cost <> it.annual_cost
order by p.policy_number;

-- 6) Примеры update/insert
insert into employees (employee_id, full_name, passport, position)
values (1004, 'новичков фёдор тимофеевич', '4003 456789', 'агент');

update employees
set position = 'ведущий агент'
where employee_id = 1004;

insert into policyholders (policy_number, passport, full_name, birth_date, insurance_type_id, employee_id, contract_date, end_date, premium_amount, policy_cost)
values ('pl00000005', '4500 555555', 'орлов игорь', date '1992-09-09', 1, 1001, current_date, current_date + interval '1 year', 9000, 10000);

insert into claims (policy_number, description, event_date, payout)
values ('pl00000005', 'первый случай через 7 дней', (select contract_date + interval '7 days' from policyholders where policy_number='pl00000005'), 1000.00);