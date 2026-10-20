from __future__ import annotations

from PyQt6.QtGui import QStandardItem, QStandardItemModel
from PyQt6.QtWidgets import (
    QComboBox, QHBoxLayout, QMessageBox, QPushButton, QTableView, QVBoxLayout, QWidget,
)

from db_manager import DatabaseError, DatabaseManager

REPORTS: list[tuple[str, str]] = [
    (
        "1. Полисы, где выплаты превышают премию",
        """
        WITH payouts AS (
            SELECT policy_number, COALESCE(SUM(payout), 0) AS total_payout
            FROM claims
            GROUP BY policy_number
        )
        SELECT p.policy_number, p.full_name, p.premium_amount,
               COALESCE(py.total_payout, 0) AS total_payout
        FROM policyholders p
        LEFT JOIN payouts py ON py.policy_number = p.policy_number
        WHERE COALESCE(py.total_payout, 0) > p.premium_amount
        ORDER BY p.policy_number
        """,
    ),
    (
        "2. Полисы, где первый случай раньше чем через 6 дней",
        """
        WITH first_claim AS (
            SELECT policy_number, MIN(event_date) AS first_event_date
            FROM claims
            GROUP BY policy_number
        )
        SELECT p.policy_number, p.full_name, p.contract_date, fc.first_event_date,
               DATE_ADD(p.contract_date, INTERVAL 6 DAY) AS must_not_be_earlier_than
        FROM policyholders p
        JOIN first_claim fc ON fc.policy_number = p.policy_number
        WHERE fc.first_event_date < DATE_ADD(p.contract_date, INTERVAL 6 DAY)
        ORDER BY p.policy_number
        """,
    ),
    (
        "3. Страхователи по видам страхования",
        """
        SELECT it.name AS insurance_type, p.policy_number, p.full_name AS policyholder_name,
               p.contract_date, p.end_date
        FROM policyholders p
        JOIN insurance_types it ON it.insurance_type_id = p.insurance_type_id
        ORDER BY insurance_type, policyholder_name
        """,
    ),
    (
        "4. Сотрудники без договоров за сегодня",
        """
        SELECT e.employee_id, e.full_name, e.position, e.count_policyholders
        FROM employees e
        WHERE NOT EXISTS (
            SELECT 1 FROM policyholders p
            WHERE p.employee_id = e.employee_id AND p.contract_date = CURDATE()
        )
        ORDER BY e.full_name
        """,
    ),
    (
        "5. Полисы, где стоимость не равна годовой стоимости вида страхования",
        """
        SELECT p.policy_number, p.full_name, it.name AS insurance_type,
               p.policy_cost, it.annual_cost
        FROM policyholders p
        JOIN insurance_types it ON it.insurance_type_id = p.insurance_type_id
        WHERE p.policy_cost <> it.annual_cost
        ORDER BY p.policy_number
        """,
    ),
]


class ReportsTab(QWidget):
    def __init__(self, db: DatabaseManager, parent=None):
        super().__init__(parent)
        self.db = db

        self.combo = QComboBox()
        for title, _ in REPORTS:
            self.combo.addItem(title)
        run_btn = QPushButton("Выполнить")
        run_btn.clicked.connect(self.run_report)

        top = QHBoxLayout()
        top.addWidget(self.combo)
        top.addWidget(run_btn)
        top.addStretch()

        self.table_view = QTableView()
        self.model = QStandardItemModel()
        self.table_view.setModel(self.model)
        self.table_view.setEditTriggers(QTableView.EditTrigger.NoEditTriggers)

        layout = QVBoxLayout(self)
        layout.addLayout(top)
        layout.addWidget(self.table_view)

    def run_report(self) -> None:
        index = self.combo.currentIndex()
        _, query = REPORTS[index]
        try:
            rows = self.db.run_query(query)
        except DatabaseError as exc:
            QMessageBox.critical(self, "Ошибка отчёта", str(exc))
            return
        self.model.clear()
        if not rows:
            self.model.setHorizontalHeaderLabels(["Нарушений не найдено"])
            return
        headers = list(rows[0].keys())
        self.model.setHorizontalHeaderLabels(headers)
        for row in rows:
            items = [QStandardItem("" if v is None else str(v)) for v in row.values()]
            for it in items:
                it.setEditable(False)
            self.model.appendRow(items)
        self.table_view.resizeColumnsToContents()
