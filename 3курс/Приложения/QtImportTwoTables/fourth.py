import sys
from PyQt6.QtWidgets import (QApplication, QMainWindow, QVBoxLayout, QHBoxLayout, 
                             QWidget, QLabel, QMessageBox, QSplitter, 
                             QComboBox, QScrollArea, QFrame)
from PyQt6.QtSql import QSqlDatabase, QSqlQueryModel, QSqlQuery
from PyQt6.QtCore import Qt
from PyQt6.QtGui import QFont

# Настройки БД
MYSQL_HOST, MYSQL_DB, MYSQL_USER, MYSQL_PASS = "127.0.0.1", "insurance", "root", "1234"
PG_HOST, PG_DB, PG_USER, PG_PASS = "127.0.0.1", "insurance", "postgres", "1234"

# SQL Запросы
SQL_EMPLOYEES = "SELECT employee_id, full_name FROM employees"

# Запрос: найти все страховые случаи по ID сотрудника (через связь с policyholders)
SQL_GET_CLAIMS_BY_EMP = """
    SELECT c.claim_id, p.full_name as client_name, c.description, c.event_date, c.payout, p.policy_number
    FROM claims c 
    JOIN policyholders p ON c.policy_number = p.policy_number 
    WHERE p.employee_id = :emp_id
"""

# --- АНАЛОГ USER CONTROL (Карточка страхового случая) ---
class ClaimCard(QFrame):
    def __init__(self, client_name, description, date, payout, policy_num):
        super().__init__()
        self.setFrameShape(QFrame.Shape.StyledPanel)
        # Стилизация карточки (CSS)
        self.setStyleSheet("""
            ClaimCard {
                background-color: #ffffff;
                border-radius: 10px;
                margin: 5px;
            }
            QLabel#ClientName { font-weight: bold; font-size: 16px; color: #2c3e50; }
            QLabel#Payout { font-weight: bold; color: #27ae60; font-size: 15px; }
            QLabel#Details { color: #7f8c8d; }
        """)

        layout = QVBoxLayout()
        
        lbl_client = QLabel(f"Клиент: {client_name}")
        lbl_client.setObjectName("ClientName")
        
        lbl_policy = QLabel(f"Полис: {policy_num}")
        lbl_policy.setObjectName("Details")

        lbl_desc = QLabel(f"Описание: {description}")
        lbl_desc.setWordWrap(True)
        
        lbl_date = QLabel(f"Дата: {date}")
        lbl_date.setObjectName("Details")

        lbl_payout = QLabel(f"Выплата: {payout} руб.")
        lbl_payout.setObjectName("Payout")

        layout.addWidget(lbl_client)
        layout.addWidget(lbl_policy)
        layout.addWidget(lbl_desc)
        layout.addWidget(lbl_date)
        layout.addSpacing(10)
        layout.addWidget(lbl_payout)
        
        self.setLayout(layout)

# --- ОСНОВНАЯ СЕКЦИЯ БД ---
class DatabaseSection(QWidget):
    def __init__(self, title, db_connection_name):
        super().__init__()
        self.db_name = db_connection_name
        self.init_ui(title)
        self.load_employees()

    def init_ui(self, title):
        layout = QVBoxLayout()
        
        # Заголовок
        lbl_title = QLabel(title)
        lbl_title.setFont(QFont("Bahnschrift", 20))
        layout.addWidget(lbl_title)

        # Выбор сотрудника (ComboBox)
        layout.addWidget(QLabel("Выберите сотрудника (Employee):"))
        self.employee_combo = QComboBox()
        self.employee_combo.setMinimumHeight(35)
        self.employee_combo.currentIndexChanged.connect(self.on_employee_selected)
        layout.addWidget(self.employee_combo)

        layout.addSpacing(20)
        layout.addWidget(QLabel("Страховые случаи (Claims) в виде карточек:"))

        # Область прокрутки для карточек
        self.scroll = QScrollArea()
        self.scroll.setWidgetResizable(True)
        self.scroll.setStyleSheet("QScrollArea { border: none; background: transparent; }")
        
        # Контейнер внутри ScrollArea
        self.cards_container = QWidget()
        self.cards_layout = QVBoxLayout(self.cards_container)
        self.cards_layout.setAlignment(Qt.AlignmentFlag.AlignTop)
        self.scroll.setWidget(self.cards_container)
        
        layout.addWidget(self.scroll)
        self.setLayout(layout)

        # Модель для комбобокса
        self.master_model = QSqlQueryModel()
        self.employee_combo.setModel(self.master_model)

    def load_employees(self):
        db = QSqlDatabase.database(self.db_name)
        query = QSqlQuery(db)
        query.exec(SQL_EMPLOYEES)
        self.master_model.setQuery(query)
        self.employee_combo.setModelColumn(1) # Показываем full_name

    def clear_cards(self):
        """ Удаляет все карточки из лэйаута """
        while self.cards_layout.count():
            item = self.cards_layout.takeAt(0)
            widget = item.widget()
            if widget:
                widget.deleteLater()

    def on_employee_selected(self, index):
        """ Загрузка данных и создание карточек """
        if index < 0: return
        self.clear_cards()

        # Получаем ID сотрудника
        employee_id = self.master_model.record(index).value("employee_id")
        
        db = QSqlDatabase.database(self.db_name)
        query = QSqlQuery(db)
        query.prepare(SQL_GET_CLAIMS_BY_EMP)
        query.bindValue(":emp_id", employee_id)
        query.exec()

        # Для каждой записи создаем "UserControl" (карточку)
        while query.next():
            card = ClaimCard(
                client_name=query.value("client_name"),
                description=query.value("description"),
                date=query.value("event_date"),
                payout=query.value("payout"),
                policy_num=query.value("policy_number")
            )
            self.cards_layout.addWidget(card)

# --- ГЛАВНОЕ ОКНО ---
class MainWindow(QMainWindow):
    def __init__(self):
        super().__init__()
        self.setWindowTitle("Insurance Cards")
        self.resize(1100, 850)

        if not self.init_databases():
            sys.exit(1)

        main_layout = QHBoxLayout()
        splitter = QSplitter(Qt.Orientation.Horizontal)

        self.mysql_section = DatabaseSection("MySQL", "mysql_conn")
        self.pg_section = DatabaseSection("PostgreSQL", "pg_conn")

        splitter.addWidget(self.mysql_section)
        splitter.addWidget(self.pg_section)
        
        self.setCentralWidget(splitter)

    def init_databases(self):
        db_ms = QSqlDatabase.addDatabase("QMARIADB", "mysql_conn")
        db_ms.setHostName(MYSQL_HOST); db_ms.setDatabaseName(MYSQL_DB)
        db_ms.setUserName(MYSQL_USER); db_ms.setPassword(MYSQL_PASS)

        db_pg = QSqlDatabase.addDatabase("QPSQL", "pg_conn")
        db_pg.setHostName(PG_HOST); db_pg.setDatabaseName(PG_DB)
        db_pg.setUserName(PG_USER); db_pg.setPassword(PG_PASS)

        if not db_ms.open():
            QMessageBox.critical(self, "MySQL Error", db_ms.lastError().text())
            return False
        if not db_pg.open():
            QMessageBox.critical(self, "PostgreSQL Error", db_pg.lastError().text())
            return False
        return True

if __name__ == "__main__":
    app = QApplication(sys.argv)
    app.setFont(QFont("Segoe UI", 10))
    window = MainWindow()
    window.show()
    sys.exit(app.exec())