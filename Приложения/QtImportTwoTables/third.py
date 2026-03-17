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

# --- SQL ЗАПРОСЫ ---
# 1. Загрузка сотрудников
SQL_GET_EMPLOYEES = "SELECT employee_id, full_name FROM employees"

# 2. Загрузка клиентов конкретного сотрудника
SQL_GET_POLICYHOLDERS = "SELECT policy_number, full_name FROM policyholders WHERE employee_id = :emp_id"

# 3. Загрузка претензий конкретного клиента
SQL_GET_CLAIMS = """
    SELECT claim_id, policy_number, description, event_date, payout 
    FROM claims 
    WHERE policy_number = :policy_num
"""

class DatabaseSection(QWidget):
    def __init__(self, title, db_connection_name):
        super().__init__()
        self.db_name = db_connection_name
        self.init_ui(title)
        self.load_employees()

    def init_ui(self, title):
        layout = QVBoxLayout()
        
        # Шрифты
        font_title = QFont("Bahnschrift", 20)
        font_label = QFont("Bahnschrift", 12)

        layout.addWidget(QLabel(title, font=font_title))

        # --- ВЫБОР СОТРУДНИКА ---
        layout.addWidget(QLabel("Выберите сотрудника:", font=font_label))
        self.combo_employee = QComboBox()
        self.combo_employee.currentIndexChanged.connect(self.on_employee_changed)
        layout.addWidget(self.combo_employee)

        # --- ВЫБОР КЛИЕНТА ---
        layout.addWidget(QLabel("Выберите застрахованное лицо (клиента):", font=font_label))
        self.combo_policyholder = QComboBox()
        self.combo_policyholder.currentIndexChanged.connect(self.on_policyholder_changed)
        layout.addWidget(self.combo_policyholder)

        # --- ТАБЛИЦА ПРЕТЕНЗИЙ ---
        layout.addWidget(QLabel("Список страховых случаев (Claims):", font=font_label))
        self.claims_view = QTableView()
        self.claims_view.horizontalHeader().setSectionResizeMode(QHeaderView.ResizeMode.Stretch)
        layout.addWidget(self.claims_view)

        self.setLayout(layout)

        # Модели данных
        self.model_emp = QSqlQueryModel()
        self.model_holder = QSqlQueryModel()
        self.model_claims = QSqlQueryModel()

        self.combo_employee.setModel(self.model_emp)
        self.combo_policyholder.setModel(self.model_holder)
        self.claims_view.setModel(self.model_claims)

    def load_employees(self):
        """Заполнение первого комбобокса (Сотрудники)"""
        db = QSqlDatabase.database(self.db_name)
        query = QSqlQuery(db)
        query.exec(SQL_GET_EMPLOYEES)
        self.model_emp.setQuery(query)
        self.combo_employee.setModelColumn(1) # Показываем full_name

    def on_employee_changed(self, index):
        """Событие: выбран сотрудник -> грузим его клиентов"""
        if index < 0: return
        
        # Получаем ID сотрудника из скрытой колонки 0
        emp_id = self.model_emp.record(index).value("employee_id")
        
        db = QSqlDatabase.database(self.db_name)
        query = QSqlQuery(db)
        query.prepare(SQL_GET_POLICYHOLDERS)
        query.bindValue(":emp_id", emp_id)
        query.exec()
        
        self.model_holder.setQuery(query)
        self.combo_policyholder.setModelColumn(1) # Показываем full_name клиента

    def on_policyholder_changed(self, index):
        """Событие: выбран клиент -> грузим его страховые случаи"""
        if index < 0:
            self.model_claims.clear()
            return
        
        # Получаем policy_number клиента из скрытой колонки 0
        policy_num = self.model_holder.record(index).value("policy_number")
        
        db = QSqlDatabase.database(self.db_name)
        query = QSqlQuery(db)
        query.prepare(SQL_GET_CLAIMS)
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
        self.setWindowTitle("Insurance Database Management")
        self.resize(1000, 700)

        if not self.init_databases():
            sys.exit(1)

        splitter = QSplitter(Qt.Orientation.Horizontal)
        splitter.addWidget(DatabaseSection("MySQL", "mysql_conn"))
        splitter.addWidget(DatabaseSection("PostgreSQL", "pg_conn"))

        self.setCentralWidget(splitter)

    def init_databases(self):
        # MySQL
        db_ms = QSqlDatabase.addDatabase("QMARIADB", "mysql_conn")
        db_ms.setHostName(MYSQL_HOST)
        db_ms.setDatabaseName(MYSQL_DB)
        db_ms.setUserName(MYSQL_USER)
        db_ms.setPassword(MYSQL_PASS)

        # PostgreSQL
        db_pg = QSqlDatabase.addDatabase("QPSQL", "pg_conn")
        db_pg.setHostName(PG_HOST)
        db_pg.setDatabaseName(PG_DB)
        db_pg.setUserName(PG_USER)
        db_pg.setPassword(PG_PASS)

        if not db_ms.open():
            QMessageBox.critical(self, "Error", f"MySQL fail: {db_ms.lastError().text()}")
            return False
        if not db_pg.open():
            QMessageBox.critical(self, "Error", f"PostgreSQL fail: {db_pg.lastError().text()}")
            return False
        return True

if __name__ == "__main__":
    app = QApplication(sys.argv)
    app.setFont(QFont("Segoe UI", 10))
    window = MainWindow()
    window.show()
    sys.exit(app.exec())