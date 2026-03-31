import sys
from PyQt6.QtWidgets import (QApplication, QMainWindow, QTableView, 
                             QVBoxLayout, QWidget, QLabel, 
                             QMessageBox, QSplitter, QHeaderView, QComboBox)
from PyQt6.QtSql import QSqlDatabase, QSqlQueryModel, QSqlQuery
from PyQt6.QtCore import Qt
from PyQt6.QtGui import QFont

# Настройки подключения
MYSQL_HOST, MYSQL_DB, MYSQL_USER, MYSQL_PASS = "127.0.0.1", "insurance", "root", "1234"
PG_HOST, PG_DB, PG_USER, PG_PASS = "127.0.0.1", "insurance", "postgres", "1234"

# --- ОБНОВЛЕННЫЕ SQL ЗАПРОСЫ ---

# 1. Сотрудники + строка "Все"
SQL_GET_EMPLOYEES = """
    SELECT 0 AS employee_id, 'Все' AS full_name
    UNION ALL
    SELECT employee_id, full_name FROM employees
"""

# 2. Клиенты (Фильтр по сотруднику или Все)
SQL_GET_POLICYHOLDERS_ALL = """
    SELECT '0' AS policy_number, 'Все' AS full_name
    UNION ALL
    SELECT policy_number, full_name FROM policyholders
"""
SQL_GET_POLICYHOLDERS_BY_EMP = """
    SELECT '0' AS policy_number, 'Все' AS full_name
    UNION ALL
    SELECT policy_number, full_name FROM policyholders WHERE employee_id = :emp_id
"""

# 3. Претензии (Разные уровни фильтрации)
SQL_CLAIMS_ALL = "SELECT claim_id, policy_number, description, event_date, payout FROM claims"

SQL_CLAIMS_BY_EMP = """
    SELECT c.claim_id, c.policy_number, c.description, c.event_date, c.payout 
    FROM claims c
    JOIN policyholders p ON c.policy_number = p.policy_number
    WHERE p.employee_id = :emp_id
"""

SQL_CLAIMS_BY_POLICY = """
    SELECT claim_id, policy_number, description, event_date, payout 
    FROM claims WHERE policy_number = :policy_num
"""

class DatabaseSection(QWidget):
    def __init__(self, title, db_connection_name):
        super().__init__()
        self.db_name = db_connection_name
        self.current_emp_id = 0 # По умолчанию "Все"
        self.init_ui(title)
        self.load_employees()

    def init_ui(self, title):
        layout = QVBoxLayout()
        font_title = QFont("Bahnschrift", 20)
        font_label = QFont("Bahnschrift", 12)

        layout.addWidget(QLabel(title, font=font_title))

        layout.addWidget(QLabel("Выберите сотрудника:", font=font_label))
        self.combo_employee = QComboBox()
        layout.addWidget(self.combo_employee)

        layout.addWidget(QLabel("Выберите клиента:", font=font_label))
        self.combo_policyholder = QComboBox()
        layout.addWidget(self.combo_policyholder)

        layout.addWidget(QLabel("Список страховых случаев (Claims):", font=font_label))
        self.claims_view = QTableView()
        self.claims_view.horizontalHeader().setSectionResizeMode(QHeaderView.ResizeMode.Stretch)
        layout.addWidget(self.claims_view)

        self.setLayout(layout)

        self.model_emp = QSqlQueryModel()
        self.model_holder = QSqlQueryModel()
        self.model_claims = QSqlQueryModel()

        self.combo_employee.setModel(self.model_emp)
        self.combo_policyholder.setModel(self.model_holder)
        self.claims_view.setModel(self.model_claims)

        # Сигналы
        self.combo_employee.currentIndexChanged.connect(self.on_employee_changed)
        self.combo_policyholder.currentIndexChanged.connect(self.on_policyholder_changed)

    def load_employees(self):
        db = QSqlDatabase.database(self.db_name)
        query = QSqlQuery(db)
        query.exec(SQL_GET_EMPLOYEES)
        self.model_emp.setQuery(query)
        self.combo_employee.setModelColumn(1)

    def on_employee_changed(self, index):
        if index < 0: return
        self.current_emp_id = self.model_emp.record(index).value("employee_id")
        
        db = QSqlDatabase.database(self.db_name)
        query = QSqlQuery(db)
        
        if self.current_emp_id == 0:
            query.exec(SQL_GET_POLICYHOLDERS_ALL)
        else:
            query.prepare(SQL_GET_POLICYHOLDERS_BY_EMP)
            query.bindValue(":emp_id", self.current_emp_id)
            query.exec()
        
        self.model_holder.setQuery(query)
        self.combo_policyholder.setModelColumn(1)
        # Сбрасываем выбор клиента на "Все" (индекс 0)
        self.combo_policyholder.setCurrentIndex(0)

    def on_policyholder_changed(self, index):
        if index < 0: return
        
        policy_num = str(self.model_holder.record(index).value("policy_number"))
        db = QSqlDatabase.database(self.db_name)
        query = QSqlQuery(db)
        
        if policy_num == "0": # Выбрано "Все" в списке клиентов
            if self.current_emp_id == 0: # И "Все" в списке сотрудников
                query.prepare(SQL_CLAIMS_ALL)
            else: # Все клиенты конкретного сотрудника
                query.prepare(SQL_CLAIMS_BY_EMP)
                query.bindValue(":emp_id", self.current_emp_id)
        else: # Конкретный клиент
            query.prepare(SQL_CLAIMS_BY_POLICY)
            query.bindValue(":policy_num", policy_num)
        
        query.exec()
        self.model_claims.setQuery(query)
        self._set_claims_headers()

    def _set_claims_headers(self):
        headers = ["ID случая", "№ Полиса", "Описание", "Дата", "Выплата"]
        for i, h in enumerate(headers):
            self.model_claims.setHeaderData(i, Qt.Orientation.Horizontal, h)

class MainWindow(QMainWindow):
    def __init__(self):
        super().__init__()
        self.setWindowTitle("Insurance Database Management (Filter: ALL)")
        self.resize(1100, 750)
        if not self.init_databases(): sys.exit(1)

        splitter = QSplitter(Qt.Orientation.Horizontal)
        splitter.addWidget(DatabaseSection("MySQL", "mysql_conn"))
        splitter.addWidget(DatabaseSection("PostgreSQL", "pg_conn"))
        self.setCentralWidget(splitter)

    def init_databases(self):
        db_ms = QSqlDatabase.addDatabase("QMARIADB", "mysql_conn")
        db_ms.setHostName(MYSQL_HOST); db_ms.setDatabaseName(MYSQL_DB)
        db_ms.setUserName(MYSQL_USER); db_ms.setPassword(MYSQL_PASS)

        db_pg = QSqlDatabase.addDatabase("QPSQL", "pg_conn")
        db_pg.setHostName(PG_HOST); db_pg.setDatabaseName(PG_DB)
        db_pg.setUserName(PG_USER); db_pg.setPassword(PG_PASS)

        if not db_ms.open() or not db_pg.open():
            QMessageBox.critical(self, "Error", "Connection failed")
            return False
        return True

if __name__ == "__main__":
    app = QApplication(sys.argv)
    app.setFont(QFont("Segoe UI", 10))
    window = MainWindow()
    window.show()
    sys.exit(app.exec())