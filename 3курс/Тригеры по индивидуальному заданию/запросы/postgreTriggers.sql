-- # 25. бд страховой компании.

-- # перед вставкой проверять, что дата заключения договора раньше даты окончания договора. 
-- # в случае нарушения отменять операцию и выдавать сообщение

create or replace function fn_insert_policyholder_date_checker()
returns trigger as $$
begin
    if new.contract_date > new.end_date then
        -- # в postgres используется raise exception вместо signal sqlstate
        raise exception 'ошибка: дата начала договора не может быть позже чем конец';
    else
        insert into logs (message, user_name)
        values (concat('добавили пользователя с id ', new.policy_number), current_user);
    end if;
    return new;
end;
$$ language plpgsql;

drop trigger if exists db_insert_policyholder_date_checker on policyholders;
create trigger db_insert_policyholder_date_checker
    before insert on policyholders
    for each row execute function fn_insert_policyholder_date_checker();

-- # перед обновлением проверять, что дата заключения договора раньше даты окончания договора. 

create or replace function fn_before_update_policyholder_date_checker()
returns trigger as $$
begin
    if new.contract_date > new.end_date then
        raise exception 'ошибка: дата начала договора не может быть позже чем конец';
    elsif new.end_date > old.end_date then
        insert into logs (message, user_name)
        values (concat('обновили пользователя с id ', new.policy_number, ' его договор продлён на ',
                       (new.end_date - old.end_date), ' дней'), current_user);
    else
        insert into logs (message, user_name)
        values (concat('обновили пользователя с id ', new.policy_number, ' его договор укорочен на ',
                       (old.end_date - new.end_date), ' дней'), current_user);
    end if;
    return new;
end;
$$ language plpgsql;

drop trigger if exists db_before_update_policyholder_date_checker on policyholders;
create trigger db_before_update_policyholder_date_checker
    before update on policyholders
    for each row execute function fn_before_update_policyholder_date_checker();


-- # перед вставкой автоматически рассчитывать сумму страховки. 

create or replace function fn_insert_policyholders_calculate_sum()
returns trigger as $$
begin
    new.policy_cost := new.premium_amount * (new.end_date - new.contract_date);
    return new;
end;
$$ language plpgsql;

drop trigger if exists db_insert_policyholders_calculate_sum on policyholders;
create trigger db_insert_policyholders_calculate_sum
    before insert on policyholders
    for each row execute function fn_insert_policyholders_calculate_sum();


-- # после заключения нового договора увеличить счетчик договоров у сотрудника.

create or replace function fn_update_employers_if_insert_policyholder()
returns trigger as $$
begin
    update employees
    set count_policyholders = coalesce(count_policyholders, 0) + 1
    where employee_id = new.employee_id;
    return new;
end;
$$ language plpgsql;

drop trigger if exists db_update_employers_if_insert_policyholder on policyholders;
create trigger db_update_employers_if_insert_policyholder
    after insert on policyholders
    for each row execute function fn_update_employers_if_insert_policyholder();


-- # при досрочном расторжении записывать старые данные записи в архив.

create or replace function fn_remove_set_log()
returns trigger as $$
begin
    insert into logs (message, user_name)
    values (concat('договор: ', old.policy_number, ' был росторжен за ', (old.end_date - old.contract_date), ' дней'),
            current_user);
    return old;
end;
$$ language plpgsql;

drop trigger if exists db_remove_set_log on policyholders;
create trigger db_remove_set_log
    before delete on policyholders
    for each row execute function fn_remove_set_log();