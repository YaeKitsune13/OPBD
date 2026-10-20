# 25. БД страховой компании.

# Перед вставкой проверять, что дата заключения договора раньше даты окончания договора. В случае нарушения отменять операцию и выдавать сообщение
delimiter //
drop trigger if exists db_insert_policyholder_date_checker;
create trigger db_insert_policyholder_date_checker
    before insert
    on policyholders
    for each row
begin
    if NEW.contract_date > new.end_date then
        # Пытался добавить логирование на ошибку но не получилось при вызове sqlstate
        # скрипт закрывается раньше чем успевает выполнить запрос
        signal sqlstate '45000'
            set message_text = "Ошибка: Дата начала договора не может быть позже чем конец";
    else
        insert into logs
        values (CONCAT('Добавили пользователя с id ', new.policy_number),
                user());
    end if;
end //

delimiter ;

# Пример ошибочного запроса с датой в на год отличающейся от основного
insert into policyholders (policy_number, passport, full_name, birth_date, insurance_type_id, employee_id,
                           contract_date, end_date, premium_amount, policy_cost)
values ('pl00000008', '4500 111111', 'кузнецов максим', date '1990-05-01', 1, 1001,
        date_add(curdate(), interval 2 year), date_add(curdate(), interval 1 year), 12000, 10000);


# Перед обновлением проверять, что дата заключения договора раньше даты окончания договора. В случае нарушения отменять операцию и выдавать сообщение

delimiter //

drop trigger if exists db_before_update_policyholder_date_checker;
create trigger db_before_update_policyholder_date_checker
    before update
    on policyholders
    for each row
begin
    if new.contract_date > new.end_date then
        signal sqlstate '45000'
            set message_text = "Ошибка: Дата начала договора не может быть позже чем конец";
    elseif new.end_date > old.end_date then
        insert into logs
        values (CONCAT('Обновили пользователя с id ', new.policy_number, ' его договор продлён на ',
                       datediff(new.end_date, old.end_date), ' дней'),
                user());
    else
        insert into logs
        values (CONCAT('Обновили пользователя с id ', new.policy_number, ' его договор укорочен на ',
                       datediff(old.end_date, new.end_date), ' дней'),
                user());
    end if;
end //
delimiter ;
# начальные значения перед обновлением
select * from policyholders;
# Успешное обновление
update policyholders
set end_date = adddate(curdate(), interval 5 year)
where policy_number = 'pl00000001';

# Ошибочное обновление
update policyholders
set end_date = adddate(curdate(), interval -10 year)
where policy_number = 'pl00000002';
# Перед вставкой автоматически рассчитывать сумму страховки. Формула: сумма страховки = Сумма премии * (Количество лет между датой заключения и датой окончания).

# проверка обновления
select * from policyholders;

delimiter //

drop trigger if exists db_insert_policyholders_calculate_sum;
create trigger db_insert_policyholders_calculate_sum
    before insert
    on policyholders
    for each row
begin
    set new.policy_cost = new.premium_amount * datediff(new.end_date, new.contract_date);
end //
delimiter ;

# Успешный ввод
insert into policyholders (policy_number, passport, full_name, birth_date, insurance_type_id, employee_id,
                           contract_date, end_date, premium_amount, policy_cost)
values ('pl00000010', '4500 111111', 'кузнецов максим', date '1990-05-01', 1, 1001,
        curdate(), date_add(curdate(), interval 1 year), 12000, 0);
#Проверка добавление
select * from policyholders;

# После заключения нового договора увеличить счетчик договоров у сотрудника.

delimiter //

drop trigger if exists db_update_employers_if_insert_policyholder;
create trigger db_update_employers_if_insert_policyholder
    after insert
    on policyholders
    for each row
begin
    update employees
    set count_policyholders = count_policyholders + 1
    where NEW.employee_id = employee_id;
end //

delimiter ;

# Пример успешного
insert into policyholders (policy_number, passport, full_name, birth_date, insurance_type_id, employee_id,
                           contract_date, end_date, premium_amount, policy_cost)
values ('pl00000012', '4500 111111', 'кузнецов максим', date '1990-05-01', 1, 1001,
        curdate(), date_add(curdate(), interval 2 year), 12000, 0);
#Проверка добалвнеия
select * from policyholders;
# При досрочном расторжении записывать старые данные записи в архив (отдельную таблицу).

delimiter //

drop trigger if exists db_remove_set_log;
create trigger db_remove_set_log
    before delete
    on policyholders
    for each row
begin
    insert into logs
    values (CONCAT('Договор: ', old.policy_number, ' был росторжен за ', datediff(old.end_date, old.contract_date), ' дней'),
                   user()
           );
end //

delimiter ;


DELETE
FROM insurance.policyholders
WHERE policy_number = 'pl00000012';

#Проверка удаления
select * from policyholders;

show triggers;