from __future__ import annotations

from PyQt6.QtWidgets import QMainWindow, QTabWidget

from crud_tab import CrudTab
from db_manager import DatabaseManager
from reports_tab import ReportsTab
from table_specs import CLAIMS_SPEC, EMPLOYEES_SPEC, INSURANCE_TYPES_SPEC, POLICYHOLDERS_SPEC


class MainWindow(QMainWindow):
    def __init__(self, db: DatabaseManager):
        super().__init__()
        self.db = db
        self.setWindowTitle("Учёт страховой компании")
        self.resize(1150, 650)

        tabs = QTabWidget()
        tabs.addTab(CrudTab(INSURANCE_TYPES_SPEC, db), INSURANCE_TYPES_SPEC.title)
        tabs.addTab(CrudTab(EMPLOYEES_SPEC, db), EMPLOYEES_SPEC.title)
        tabs.addTab(CrudTab(POLICYHOLDERS_SPEC, db), POLICYHOLDERS_SPEC.title)
        tabs.addTab(CrudTab(CLAIMS_SPEC, db), CLAIMS_SPEC.title)
        tabs.addTab(ReportsTab(db), "Отчёты")
        self.setCentralWidget(tabs)
