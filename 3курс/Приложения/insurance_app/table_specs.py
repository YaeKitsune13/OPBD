"""
Метаданные таблиц: какие колонки показывать, какого типа, и где брать
варианты для выпадающих списков (внешние ключи). На основе этого
crud_tab.py строит и таблицу, и форму добавления/изменения автоматически.
"""
from __future__ import annotations

from dataclasses import dataclass
from enum import Enum, auto


class FieldType(Enum):
    TEXT = auto()
    TEXTAREA = auto()
    INT = auto()
    DECIMAL = auto()
    DATE = auto()
    LOOKUP = auto()  # выпадающий список со значениями из другой таблицы (внешний ключ)


@dataclass
class FieldSpec:
    name: str
    label: str
    field_type: FieldType
    editable: bool = True
    required: bool = True
    lookup_table: str | None = None
    lookup_key: str | None = None
    lookup_label: str | None = None


@dataclass
class TableSpec:
    table: str
    pk: str
    title: str
    fields: list[FieldSpec]
    pk_editable: bool = True  # False для автоинкрементного pk (claim_id)


INSURANCE_TYPES_SPEC = TableSpec(
    table="insurance_types",
    pk="insurance_type_id",
    title="Виды страхования",
    pk_editable=True,
    fields=[
        FieldSpec("insurance_type_id", "Код вида", FieldType.INT),
        FieldSpec("name", "Наименование", FieldType.TEXT),
        FieldSpec("description", "Описание", FieldType.TEXTAREA, required=False),
        FieldSpec("annual_cost", "Годовая стоимость", FieldType.DECIMAL),
    ],
)

EMPLOYEES_SPEC = TableSpec(
    table="employees",
    pk="employee_id",
    title="Сотрудники",
    pk_editable=True,
    fields=[
        FieldSpec("employee_id", "Табельный номер", FieldType.INT),
        FieldSpec("full_name", "ФИО", FieldType.TEXT),
        FieldSpec("passport", "Паспорт", FieldType.TEXT),
        FieldSpec("position", "Должность", FieldType.TEXT),
        FieldSpec(
            "count_policyholders", "Кол-во страхователей", FieldType.INT,
            editable=False, required=False,
        ),
    ],
)

POLICYHOLDERS_SPEC = TableSpec(
    table="policyholders",
    pk="policy_number",
    title="Страхователи (полисы)",
    pk_editable=True,
    fields=[
        FieldSpec("policy_number", "Номер полиса", FieldType.TEXT),
        FieldSpec("passport", "Паспорт", FieldType.TEXT),
        FieldSpec("full_name", "ФИО", FieldType.TEXT),
        FieldSpec("birth_date", "Дата рождения", FieldType.DATE),
        FieldSpec(
            "insurance_type_id", "Вид страхования", FieldType.LOOKUP,
            lookup_table="insurance_types", lookup_key="insurance_type_id", lookup_label="name",
        ),
        FieldSpec(
            "employee_id", "Сотрудник", FieldType.LOOKUP,
            lookup_table="employees", lookup_key="employee_id", lookup_label="full_name",
        ),
        FieldSpec("contract_date", "Дата договора", FieldType.DATE),
        FieldSpec("end_date", "Дата окончания", FieldType.DATE),
        FieldSpec("premium_amount", "Сумма премии", FieldType.DECIMAL),
        FieldSpec("policy_cost", "Стоимость полиса", FieldType.DECIMAL),
    ],
)

CLAIMS_SPEC = TableSpec(
    table="claims",
    pk="claim_id",
    title="Страховые случаи",
    pk_editable=False,
    fields=[
        FieldSpec(
            "policy_number", "Номер полиса", FieldType.LOOKUP,
            lookup_table="policyholders", lookup_key="policy_number", lookup_label="full_name",
        ),
        FieldSpec("description", "Описание", FieldType.TEXTAREA),
        FieldSpec("event_date", "Дата случая", FieldType.DATE),
        FieldSpec("payout", "Сумма выплаты", FieldType.DECIMAL),
    ],
)

ALL_SPECS = [INSURANCE_TYPES_SPEC, EMPLOYEES_SPEC, POLICYHOLDERS_SPEC, CLAIMS_SPEC]
